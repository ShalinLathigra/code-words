package terminal

import (
	"fmt"
	"os"
	"slices"

	"golang.org/x/term"
	"sl.com/log"
	"sl.com/math"
)

func Log() *log.LogBuilder { return log.CreateLogger("terminal") }

var termFd int
var oldState *term.State
var Size math.Vec2
var terminal *term.Terminal
var clearDebugLineBuffer []byte
var buf = make([]byte, 3)

func Init() {
	var err error
	termFd = int(os.Stdin.Fd())
	if Size.X, Size.Y, err = term.GetSize(termFd); err != nil {
		panic(err)
	}
	for range Size.Y {
		fmt.Println("")
	}
	if oldState, err = term.MakeRaw(termFd); err != nil {
		panic(err)
	}
	clearDebugLineBuffer = slices.Repeat([]byte{' '}, Size.X)
	terminal = term.NewTerminal(os.Stdout, "")
	go readInputContinuous()
}

const (
	ESCAPE byte = 27
	ARROW  byte = 91
	UP     byte = 65
	DOWN   byte = 66
	RIGHT  byte = 67
	LEFT   byte = 68
)

var inputChan = make(chan byte)

func ReadInput() (byte, bool) {
	select {
	case input := <-inputChan:
		return input, true
	default:
		return 0, false
	}
}
func readInputContinuous() {
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}
		if n == 3 {
			if buf[0] != ESCAPE || buf[1] != ARROW {
				continue
			}
			logArrowInput(buf[2])
			inputChan <- buf[2]
		}
	}
}

func logArrowInput(input byte) {
	if input < UP || input > LEFT {
		panic(fmt.Errorf("logging an invalid arrow input %d (%c)", input, input))
	}
	debugStr := fmt.Sprintf("%d %d %d", ESCAPE, ARROW, input)
	direction := ""
	switch input {
	case UP:
		direction = "UP"
	case DOWN:
		direction = "DOWN"
	case RIGHT:
		direction = "RIGHT"
	case LEFT:
		direction = "LEFT"
	}
	Log().String("bytes", debugStr).String("dir", direction).Msg("Received Input")
	DebugWriteString(fmt.Sprintf("Input: %s Dir: %s", debugStr, direction))
}

func Clean() {
	// switch back from raw mode
	if err := term.Restore(termFd, oldState); err != nil {
		panic(err)
	}
}

func WriteBytes(bytes []byte) {
	// write bytes
	if _, err := terminal.Write(bytes); err != nil {
		panic(err)
	}
}

func CursorTo(x int, y int) {
	WriteBytes([]byte(fmt.Sprintf("\033[%d;%dH", y+1, x+1)))
}

func DebugWriteString(msg string) {
	CursorTo(0, Size.Y)
	endBuf := slices.Clone(clearDebugLineBuffer)
	bufLen := math.Min(len(msg), len(endBuf))
	copy(endBuf[0:bufLen], msg)
	if _, err := terminal.Write(endBuf); err != nil {
		panic(err)
	}
}
