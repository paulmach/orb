package geo

import (
	"math"

	"github.com/paulmach/orb"
)

// Ring represents a set of ring on the earth.
type Ring LineString

// NewRing creates a new ring.
func NewRing() Ring {
	return Ring{}
}

// GeoJSONType returns the GeoJSON type for the object.
func (r Ring) GeoJSONType() string {
	return "Polygon"
}

// Valid will return true if the ring is a real ring.
// ie. 4+ points and the first and last points match.
// NOTE: this will not check for self-intersection.
func (r Ring) Valid() bool {
	if len(r) < 4 {
		return false
	}

	// first must equal last
	return r[0] == r[len(r)-1]
}

// Reverse changes the direction of the ring.
// This is done inplace, ie. it modifies the original data.
func (r Ring) Reverse() {
	LineString(r).Reverse()
}

// Bound returns a rect around the ring. Uses rectangular coordinates.
func (r Ring) Bound() Bound {
	return MultiPoint(r).Bound()
}

// Area calculate the approximate area of the ring.
func (r Ring) Area() float64 {
	return math.Abs(r.SignedArea())
}

// SignedArea calculate the approximate area of the ring.
// Area will be positive if ring is oriented counter-clockwise,
// otherwise it will be negative.
func (r Ring) SignedArea() float64 {
	if len(r) < 3 {
		return 0
	}
	var lo, mi, hi int

	l := len(r)
	if r[0] != r[len(r)-1] {
		// if not a closed ring, add an implicit calc for that last point.
		l++
	}

	// To support implicit closing of ring, replace references to
	// the last point in r to the first 1.

	area := 0.0
	for i := 0; i < l; i++ {
		if i == l-3 { // i = N-3
			lo = l - 3
			mi = l - 2
			hi = 0
		} else if i == l-2 { // i = N-2
			lo = l - 2
			mi = 0
			hi = 0
		} else if i == l-1 { // i = N-1
			lo = 0
			mi = 0
			hi = 1
		} else { // i = 0 to N-3
			lo = i
			mi = i + 1
			hi = i + 2
		}

		area += (rad(r[hi][0]) - rad(r[lo][0])) * math.Sin(rad(r[mi][1]))
	}

	return -area * orb.EarthRadius * orb.EarthRadius / 2
}

// Orientation returns 1 if the the ring is in couter-clockwise order,
// return -1 if the ring is the clockwise order and 0 if the ring is
// degenerate and had no area.
func (r Ring) Orientation() orb.Orientation {
	area := 0.0

	// This is a fast planar area computation, which is okay for this use.
	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		area += (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
	}

	if area > 0 {
		return orb.CCW
	}

	if area < 0 {
		return orb.CW
	}

	// degenerate case, no area
	return 0
}

// Equal compares two rings. Returns true if lengths are the same
// and all points are Equal.
func (r Ring) Equal(ring Ring) bool {
	return MultiPoint(r).Equal(MultiPoint(ring))
}

// Clone returns a new copy of the ring.
func (r Ring) Clone() Ring {
	ps := MultiPoint(r)
	return Ring(ps.Clone())
}

// WKT returns the ring in WKT format, eg. POLYGON((30 10,10 30,40 40))
// For empty rings the result will be 'EMPTY'.
func (r Ring) WKT() string {
	polygon := Polygon{r}
	return polygon.WKT()
}

// String returns the wkt representation of the ring.
func (r Ring) String() string {
	return r.WKT()
}

func rad(n float64) float64 {
	return n * math.Pi / 180
}
