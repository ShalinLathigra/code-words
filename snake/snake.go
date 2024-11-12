package snake

import (
	"fmt"
	"slices"

	"sl.com/log"
	"sl.com/math"
	"sl.com/render"
	"sl.com/terminal"
)

func Log() *log.LogBuilder { return log.CreateLogger("game") }

type TileState int

const (
	Empty     byte = ' '
	SnakeHead byte = '+'
	SnakeBody byte = '#'
	Apple     byte = '&'
)

type Grid struct {
	math.Rect
	tiles [][]byte
}

type SnakeCell struct {
	math.Vec2
	prev *SnakeCell
	next  *SnakeCell
}

func (c *SnakeCell) clone () *SnakeCell {
	return &SnakeCell{Vec2: c.Vec2}
}

type Snake struct {
	head *SnakeCell
	tail *SnakeCell
	current *SnakeCell
	dir  math.Vec2
}

type GridFrame struct {
	frame *render.BufferedFrame // Similar Idea, Just have it contain a buffered node reference. then during the game we just trigger the Buffer update trhough this link
	Grid
}

func (gf *GridFrame) UpdateBuffer() {
	bufIndex := 0
	for y := 0; y < gf.Size.Y; y++ {
		bufIndex = (y+gf.Y)*gf.frame.Aabb.Size.X + gf.X
		copy(gf.frame.Buffer[bufIndex:bufIndex+gf.Size.X], grid.tiles[y])
	}
}

func NewGridFrame(rect math.Rect) GridFrame {
	if rect.Size.X <= 0 || rect.Size.Y <= 0 {
		panic("instancing a grid with no size")
	}
	gridRect, ok := math.Shrink(rect, math.One)
	if !ok {
		panic(fmt.Errorf("failed to shrink rect: %s", rect))
	}
	gridRect.Vec2 = math.Vec2{X: 1, Y: 1}
	grid := Grid{
		Rect: gridRect,
	}
	for range grid.Size.Y {
		grid.tiles = append(grid.tiles, slices.Repeat([]byte{Empty}, grid.Size.X))
	}
	gf := GridFrame{
		Grid:  grid,
		frame: render.CreatePatternedBufferedFrame("Grid", rect),
	}
	render.Root.AddChild(gf.frame)
	Log().String("frame size", gf.frame.Aabb.String()).String("grid size", grid.Size.String()).Msg("Instancing Grid")
	return gf
}

func (g Grid) String() string {
	return fmt.Sprintf("Grid: XY(%s) WH(%s)", g.Vec2, g.Size)
}

type GameState interface {
	Tick()
}

func CreateSnake(grid *Grid, pos math.Vec2, len uint) Snake {
	if !math.VecContains(grid.Size, pos) {
		panic(fmt.Errorf("position outside of bounds %s %s", pos, grid.Size))
	}
	if len < 1 {
		panic(fmt.Errorf("snake must have positive length: %d", len))
	}
	head := &SnakeCell{
		Vec2: pos,
	}
	tail := head
	snake := Snake{
		head: head,
		tail: tail,
		dir:  math.Zero,
	}
	return snake
}

func die() {
	panic("death not implemented")
}

func setTile(grid *Grid, at math.Vec2, state byte) {
	if !math.VecContains(grid.Size, at) {
		panic(fmt.Sprintf("setting out-of-bounds index: %s in %s", at, grid))
	}
	grid.tiles[at.Y][at.X] = state
}

var grid GridFrame
var snake Snake

func Init(insets math.Vec2) {
	rect, ok := math.Shrink(math.Rect{Size: render.Root.Aabb.Size}, insets)
	// rect should be odd on x axis
	if rect.Size.X % 2 == 0 {
		rect.Size.X += 1
	}
	if !ok {
		panic("game area too small")
	}
	grid = NewGridFrame(rect)
	snake = CreateSnake(&grid.Grid, math.Vec2{X: rect.Size.X / 2, Y: rect.Size.Y / 2}, 1)
	setTile(&grid.Grid, math.Vec2{4*1,3}, Apple)
	setTile(&grid.Grid, math.Vec2{4*2, 3}, Apple)
	setTile(&grid.Grid, math.Vec2{4*3, 3}, Apple)
	setTile(&grid.Grid, math.Vec2{4*4, 3}, Apple)
	setTile(&grid.Grid, math.Vec2{4*5, 3}, Apple)
	setTile(&grid.Grid, math.Vec2{4*6, 3}, Apple)
	Log().String("grid", grid.String()).String("snake.head", snake.head.String()).Msg("Game Start")
}

func Tick(_ int64, frame int) bool {
	lastPos := snake.head.Vec2
	snake.head.Vec2 = math.Add(snake.head.Vec2, snake.dir)
	if !math.VecContains(grid.Size, snake.head.Vec2) {
		// If here, die
		terminal.DebugWriteString("Out of bounds")
		die()
	}
	switch grid.tiles[snake.head.Y][snake.head.X] {
	case Apple:
		extendSnake(&snake)
	case SnakeBody:
		terminal.DebugWriteString(fmt.Sprintf("Hit yourself at %s", snake.head.Vec2))
	default:
	}
	// update the grid display
	if snake.current == nil {
		setTile(&grid.Grid, lastPos, Empty)
	} else {
		// need to move snake.current to lastPos
		setTile(&grid.Grid, snake.current.Vec2, Empty)
		setTile(&grid.Grid, lastPos, SnakeBody)
		snake.current.Vec2 = lastPos
		snake.current = snake.current.next
		if snake.current == nil {
			snake.current = snake.head.next
		}
	}
	setTile(&grid.Grid, snake.head.Vec2, SnakeHead)

	grid.UpdateBuffer()
	return true
}

func extendSnake(s *Snake) {
	newCell := s.tail.clone()
	s.tail.next = newCell
	newCell.prev = s.tail
	s.tail = newCell
	if s.current == nil {
		s.current = s.tail
	}
}

func HandleInputByte(b byte) {
	switch b {
	case 65:
		snake.dir = math.Up
	case 66:
		snake.dir = math.Down
	case 67:
		snake.dir = math.Scale(math.Right, 2)
	case 68:
		snake.dir = math.Scale(math.Left, 2)
	}
}
