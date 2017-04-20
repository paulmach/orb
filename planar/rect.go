package planar

import "github.com/paulmach/orb/internal/rect"

// A Rect represents an enclosed "box" in the 2D Euclidean plane.
type Rect struct {
	rect.Rect
}

// NewRect creates a new rectangle given the parameters.
func NewRect(left, right, bottom, top float64) Rect {
	return Rect{
		Rect: rect.New(left, right, bottom, top),
	}
}

// NewRectFromPoints creates a new rectangle given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewRectFromPoints(corner, oppositeCorner Point) Rect {
	return Rect{
		Rect: rect.FromPoints(rect.Point(corner), rect.Point(oppositeCorner)),
	}
}

// NewRectAroundPoint creates a new rectangle given a center point,
// and a distance from the center point.
func NewRectAroundPoint(center Point, distance float64) Rect {
	r := NewRectFromPoints(center, center)

	r.SW[0] -= distance
	r.SW[0] -= distance
	r.NE[1] += distance
	r.NE[1] += distance

	return r
}

// Extend grows the rectable to include the new point.
func (r Rect) Extend(point Point) Rect {
	r.Rect = r.Rect.Extend(rect.Point(point))
	return r
}

// Union extends this rectable to contain the union of this and the given rectangle.
func (r Rect) Union(other Rect) Rect {
	r.Rect = r.Rect.Union(other.Rect)
	return r
}

// Contains determines if the point is within the rectangle.
// Points on the boundary are considered within.
func (r Rect) Contains(point Point) bool {
	return r.Rect.Contains(rect.Point(point))
}

// Intersects determines if two rectangles intersect.
// Returns true if they are touching.
func (r Rect) Intersects(rectangle Rect) bool {
	return r.Rect.Intersects(rectangle.Rect)
}

// Center returns the center of the rectangle.
func (r Rect) Center() Point {
	return Point(r.Rect.Center())
}

// Pad expands the rectangle in all directions by the amount given. The amount must be
// in the units of the rectangle. Technically one can pad with negative value,
// but no error checking is done.
func (r Rect) Pad(amount float64) Rect {
	r.SW[0] -= amount
	r.SW[1] -= amount

	r.NE[0] += amount
	r.NE[1] += amount

	return r
}

// Height returns just the difference in the point's Y/Latitude.
func (r Rect) Height() float64 {
	return r.NE.Y() - r.SW.Y()
}

// Width returns just the difference in the point's X/Longitude.
func (r Rect) Width() float64 {
	return r.NE.X() - r.SW.X()
}

// Top returns the north side of the rectangle.
func (r Rect) Top() float64 {
	return r.NE[1]
}

// Bottom returns the south side of the rectangle.
func (r Rect) Bottom() float64 {
	return r.SW[1]
}

// Right returns the east side of the rectangle.
func (r Rect) Right() float64 {
	return r.NE[0]
}

// Left returns the west side of the rectangle.
func (r Rect) Left() float64 {
	return r.SW[0]
}

// IsEmpty returns true if it's in some malformed negative state
// where the left point is larger than the right.
// This can be caused by padding too much negative.
func (r Rect) IsEmpty() bool {
	return r.Rect.IsEmpty()
}

// Equal returns if two rectangles are equal.
func (r Rect) Equal(c Rect) bool {
	return r.SW == c.SW && r.NE == c.NE
}

// String returns the string respentation of the rectangle in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) String() string {
	return r.Rect.WKT()
}
