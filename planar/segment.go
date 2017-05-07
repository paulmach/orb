package planar

import (
	"errors"
	"fmt"
	"math"
)

// A Segment represents the shortest path between A and B.
type Segment [2]Point

// NewSegment creates a new segment by cloning the provided points.
func NewSegment(a, b Point) Segment {
	return Segment{a, b}
}

// DistanceFrom finds the distance from the segment using Euclidean
// geometry, using the units the points are in.
func (s Segment) DistanceFrom(point Point) float64 {
	// yes duplicate code, but saw a 15% performance increase by removing the function call
	// return math.Sqrt(s.SquaredDistanceFrom(point))
	x := s[0][0]
	y := s[0][1]
	dx := s[1][0] - x
	dy := s[1][1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = s[1][0]
			y = s[1][1]
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
// the segment using standard Euclidean geometry, using the units the points are in.
func (s Segment) DistanceFromSquared(point Point) float64 {
	x := s[0][0]
	y := s[0][1]
	dx := s[1][0] - x
	dy := s[1][1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = s[1][0]
			y = s[1][1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y

	return dx*dx + dy*dy
}

// Distance computes the distance of the segment, ie. its length, in Euclidian space.
func (s Segment) Distance() float64 {
	return s[0].DistanceFrom(s[1])
}

// DistanceSquared computes the squared distance of the segment, ie. its length, in Euclidian space.
// This can save a sqrt computation.
func (s Segment) DistanceSquared() float64 {
	return s[0].DistanceFromSquared(s[1])
}

// Project returns the normalized distance of the point on the segment nearest the given point.
// Returned values may be outside of [0,1]. This function is the opposite of Interpolate.
func (s Segment) Project(point Point) float64 {
	if point.Equal(s[0]) {
		return 0.0
	}

	if point.Equal(s[1]) {
		return 1.0
	}

	dx := s[1][0] - s[0][0]
	dy := s[1][1] - s[0][1]
	d := (dx*dx + dy*dy)
	if d == 0 {
		return 0
	}

	return ((point[0]-s[0][0])*dx + (point[1]-s[0][1])*dy) / d
}

// Interpolate performs a linear interpolation, from A to B.
// This function is the opposite of Project.
func (s Segment) Interpolate(percent float64) Point {
	return Point{
		s[0][0] + percent*(s[1][0]-s[0][0]),
		s[0][1] + percent*(s[1][1]-s[0][1]),
	}
}

// Centroid computes the midpoint or centroid of the segment.
func (s Segment) Centroid() Point {
	return Point{
		s[0][0] + 0.5*(s[1][0]-s[0][0]),
		s[0][1] + 0.5*(s[1][1]-s[0][1]),
	}
}

// Side returns -1 if the point is on the right side, 1 if on the left side, and 0 if collinear.
func (s Segment) Side(p Point) int {
	val := (s[1][0]-s[0][0])*(p[1]-s[1][1]) - (s[1][1]-s[0][1])*(p[0]-s[1][0])

	if val < 0 {
		return -1 // right
	} else if val > 0 {
		return 1 // left
	} else {
		return 0 // collinear
	}
}

// Intersection finds the intersection of the two segments or nil,
// if the segments are collinear will return NewPoint(math.Inf(1), math.Inf(1)) == InfinityPoint
func (s Segment) Intersection(seg Segment) (Point, error) {
	den := (seg[1][1]-seg[0][1])*(s[1][0]-s[0][0]) - (seg[1][0]-seg[0][0])*(s[1][1]-s[0][1])
	U1 := (seg[1][0]-seg[0][0])*(s[0][1]-seg[0][1]) - (seg[1][1]-seg[0][1])*(s[0][0]-seg[0][0])
	U2 := (s[1][0]-s[0][0])*(s[0][1]-seg[0][1]) - (s[1][1]-s[0][1])*(s[0][0]-seg[0][0])

	if den == 0 {
		// collinear, all bets are off
		if U1 == 0 && U2 == 0 {
			return Point{}, errors.New("collinear") // TODO: improve
		}

		return Point{}, errors.New("nointersection") // TODO: improve
	}

	if U1/den < 0 || U1/den > 1 || U2/den < 0 || U2/den > 1 {
		return Point{}, errors.New("nointersection") // TODO: improve
	}

	return s.Interpolate(U1 / den), nil
}

// Intersects will return true if the segments are collinear AND intersect.
// Based on: http://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/
func (s Segment) Intersects(seg Segment) bool {
	s1 := s.Side(seg[0])
	s2 := s.Side(seg[1])
	s3 := seg.Side(s[0])
	s4 := seg.Side(s[1])

	if s1 != s2 && s3 != s4 {
		return true
	}

	// Special Cases
	// l1 and l2[0] collinear, check if l2[0] is on l1
	lBound := s.Bound()
	if s1 == 0 && lBound.Contains(seg[0]) {
		return true
	}

	// l1 and l2[1] collinear, check if l2[1] is on l1
	if s2 == 0 && lBound.Contains(seg[1]) {
		return true
	}

	// TODO: are these next two tests redudant give the test above.
	// Thinking yes if there is round off magic.

	// s2 and s1[0] collinear, check if s1[0] is on s2
	segBound := seg.Bound()
	if s3 == 0 && segBound.Contains(s[0]) {
		return true
	}

	// s2 and s1[1] collinear, check if s1[1] is on s2
	if s4 == 0 && segBound.Contains(s[1]) {
		return true
	}

	return false
}

// Midpoint returns the Euclidean midpoint of the segment.
func (s Segment) Midpoint() Point {
	return Point{(s[0][0] + s[1][0]) / 2, (s[0][1] + s[1][1]) / 2}
}

// Bound returns a rectangle bound around the segment.
func (s Segment) Bound() Bound {
	return NewBound(
		math.Min(s[0][0], s[1][0]),
		math.Max(s[0][0], s[1][0]),
		math.Min(s[0][1], s[1][1]),
		math.Max(s[0][1], s[1][1]),
	)
}

// Reverse swaps the start and end of the segment.
func (s Segment) Reverse() Segment {
	s[0], s[1] = s[1], s[0]
	return s
}

// Equal returns the segments equality. Segments of different direction
// will not be equal
func (s Segment) Equal(seg Segment) bool {
	return (s[0].Equal(seg[0]) && s[1].Equal(seg[1]))
}

// A returns the first point in the segment.
func (s Segment) A() Point {
	return s[0]
}

// B returns the second point in the segment.
func (s Segment) B() Point {
	return s[1]
}

// WKT returns the segment in WKT format, eg. LINESTRING(30 10,10 30)
func (s Segment) WKT() string {
	return fmt.Sprintf("LINESTRING(%g %g,%g %g)", s[0][0], s[0][1], s[1][0], s[1][1])
}

// String returns a string representation of the segment.
// The format is WKT, e.g. LINESTRING(30 10,10 30)
func (s Segment) String() string {
	return s.WKT()
}
