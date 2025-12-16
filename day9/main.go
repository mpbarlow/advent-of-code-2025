package main

import (
	"advent-of-code-2025/support"
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"
)

func main() {
	points := support.Map(
		strings.Split(support.LoadInput(), "\n"),
		func(line string) support.Point2 {
			coords := support.SliceOfNumericStringsToSliceOfInts(strings.Split(line, ","))
			if len(coords) != 2 {
				panic("Expected co-ordinate with two points")
			}

			return support.Point2{X: coords[0], Y: coords[1]}
		},
	)

	fmt.Println(partOne(points))
	fmt.Println(partTwo(points))
}

func partOne(points []support.Point2) int {
	largestArea := -1

	// Brute force is plenty fast enough here
	for _, left := range points {
		for _, right := range points {
			if left == right {
				continue
			}

			pair := pointPair{from: left, to: right}
			area := pair.Area()

			if area > largestArea {
				largestArea = area
			}
		}
	}

	return largestArea
}

type pointPair struct {
	from, to support.Point2
}

func (p *pointPair) Area() int {
	// +1 because the edges are included in the area
	return int(math.Abs(float64(p.from.X-p.to.X))+1) * int(math.Abs(float64(p.from.Y-p.to.Y))+1)
}

func partTwo(points []support.Point2) int {
	horizontalWalls, verticalWalls := collectWalls(points)

	for _, pair := range collectPointPairs(points) {
		// The first valid rectangle formed must be the largest
		if validRectangle(pair, horizontalWalls, verticalWalls) {
			return pair.Area()
		}
	}

	return -1
}

// Returns a slice of each unique pair of points, ordered by the area of the rectange they form, largest first.
func collectPointPairs(points []support.Point2) []pointPair {
	uniquePairs := support.NewSet[pointPair]()

	for _, left := range points {
		for _, right := range points {
			if left == right {
				continue
			}

			// Use consistent ordering so we only get unique pairs
			left, right := left, right
			if !aSmallerThanB(left, right) {
				left, right = right, left
			}

			uniquePairs.Add(pointPair{from: left, to: right})
		}
	}

	pairs := slices.Collect(maps.Keys(uniquePairs))
	slices.SortFunc(pairs, func(a, b pointPair) int {
		if b.Area() < a.Area() {
			return -1
		}

		if b.Area() > a.Area() {
			return 1
		}

		return 0
	})

	return pairs
}

// Returns true if point a is lexographically smaller than b
func aSmallerThanB(a, b support.Point2) bool {
	if a.X != b.X {
		return a.X < b.X
	}

	return a.Y < b.Y
}

// Returns two slices representing the horizontal and vertical "walls" of the shape formed by subsequent point pairs.
func collectWalls(points []support.Point2) ([]pointPair, []pointPair) {
	horizontalWalls := make([]pointPair, 0)
	verticalWalls := make([]pointPair, 0)

	for i, curr := range points {
		var next support.Point2

		// Handle wrap around at the end
		if i == len(points)-1 {
			next = points[0]
		} else {
			next = points[i+1]
		}

		if curr.Y == next.Y {
			horizontalWalls = append(horizontalWalls, pointPair{from: curr, to: next})
		} else if curr.X == next.X {
			verticalWalls = append(verticalWalls, pointPair{from: curr, to: next})
		} else {
			panic("Both x and y changed between subsequent points")
		}
	}

	return horizontalWalls, verticalWalls
}

// Determines if the rectangle formed by `pair` is valid by making sure it doesn't intersect any walls of the shape.
func validRectangle(pair pointPair, horizontalWalls []pointPair, verticalWalls []pointPair) bool {
	startX := support.MinInt(pair.from.X, pair.to.X)
	endX := support.MaxInt(pair.from.X, pair.to.X)
	startY := support.MinInt(pair.from.Y, pair.to.Y)
	endY := support.MaxInt(pair.from.Y, pair.to.Y)

	// We check the space inside the rectangle rather than the rectangle itself because the edges of the rectangle may
	// be formed of multiple walls, and so "intersects" walls that do not interrupt the rectangle.
	topAndBottom := []pointPair{
		{from: support.Point2{X: startX + 1, Y: startY + 1}, to: support.Point2{X: endX - 1, Y: startY + 1}},
		{from: support.Point2{X: startX + 1, Y: endY - 1}, to: support.Point2{X: endX - 1, Y: endY - 1}},
	}

	for _, a := range topAndBottom {
		for _, b := range verticalWalls {
			if intersects(a, b) {
				return false
			}
		}
	}

	leftAndRight := []pointPair{
		{from: support.Point2{X: startX + 1, Y: startY + 1}, to: support.Point2{X: startX + 1, Y: endY - 1}},
		{from: support.Point2{X: endX - 1, Y: startY + 1}, to: support.Point2{X: endX - 1, Y: endY - 1}},
	}

	for _, a := range leftAndRight {
		for _, b := range horizontalWalls {
			if intersects(a, b) {
				return false
			}
		}
	}

	return true
}

// Check if two perpendicular lines intersect. No idea what it does if they aren't perpendicular so...don't do that.
func intersects(a, b pointPair) bool {
	// Always use the horizontal line as `a` to make life easier
	if a.from.X == a.to.X {
		a, b = b, a
	}

	aMinX := support.MinInt(a.from.X, a.to.X)
	aMaxX := support.MaxInt(a.from.X, a.to.X)
	aMinY := support.MinInt(a.from.Y, a.to.Y)
	bMinX := support.MinInt(b.from.X, b.to.X)
	bMinY := support.MinInt(b.from.Y, b.to.Y)
	bMaxY := support.MaxInt(b.from.Y, b.to.Y)

	return aMinX <= bMinX && aMaxX >= bMinX && aMinY >= bMinY && aMinY <= bMaxY
}
