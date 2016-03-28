package geo

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

// Centroid uses an algorithm to compute the centroid of points
// on the earth's surface. The points are first projected into 3D space then
// averaged. The result is projected back onto the sphere. This method is about
func (ps PointSet) Centroid() Point {

	// Implementation sourced from Geolib
	// https://github.com/manuelbieh/Geolib/blob/74593bf93f9a99d5ce7e6bcefa367c5a78f5321b/src/geolib.js#L416
	var x, y, z float64

	for _, p := range ps {
		lngSin, lngCos := math.Sincos(deg2rad(p[0]))
		latSin, latCos := math.Sincos(deg2rad(p[1]))

		x += latCos * lngCos
		y += latCos * lngSin
		z += latSin
	}

	np := float64(len(ps))
	x /= np
	y /= np
	z /= np

	return Point{
		rad2deg(math.Atan2(y, x)),
		rad2deg(math.Atan2(z, math.Sqrt(x*x+y*y))),
	}
}

// DistanceFrom returns the minimum geo distance from the point set,
// along with the index of the point with minimum index.
func (ps PointSet) DistanceFrom(point Point) (float64, int) {
	dist := math.Inf(1)
	index := 0

	for i := range ps {
		if d := ps[i].DistanceFrom(point); d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
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
		minX = math.Min(minX, v.Lng())
		minY = math.Min(minY, v.Lat())

		maxX = math.Max(maxX, v.Lng())
		maxY = math.Max(maxY, v.Lat())
	}

	return NewBound(maxX, minX, maxY, minY)
}

// Equals compares two point sets. Returns true if lengths are the same
// and all points are Equal
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
