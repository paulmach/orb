package planar

import "math"

// A Bound represents an enclosed "box" in the 2D Euclidean plane.
type Bound [2]Point

// NewBound creates a new bound given the parameters.
func NewBound(left, right, bottom, top float64) Bound {
	return Bound{
		Point{math.Min(left, right), math.Min(bottom, top)},
		Point{math.Max(left, right), math.Max(bottom, top)},
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{corner, corner}.Extend(oppositeCorner)
}

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point.
func NewBoundAroundPoint(center Point, distance float64) Bound {
	b := NewBoundFromPoints(center, center)

	b[0][0] -= distance
	b[0][0] -= distance
	b[1][1] += distance
	b[1][1] += distance

	return b
}

// ToRing converts the bound into a ring defined
// by the boundary of the box.
func (b Bound) ToRing() Ring {
	r := make(Ring, 5)
	r[0] = b[0]
	r[1] = Point{b[0][0], b[1][1]}
	r[2] = b[1]
	r[3] = Point{b[1][0], b[0][1]}
	r[4] = b[0]

	return r
}

// ToPolygon converts the bound into a rectangle Polygon object.
func (b Bound) ToPolygon() Polygon {
	return Polygon{b.ToRing()}
}

// DistanceFrom return the distance from the bound.
// Will return 0 if inside.
func (b Bound) DistanceFrom(p Point) float64 {
	if b.Contains(p) {
		return 0
	}

	return b.ToRing().DistanceFrom(p)
}

// IsZero will return true if the bound is the empty [0 0, 0 0] bound.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b[0][0] = math.Min(b[0].X(), point.X())
	b[1][0] = math.Max(b[1].X(), point.X())

	b[0][1] = math.Min(b[0].Y(), point.Y())
	b[1][1] = math.Max(b[1].Y(), point.Y())

	return b
}

// Union extends this bound to contain the union of this and the given bound.
func (b Bound) Union(other Bound) Bound {
	b = b.Extend(other[0])
	b = b.Extend(other[1])
	b = b.Extend(Point{other[0][0], other[1][1]})
	b = b.Extend(Point{other[1][0], other[0][1]})

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	if point.Y() < b[0].Y() || b[1].Y() < point.Y() {
		return false
	}

	if point.X() < b[0].X() || b[1].X() < point.X() {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	if (b[1][0] < bound[0][0]) ||
		(b[0][0] > bound[1][0]) ||
		(b[1][1] < bound[0][1]) ||
		(b[0][1] > bound[1][1]) {
		return false
	}

	return true
}

// Centroid returns the center of the bound.
func (b Bound) Centroid() Point {
	return Point{
		(b[1][0] + b[0][0]) / 2.0,
		(b[1][1] + b[0][1]) / 2.0,
	}
}

// Pad expands the bound in all directions by the amount given. The amount must be
// in the units of the bound. Technically one can pad with negative value,
// but no error checking is done.
func (b Bound) Pad(amount float64) Bound {
	b[0][0] -= amount
	b[0][1] -= amount

	b[1][0] += amount
	b[1][1] += amount

	return b
}

// Height returns just the difference in the point's Y/Latitude.
func (b Bound) Height() float64 {
	return b[1].Y() - b[0].Y()
}

// Width returns just the difference in the point's X/Longitude.
func (b Bound) Width() float64 {
	return b[1].X() - b[0].X()
}

// Top returns the north side of the bound.
func (b Bound) Top() float64 {
	return b[1][1]
}

// Bottom returns the south side of the bound.
func (b Bound) Bottom() float64 {
	return b[0][1]
}

// Right returns the east side of the bound.
func (b Bound) Right() float64 {
	return b[1][0]
}

// Left returns the west side of the bound.
func (b Bound) Left() float64 {
	return b[0][0]
}

// IsEmpty returns true if it's in some malformed negative state
// where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b[0].X() > b[1].X() || b[0].Y() > b[1].Y()
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b[0] == c[0] && b[1] == c[1]
}

// WKT returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) WKT() string {
	return b.ToPolygon().WKT()
}

// String returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) String() string {
	return b.WKT()
}
