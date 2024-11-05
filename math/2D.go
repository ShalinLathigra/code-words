package math

import "fmt"

type Vec2 struct {
	X int
	Y int
}

func (v Vec2) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func add(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

func subtract(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X - b.X, a.Y - b.Y}
}

func dot(a Vec2, b Vec2) int {
	return a.X*b.X + a.Y*b.Y
}

type Rect struct {
	Vec2
	Size Vec2
}

func (r Rect) String() string {
	return fmt.Sprintf("%s*%s", r.Vec2, r.Size)
}

func contains(rect Rect, point Vec2) bool {
	min, max := rect.Vec2, add(rect.Vec2, rect.Size)
	return point.X >= min.X && point.X < max.X && point.Y >= min.Y && point.Y < max.Y
}

func overlap(rect Rect, other Rect) (inter Rect, ok bool) {
	minR, maxR := rect.Vec2, add(rect.Vec2, rect.Size)
	minO, maxO := other.Vec2, add(other.Vec2, other.Size)

	inter.X = max(minR.X, minO.X)
	inter.Y = max(minR.Y, minO.Y)
	inter.Size.X = min(maxR.X, maxO.X)
	inter.Size.Y = min(maxR.Y, maxO.Y)
	inter.Size = subtract(inter.Size, inter.Vec2)
	ok = inter.Size.X > 0 && inter.Size.Y > 0
	if ok {
		return inter, true
	}
	return Rect{}, false
}
