package snake

import (
	"fmt"
	"math/rand"
	"slices"

	"sl.com/log"
	"sl.com/math"
	"sl.com/render"
	"sl.com/terminal"
)

func Log() *log.LogBuilder { return log.CreateLogger("game") }

type TileState int

const (
	Empty              byte = ' '
	SnakeHead          byte = '+'
	SnakeBody          byte = '#'
	Apple              byte = '&'
	SectorW            int  = 16
	SectorH            int  = 4
	MaxApplesPerSecond int  = 2
)

type Grid struct {
	math.Rect
	tiles      [][]byte
	sectors    []Sector
	numSectors math.Vec2
}

type Sector struct {
	math.Rect
	occupancy int
}

type SnakeCell struct {
	math.Vec2
	prev *SnakeCell
	next *SnakeCell
}

func (c *SnakeCell) clone() *SnakeCell {
	return &SnakeCell{Vec2: c.Vec2}
}

type Snake struct {
	head    *SnakeCell
	tail    *SnakeCell
	current *SnakeCell
	dir     math.Vec2
}

type GridFrame struct {
	frame *render.BufferedFrame // Similar Idea, Just have it contain a buffered node reference. then during the game we just trigger the Buffer update trhough this link
	Grid
}

func (gf *GridFrame) updateBuffer() {
	bufIndex := 0
	for y := 0; y < gf.Size.Y; y++ {
		bufIndex = (y+gf.Y)*gf.frame.Aabb.Size.X + gf.X
		copy(gf.frame.Buffer[bufIndex:bufIndex+gf.Size.X], grid.tiles[y])
	}
}

func newGridFrame(rect math.Rect) GridFrame {
	if rect.Size.X <= 0 || rect.Size.Y <= 0 {
		panic("instancing a grid with no size")
	}
	gridRect, ok := math.Shrink(rect, math.One)
	if !ok {
		panic(fmt.Errorf("failed to shrink rect: %s", rect))
	}
	gridRect.Vec2 = math.Vec2{X: 1, Y: 1}
	numSectors := math.Vec2{X: gridRect.Size.X/SectorW + 1, Y: gridRect.Size.Y/SectorH + 1}
	sectors := make([]Sector, numSectors.X*numSectors.Y)
	for y := range numSectors.Y {
		for x := range numSectors.X {
			newSector := Sector{
				Rect: math.Rect{
					Vec2: math.Vec2{X: x * SectorW, Y: y * SectorH},
					Size: math.Vec2{
						X: math.Min(SectorW, gridRect.Size.X-SectorW*x),
						Y: math.Min(SectorH, gridRect.Size.Y-SectorH*y),
					},
				},
			}
			if newSector.Size.X*newSector.Size.Y <= 0 {
				panic("bad sector")
			}
			sectors[x+y*numSectors.X] = newSector
		}
	}
	grid := Grid{
		Rect:       gridRect,
		sectors:    sectors,
		numSectors: numSectors,
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

func createSnake(grid *Grid, pos math.Vec2, len uint) Snake {
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

func setTile(grid *Grid, at math.Vec2, state byte) (ret bool) {
	if !math.VecContains(grid.Size, at) {
		panic(fmt.Sprintf("setting out-of-bounds index: %s in %s", at, grid))
	}
	ret = grid.tiles[at.Y][at.X] != state
	grid.tiles[at.Y][at.X] = state
	sectorIndex := at.X/SectorW + at.Y/SectorH*grid.numSectors.X
	if state == Empty {
		grid.sectors[sectorIndex].occupancy -= 1
	} else {
		grid.sectors[sectorIndex].occupancy += 1
	}
	return
}

var grid GridFrame
var snake Snake
var expectedNumApples int
var numAddedApples int
var fps int64

func Init(insets math.Vec2, numApples int, FPS int64) {
	rect, ok := math.Shrink(math.Rect{Size: render.Root.Aabb.Size}, insets)
	// rect should be odd on x axis
	if rect.Size.X%2 == 0 {
		rect.Size.X += 1
	}
	if !ok {
		panic("game area too small")
	}
	if FPS <= 0 {
		panic("framerate must be >= 0")
	}
	grid = newGridFrame(rect)
	snake = createSnake(&grid.Grid, math.Vec2{X: rect.Size.X/2 - (rect.Size.X/2)%2, Y: rect.Size.Y/2 - (rect.Size.Y/2)%2}, 1)
	expectedNumApples = numApples
	fps = FPS
	Log().String("grid", grid.String()).String("snake.head", snake.head.String()).Msg("Game Start")
}

func findRandomEmptyPoints(area *Grid, numPoints int) (indices []math.Vec2) {
	indices = make([]math.Vec2, 0, numPoints)
	// Basically, if we pick one that's at above 50% occupancy, grab another sector and take the lowest one
	// If we still fail to find a working sector, then we wait till next frame and try again
	for range numPoints {
		// pick a random sector
		sectorIndex := rand.Int() % (area.numSectors.X * area.numSectors.Y)
		terminal.DebugWriteString(fmt.Sprintf("processing point: %d/%d - %s", sectorIndex, len(area.sectors), area.sectors[sectorIndex].Size))
		point := math.PointWithin(area.sectors[sectorIndex].Rect)
		point.X = math.Clamp(point.X, area.sectors[sectorIndex].X, area.sectors[sectorIndex].X+area.sectors[sectorIndex].Size.X-2)
		if point.X%2 == 1 {
			point.X -= 1
		}
		indices = append(indices, point)
	}
	return indices
}

func Tick(_ int64, frame int) bool {
	// Handle spawning apples
	if int64(frame)%fps == 0 {
		for _, point := range findRandomEmptyPoints(&grid.Grid, math.Min(MaxApplesPerSecond, expectedNumApples-numAddedApples)) {
			if grid.Grid.tiles[point.Y][point.Y] == Empty && setTile(&grid.Grid, point, Apple) {
				numAddedApples += 1
			}
		}
	}

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
		numAddedApples -= 1
	case SnakeBody:
		terminal.DebugWriteString(fmt.Sprintf("Hit yourself at %s", snake.head.Vec2))
		die()
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

	grid.updateBuffer()
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
