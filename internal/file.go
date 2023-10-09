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
