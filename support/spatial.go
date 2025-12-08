package support

import "math"

type Point2 struct {
	X int
	Y int
}

// The 8 relative positions around a cell in a 2D grid
var RelativeCardinalDirections = []Point2{
	{X: -1, Y: 0},
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 0, Y: 1},
	{X: -1, Y: 1},
}

type Point3 struct {
	X int
	Y int
	Z int
}

func (p Point3) DistanceTo(other Point3) float64 {
	return math.Sqrt(
		math.Pow(float64(other.X-p.X), 2) + math.Pow(float64(other.Y-p.Y), 2) + math.Pow(float64(other.Z-p.Z), 2),
	)
}
