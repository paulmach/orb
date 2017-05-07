package planar

import (
	"bytes"
	"fmt"
	"math"
)

// A MultiPoint represents a set of points in the 2D Eucledian or Cartesian plane.
type MultiPoint []Point

// NewMultiPoint creates a new MultiPoint object.
func NewMultiPoint() MultiPoint {
	return MultiPoint{}
}

// Clone returns a new copy of the MultiPoint object.
func (mp MultiPoint) Clone() MultiPoint {
	points := make([]Point, len(mp))
	copy(points, mp)

	return MultiPoint(points)
}

// Centroid returns the average x and y coordinate of the points.
// This can also be used for small clusters of lat/lng points.
func (mp MultiPoint) Centroid() Point {
	x, y := 0.0, 0.0
	for _, p := range mp {
		x += p[0]
		y += p[1]
	}

	num := float64(len(mp))
	return Point{x / num, y / num}
}

// DistanceFrom returns the minimum euclidean distance from the points.
func (mp MultiPoint) DistanceFrom(point Point) float64 {
	d, _ := mp.DistanceFromWithIndex(point)
	return d
}

// DistanceFromWithIndex returns the minimum euclidean distance
// from the points plus the index of that point.
func (mp MultiPoint) DistanceFromWithIndex(point Point) (float64, int) {
	dist := math.Inf(1)
	index := 0

	for i := range mp {
		if d := mp[i].DistanceFromSquared(point); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}

// Bound returns a rectangle bound around the points. Uses rectangular coordinates.
func (mp MultiPoint) Bound() Bound {
	if len(mp) == 0 {
		return NewBound(0, 0, 0, 0)
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range mp {
		minX = math.Min(minX, v.X())
		minY = math.Min(minY, v.Y())

		maxX = math.Max(maxX, v.X())
		maxY = math.Max(maxY, v.Y())
	}

	return NewBound(maxX, minX, maxY, minY)
}

// Equal compares two sets of points. Returns true if lengths are the same
// and all points are Equal.
func (mp MultiPoint) Equal(multiPoint MultiPoint) bool {
	if len(mp) != len(multiPoint) {
		return false
	}

	for i := range mp {
		if !mp[i].Equal(multiPoint[i]) {
			return false
		}
	}

	return true
}

// WKT returns the points in WKT format,
// eg. MULTIPOINT(30 10, 10 30, 40 40)
func (mp MultiPoint) WKT() string {
	if len(mp) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTIPOINT")
	wktPoints(buff, mp)

	return buff.String()
}

// String returns a string representation of the points.
// The format is WKT, e.g. MULTIPOINT(30 10,10 30,40 40)
// For empty sets the result will be 'EMPTY'.
func (mp MultiPoint) String() string {
	return mp.WKT()
}
