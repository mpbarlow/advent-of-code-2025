package support

import "strings"

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
