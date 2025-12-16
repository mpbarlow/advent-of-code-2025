package support

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func AbsInt[T int | int64](i T) int {
	if i < 0 {
		return int(-i)
	}

	return int(i)
}

func MinInt(ints ...int) int {
	min := math.MaxInt

	for _, i := range ints {
		if i < min {
			min = i
		}
	}

	return min
}

func MaxInt(ints ...int) int {
	max := math.MinInt

	for _, i := range ints {
		if i > max {
			max = i
		}
	}

	return max
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

func SliceOfNumericStringsToSliceOfInts(in []string) []int {
	out := make([]int, len(in))

	for i, numericStr := range in {
		asInt, err := strconv.Atoi(numericStr)
		if err != nil {
			panic(fmt.Sprintf("Could not convert %s to int", numericStr))
		}

		out[i] = asInt
	}

	return out
}
