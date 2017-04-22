package planar

import "math"

// A Rect represents an enclosed "box" in the 2D Euclidean plane.
type Rect [2]Point

// NewRect creates a new rectangle given the parameters.
func NewRect(left, right, bottom, top float64) Rect {
	return Rect{
		Point{math.Min(left, right), math.Min(bottom, top)},
		Point{math.Max(left, right), math.Max(bottom, top)},
	}
}

// NewRectFromPoints creates a new rectangle given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewRectFromPoints(corner, oppositeCorner Point) Rect {
	return Rect{corner, corner}.Extend(oppositeCorner)
}

// NewRectAroundPoint creates a new rectangle given a center point,
// and a distance from the center point.
func NewRectAroundPoint(center Point, distance float64) Rect {
	r := NewRectFromPoints(center, center)

	r[0][0] -= distance
	r[0][0] -= distance
	r[1][1] += distance
	r[1][1] += distance

	return r
}

// ToLineString converts the Rect into a loop defined
// by the Rect boundary.
func (r Rect) ToLineString() LineString {
	return append(NewLineStringPreallocate(0, 5),
		r[0],
		NewPoint(r[0][0], r[1][1]),
		r[1],
		NewPoint(r[1][0], r[0][1]),
		r[0],
	)
}

// ToPolygon converts the Rect into a rectangle Polygon object.
func (r Rect) ToPolygon() Polygon {
	return Polygon{r.ToLineString()}
}

// DistanceFrom return the distance from the rectangle.
// Will return 0 if inside.
func (r Rect) DistanceFrom(p Point) float64 {
	if r.Contains(p) {
		return 0
	}

	return r.ToLineString().DistanceFrom(p)
}

// IsZero will return true if the rectangle is the empty [0 0, 0 0] rectangle.
func (r Rect) IsZero() bool {
	return r == Rect{}
}

// Bound returns the the same rectangle.
func (r Rect) Bound() Rect {
	return r
}

// Extend grows the rectable to include the new point.
func (r Rect) Extend(point Point) Rect {
	// already included, no big deal
	if r.Contains(point) {
		return r
	}

	r[0][0] = math.Min(r[0].X(), point.X())
	r[1][0] = math.Max(r[1].X(), point.X())

	r[0][1] = math.Min(r[0].Y(), point.Y())
	r[1][1] = math.Max(r[1].Y(), point.Y())

	return r
}

// Union extends this rectable to contain the union of this and the given rectangle.
func (r Rect) Union(other Rect) Rect {
	r = r.Extend(other[0])
	r = r.Extend(other[1])
	r = r.Extend(Point{other[0][0], other[1][1]})
	r = r.Extend(Point{other[1][0], other[0][1]})

	return r
}

// Contains determines if the point is within the rectangle.
// Points on the boundary are considered within.
func (r Rect) Contains(point Point) bool {
	if point.Y() < r[0].Y() || r[1].Y() < point.Y() {
		return false
	}

	if point.X() < r[0].X() || r[1].X() < point.X() {
		return false
	}

	return true
}

// Intersects determines if two rectangles intersect.
// Returns true if they are touching.
func (r Rect) Intersects(rect Rect) bool {
	if (r[1][0] < rect[0][0]) ||
		(r[0][0] > rect[1][0]) ||
		(r[1][1] < rect[0][1]) ||
		(r[0][1] > rect[1][1]) {
		return false
	}

	return true
}

// Centroid returns the center of the rectangle.
func (r Rect) Centroid() Point {
	return Point{
		(r[1][0] + r[0][0]) / 2.0,
		(r[1][1] + r[0][1]) / 2.0,
	}
}

// Pad expands the rectangle in all directions by the amount given. The amount must be
// in the units of the rectangle. Technically one can pad with negative value,
// but no error checking is done.
func (r Rect) Pad(amount float64) Rect {
	r[0][0] -= amount
	r[0][1] -= amount

	r[1][0] += amount
	r[1][1] += amount

	return r
}

// Height returns just the difference in the point's Y/Latitude.
func (r Rect) Height() float64 {
	return r[1].Y() - r[0].Y()
}

// Width returns just the difference in the point's X/Longitude.
func (r Rect) Width() float64 {
	return r[1].X() - r[0].X()
}

// Top returns the north side of the rectangle.
func (r Rect) Top() float64 {
	return r[1][1]
}

// Bottom returns the south side of the rectangle.
func (r Rect) Bottom() float64 {
	return r[0][1]
}

// Right returns the east side of the rectangle.
func (r Rect) Right() float64 {
	return r[1][0]
}

// Left returns the west side of the rectangle.
func (r Rect) Left() float64 {
	return r[0][0]
}

// IsEmpty returns true if it's in some malformed negative state
// where the left point is larger than the right.
// This can be caused by padding too much negative.
func (r Rect) IsEmpty() bool {
	return r[0].X() > r[1].X() || r[0].Y() > r[1].Y()
}

// Equal returns if two rectangles are equal.
func (r Rect) Equal(c Rect) bool {
	return r[0] == c[0] && r[1] == c[1]
}

// WKT returns the string respentation of the rectangle in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) WKT() string {
	return r.ToPolygon().WKT()
}

// String returns the string respentation of the rectangle in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) String() string {
	return r.WKT()
}
