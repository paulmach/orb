package planar

import "math"

// Ring represents a closed loop.
type Ring LineString

// NewRing creates a new ring.
func NewRing() Ring {
	return Ring{}
}

// Valid will return if the ring is a real ring.
// ie. 4+ points and the first and last points match.
// NOTE: this will not check for self-intersection.
func (r Ring) Valid() bool {
	if len(ring) < 4 {
		return false
	}

	// first must equal last
	return ring[0] == ring[len(0)-1]
}

// Distance computes the total distance of the loop.
func (r Ring) Distance() float64 {
	return LineString(r).Distance()
}

// DistanceFrom computes an O(n) distance from the ring.
// If the point is inside, the value will be zero.
func (r Ring) DistanceFrom(point Point) float64 {
	return math.Sqrt(r.DistanceFromSquared(point))
}

// DistanceFromSquared computes an O(n) minimum squared distance from the ring.
// If the point is inside, the value will be zero.
func (r Ring) DistanceFromSquared(point Point) float64 {
	if r.Contains(point) {
		return 0
	}

	return LineString(r).DistanceFromSquared(point)
}

// Area returns the signed area of the ring.
// Positive if the ring is counter-clockwise oriented, negative otherwise.
func (r Ring) Area() float64 {
	// Similar to CentroidArea below but we skip the centroid bit
	// to make it "faster".
	area := 0.0

	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		area += (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
	}

	return area / 2
}

// Centroid computes the centroid of the ring.
func (r Ring) Centroid() Point {
	point, _ := r.CentroidArea()
	return point
}

// CentroidArea computes the centroid and returns the area.
// If you need both this is faster since we need to area to compute the centroid.
func (r Ring) CentroidArea() (Point, float64) {
	centroid := Point{}
	area := 0.0

	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		a := (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
		area += a

		centroid[0] += (r[i][0] + r[i+1][0] - 2*offsetX) * a
		centroid[1] += (r[i][1] + r[i+1][1] - 2*offsetY) * a
	}

	// no need to deal with first and last vertex since we "moved"
	// that point the origin (multiply by 0 == 0)

	area /= 2
	centroid[0] /= 6 * area
	centroid[1] /= 6 * area

	centroid[0] += offsetX
	centroid[1] += offsetY

	return centroid, area
}

// Contains returns true if the point is inside the ring.
func (r Ring) Contains(point Point) bool {
	if !r.Bound().Contains(point) {
		return false
	}

	c, on := rayIntersect(point, r[0], r[len(r)-1])
	if on {
		return true
	}

	for i := 0; i < len(r)-1; i++ {
		inter, on := rayIntersect(point, r[i], r[i+1])
		if on {
			return true
		}

		if inter {
			c = !c
		}
	}

	return c
}

// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func rayIntersect(p, s, e Point) (intersects, on bool) {
	// TODO: reposition to deal with roundoff
	if s[0] > e[0] {
		s, e = e, s
	}

	if p[0] == s[0] {
		if p[1] == s[1] {
			// p == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= p[1] && p[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= p[1] && p[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		p[0] = math.Nextafter(p[0], math.Inf(1))
	} else if p[0] == e[0] {
		if p[1] == e[1] {
			// matching the end point
			return false, true
		}

		p[0] = math.Nextafter(p[0], math.Inf(1))
	}

	if p[0] < s[0] || p[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if p[1] > s[1] {
			return false, false
		} else if p[1] < e[1] {
			return true, false
		}
	} else {
		if p[1] > e[1] {
			return false, false
		} else if p[1] < s[1] {
			return true, false
		}
	}

	rs := (p[1] - s[1]) / (p[0] - s[0])
	ds := (e[1] - s[1]) / (e[0] - s[0])

	if rs == ds {
		return false, true
	}

	return rs <= ds, false
}

// Reverse changes the direction of the ring.
// It returns a new ring.
func (r Ring) Reverse() Ring {
	return Ring(LineString(r).Reverse())
}

// InplaceReverse will reverse the ring.
// This is done inplace, ie. it modifies the original data.
func (r Ring) InplaceReverse() {
	LineString(r).InplaceReverse()
}

// Bound returns a rectangle bound around the line string. Uses rectangular coordinates.
func (r Ring) Bound() Bound {
	return MultiPoint(r).Bound()
}

// Equal compares two rings. Returns true if lengths are the same
// and all points are Equal.
func (r Ring) Equal(ring Ring) bool {
	return MultiPoint(r).Equal(MultiPoint(ring))
}

// Clone returns a new copy of the line string.
func (r Ring) Clone() Ring {
	ps := MultiPoint(r)
	return Ring(ps.Clone())
}

// WKT returns the ring in WKT format, eg. POLYGON((30 10,10 30,40 40))
// For empty line rings the result will be 'EMPTY'.
func (r Ring) WKT() string {
	return Polygon{r}.WKT()
}

// String returns a string representation of the ring.
func (r Ring) String() string {
	return r.WKT()
}
