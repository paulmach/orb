package planar

import (
	"bytes"
	"fmt"
	"math"
)

// Geometry is an interface that represents the shared attributes
// of a geometry.
type Geometry interface {
	Bound() Rect
	DistanceFrom(Point) float64
	Centroid() Point
	WKT() string
}

// compile time checks
var (
	_ Geometry = Point{}
	_ Geometry = Segment{}
	_ Geometry = MultiPoint{}
	_ Geometry = LineString{}
	_ Geometry = Rect{}
	_ Geometry = Polygon{}
)

// A Collection is a collection of geometries that is also a Geometry.
type Collection []Geometry

// Bound returns the bounding box of all the Geometries combined.
func (c Collection) Bound() Rect {
	r := c[0].Bound()
	for i := 1; i < len(c); i++ {
		r = r.Union(c[i].Bound())
	}

	return r
}

// DistanceFrom computes the min distance to the collection.
func (c Collection) DistanceFrom(p Point) float64 {
	dist := math.MaxFloat64
	for i := range c {
		if d := c[i].DistanceFrom(p); d < dist {
			dist = d
		}

		if dist == 0 {
			return 0
		}
	}

	return dist
}

// WKT returns the wkt of the geometry collection.
func (c Collection) WKT() string {
	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "GEOMETRYCOLLECTION(")

	for i := range c {
		buff.WriteString(c[i].WKT())
	}

	fmt.Fprintf(buff, ")")
	return buff.String()
}
