package planar

import (
	"errors"
	"fmt"
	"math"
)

// Line represents the shortest path between A and B.
type Line struct {
	a, b Point
}

// NewLine creates a new line by cloning the provided points.
func NewLine(a, b Point) Line {
	return Line{a, b}
}

// DistanceFrom does NOT use spherical geometry. It finds the distance from
// the line using standard Euclidean geometry, using the units the points are in.
func (l Line) DistanceFrom(point Point) float64 {
	// yes duplicate code, but saw a 15% performance increase by removing the function call
	// return math.Sqrt(l.SquaredDistanceFrom(point))
	x := l.a[0]
	y := l.a[1]
	dx := l.b[0] - x
	dy := l.b[1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = l.b[0]
			y = l.b[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y

	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceFromSquared does NOT use spherical geometry. It finds the squared distance from
// the line using standard Euclidean geometry, using the units the points are in.
func (l Line) DistanceFromSquared(point Point) float64 {
	x := l.a[0]
	y := l.a[1]
	dx := l.b[0] - x
	dy := l.b[1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = l.b[0]
			y = l.b[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y

	return dx*dx + dy*dy
}

// Distance computes the distance of the line, ie. its length, in Euclidian space.
func (l Line) Distance() float64 {
	return l.a.DistanceFrom(l.b)
}

// DistanceSquared computes the squared distance of the line, ie. its length, in Euclidian space.
// This can save a sqrt computation.
func (l Line) DistanceSquared() float64 {
	return l.a.DistanceFromSquared(l.b)
}

// Project returns the normalized distance of the point on the line nearest the given point.
// Returned values may be outside of [0,1]. This function is the opposite of Interpolate.
func (l Line) Project(point Point) float64 {
	if point.Equal(l.a) {
		return 0.0
	}

	if point.Equal(l.b) {
		return 1.0
	}

	dx := l.b[0] - l.a[0]
	dy := l.b[1] - l.a[1]
	d := (dx*dx + dy*dy)
	if d == 0 {
		return 0
	}

	return ((point[0]-l.a[0])*dx + (point[1]-l.a[1])*dy) / d
}

// Interpolate performs a simple linear interpolation, from A to B.
// This function is the opposite of Project.
func (l Line) Interpolate(percent float64) Point {
	return Point{
		l.a[0] + percent*(l.b[0]-l.a[0]),
		l.a[1] + percent*(l.b[1]-l.a[1]),
	}
}

// Side returns -1 if the point is on the right side, 1 if on the left side, and 0 if collinear.
func (l Line) Side(p Point) int {
	val := (l.b[0]-l.a[0])*(p[1]-l.b[1]) - (l.b[1]-l.a[1])*(p[0]-l.b[0])

	if val < 0 {
		return -1 // right
	} else if val > 0 {
		return 1 // left
	}

	return 0 // collinear
}

// Intersection finds the intersection of the two lines or nil,
// if the lines are collinear will return NewPoint(math.Inf(1), math.Inf(1)) == InfinityPoint
func (l Line) Intersection(line Line) (Point, error) {
	den := (line.b[1]-line.a[1])*(l.b[0]-l.a[0]) - (line.b[0]-line.a[0])*(l.b[1]-l.a[1])
	U1 := (line.b[0]-line.a[0])*(l.a[1]-line.a[1]) - (line.b[1]-line.a[1])*(l.a[0]-line.a[0])
	U2 := (l.b[0]-l.a[0])*(l.a[1]-line.a[1]) - (l.b[1]-l.a[1])*(l.a[0]-line.a[0])

	if den == 0 {
		// collinear, all bets are off
		if U1 == 0 && U2 == 0 {
			return Point{}, errors.New("collinear") // TODO: improve
		}

		return Point{}, errors.New("nointersection") // TOOD: improve
	}

	if U1/den < 0 || U1/den > 1 || U2/den < 0 || U2/den > 1 {
		return Point{}, errors.New("nointersection") // TOOD: improve
	}

	return l.Interpolate(U1 / den), nil
}

// Intersects will return true if the lines are collinear AND intersect.
// Based on: http://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/
func (l Line) Intersects(line Line) bool {
	s1 := l.Side(line.a)
	s2 := l.Side(line.b)
	s3 := line.Side(l.a)
	s4 := line.Side(l.b)

	if s1 != s2 && s3 != s4 {
		return true
	}

	// Special Cases
	// l1 and l2.a collinear, check if l2.a is on l1
	lBound := l.Bound()
	if s1 == 0 && lBound.Contains(line.a) {
		return true
	}

	// l1 and l2.b collinear, check if l2.b is on l1
	if s2 == 0 && lBound.Contains(line.b) {
		return true
	}

	// TODO: are these next two tests redudant give the test above.
	// Thinking yes if there is round off magic.

	// l2 and l1.a collinear, check if l1.a is on l2
	lineBound := line.Bound()
	if s3 == 0 && lineBound.Contains(l.a) {
		return true
	}

	// l2 and l1.b collinear, check if l1.b is on l2
	if s4 == 0 && lineBound.Contains(l.b) {
		return true
	}

	return false
}

// Midpoint returns the Euclidean midpoint of the line.
func (l Line) Midpoint() Point {
	return Point{(l.a[0] + l.b[0]) / 2, (l.a[1] + l.b[1]) / 2}
}

// Bound returns a rectangle bound around the line. Simply uses rectangular coordinates.
func (l Line) Bound() Rect {
	return NewRect(math.Max(l.a[0], l.b[0]), math.Min(l.a[0], l.b[0]),
		math.Max(l.a[1], l.b[1]), math.Min(l.a[1], l.b[1]))
}

// Reverse swaps the start and end of the line.
func (l Line) Reverse() Line {
	l.a, l.b = l.b, l.a
	return l
}

// Equal returns the line equality and is irrespective of direction,
// i.e. true if one is the reverse of the other.
func (l Line) Equal(line Line) bool {
	return (l.a.Equal(line.a) && l.b.Equal(line.b)) || (l.a.Equal(line.b) && l.b.Equal(line.a))
}

// A returns the first point in the line.
func (l Line) A() Point {
	return l.a
}

// B returns the second point in the line.
func (l Line) B() Point {
	return l.b
}

// WKT returns the line in WKT format, eg. LINESTRING(30 10,10 30)
func (l Line) WKT() string {
	return fmt.Sprintf("LINESTRING(%g %g,%g %g)", l.a[0], l.a[1], l.b[0], l.b[1])
}

// String returns a string representation of the line.
// The format is WKT, e.g. LINESTRING(30 10,10 30)
func (l Line) String() string {
	return l.WKT()
}
