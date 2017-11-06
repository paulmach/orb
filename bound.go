package orb

import (
	"math"
)

// A Bound represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Bound [2]Point

// NewBound creates a new bound given the parameters.
func NewBound(left, right, bottom, top float64) Bound {
	return Bound{
		Point{left, bottom},
		Point{right, top},
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{corner, corner}.Extend(oppositeCorner)
}

// GeoJSONType returns the GeoJSON type for the object.
func (b Bound) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Bound is a 2d object.
func (b Bound) Dimensions() int {
	return 2
}

// ToPolygon converts the bound into a Polygon object.
func (b Bound) ToPolygon() Polygon {
	return Polygon{b.ToRing()}
}

// ToRing converts the bound into a loop defined
// by the boundary of the box.
func (b Bound) ToRing() Ring {
	r := make(Ring, 5)
	r[0] = b[0]
	r[1] = Point{b[1][0], b[0][1]}
	r[2] = b[1]
	r[3] = Point{b[0][0], b[1][1]}
	r[4] = b[0]

	return r
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b[0][0] = math.Min(b[0][0], point[0])
	b[1][0] = math.Max(b[1][0], point[0])

	b[0][1] = math.Min(b[0][1], point[1])
	b[1][1] = math.Max(b[1][1], point[1])

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
	if point[1] < b[0][1] || b[1][1] < point[1] {
		return false
	}

	if point[0] < b[0][0] || b[1][0] < point[0] {
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

// Center returns the center of the bounds by "averaging" the x and y coords.
func (b Bound) Center() Point {
	return Point{
		(b[0][0] + b[1][0]) / 2.0,
		(b[0][1] + b[1][1]) / 2.0,
	}
}

// Top returns the top of the bound.
func (b Bound) Top() float64 {
	return b[1][1]
}

// Bottom returns the bottom of the bound.
func (b Bound) Bottom() float64 {
	return b[0][1]
}

// Right returns the right of the bound.
func (b Bound) Right() float64 {
	return b[1][0]
}

// Left returns the left of the bound.
func (b Bound) Left() float64 {
	return b[0][0]
}

// SouthWest returns the lower left point of the bound.
func (b Bound) SouthWest() Point {
	return NewPoint(b[0][0], b[0][1])
}

// NorthEast return the upper right point of the bound.
func (b Bound) NorthEast() Point {
	return NewPoint(b[1][0], b[1][1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b[0][0] > b[1][0] || b[0][1] > b[1][1]
}

// IsZero return true if the bound just includes just null island.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b[0] == c[0] && b[1] == c[1]
}
