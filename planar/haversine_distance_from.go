package planar

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
)

// HaversineDistanceFromSegment returns point's haversine distance on earth from the segment [a, b] in kilometers.
func HaversineDistanceFromSegment(a, b, point orb.Point) float64 {
	x := a[0]
	y := a[1]
	dx := b[0] - x
	dy := b[1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = b[0]
			y = b[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y

	return HaversineDistance(point,orb.Point{x,y})
}

// HaversineDistanceFrom returns the distance on earth from the boundary of the geometry in
// kilometers
func HaversineDistanceFrom(g orb.Geometry, p orb.Point) float64 {
	d, _ := HaversineDistanceFromWithIndex(g, p)
	return d
}


// HaversineDistanceFromWithIndex returns the minimum haversine distance on earth in kilometers
// from the boundary of the geometry plus the index of the sub-geometry
// that was the match.
func HaversineDistanceFromWithIndex(g orb.Geometry, p orb.Point) (float64, int) {
	if g == nil {
		return math.Inf(1), -1
	}

	switch g := g.(type) {
	case orb.Point:
		return HaversineDistance(g, p), 0
	case orb.MultiPoint:
		return multiPointHaversineDistanceFrom(g, p)
	case orb.LineString:
		return lineStringHaversineDistanceFrom(g, p)
	case orb.MultiLineString:
		dist := math.Inf(1)
		index := -1
		for i, ls := range g {
			if d, _ := lineStringHaversineDistanceFrom(ls, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case orb.Ring:
		return lineStringHaversineDistanceFrom(orb.LineString(g), p)
	case orb.Polygon:
		return polygonHaversineDistanceFrom(g, p)
	case orb.MultiPolygon:
		dist := math.Inf(1)
		index := -1
		for i, poly := range g {
			if d, _ := polygonHaversineDistanceFrom(poly, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case orb.Collection:
		dist := math.Inf(1)
		index := -1
		for i, ge := range g {
			if d, _ := HaversineDistanceFromWithIndex(ge, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case orb.Bound:
		return HaversineDistanceFromWithIndex(g.ToRing(), p)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

func multiPointHaversineDistanceFrom(mp orb.MultiPoint, p orb.Point) (float64, int) {
	dist := math.Inf(1)
	index := -1

	for i := range mp {
		if d := HaversineDistance(mp[i], p); d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
}

func lineStringHaversineDistanceFrom(ls orb.LineString, p orb.Point) (float64, int) {
	dist := math.Inf(1)
	index := -1

	for i := 0; i < len(ls)-1; i++ {
		if d := segmentHaversineDistanceFrom(ls[i], ls[i+1], p); d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
}

func polygonHaversineDistanceFrom(p orb.Polygon, point orb.Point) (float64, int) {
	if len(p) == 0 {
		return math.Inf(1), -1
	}

	dist, index := lineStringHaversineDistanceFrom(orb.LineString(p[0]), point)
	for i := 1; i < len(p); i++ {
		d, i := lineStringHaversineDistanceFrom(orb.LineString(p[i]), point)
		if d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
}

func segmentHaversineDistanceFrom(p1, p2, point orb.Point) float64 {
	x := p1[0]
	y := p1[1]
	dx := p2[0] - x
	dy := p2[1] - y

	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = p2[0]
			y = p2[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y

	return HaversineDistance(point,orb.Point{x,y})
}
