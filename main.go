package main

import (
	"time"

	"sl.com/log"
	"sl.com/math"
	"sl.com/render"
	"sl.com/server"
	game "sl.com/snake"
	"sl.com/terminal"
)

func Log() *log.LogBuilder { return log.CreateLogger("main") }

var frame int = 0

const (
	FPS int64 = 12
)

func main() {
	server.Test()
	return
	var delta int64
	defer render.Clean()
	defer terminal.Clean()
	defer log.Clean()
	Log().Msg("Initializing game")
	lastTime := time.Now()
	terminal.Init() // must happen first
	render.Init()
	game.Init(math.Vec2{12, 3}, 8, FPS)

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
		if !game.Tick(delta, frame) {
			Log().Msg("Exiting Game")
			return
		}
		tickDuration := time.Since(lastTime).Microseconds()
		Log().Any("delta", delta).Any("us", tickDuration).Msg("Tick Duration")
		render.ReDraw()
		frame = frame + 1
		Log().Any("delta", delta).Any("us", time.Since(lastTime).Microseconds()-tickDuration).Msg("Render Duration")
	}
}
