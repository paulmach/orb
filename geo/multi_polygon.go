package geo

import (
	"bytes"
	"fmt"
)

// MultiPolygon is a set of polygons.
type MultiPolygon []Polygon

// NewMultiPolygon creates a new MultiPolygon.
func NewMultiPolygon() MultiPolygon {
	return MultiPolygon{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (mp MultiPolygon) GeoJSONType() string {
	return "MultiPolygon"
}

// Dimensions returns 2 because a MultiPolygon is a 2d object.
func (mp MultiPolygon) Dimensions() int {
	return 2
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
