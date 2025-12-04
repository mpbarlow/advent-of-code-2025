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
