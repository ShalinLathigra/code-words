package math

func min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

// if value < a, return a
// if value > b, return b
// if value is in between, return
func clamp(v int, a int, b int) int {
	return min(max(v, a), b)
}
