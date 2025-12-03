package support

import (
	"fmt"
	"strconv"
	"strings"
)

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
