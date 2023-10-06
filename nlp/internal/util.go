package internal

import (
	"os"
	"path/filepath"
)

func ReadDataFile(path string) []byte {
	p, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	data, ferr := os.ReadFile(p)
	if err != nil {
		panic(ferr)
	}

	return data
}

// min returns the minimum of `a` and `b`.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// stringInSlice determines if `slice` contains the string `a`.
func StringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}
