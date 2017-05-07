package geo

import (
	"bytes"
	"fmt"
)

// MultiLineString is a set of polylines.
type MultiLineString []LineString

// NewMultiLineString creates a new multi line string.
func NewMultiLineString() MultiLineString {
	return MultiLineString{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (mls MultiLineString) GeoJSONType() string {
	return "MultiLineString"
}

// Bound returns a bound around all the line strings.
func (mls MultiLineString) Bound() Bound {
	bound := mls[0].Bound()
	for i := 1; i < len(mls); i++ {
		bound = bound.Union(mls[i].Bound())
	}

	return bound
}

// Equal compares two multi line strings. Returns true if lengths are the same
// and all points are Equal.
func (mls MultiLineString) Equal(multiLineString MultiLineString) bool {
	if len(mls) != len(multiLineString) {
		return false
	}

	for i, ls := range mls {
		if !ls.Equal(multiLineString[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the multi line string.
func (mls MultiLineString) Clone() MultiLineString {
	nmls := make(MultiLineString, 0, len(mls))
	for _, ls := range mls {
		nmls = append(nmls, ls.Clone())
	}

	return nmls
}

// WKT returns the line string in WKT format, eg. MULTILINESTRING((30 10,10 30,40 40))
// For empty line strings the result will be 'EMPTY'.
func (mls MultiLineString) WKT() string {
	if len(mls) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "MULTILINESTRING(")
	wktPoints(buff, mls[0])

	for i := 1; i < len(mls); i++ {
		buff.Write([]byte(","))
		wktPoints(buff, mls[i])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns a string representation of the line string.
func (mls MultiLineString) String() string {
	return mls.WKT()
}
