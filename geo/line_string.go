package geo

import (
	"bytes"
	"fmt"
	"io"
)

// LineString represents a set of points to be thought of as a polyline.
type LineString []Point

// NewLineString creates a new line string.
func NewLineString() LineString {
	return LineString{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (ls LineString) GeoJSONType() string {
	return "LineString"
}

// Dimensions returns 1 because a LineString is a 1d object.
func (ls LineString) Dimensions() int {
	return 1
}

// Distance computes the total distance using spherical geometry.
func (ls LineString) Distance(haversine ...bool) float64 {
	yesgeo := yesHaversine(haversine)
	sum := 0.0

	loopTo := len(ls) - 1
	for i := 0; i < loopTo; i++ {
		sum += ls[i].DistanceFrom(ls[i+1], yesgeo)
	}

	return sum
}

// Reverse will reverse the line string.
// This is done inplace, ie. it modifies the original data.
func (ls LineString) Reverse() {
	l := len(ls) - 1
	for i := 0; i <= l/2; i++ {
		ls[i], ls[l-i] = ls[l-i], ls[i]
	}
}

// Bound returns a rect around the line string. Uses rectangular coordinates.
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

// String returns the wkt representation of the line string.
func (ls LineString) String() string {
	return ls.WKT()
}

func wktPoints(w io.Writer, ps []Point) {
	if len(ps) == 0 {
		w.Write([]byte(`EMPTY`))
		return
	}

	fmt.Fprintf(w, "(%g %g", ps[0][0], ps[0][1])
	for i := 1; i < len(ps); i++ {
		fmt.Fprintf(w, ",%g %g", ps[i][0], ps[i][1])
	}

	fmt.Fprintf(w, ")")
}
