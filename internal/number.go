package internal

import "fmt"

// min returns the minimum of `a` and `b`.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func EqualFloat(expected, observed float64) bool {
	return fmt.Sprintf("%0.2f", expected) == fmt.Sprintf("%0.2f", observed)
}
