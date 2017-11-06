package orb

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

// Dimensions returns 1 because a MultiLineString is a 2d object.
func (mls MultiLineString) Dimensions() int {
	return 2
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
