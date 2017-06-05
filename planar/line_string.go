package planar

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

// LineString represents a set of points to be thought of as a polyline.
type LineString []Point

// NewLineString creates a new line string.
func NewLineString() LineString {
	return LineString{}
}

// NewLineStringFromXYData creates a line string from a slice of [2]float64 values
// representing [horizontal, vertical] type data.
func NewLineStringFromXYData(data [][2]float64) LineString {
	ls := make(LineString, len(data))
	for i := range data {
		ls[i] = Point{data[i][0], data[i][1]}
	}

	return ls
}

// NewLineStringFromYXData creates a line string from a slice of [2]float64 values
// representing [vertical, horizontal] type data, for example typical lat/lng data.
func NewLineStringFromYXData(data [][2]float64) LineString {
	ls := make(LineString, len(data))
	for i := range data {
		ls[i] = Point{data[i][1], data[i][0]}
	}

	return ls
}

// NewLineStringFromXYSlice creates a line string from a slice of []float64 values.
// The first two elements are taken to be horizontal and vertical components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewLineStringFromXYSlice(data [][]float64) LineString {
	ls := make(LineString, 0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			ls = append(ls, Point{data[i][0], data[i][1]})
		}
	}

	return ls
}

// NewLineStringFromYXSlice creates a line string from a slice of []float64 values.
// The first two elements are taken to be vertical and horizontal components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewLineStringFromYXSlice(data [][]float64) LineString {
	ls := make(LineString, 0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			ls = append(ls, Point{data[i][1], data[i][0]})
		}
	}

	return ls
}

// Distance computes the total distance in the units of the points.
func (ls LineString) Distance() float64 {
	sum := 0.0

	loopTo := len(ls) - 1
	for i := 0; i < loopTo; i++ {
		sum += ls[i].DistanceFrom(ls[i+1])
	}

	return sum
}

// DistanceFrom computes an O(n) distance from the line string. Loops over every
// subline to find the minimum distance.
func (ls LineString) DistanceFrom(point Point) float64 {
	return math.Sqrt(ls.DistanceFromSquared(point))
}

// DistanceFromSquared computes an O(n) minimum squared distance from the line string.
// Loops over every subline to find the minimum distance.
func (ls LineString) DistanceFromSquared(point Point) float64 {
	dist := math.Inf(1)

	seg := segment{}
	loopTo := len(ls) - 1
	for i := 0; i < loopTo; i++ {
		seg[0] = ls[i]
		seg[1] = ls[i+1]

		// TODO: would checking the distance to the bound first be faster?
		dist = math.Min(seg.DistanceFromSquared(point), dist)
	}

	return dist
}

// Interpolate interpolates the line string by distance.
func (ls LineString) Interpolate(percent float64) Point {
	if percent <= 0 {
		return ls[0]
	} else if percent >= 1 {
		return ls[len(ls)-1]
	}

	destination := ls.Distance() * percent
	travelled := 0.0

	for i := 0; i < len(ls)-1; i++ {
		dist := ls[i].DistanceFrom(ls[i+1])
		if (travelled + dist) > destination {
			factor := (destination - travelled) / dist
			return Point{
				ls[i][0]*(1-factor) + ls[i+1][0]*factor,
				ls[i][1]*(1-factor) + ls[i+1][1]*factor,
			}
		}
		travelled += dist
	}

	return ls[0]
}

// Project computes the percent along this line string closest to the given point,
// normalized to the length of the line string.
func (ls LineString) Project(point Point) float64 {
	minDistance := math.Inf(1)
	measure := math.Inf(-1)
	sum := 0.0

	seg := segment{}
	for i := 0; i < len(ls)-1; i++ {
		seg[0] = ls[i]
		seg[1] = ls[i+1]

		distanceToLine := seg.DistanceFromSquared(point)
		segDistance := seg.Distance()

		if distanceToLine < minDistance {
			minDistance = distanceToLine

			proj := seg.Project(point)
			if proj < 0 {
				proj = 0
			} else if proj > 1 {
				proj = 1
			}

			measure = sum + proj*segDistance
		}
		sum += segDistance
	}
	return measure / sum
}

// Centroid computes the centroid of the line string.
func (ls LineString) Centroid() Point {
	point, _ := ls.CentroidDistance()
	return point
}

// CentroidDistance computes the centroid and the distance/length
// of the line. If you need both this is faster since we need the
// length to compute the centroid.
func (ls LineString) CentroidDistance() (Point, float64) {
	dist := 0.0
	point := Point{}

	seg := segment{}

	// implicitly move everything to near the origin to help with roundoff
	offset := Vector(ls[0]).Negative()
	for i := 0; i < len(ls)-1; i++ {
		seg[0] = ls[i].Add(offset)
		seg[1] = ls[i+1].Add(offset)

		d := seg.Distance()
		centroid := seg.Centroid()

		point[0] += centroid[0] * d
		point[1] += centroid[1] * d

		dist += d
	}

	point[0] /= dist
	point[1] /= dist

	point = point.Add(offset.Negative())
	return point, dist
}

// Reverse changes the direction of the line string.
// This is done inplace, ie. it modifies the original data.
func (ls LineString) Reverse() {
	l := len(ls) - 1
	for i := 0; i <= l/2; i++ {
		ls[i], ls[l-i] = ls[l-i], ls[i]
	}
}

// Bound returns a rectangle bound around the line string. Uses rectangular coordinates.
func (ls LineString) Bound() Bound {
	return MultiPoint(ls).Bound()
}

// Equal compares two line strings. Returns true if lengths are the same
// and all points are Equal.
func (ls LineString) Equal(lineString LineString) bool {
	return MultiPoint(ls).Equal(MultiPoint(lineString))
}

// Clone returns a new copy of the line string.
func (ls LineString) Clone() LineString {
	ps := MultiPoint(ls)
	return LineString(ps.Clone())
}

// WKT returns the line string in WKT format, eg. LINESTRING(30 10,10 30,40 40)
// For empty line strings the result will be 'EMPTY'.
func (ls LineString) WKT() string {
	if len(ls) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "LINESTRING")
	wktPoints(buff, ls)

	return buff.String()
}

// String returns a string representation of the line string.
func (ls LineString) String() string {
	return ls.WKT()
}

func wktPoints(w io.Writer, ps []Point) {
	fmt.Fprintf(w, "(%g %g", ps[0][0], ps[0][1])

	for i := 1; i < len(ps); i++ {
		fmt.Fprintf(w, ",%g %g", ps[i][0], ps[i][1])
	}

	fmt.Fprintf(w, ")")
}
