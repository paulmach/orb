package rect

import (
	"fmt"
	"math"
)

// Point is rect point specific to internal use.
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

// Rect implements shared functionality between geo and planar bounds.
type Rect struct {
	SW Point
	NE Point
}

// New creates a new rectangle given the parameters.
func New(west, east, south, north float64) Rect {
	return Rect{
		SW: Point{math.Min(east, west), math.Min(north, south)},
		NE: Point{math.Max(east, west), math.Max(north, south)},
	}
}

// FromPoints creates a new rectangle given two opposite corners.
// These corners can be either sw/ne or se/nw.
func FromPoints(corner, oppositeCorner Point) Rect {
	r := Rect{
		SW: corner,
		NE: corner,
	}

	return r.Extend(oppositeCorner)
}

// Extend grows the rectangle to include the new point.
func (r Rect) Extend(point Point) Rect {

	// already included, no big deal
	if r.Contains(point) {
		return r
	}

	r.SW[0] = math.Min(r.SW.X(), point.X())
	r.NE[0] = math.Max(r.NE.X(), point.X())

	r.SW[1] = math.Min(r.SW.Y(), point.Y())
	r.NE[1] = math.Max(r.NE.Y(), point.Y())

	return r
}

// Union extends this rectangles to contain the union of this and the given rectangles.
func (r Rect) Union(other Rect) Rect {
	r = r.Extend(other.SW)
	r = r.Extend(other.NE)
	r = r.Extend(Point{other.SW[0], other.NE[1]})
	r = r.Extend(Point{other.NE[0], other.SW[1]})

	return r
}

// Contains determines if the point is within the rectangle.
// Points on the boundary are considered within.
func (r Rect) Contains(point Point) bool {

	if point.Y() < r.SW.Y() || r.NE.Y() < point.Y() {
		return false
	}

	if point.X() < r.SW.X() || r.NE.X() < point.X() {
		return false
	}

	return true
}

// Intersects determines if two rectangles intersect.
// Returns true if they are touching.
func (r Rect) Intersects(rect Rect) bool {

	if (r.NE[0] < rect.SW[0]) ||
		(r.SW[0] > rect.NE[0]) ||
		(r.NE[1] < rect.SW[1]) ||
		(r.SW[1] > rect.NE[1]) {
		return false
	}

	return true
}

// Center returns the center of the rect.
func (r Rect) Center() Point {
	return Point{
		(r.NE[0] + r.SW[0]) / 2.0,
		(r.NE[1] + r.SW[1]) / 2.0,
	}
}

// IsEmpty returns true if it's in some malformed negative state
// where the left point is larger than the right.
// This can be caused by padding too much negative.
func (r Rect) IsEmpty() bool {
	return r.SW.X() > r.NE.X() || r.SW.Y() > r.NE.Y()
}

// IsZero will return true if the rectangle is the empty undefined rectangle.
func (r Rect) IsZero() bool {
	return r == Rect{}
}

// Equal returns if two rectangles are equal.
func (r Rect) Equal(c Rect) bool {
	return r.SW == c.SW && r.NE == c.NE
}

// WKT returns the string respentation of the rectangle in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) WKT() string {
	return fmt.Sprintf("POLYGON((%g %g, %g %g, %g %g, %g %g, %g %g))",
		r.SW[0], r.SW[1], r.SW[0], r.NE[1], r.NE[0], r.NE[1], r.NE[0], r.SW[1], r.SW[0], r.SW[1])
}
