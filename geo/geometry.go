package geo

import (
	"bytes"
	"fmt"
)

// Geometry is an interface that represents the shared attributes
// of a geometry.
type Geometry interface {
	GeoJSONType() string
	Dimensions() int // e.g. 0d, 1d, 2d

	Bound() Bound
	WKT() string
}

// compile time checks
var (
	_ Geometry = Point{}
	_ Geometry = MultiPoint{}
	_ Geometry = LineString{}
	_ Geometry = MultiLineString{}
	_ Geometry = Bound{}
	_ Geometry = Ring{}
	_ Geometry = Polygon{}
	_ Geometry = MultiPolygon{}

	_ Geometry = Collection{}
)

// A Collection is a collection of geometries that is also a Geometry.
type Collection []Geometry

// GeoJSONType returns the geometry collection type.
func (c Collection) GeoJSONType() string {
	return "GeometryCollection"
}

// Dimensions returns the max of the dimensions of the collection.
func (c Collection) Dimensions() int {
	max := -1
	for _, g := range c {
		if d := g.Dimensions(); d > max {
			max = d
		}
	}

	return max
}

// Bound returns the bounding box of all the Geometries combined.
func (c Collection) Bound() Bound {
	r := c[0].Bound()
	for i := 1; i < len(c); i++ {
		r = r.Union(c[i].Bound())
	}

	return r
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
