package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	boxes := parsePositions(support.LoadInput())
	distances := getOrderedBoxPairDistances(boxes)

	fmt.Println(partOne(NewCircuitSet(boxes), distances))
	fmt.Println(partTwo(NewCircuitSet(boxes), distances))
}

const partOneIterations = 1000

func partOne(circuitSet *CircuitSet, distances []BoxPairDistance) int {
	for _, distance := range distances[:partOneIterations] {
		circuitSet.ConnectCircuits(distance.left, distance.right)
	}

	result := 1

	for _, v := range circuitSet.LargestCircuits(3) {
		result *= v
	}

	return result
}

func partTwo(circuitSet *CircuitSet, distances []BoxPairDistance) int {
	for _, distance := range distances {
		circuitSet.ConnectCircuits(distance.left, distance.right)

		if circuitSet.IsThereOnlyOneCircuitYet() {
			return distance.left.X * distance.right.X
		}
	}

	panic("There's still more than one circuit")
}

type BoxPairDistance struct {
	left     support.Point3
	right    support.Point3
	distance float64
}

// Given a slice of points in 3D space, return a slice of each pair and their distance, ordered by distance ascending.
// Each pair of points only appears in one order.
func getOrderedBoxPairDistances(boxes []support.Point3) []BoxPairDistance {
	distances := make(map[support.Point3]map[support.Point3]float64)

	for _, left := range boxes {
		for _, right := range boxes {
			if left == right {
				continue
			}

			// If we've already recorded this distance in the opposite direction don't do it again. There is probably a
			// better way to do this but meh.
			if _, ok := distances[right]; ok {
				if _, ok := distances[right][left]; ok {
					continue
				}
			}

			if _, ok := distances[left]; !ok {
				distances[left] = make(map[support.Point3]float64)
			}

			distances[left][right] = left.DistanceTo(right)
		}
	}

	orderedDistances := make([]BoxPairDistance, 0, len(distances))

	for left, rights := range distances {
		for right, distance := range rights {
			orderedDistances = append(orderedDistances, BoxPairDistance{left: left, right: right, distance: distance})
		}
	}

	slices.SortFunc(orderedDistances, func(a BoxPairDistance, b BoxPairDistance) int {
		if a.distance < b.distance {
			return -1
		}

		if a.distance > b.distance {
			return 1
		}

		return 0
	})

	return orderedDistances
}

type CircuitSet struct {
	circuitMap map[support.Point3]int
}

func NewCircuitSet(junctionBoxes []support.Point3) *CircuitSet {
	set := CircuitSet{
		circuitMap: make(map[support.Point3]int, len(junctionBoxes)),
	}

	for i, box := range junctionBoxes {
		// Start with every junction box being on its own circuit
		set.circuitMap[box] = i
	}

	return &set
}

func (c *CircuitSet) ConnectCircuits(left support.Point3, right support.Point3) {
	fromCircuit, fromOk := c.circuitMap[right]
	toCircuit, toOk := c.circuitMap[left]

	if !fromOk || !toOk {
		panic("Received a point I don't know about!")
	}

	for k, v := range c.circuitMap {
		if v == fromCircuit {
			c.circuitMap[k] = toCircuit
		}
	}
}

func (c *CircuitSet) LargestCircuits(n int) []int {
	histogram := make(map[int]int, 0)

	for _, circuitNo := range c.circuitMap {
		if v, ok := histogram[circuitNo]; !ok {
			histogram[circuitNo] = 1
		} else {
			histogram[circuitNo] = v + 1
		}
	}

	counts := slices.Collect(maps.Values(histogram))
	slices.Sort(counts)
	slices.Reverse(counts)

	return counts[:n]
}

func (c *CircuitSet) IsThereOnlyOneCircuitYet() bool {
	circuit := -1

	for _, v := range c.circuitMap {
		if circuit == -1 {
			circuit = v
		}

		if v != circuit {
			return false
		}
	}

	return true
}

func parsePositions(input string) []support.Point3 {
	lines := strings.Split(input, "\n")
	connections := make([]support.Point3, 0, len(lines))

	for _, line := range lines {
		point := support.SliceOfNumericStringsToSliceOfInts(strings.Split(line, ","))
		if len(point) != 3 {
			panic("Expected 3D point")
		}

		connections = append(connections, support.Point3{X: point[0], Y: point[1], Z: point[2]})
	}

	return connections
}
