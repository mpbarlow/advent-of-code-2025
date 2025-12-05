package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"strings"
)

func main() {
	rangesAndIngredientIds := strings.Split(support.LoadInput(), "\n\n")

	ingredientMap := buildIngredientMap(strings.Split(rangesAndIngredientIds[0], "\n"))
	ingredientIds := support.SliceOfNumericStringsToSliceOfInts(strings.Split(rangesAndIngredientIds[1], "\n"))

	fmt.Println(partOne(ingredientMap, ingredientIds))
	fmt.Println(partTwo(ingredientMap))
}

func partOne(ingredientMap map[int]int, ingredientIds []int) int {
	freshCount := 0

	for _, id := range ingredientIds {
		for k, v := range ingredientMap {
			if k <= id && v >= id {
				freshCount++
				break
			}
		}
	}

	return freshCount
}

func partTwo(ingredientMap map[int]int) int {
	total := 0

	for k, v := range ingredientMap {
		total += (v - k) + 1
	}

	return total
}

// Given a slice of raw ingredient range lines e.g. {"3-5", "6-8"}, return a map of lower bounds to upper bounds.
// Overlapping ranges are coalesced.
func buildIngredientMap(ranges []string) map[int]int {
	ingredientMap := make(map[int]int, len(ranges))

	for _, r := range ranges {
		lowerBound, upperBound := rangeFromString(r)

		if prevUpper, ok := ingredientMap[lowerBound]; !ok || upperBound > prevUpper {
			ingredientMap[lowerBound] = upperBound
		}
	}

	return coalesceRanges(ingredientMap)
}

// Combine any overlapping ranges together. This is not necessary for part one but makes part two super easy.
func coalesceRanges(ranges map[int]int) map[int]int {
	initialCount := len(ranges)

	for {
	outer:
		for kOther, vOther := range ranges {
			for k, v := range ranges {
				// Don't try to coalesce something with itself
				if kOther == k && vOther == v {
					continue
				}

				// Completely replace existing range with a larger one fully overlapping it
				if k <= kOther && v >= vOther {
					delete(ranges, kOther)
					ranges[k] = v

					// We're messing with the map every time we coalesce a range, so we need to restart the iteration
					break outer
				}

				// Extend existing lower bound downwards
				if k <= kOther && v >= kOther && v <= vOther {
					delete(ranges, kOther)
					ranges[k] = vOther

					break outer
				}

				// Extend existing upper bound upwards
				if k >= kOther && k <= vOther && v >= vOther {
					delete(ranges, k)
					ranges[kOther] = v

					break outer
				}
			}
		}

		// Once we are no longer able to coalesce any ranges, we're done
		if len(ranges) == initialCount {
			return ranges
		}

		initialCount = len(ranges)
	}
}

// Given a string e.g. "3-5" return (3, 5)
func rangeFromString(str string) (int, int) {
	upperAndLowerBound := support.SliceOfNumericStringsToSliceOfInts(strings.Split(str, "-"))
	if len(upperAndLowerBound) != 2 {
		panic(fmt.Sprintf("Could not get two parts from %s", str))
	}

	return upperAndLowerBound[0], upperAndLowerBound[1]
}
