package render

import (
	"fmt"
	"slices"

	"sl.com/log"
	"sl.com/math"
)

/*
		basically, we assemble a tree starting with some arbitrary Root node with some size and width
		We attach additional nodes as Children of this one to define regions of the screen
	 		render order is determined by order in the list of Children
			Can have an arbitrary amount of Children
			if transparent, then we need to step through every byte and only copy over the non-zero ones

		Render is called to start a chain of RenderOntos by all the Children
		Tree Root generates its buffer first, then for every child node we call RenderOnto(buffer, source rect)
			which recursively draws all of the Children in

		Expected scene layout is going to be:

		Root:
			Score Display
			GameArea
				(Special case node implemented in the game logic that will handle creating and writing that whole thing)
			End Game Banner
*/
func Log() *log.LogBuilder { return log.CreateLogger("render") }

type Node interface {
	GetParent() Node
	GetChildren() []Node
	Name() string
	HasChild(Node) int // Check if part
	AddChild(Node)
	RemoveChild(Node) // Remove child, if no action needed return false
	SetParent(Node)
	RenderOnto(*[]byte, math.Rect)
}

type LinkedNode struct {
	Children []Node
	Parent   Node
	Name     string
}

type BufferedFrame struct {
	LinkedNode
	Buffer []byte
	Aabb   math.Rect
}

func (f *BufferedFrame) GetParent() Node {
	return f.Parent
}

func (f *BufferedFrame) GetChildren() []Node {
	return f.Children
}

func (f *BufferedFrame) HasChild(n Node) int {
	for i, child := range f.Children {
		if child == n {
			return i
		}
	}
	return -1
}

func (f *BufferedFrame) AddChild(n Node) {
	if index := f.HasChild(n); index >= 0 || n.GetParent() == f {
		return
	}
	f.Children = append(f.Children, n)
	if p := n.GetParent(); p != nil {
		p.RemoveChild(n)
	}
	n.SetParent(f)
	f.Children = append(f.Children, n)
}

func (f *BufferedFrame) RemoveChild(n Node) {
	index := f.HasChild(n)
	if index >= 0 {
		f.Children = slices.Delete(f.Children, index, index+1)
	}
	if n.GetParent() == f {
		n.SetParent(nil)
	}
}

func (f *BufferedFrame) SetParent(n Node) {
	if f.GetParent() == n {
		return
	}
	f.Parent = n
}

func (f *BufferedFrame) RenderOnto(buffer *[]byte, rect math.Rect) {
	inter, ok := math.Overlap(rect, f.Aabb)
	if !ok {
		Log().Msg("Failed to overlap")
		return
	}
	fOffset := math.Subtract(f.Aabb.Vec2, inter.Vec2)
	for y := 0; y < inter.Size.Y; y++ {
		bufStartIndex := (inter.Y+y)*rect.Size.X + inter.X
		fStartIndex := (y+fOffset.Y)*f.Aabb.Size.X + fOffset.X
		copy((*buffer)[bufStartIndex:bufStartIndex+inter.Size.X], f.Buffer[fStartIndex:fStartIndex+inter.Size.X])
	}
	for _, child := range f.Children {
		child.RenderOnto(buffer, math.Rect{Vec2: fOffset, Size: inter.Size})
	}
}

func (f *BufferedFrame) Name() string {
	return f.LinkedNode.Name
}

func CreatePatternedBufferedFrame(name string, aabb math.Rect) *BufferedFrame {
	return &BufferedFrame{
		LinkedNode: LinkedNode{
			Children: make([]Node, 0, 8),
			Name:     name,
		},
		Buffer: makePatternedFrameBuffer(aabb.Size),
		Aabb:   aabb,
	}
}

func makePatternedFrameBuffer(size math.Vec2) []byte {
	if size.X*size.Y == 0 {
		panic("cannot init frame with zero area")
	}
	ret := slices.Repeat([]byte{' '}, size.X*size.Y)
	// What are the steps?
	i := 0
	for y := range size.Y {
		for x := range size.X {
			i = x + size.X*y
			if x == 0 || x == size.X-1 {
				if y == 0 || y == size.Y-1 || y%3 == 0 {
					ret[i] = '+'
				} else {
					ret[i] = '|'
				}
			} else if y == 0 {
				if x%5 == 0 {
					ret[i] = '+'
				} else {
					ret[i] = '-'
				}
			} else if y == size.Y-1 {
				if x%5 == 0 {
					ret[i] = '+'
				} else {
					ret[i] = '='
				}
			}
		}
	}
	return ret
}

func render() ([]byte, math.Rect) {
	if Root == nil {
		panic("cannot render with empty Root")
	}
	copy(dirty, Root.Buffer)
	for _, child := range Root.Children {
		child.RenderOnto(&dirty, Root.Aabb)
	}
	return dirty, Root.Aabb
}

var Root *BufferedFrame
var dirty []byte

func Init() {
	InitTerminal()
	Root = CreatePatternedBufferedFrame("Root", math.Rect{Size: math.Subtract(Size, math.Down)})
	dirty = make([]byte, Root.Aabb.Size.X*Root.Aabb.Size.Y)
	Log().String("Root", Root.Aabb.String()).Msg("Init root render frame")
}

func Clean() {
	CleanTerminal()
	fmt.Println("")
}

func ReDraw() {
	// clear()
	buf, box := render()
	for y := 0; y < box.Size.Y; y++ {
		cursorTo(0, y)
		WriteBytes(buf[y*box.Size.X : (y+1)*box.Size.X])
	}
}
