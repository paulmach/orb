package orb

import (
	"math"
)

// A MultiPoint represents a set of points in the 2D Eucledian or Cartesian plane.
type MultiPoint []Point

// NewMultiPoint simply creates a new MultiPoint object.
func NewMultiPoint() MultiPoint {
	return MultiPoint{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (mp MultiPoint) GeoJSONType() string {
	return "MultiPoint"
}

// Dimensions returns 0 because a MultiPoint is a 0d object.
func (mp MultiPoint) Dimensions() int {
	return 2
}

// Clone returns a new copy of the points.
func (mp MultiPoint) Clone() MultiPoint {
	points := make([]Point, len(mp))
	copy(points, mp)

	return MultiPoint(points)
}

// Bound returns a bound around the points. Uses rectangular coordinates.
func (mp MultiPoint) Bound() Bound {
	if len(mp) == 0 {
		return Bound{}
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range mp {
		minX = math.Min(minX, v[0])
		minY = math.Min(minY, v[1])

		maxX = math.Max(maxX, v[0])
		maxY = math.Max(maxY, v[1])
	}

	return NewBound(minX, maxX, minY, maxY)
}

// Equal compares two MultiPoint objects. Returns true if lengths are the same
// and all points are Equal, and in the same order.
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
