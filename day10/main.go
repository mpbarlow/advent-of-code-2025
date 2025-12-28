package main

import (
	"advent-of-code-2025/day10/augmentedmatrix"
	"advent-of-code-2025/support"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	machines := support.Map(strings.Split(support.LoadInput(), "\n"), parseMachine)

	fmt.Println(partOne(machines))
	fmt.Println(partTwo(machines))
}

func partOne(machines []machine) int {
	// Keep generating all combinations of switches starting at 1 and moving upwards; the first valid solution
	// will therefore be the smallest.
	totalPresses := 0

	for _, m := range machines {
		totalPresses += findSmallestSequence(m)
	}

	return totalPresses
}

// Part two forms a system of vector linear equations where variables a...f are the number of times pressing a button:
//
// a(3) + b(1,3) + c(2) + d(2,3) + e(0,2) + f(0,1) = {3,5,4,7}
// -> a(0, 0, 0, 1) + b(0, 1, 0, 1) + c(0, 0, 1, 0) + d(0, 0, 1, 1) + e(1, 0, 1, 0) + f(1, 1, 0, 0) = {3,5,4,7}
//
// These vector equations can then be rewritten as four scalar (?) equations:
//
// 0a + 0b + 0c + 0d + 1e + 1f = 3
// 0a + 1b + 0c + 0d + 0e + 1f = 5
// 0a + 0b + 1c + 1d + 1e + 0f = 4
// 1a + 1b + 0c + 1d + 0e + 0f = 7
//
// Which can then be represented as an augmented matrix:
//
// | 0 0 0 0 1 1 || 3 |
// | 0 1 0 0 0 1 || 5 |
// | 0 0 1 1 1 0 || 4 |
// | 1 1 0 1 0 0 || 7 |
//
// We can then convert that to row echelon form to get solutions (see augmentedmatrix.AugmentedMatrix.toRowEchelonForm)
func partTwo(machines []machine) int {
	presses := 0

	for _, m := range machines {
		// Transform the bitfield version of the switches into a regular slice of ints 0...1
		switches := support.Map(m.switches, func(s int) []int {
			width := len(m.joltageLevels)
			result := make([]int, width)
			for i := width - 1; i >= 0; i-- {
				result[i] = s & 1
				s >>= 1
			}
			return result
		})

		augmentedMatrix := augmentedmatrix.NewRefAugmentedMatrix(switches, m.joltageLevels)
		presses += augmentedMatrix.Solve()
	}

	return presses
}

type machine struct {
	lights        int   // Bitfield of the target light pattern
	switches      []int // Slice of bitfields that can be XORed against the lights to model pressing that switch
	joltageLevels []int // Slice of ints representing the target final joltage levels
}

// Toggling switches is an XOR.
// XOR is commutative, associative, and its own inverse, meaning in the shortest sequence no button is pressed more than
// once: pressing it again would undo its effect, regardless of any other buttons pressed in between.
// I did not know this off the top of my head and had to look it up following a reddit hint.
// Following this, we can represent the lights as a bitfield and the switches as bitfields to XOR onto the lights.
func parseMachine(m string) machine {
	lightRegex := regexp.MustCompile(`([.#]+)`)
	lightsString := lightRegex.FindStringSubmatch(m)[0]

	// Convert the target light pattern to a bitfield where 1 is on and 0 is off
	lightsInt := 0
	for i, lightChar := range lightsString {
		if lightChar == '.' {
			continue
		}

		lightsInt |= 1 << ((len(lightsString) - 1) - i)
	}

	// Convert each switch to a bitfield representing the XOR that pressing it would perform
	switchRegex := regexp.MustCompile(`\([\d,]+\)`)
	switches := support.Map(
		switchRegex.FindAllStringSubmatch(m, -1),
		func(s []string) int {
			toggles := support.SliceOfNumericStringsToSliceOfInts(strings.Split(strings.Trim(s[0], "()"), ","))

			switchBitfield := 0
			for _, toggle := range toggles {
				switchBitfield |= 1 << ((len(lightsString) - 1) - toggle)
			}

			return switchBitfield
		},
	)

	joltageLevels := support.SliceOfNumericStringsToSliceOfInts(
		strings.Split(
			strings.Trim(
				string(regexp.MustCompile(`{[\d,]+}`).Find([]byte(m))),
				"{}",
			),
			",",
		),
	)

	return machine{lights: lightsInt, switches: switches, joltageLevels: joltageLevels}
}

func findSmallestSequence(m machine) int {
	presses := 1

	for {
		combinations := support.GenerateCombinations(m.switches, presses)

		if len(combinations) == 0 {
			panic("Could not find a valid combination!")
		}

		for _, combination := range combinations {
			lights := 0

			for _, press := range combination {
				lights ^= press
			}

			if lights == m.lights {
				return presses
			}
		}

		presses++
	}
}
