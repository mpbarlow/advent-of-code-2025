package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := strings.Split(support.LoadInput(), "\n")
	input := support.Map(lines, support.StringOfDigitsAsSliceOfInts)

	fmt.Println(solvePart(input, 2))
	fmt.Println(solvePart(input, 12))
}

func solvePart(input [][]int, batteryCount int) int {
	results := 0

	for _, line := range input {
		result := solveLine(line, batteryCount)

		resultAsInt, err := strconv.Atoi(result)
		if err != nil {
			panic(fmt.Sprintf("Could not convert %s to int", result))
		}

		results += resultAsInt
	}

	return results
}

func solveLine(line []int, targetLength int) string {
	if targetLength == 0 {
		return ""
	}

	// Given we can't rearrange batteries, we have to choose the largest digit that has enough digits to the right to
	// make up to target length.
	leftIndex, leftDigit := findMax(line[:(len(line)-targetLength)+1])

	// Then just do this recursively with those digits to the right.
	return strconv.Itoa(leftDigit) + solveLine(line[leftIndex+1:], targetLength-1)
}

// Returns (indexOfLargestValue, largestValue) from line; (-1, -1) for empty slices.
// Assumes all values in line are positive.
func findMax(line []int) (int, int) {
	largestValue := -1
	indexOfLargestValue := -1

	for index, value := range line {
		if value > largestValue {
			largestValue = value
			indexOfLargestValue = index
		}
	}

	return indexOfLargestValue, largestValue
}
