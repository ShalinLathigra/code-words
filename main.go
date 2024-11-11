package main

import (
	"time"

	"sl.com/log"
	"sl.com/math"
	"sl.com/render"
	game "sl.com/snake"
	"sl.com/terminal"
)

func Log() *log.LogBuilder { return log.CreateLogger("main") }

var frame int = 0

const (
	FPS int64 = 20
)

func main() {
	defer render.Clean()
	defer terminal.Clean()
	defer log.Clean()
	var delta int64
	lastTime := time.Now()
	terminal.Init() // must happen first
	render.Init()
	game.Init(math.Vec2{24, 3})

	// next step, add contextual printing. Create "BufferedFrame" and "GridFrame", BufferedFrame is just a region and a pre-programmed byte Buffer, GridFrame is a region in which
	// Original frame can be set to use the fancy background and such
	// GridFrame is implemented here, takes a reference to a Grid, and generates a Buffer
	// Frames use the same rect math to determine when/what to render
	// 	First in, last drawn

	for {
		delta = time.Since(lastTime).Milliseconds()
		if delta < 1000/FPS {
			continue
		}
		lastTime = time.Now()
		input, ok := terminal.ReadInput()
		if ok {
			game.HandleInputByte(input)
		}
		if !game.Tick() {
			Log().Msg("Exiting Game")
			return
		}
		tickDuration := time.Since(lastTime).Microseconds()
		Log().Any("delta", delta).Any("us", tickDuration).Msg("Tick Duration")
		render.ReDraw()
		// terminal.DebugWriteString(fmt.Sprintf("Hello World: %d", frame))
		frame = frame + 1
		Log().Any("delta", delta).Any("us", time.Since(lastTime).Microseconds()-tickDuration).Msg("Render Duration")
	}
}
