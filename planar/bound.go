package planar

import "github.com/paulmach/orb/internal/bound"

// A Bound represents an enclosed "box" in the 2D Euclidean plane.
type Bound struct {
	bound.Bound
}

// NewBound creates a new bound given the parameters.
func NewBound(left, right, bottom, top float64) Bound {
	return Bound{
		Bound: bound.New(left, right, bottom, top),
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{
		Bound: bound.FromPoints(bound.Point(corner), bound.Point(oppositeCorner)),
	}
}

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point.
func NewBoundAroundPoint(center Point, distance float64) Bound {
	b := NewBoundFromPoints(center, center)

	b.SW[0] -= distance
	b.SW[0] -= distance
	b.NE[1] += distance
	b.NE[1] += distance

	return b
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	b.Bound = b.Bound.Extend(bound.Point(point))
	return b
}

// Union extends this bounds to contain the union of this and the given bounds.
func (b Bound) Union(other Bound) Bound {
	b.Bound = b.Bound.Union(other.Bound)
	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	return b.Bound.Contains(bound.Point(point))
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	return b.Bound.Intersects(bound.Bound)
}

// Center returns the center of the bound.
func (b Bound) Center() Point {
	return Point(b.Bound.Center())
}

// Pad expands the bound in all directions by the amount given. The amount must be
// in the units of the bounds. Technically one can pad with negative value,
// but no error checking is done.
func (b Bound) Pad(amount float64) Bound {
	b.SW[0] -= amount
	b.SW[1] -= amount

	b.NE[0] += amount
	b.NE[1] += amount

	return b
}

// Height returns just the difference in the point's Y/Latitude.
func (b Bound) Height() float64 {
	return b.NE.Y() - b.SW.Y()
}

// Width returns just the difference in the point's X/Longitude.
func (b Bound) Width() float64 {
	return b.NE.X() - b.SW.X()
}

// Top returns the north side of the bound.
func (b Bound) Top() float64 {
	return b.NE[1]
}

// Bottom returns the south side of the bound.
func (b Bound) Bottom() float64 {
	return b.SW[1]
}

// Right returns the east side of the bound.
func (b Bound) Right() float64 {
	return b.NE[0]
}

// Left returns the west side of the bound.
func (b Bound) Left() float64 {
	return b.SW[0]
}

// Empty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) Empty() bool {
	return b.Bound.Empty()
}

// Equals returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b.SW == c.SW && b.NE == c.NE
}

// String returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) String() string {
	return b.Bound.WKT()
}

// MysqlIntersectsCondition returns a condition defining the intersection
// of the column and the bound. To be used in a MySQL query.
func (b Bound) MysqlIntersectsCondition(column string) string {
	return b.Bound.MysqlIntersectsCondition(column)
}
