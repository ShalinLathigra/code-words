package render

import (
	"os"

	"golang.org/x/term"
	"sl.com/math"
)

var termFd int
var Size math.Vec2

func Init() {
	termFd = int(os.Stdin.Fd())
	if w, h, err := term.GetSize(termFd); err != nil {
		panic(err)
	} else {
		Size = math.Vec2{X: w, Y: h}
	}
}
