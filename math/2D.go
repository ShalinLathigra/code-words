package math

import (
	"fmt"
)

type Vec2 struct {
	X int
	Y int
}

var Left = Vec2{-1, 0}
var Right = Vec2{1, 0}
var Up = Vec2{0, -1}
var Down = Vec2{0, 1}
var One = Vec2{1, 1}

func (v Vec2) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func Add(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

func Subtract(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X - b.X, a.Y - b.Y}
}

func Scale(a Vec2, b int) Vec2 {
	return Vec2{a.X * b, a.Y * b}
}

func Dot(a Vec2, b Vec2) int {
	return a.X*b.X + a.Y*b.Y
}

type Rect struct {
	Vec2
	Size Vec2
}

func (r Rect) String() string {
	return fmt.Sprintf("%s*%s", r.Vec2, r.Size)
}

func Shrink(rect Rect, skinWidth Vec2) (ret Rect, ok bool) {
	if skinWidth.X*2 >= rect.Size.X || skinWidth.Y*2 >= rect.Size.Y {
		return Rect{}, false
	}
	ok = true
	ret.Vec2 = Add(rect.Vec2, skinWidth)
	ret.Size = Add(rect.Size, Scale(skinWidth, -2))
	return
}

func Contains(rect Rect, point Vec2) bool {
	min, max := rect.Vec2, Add(rect.Vec2, rect.Size)
	return point.X >= min.X && point.X < max.X && point.Y >= min.Y && point.Y < max.Y
}

func VecContains(area Vec2, point Vec2) bool {
	return point.X >= 0 && point.X < area.X && point.Y >= 0 && point.Y < area.Y
}

// inter.X,Y is alwyas relative to rect. If there is total overlap, it'll be at 0,0
func Overlap(rect Rect, other Rect) (inter Rect, ok bool) {
	minR, maxR := rect.Vec2, Add(rect.Vec2, rect.Size)
	minO, maxO := other.Vec2, Add(other.Vec2, other.Size)

	inter.X = Max(minR.X, minO.X)
	inter.Y = Max(minR.Y, minO.Y)
	inter.Size.X = Min(maxR.X, maxO.X)
	inter.Size.Y = Min(maxR.Y, maxO.Y)
	inter.Size = Subtract(inter.Size, inter.Vec2)
	ok = inter.Size.X > 0 && inter.Size.Y > 0
	if ok {
		return inter, true
	}
	return Rect{}, false
}
