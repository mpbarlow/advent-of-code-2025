package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := support.LoadInput()
	moves := strings.Split(input, "\n")

	fmt.Println(partOne(moves))
	fmt.Println(partTwo(moves))
}

func partOne(moves []string) int {
	dial := 50
	zeroes := 0

	for _, line := range moves {
		dir := line[:1]
		rotateBy, _ := strconv.ParseInt(line[1:], 10, 64)

		if dir == "L" {
			rotateBy *= -1
		}

		dial += int(rotateBy % 100)

		// Normalise back to 0...99
		if dial < 0 {
			dial += 100
		} else if dial > 99 {
			dial -= 100
		}

		if dial == 0 {
			zeroes += 1
		}
	}

	return zeroes
}

func partTwo(moves []string) int {
	dial := 50
	zeroes := 0

	for _, line := range moves {
		dialStart := dial

		dir := line[:1]
		rotateBy, _ := strconv.ParseInt(line[1:], 10, 64)

		if dir == "L" {
			rotateBy *= -1
		}

		dial += int(rotateBy % 100)

		// Calculate the number of complete rotations
		zeroes += support.AbsInt(rotateBy / 100)

		// Normalise back to 0...99, counting any partial rotations that crossed 0
		if dial < 0 {
			dial += 100
			// e.g. L5 from 0 would not count as crossing zero
			if dialStart > 0 {
				zeroes += 1
			}
		} else if dial > 99 {
			dial -= 100
			zeroes += 1
		} else if dial == 0 {
			zeroes += 1
		}
	}

	return zeroes
}
