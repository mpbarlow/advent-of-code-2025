package support

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

// The 8 relative positions around a cell in a 2D grid
var RelativeCardinalDirections = []Point{
	{X: -1, Y: 0},
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 0, Y: 1},
	{X: -1, Y: 1},
}

func AbsInt[T int | int64](i T) int {
	if i < 0 {
		return int(-i)
	}

	return int(i)
}

func StringOfDigitsAsSliceOfInts(in string) []int {
	digits := strings.Split(in, "")
	out := make([]int, 0, len(digits))

	for _, digit := range digits {
		intDigit, err := strconv.Atoi(digit)
		if err != nil {
			panic(fmt.Sprintf("Could not convert %s to int", digit))
		}

		out = append(out, intDigit)
	}

	return out
}
