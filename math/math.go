package math

func Min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func Max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

// if value < a, return a
// if value > b, return b
// if value is in between, return
func Clamp(v int, a int, b int) int {
	return Min(Max(v, a), b)
}
