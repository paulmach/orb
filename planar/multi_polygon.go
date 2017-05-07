package planar

import (
	"bytes"
	"fmt"
	"math"
)

// MultiPolygon is a set of polygons.
type MultiPolygon []Polygon

// NewMultiPolygon creates a new MultiPolygon.
func NewMultiPolygon() MultiPolygon {
	return MultiPolygon{}
}

// DistanceFrom will return the distance from the point to
// the closest polygon. Returns 0 if the point is within the polygon.
func (mp MultiPolygon) DistanceFrom(point Point) float64 {
	dist := math.MaxFloat64

	for _, p := range mp {
		d := p.DistanceFrom(point)
		if d == 0 {
			return 0
		}

		if d < dist {
			dist = d
		}
	}

	return dist
}

// Centroid computes the area based centroid of the polygon.
// The algorithm removes the contribution of the holes.
func (mp MultiPolygon) Centroid() Point {
	point := Point{}
	area := 0.0

	for _, p := range mp {
		c, a := p.CentroidArea()

		point[0] += c[0] * a
		point[1] += c[1] * a

		area += a
	}

	point[0] /= area
	point[1] /= area

	return point
}

// Contains checks if the point is within any of the polygons.
func (mp MultiPolygon) Contains(point Point) bool {
	for _, p := range mp {
		if p.Contains(point) {
			return true
		}
	}

	return false
}

// Area computes the sum of all the polygons.
func (mp MultiPolygon) Area() float64 {
	area := 0.0

	for _, p := range mp {
		area += p.Area()
	}

	return area
}

// Bound returns a bound around the multi-polygon.
func (mp MultiPolygon) Bound() Bound {
	bound := mp[0].Bound()
	for i := 1; i < len(mp); i++ {
		bound = bound.Union(mp[i].Bound())
	}

	return bound
}

// WKT returns the polygon in WKT format, eg. POlYGON((0 0,1 0,1 1,0 0))
// For empty polygons the result will be 'EMPTY'.
func (mp MultiPolygon) WKT() string {
	if len(mp) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTIPOLYGON(")

	for i, p := range mp {
		if i != 0 {
			buff.Write([]byte(","))
		}

		buff.Write([]byte("("))
		for j, r := range p {
			if j != 0 {
				buff.Write([]byte(","))

			}
			wktPoints(buff, r)
		}
		buff.Write([]byte(")"))
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns the wkt representation of the multi polygon.
func (mp MultiPolygon) String() string {
	return mp.WKT()
}

// Equal compares two multi-polygons.
func (mp MultiPolygon) Equal(multiPolygon MultiPolygon) bool {
	if len(mp) != len(multiPolygon) {
		return false
	}

	for i, p := range mp {
		if !p.Equal(multiPolygon[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the multi-polygon.
func (mp MultiPolygon) Clone() MultiPolygon {
	nmp := make(MultiPolygon, 0, len(mp))
	for _, p := range mp {
		nmp = append(nmp, p.Clone())
	}

	return nmp
}
