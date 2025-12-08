package main

import (
	"advent-of-code-2025/support"
	"fmt"
)

const roll rune = '@'
const emptySpace rune = '.'

func main() {
	input := support.LoadInput()

	fmt.Println(partOne(support.InputTo2DGrid(input)))
	fmt.Println(partTwo(support.InputTo2DGrid(input)))
}

func partOne(grid [][]rune) int {
	return len(getReachableRolls(grid))
}

func partTwo(grid [][]rune) int {
	totalRemovableRolls := 0

	for {
		removableRolls := getReachableRolls(grid)

		if len(removableRolls) == 0 {
			break
		}

		totalRemovableRolls += len(removableRolls)
		grid = removeRolls(grid, removableRolls)
	}

	return totalRemovableRolls
}

// Return the coordinates of all removable rolls
func getReachableRolls(grid [][]rune) []support.Point2 {
	reachableRolls := make([]support.Point2, 0)

	for y, line := range grid {
		for x := range line {
			if grid[y][x] != roll {
				continue
			}

			neighbouringRolls := 0

			for _, offset := range support.RelativeCardinalDirections {
				if x+offset.X < 0 || x+offset.X >= len(line) {
					continue
				}

				if y+offset.Y < 0 || y+offset.Y >= len(grid) {
					continue
				}

				if grid[y+offset.Y][x+offset.X] == roll {
					neighbouringRolls++
				}
			}

			if neighbouringRolls < 4 {
				reachableRolls = append(reachableRolls, support.Point2{X: x, Y: y})
			}
		}
	}

	return reachableRolls
}

// Remove the rolls specified by rollCoords from grid, returning the new grid
func removeRolls(grid [][]rune, rollCoords []support.Point2) [][]rune {
	for _, coord := range rollCoords {
		grid[coord.Y][coord.X] = emptySpace
	}

	return grid
}
