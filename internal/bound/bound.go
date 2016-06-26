package bound

import (
	"fmt"
	"math"
)

// Point is bound point specific to internal use.
type Point [2]float64

// NewPoint initializes a point.
func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// X returns the horizontal value.
func (p Point) X() float64 {
	return p[0]
}

// Y returns the vertical value.
func (p Point) Y() float64 {
	return p[1]
}

// Bound implements shared functionality between geo and planar bounds.
type Bound struct {
	SW Point
	NE Point
}

// New creates a new bound given the parameters.
func New(west, east, south, north float64) Bound {
	return Bound{
		SW: Point{math.Min(east, west), math.Min(north, south)},
		NE: Point{math.Max(east, west), math.Max(north, south)},
	}
}

// FromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func FromPoints(corner, oppositeCorner Point) Bound {
	b := Bound{
		SW: corner,
		NE: corner,
	}

	return b.Extend(oppositeCorner)
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {

	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b.SW[0] = math.Min(b.SW.X(), point.X())
	b.NE[0] = math.Max(b.NE.X(), point.X())

	b.SW[1] = math.Min(b.SW.Y(), point.Y())
	b.NE[1] = math.Max(b.NE.Y(), point.Y())

	return b
}

// Union extends this bounds to contain the union of this and the given bounds.
func (b Bound) Union(other Bound) Bound {
	b = b.Extend(other.SW)
	b = b.Extend(other.NE)
	b = b.Extend(Point{other.SW[0], other.NE[1]})
	b = b.Extend(Point{other.NE[0], other.SW[1]})

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {

	if point.Y() < b.SW.Y() || b.NE.Y() < point.Y() {
		return false
	}

	if point.X() < b.SW.X() || b.NE.X() < point.X() {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {

	if (b.NE[0] < bound.SW[0]) ||
		(b.SW[0] > bound.NE[0]) ||
		(b.NE[1] < bound.SW[1]) ||
		(b.SW[1] > bound.NE[1]) {
		return false
	}

	return true
}

// Center returns the center of the bound.
func (b Bound) Center() Point {
	return Point{
		(b.NE[0] + b.SW[0]) / 2.0,
		(b.NE[1] + b.SW[1]) / 2.0,
	}
}

// IsEmpty returns true if it's in some malformed negative state
// where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b.SW.X() > b.NE.X() || b.SW.Y() > b.NE.Y()
}

// IsZero will return true if the bound is the empty undefined bound.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b.SW == c.SW && b.NE == c.NE
}

// WKT returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) WKT() string {
	return fmt.Sprintf("POLYGON((%g %g, %g %g, %g %g, %g %g, %g %g))", b.SW[0], b.SW[1], b.SW[0], b.NE[1], b.NE[0], b.NE[1], b.NE[0], b.SW[1], b.SW[0], b.SW[1])
}

// MysqlIntersectsCondition returns a condition defining the intersection
// of the column and the bound. To be used in a MySQL query.
func (b Bound) MysqlIntersectsCondition(column string) string {
	return fmt.Sprintf("INTERSECTS(%s, GEOMFROMTEXT('%s'))", column, b.WKT())
}
