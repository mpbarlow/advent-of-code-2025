package support

import (
	"os"
	"path/filepath"
	"runtime"
)

func LoadInput() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not determine running directory!")
	}

	bytes, err := os.ReadFile(filepath.Dir(file) + "/input.txt")
	if err != nil {
		panic("Could not read local input file!")
	}

	return string(bytes)
}

func AbsInt[T int | int64](i T) int {
	if i < 0 {
		return int(-i)
	}

	return int(i)
}
