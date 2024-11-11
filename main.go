package main

import (
	"fmt"
	"time"

	"sl.com/log"
	"sl.com/render"
	game "sl.com/snake"
)

func Log() *log.LogBuilder { return log.CreateLogger("main") }

var frame int = 0

func main() {
	defer render.Clean()
	defer log.Clean()
	var delta int64
	lastTime := time.Now()
	render.Init()
	game.Init(6)

	// next step, add contextual printing. Create "BufferedFrame" and "GridFrame", BufferedFrame is just a region and a pre-programmed byte Buffer, GridFrame is a region in which
	// Original frame can be set to use the fancy background and such
	// GridFrame is implemented here, takes a reference to a Grid, and generates a Buffer
	// Frames use the same rect math to determine when/what to render
	// 	First in, last drawn
	// grid := NewGrid(math.Add(render.Size, math.Up))

	for {
		delta = time.Since(lastTime).Milliseconds()
		if delta < 1000/24 {
			continue
		}
		lastTime = time.Now()
		if !game.Tick() {
			Log().Msg("Exiting Game")
			return
		}
		tickDuration := time.Since(lastTime).Microseconds()
		Log().Any("delta", delta).Any("us", tickDuration).Msg("Tick Duration")
		render.ReDraw()
		render.DebugWriteString(fmt.Sprintf("Hello World: %d", frame))
		frame = frame + 1
		Log().Any("delta", delta).Any("us", time.Since(lastTime).Microseconds()-tickDuration).Msg("Render Duration")
	}
}
