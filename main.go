package main

import (
	"fmt"
	"slices"
	"time"

	"sl.com/math"
	"sl.com/render"
)

type TileState int

const (
	Empty TileState = iota
	Border
	SnakeHead
	SnakeBody
	Apple
)

type Grid struct {
	math.Rect
	tiles [][]TileState
}

type SnakeCell struct {
	math.Vec2
	parent *SnakeCell
	child  *SnakeCell
}

type Snake struct {
	head *SnakeCell
	tail *SnakeCell
	dir  math.Vec2
}

func NewGrid(size math.Vec2) (grid Grid) {
	if size.X <= 0 || size.Y <= 0 {
		panic("instancing a grid with no size")
	}
	grid.Size = size
	for y := range size.Y {
		var line []TileState
		if y == 0 || y == size.Y-1 {
			line = slices.Repeat([]TileState{Border}, size.X)
		} else {
			line = slices.Repeat([]TileState{Empty}, size.X)
			line[0], line[size.X-1] = Border, Border
		}
		grid.tiles = append(grid.tiles, line)
	}
	return
}

func (g Grid) Print() {
	fmt.Printf("Grid: (%s)\n", g.Vec2)
}

type GameState interface {
	Tick()
}

type PlayData struct {
	grid  Grid
	snake Snake
}

func CreateSnake(grid Grid, pos math.Vec2, len uint) Snake {
	if !math.Contains(grid.Rect, pos) {
		panic(fmt.Errorf("position outside of bounds %s %s", pos, grid.Rect))
	}
	if len < 1 {
		panic(fmt.Errorf("snake must have positive length: %d", len))
	}
	setTile(grid, pos, SnakeHead)
	head := &SnakeCell{
		math.Vec2{
			render.Size.X / 2, render.Size.Y / 2,
		},
		nil, nil,
	}
	tail := head
	if len > 0 {
		fmt.Println("variable length snake not supported")
	}
	snake := Snake{
		head: head,
		tail: tail,
		dir:  math.Left,
	}
	return snake
}

func main() {
	var delta int64
	lastTime := time.Now()
	render.Init()

	// next step, add contextual printing. Create "Frame" and "GridFrame", Frame is just a region and a pre-programmed byte buffer, GridFrame is a region in which
	// Original frame can be set to use the fancy background and such
	// GridFrame is implemented here, takes a reference to a Grid, and generates a buffer
	// Frames use the same rect math to determine when/what to render
	// 	First in, last drawn
	grid := NewGrid(math.Add(render.Size, math.Up))
	grid.Print()
	snake := CreateSnake(grid, math.Vec2{grid.Size.X / 2, grid.Size.Y / 2}, 2)

	for {
		delta = time.Since(lastTime).Milliseconds()
		if delta < 1000/24 {
			continue
		}
		lastTime = time.Now()

		// Update grid
		for cell := snake.tail; cell != nil; cell = cell.parent {
			grid.tiles[cell.Y][cell.X] = Empty
			if cell.parent == nil {
				cell.Vec2 = math.Add(cell.Vec2, snake.dir)
				if !math.Contains(grid.Rect, cell.Vec2) {
					// If here, die
					die()
				}
				setTile(grid, cell.Vec2, SnakeHead)
			} else {
				cell.X, cell.Y = cell.parent.X, cell.parent.Y
				setTile(grid, cell.Vec2, SnakeBody)
			}
		}

		// Print the grid
		fmt.Println("Snake at", snake.head.Vec2)
		// go debugPrintGrid(grid)
	}
}

func die() {
	panic("death not implemented")
}

func setTile(grid Grid, at math.Vec2, state TileState) {
	if !math.Contains(grid.Rect, at) {
		panic(fmt.Sprintf("setting out-of-bounds index: %s in %s", at, grid))
	}
	// do not override the head
	if grid.tiles[at.Y][at.X] == SnakeHead {
		return
	}
	grid.tiles[at.Y][at.X] = state
}

func debugPrintGrid(grid Grid) {
	debug := ""
	for j, line := range grid.tiles {
		for i, col := range line {
			debug = fmt.Sprintf("%s%c", debug, tileToRune(i, j, col))
		}
		debug = fmt.Sprintf("%s\n", debug)
	}
	fmt.Printf("%s", debug)
}

func tileToRune(x int, y int, t TileState) rune {
	switch t {
	case Empty:
		return ' '
	case Border:
		if (x+y)%2 == 0 {
			return '+'
		}
		return '='
	case SnakeHead:
		return '+'
	case SnakeBody:
		return '#'
	default:
		return '?'
	}
}
