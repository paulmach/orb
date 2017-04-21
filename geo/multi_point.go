package geo

import (
	"bytes"
	"fmt"
	"math"

	"github.com/paulmach/go.geojson"
)

// A MultiPoint represents a set of points in the 2D Eucledian or Cartesian plane.
type MultiPoint []Point

// NewMultiPoint simply creates a new MultiPoint object.
func NewMultiPoint() MultiPoint {
	return MultiPoint{}
}

// NewMultiPointPreallocate creates a new MultiPoint with points array of the given size.
func NewMultiPointPreallocate(length, capacity int) MultiPoint {
	if length > capacity {
		capacity = length
	}

	mp := make([]Point, length, capacity)
	return MultiPoint(mp)
}

// Clone returns a new copy of the points.
func (mp MultiPoint) Clone() MultiPoint {
	points := make([]Point, len(mp))
	copy(points, mp)

	return MultiPoint(points)
}

// Centroid uses an algorithm to compute the centroid of points
// on the earth's surface. The points are first projected into 3D space then
// averaged. The result is projected back onto the sphere. This method is about
func (mp MultiPoint) Centroid() Point {

	// Implementation sourced from Geolib
	// https://github.com/manuelbieh/Geolib/blob/74593bf93f9a99d5ce7e6bcefa367c5a78f5321b/src/geolib.js#L416
	var x, y, z float64

	for _, p := range mp {
		lonSin, lonCos := math.Sincos(deg2rad(p[0]))
		latSin, latCos := math.Sincos(deg2rad(p[1]))

		x += latCos * lonCos
		y += latCos * lonSin
		z += latSin
	}

	np := float64(len(mp))
	x /= np
	y /= np
	z /= np

	return Point{
		rad2deg(math.Atan2(y, x)),
		rad2deg(math.Atan2(z, math.Sqrt(x*x+y*y))),
	}
}

// DistanceFrom returns the minimum geo distance from the points,
// along with the index of the point with minimum index.
func (mp MultiPoint) DistanceFrom(point Point) (float64, int) {
	dist := math.Inf(1)
	index := 0

	for i := range mp {
		if d := mp[i].DistanceFrom(point); d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
}

// Bound returns a rectangle bound around the points. Uses rectangular coordinates.
func (mp MultiPoint) Bound() Rect {
	if len(mp) == 0 {
		return Rect{}
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range mp {
		minX = math.Min(minX, v.Lon())
		minY = math.Min(minY, v.Lat())

		maxX = math.Max(maxX, v.Lon())
		maxY = math.Max(maxY, v.Lat())
	}

	return NewRect(maxX, minX, maxY, minY)
}

// Equal compares two MultiPoint objects. Returns true if lengths are the same
// and all points are Equal
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

// GeoJSON creates a new geojson feature with a multipoint geometry
// containing all the points.
func (mp MultiPoint) GeoJSON() *geojson.Feature {
	f := geojson.NewMultiPointFeature()
	for _, p := range mp {
		f.Geometry.MultiPoint = append(f.Geometry.MultiPoint, []float64{p[0], p[1]})
	}

	return f
}

// WKT returns the points in WKT format,
// eg. MULTIPOINT(30 10, 10 30, 40 40)
func (mp MultiPoint) WKT() string {
	if len(mp) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTIPOINT(%g %g", mp[0][0], mp[0][1])

	for i := 1; i < len(mp); i++ {
		fmt.Fprintf(buff, ",%g %g", mp[i][0], mp[i][1])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns a string representation of the path.
// The format is WKT, e.g. MULTIPOINT(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (mp MultiPoint) String() string {
	return mp.WKT()
}
