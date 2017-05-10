package geo

import (
	"bytes"
	"fmt"
)

// Polygon is a closed area. The first LineString is the outer ring.
// The others are the holes. Each LineString is expected to be closed
// ie. the first point matches the last.
type Polygon []Ring

// NewPolygon creates a new Polygon.
func NewPolygon() Polygon {
	return Polygon{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (p Polygon) GeoJSONType() string {
	return "Polygon"
}

// Bound returns a bound around the polygon.
func (p Polygon) Bound() Bound {
	return p[0].Bound()
}

// WKT returns the polygon in WKT format, eg. POlYGON((0 0,1 0,1 1,0 0))
// For empty polygons the result will be 'EMPTY'.
func (p Polygon) WKT() string {
	if len(p) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "POLYGON(")
	wktPoints(buff, p[0])

	for i := 1; i < len(p); i++ {
		buff.Write([]byte(","))
		wktPoints(buff, p[i])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// Equal compares two polygons. Returns true if lengths are the same
// and all points are Equal.
func (p Polygon) Equal(polygon Polygon) bool {
	if len(p) != len(polygon) {
		return false
	}

	for i := range p {
		if !p[i].Equal(polygon[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the polygon.
// All of the rings are also cloned.
func (p Polygon) Clone() Polygon {
	np := make(Polygon, 0, len(p))
	for _, r := range p {
		np = append(np, r.Clone())
	}

	return np
}

// String returns the wkt representation of the polygon.
func (p Polygon) String() string {
	return p.WKT()
}
