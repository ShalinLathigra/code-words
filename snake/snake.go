package snake

import (
	"fmt"
	"slices"

	"sl.com/log"
	"sl.com/math"
	"sl.com/render"
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
	parent *SnakeCell
	// child  *SnakeCell
}

type Snake struct {
	head *SnakeCell
	tail *SnakeCell
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
		Vec2: math.Vec2{
			X: grid.Size.X / 2, Y: grid.Size.Y / 2,
		},
	}
	tail := head
	snake := Snake{
		head: head,
		tail: tail,
		dir:  math.Left,
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
	// do not override the head
	if grid.tiles[at.Y][at.X] == SnakeHead {
		return
	}
	grid.tiles[at.Y][at.X] = state
}

var grid GridFrame
var snake Snake

func Init(insets math.Vec2) {
	rect, ok := math.Shrink(math.Rect{Size: render.Root.Aabb.Size}, insets)
	if !ok {
		panic("game area too small")
	}
	grid = NewGridFrame(rect)
	snake = CreateSnake(&grid.Grid, math.Vec2{X: grid.Size.X / 2, Y: grid.Size.Y / 2}, 1)
	Log().String("grid", grid.String()).String("snake.head", snake.head.String()).Msg("Game Start")
}

func Tick() bool {
	// UpdateBuffer grid
	for cell := snake.tail; cell != nil; cell = cell.parent {
		grid.tiles[cell.Y][cell.X] = Empty
		if cell.parent == nil {
			cell.Vec2 = math.Add(cell.Vec2, snake.dir)
			if !math.VecContains(grid.Size, cell.Vec2) {
				// If here, die
				die()
			}
			setTile(&grid.Grid, cell.Vec2, SnakeHead)
		} else {
			cell.X, cell.Y = cell.parent.X, cell.parent.Y
			setTile(&grid.Grid, cell.Vec2, SnakeBody)
		}
	}
	grid.UpdateBuffer()
	return true
}

func HandleInputByte(b byte) {
	switch b {
	case 65:
		snake.dir = math.Up
	case 66:
		snake.dir = math.Down
	case 67:
		snake.dir = math.Right
	case 68:
		snake.dir = math.Left
	}
}
