package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := strings.Split(support.LoadInput(), ",")
	ranges := make([][]int, len(input))

	for idx, rangeString := range input {
		r := strings.Split(rangeString, "-")
		lowerBound, _ := strconv.Atoi(r[0])
		upperBound, _ := strconv.Atoi(r[1])

		ranges[idx] = []int{lowerBound, upperBound}
	}

	fmt.Println(partOne(ranges))
	fmt.Println(partTwo(ranges))
}

func partOne(ranges [][]int) int {
	total := 0

	for _, r := range ranges {
		// Check every number in the range
		for i := r[0]; i <= r[1]; i++ {
			iStr := strconv.Itoa(i)

			if len(iStr)%2 != 0 {
				continue
			}

			// If it has an even number of digits, it's invalid if the first half of the number == the second half
			if iStr[:len(iStr)/2] != iStr[len(iStr)/2:] {
				continue
			}

			total += i
		}
	}

	return total
}

func partTwo(ranges [][]int) int {
	total := 0

	for _, r := range ranges {
		for i := r[0]; i <= r[1]; i++ {
			if checkPartTwo(strconv.Itoa(i)) {
				total += i
			}
		}
	}

	return total
}

func checkPartTwo(pattern string) bool {
	chars := strings.Split(pattern, "")

	// Given e.g. 12345678, first check
	// - if 1234 1234 == 1234 5678,
	// - then if 12 12 12 12 == 12 34 56 78
	// - finally if 1 1 1 1 1 1 1 1 == 1 2 3 4 5 6 7 8
	// Big chunks first is marginally faster but that might just be for my input
	for j := len(pattern) / 2; j >= 1; j-- {
		if len(pattern)%j != 0 {
			continue
		}

		repeatTest := strings.Join(chars[:j], "")

		if strings.Repeat(repeatTest, len(pattern)/j) == pattern {
			return true
		}
	}

	return false
}
