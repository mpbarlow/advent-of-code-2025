package main

import (
	"advent-of-code-2025/support"
	"fmt"
)

const start rune = 'S'
const splitter rune = '^'
const space rune = '.'

func main() {
	input := support.InputTo2DGrid(support.LoadInput())

	fmt.Println(countBeamSplits(input, support.Point2{X: findStartX(input), Y: 0}))
	fmt.Println(countBeamPaths(input, support.Point2{X: findStartX(input), Y: 0}))
}

func findStartX(input [][]rune) int {
	for x, value := range input[0] {
		if value == start {
			return x
		}
	}

	panic("Could not find start in first line of input")
}

// Beams can merge again, so we want to make sure we only count each splitter once
var encounteredSplitters = support.NewSet[support.Point2]()

// Count the number of unique beam splits as it progresses downwards
func countBeamSplits(input [][]rune, pos support.Point2) int {
	if pos.Y == len(input)-1 {
		return 0
	}

	switch input[pos.Y][pos.X] {
	case space, start:
		return countBeamSplits(input, support.Point2{X: pos.X, Y: pos.Y + 1})

	case splitter:
		if encounteredSplitters.Has(pos) {
			return 0
		}

		encounteredSplitters.Add(pos)

		return 1 +
			countBeamSplits(input, support.Point2{X: pos.X - 1, Y: pos.Y + 1}) +
			countBeamSplits(input, support.Point2{X: pos.X + 1, Y: pos.Y + 1})

	default:
		panic("Encountered unexpected item in grid")
	}
}

var cache = make(map[support.Point2]int)

// Count the number of unique paths the beam could take across all splitters. Memoise so it's actually computable.
func countBeamPaths(input [][]rune, pos support.Point2) int {
	if pos.Y == len(input)-1 {
		return 1
	}

	switch input[pos.Y][pos.X] {
	case space, start:
		return countBeamPaths(input, support.Point2{X: pos.X, Y: pos.Y + 1})

	case splitter:
		if val, ok := cache[pos]; ok {
			return val
		}

		paths := countBeamPaths(input, support.Point2{X: pos.X - 1, Y: pos.Y + 1}) +
			countBeamPaths(input, support.Point2{X: pos.X + 1, Y: pos.Y + 1})

		cache[pos] = paths

		return paths

	default:
		panic("Encountered unexpected item in grid")
	}
}
