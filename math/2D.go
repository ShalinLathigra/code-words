package math

import "fmt"

type Vec2 struct {
	X int
	Y int
}

var Left = Vec2{-1, 0}
var Right = Vec2{1, 0}
var Up = Vec2{0, -1}
var Down = Vec2{0, 1}

func (v Vec2) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func Add(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

func Subtract(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X - b.X, a.Y - b.Y}
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

func Shrink(rect Rect, skin_width int) (ret Rect, ok bool) {
	if skin_width*2 >= rect.Size.X || skin_width*2 >= rect.Size.Y {
		return Rect{}, false
	}
	ok = true
	ret.Vec2 = Add(rect.Vec2, Vec2{skin_width, skin_width})
	ret.Size = Add(rect.Size, Vec2{-2 * skin_width, -2 * skin_width})
	return
}

func Contains(rect Rect, point Vec2) bool {
	min, max := rect.Vec2, Add(rect.Vec2, rect.Size)
	return point.X >= min.X && point.X < max.X && point.Y >= min.Y && point.Y < max.Y
}

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
