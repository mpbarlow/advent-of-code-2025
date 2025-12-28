package support

import (
	"slices"
	"strings"
)

func Map[I any, O any](input []I, fn func(I) O) []O {
	output := make([]O, 0, len(input))

	for _, item := range input {
		output = append(output, fn(item))
	}

	return output
}

func InputTo2DGrid(input string) [][]rune {
	lines := strings.Split(input, "\n")
	grid := make([][]rune, len(lines))

	for i, line := range lines {
		grid[i] = []rune(line)
	}

	return grid
}

func Transpose[T any](input [][]T) [][]T {
	rows := len(input)
	if rows == 0 {
		return nil
	}

	cols := len(input[0])

	output := make([][]T, cols)
	for i := range output {
		output[i] = make([]T, rows)
	}

	for i := range input {
		for j := range input[i] {
			output[j][i] = input[i][j]
		}
	}

	return output
}

// Generate n choose k combinations. Returns an empty slice if k > len(n).
func GenerateCombinations[T any](n []T, k int) [][]T {
	combinations := make([][]T, 0)

	for i, val := range n {
		partialCombo := []T{val}

		if k > 1 {
			// Recursively generate the other possibilities for the combination. Exclude any element we've already
			// considered to prevent duplicates.
			nextN := make([]T, 0, len(n)-1)
			nextN = append(nextN, n[i+1:]...)

			for _, next := range GenerateCombinations(nextN, k-1) {
				combinations = append(combinations, slices.Concat(partialCombo, next))
			}
		} else {
			combinations = append(combinations, partialCombo)
		}
	}

	return combinations
}

// Given n, return a slice containing min..<max
func Range(min, max int) []int {
	r := make([]int, 0, max-min)

	for i := min; i < max; i++ {
		r = append(r, i)
	}

	return r
}
