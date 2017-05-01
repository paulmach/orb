package planar

import (
	"bytes"
	"fmt"
	"math"
)

// MultiLineString is a set of polylines.
type MultiLineString []LineString

// NewMultiLineString creates a new multi line string.
func NewMultiLineString() MultiLineString {
	return MultiLineString{}
}

// Distance computes the sum of the distances of each line string.
func (mls MultiLineString) Distance() float64 {
	sum := 0.0
	for _, ls := range mls {
		sum += ls.Distance()
	}

	return sum
}

// DistanceFrom computes an O(n) distance from the multi line string.
// Loops over every subline to find the minimum distance.
func (mls MultiLineString) DistanceFrom(point Point) float64 {
	min := math.MaxFloat64

	// TODO: would checking the distance to the bound first be faster?
	for _, ls := range mls {
		if d := ls.DistanceFromSquared(point); d < min {
			min = d
		}
	}

	return math.Sqrt(min)
}

// Centroid computes the centroid of the line string.
func (mls MultiLineString) Centroid() Point {
	point := Point{}
	dist := 0.0

	for _, ls := range mls {
		c, d := ls.CentroidDistance()

		point[0] += c[0] * d
		point[1] += c[1] * d

		dist += d
	}

	point[0] /= dist
	point[1] /= dist

	return point
}

// Bound returns a rectangle bound around all the line strings.
func (mls MultiLineString) Bound() Rect {
	bound := mls[0].Bound()
	for i := 0; i < len(mls); i++ {
		bound = bound.Union(mls[i].Bound())
	}

	return bound
}

// Equal compares two multi line strings. Returns true if lengths are the same
// and all points are Equal.
func (mls MultiLineString) Equal(multiLineString MultiLineString) bool {
	if len(mls) != len(multiLineString) {
		return false
	}

	for i, ls := range mls {
		if !ls.Equal(multiLineString[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the multi line string.
func (mls MultiLineString) Clone() MultiLineString {
	nmls := make(MultiLineString, 0, len(mls))
	for _, ls := range mls {
		nmls = append(nmls, ls.Clone())
	}

	return nmls
}

// WKT returns the line string in WKT format, eg. LINESTRING(30 10,10 30,40 40)
// For empty line strings the result will be 'EMPTY'.
func (mls MultiLineString) WKT() string {
	if len(mls) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTILINESTRING(")
	wktPoints(buff, mls[0])

	for i := 1; i < len(mls); i++ {
		buff.Write([]byte(","))
		wktPoints(buff, mls[i])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns a string representation of the line string.
func (mls MultiLineString) String() string {
	return mls.WKT()
}
