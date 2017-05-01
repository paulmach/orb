package geo

import (
	"bytes"
	"fmt"
)

// Geometry is an interface that represents the shared attributes
// of a geometry.
type Geometry interface {
	GeoJSONType() string

	Bound() Rect
	WKT() string
}

// compile time checks
var (
	_ Geometry = Point{}
	_ Geometry = MultiPoint{}
	_ Geometry = LineString{}
	_ Geometry = MultiLineString{}
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
