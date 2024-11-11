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
	go ReadInput()
}

func ReadInput() {
	for {
		os.Stdin.Read(buf)
		Log().String("input", string(buf)).String("asBytes", fmt.Sprintf("%s", buf)).Msg("Received ")
		DebugWriteString(fmt.Sprintf("Hello World: %s", buf))
	}
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
