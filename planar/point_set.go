package planar

import (
	"bytes"
	"fmt"
	"math"

	"github.com/paulmach/go.geojson"
)

// A PointSet represents a set of points in the 2D Eucledian or Cartesian plane.
type PointSet []Point

// NewPointSet simply creates a new point set with points array of the given size.
func NewPointSet() PointSet {
	return PointSet{}
}

// NewPointSetPreallocate simply creates a new point set with points array of the given size.
func NewPointSetPreallocate(length, capacity int) PointSet {
	if length > capacity {
		capacity = length
	}

	ps := make([]Point, length, capacity)
	return PointSet(ps)
}

// Clone returns a new copy of the point set.
func (ps PointSet) Clone() PointSet {
	points := make([]Point, len(ps))
	copy(points, ps)

	return PointSet(points)
}

// Centroid returns the average x and y coordinate of the point set.
// This can also be used for small clusters of lat/lng points.
func (ps PointSet) Centroid() Point {
	x := 0.0
	y := 0.0
	numPoints := float64(len(ps))
	for _, point := range ps {
		x += point[0]
		y += point[1]
	}
	return Point{x / numPoints, y / numPoints}
}

// DistanceFrom returns the minimum euclidean distance from the point set.
func (ps PointSet) DistanceFrom(point Point) (float64, int) {
	dist := math.Inf(1)
	index := 0

	for i := range ps {
		if d := ps[i].DistanceFromSquared(point); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}

// Bound returns a bound around the point set. Simply uses rectangular coordinates.
func (ps PointSet) Bound() Bound {
	if len(ps) == 0 {
		return NewBound(0, 0, 0, 0)
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range ps {
		minX = math.Min(minX, v.X())
		minY = math.Min(minY, v.Y())

		maxX = math.Max(maxX, v.X())
		maxY = math.Max(maxY, v.Y())
	}

	return NewBound(maxX, minX, maxY, minY)
}

// Equal compares two point sets. Returns true if lengths are the same
// and all points are Equal.
func (ps PointSet) Equal(pointSet PointSet) bool {
	if len(ps) != len(pointSet) {
		return false
	}

	for i := range ps {
		if !ps[i].Equal(pointSet[i]) {
			return false
		}
	}

	return true
}

// GeoJSON creates a new geojson feature with a multipoint geometry
// containing all the points.
func (ps PointSet) GeoJSON() *geojson.Feature {
	f := geojson.NewMultiPointFeature()
	for _, v := range ps {
		f.Geometry.MultiPoint = append(f.Geometry.MultiPoint, []float64{v[0], v[1]})
	}

	return f
}

// WKT returns the point set in WKT format,
// eg. MULTIPOINT(30 10, 10 30, 40 40)
func (ps PointSet) WKT() string {
	return ps.String()
}

// String returns a string representation of the path.
// The format is WKT, e.g. MULTIPOINT(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (ps PointSet) String() string {
	if len(ps) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTIPOINT(%g %g", ps[0][0], ps[0][1])

	for i := 1; i < len(ps); i++ {
		fmt.Fprintf(buff, ",%g %g", ps[i][0], ps[i][1])
	}

	buff.Write([]byte(")"))
	return buff.String()
}
