package planar

import (
	"fmt"
	"math"

	"github.com/paulmach/orb/geo"
)

// DistanceFrom returns the distance from the boundary of the geometry in
// the units of the geometry.
func DistanceFrom(g geo.Geometry, p geo.Point) float64 {
	d, _ := DistanceFromWithIndex(g, p)
	return d
}

// DistanceFromWithIndex returns the minimum euclidean distance
// from the boundary of the geometry plus the index of the sub-geometry
// that was the match.
func DistanceFromWithIndex(g geo.Geometry, p geo.Point) (float64, int) {
	switch g := g.(type) {
	case geo.Point:
		return Distance(g, p), 0
	case geo.MultiPoint:
		return multiPointDistanceFrom(g, p)
	case geo.LineString:
		return lineStringDistanceFrom(g, p)
	case geo.MultiLineString:
		dist := math.Inf(1)
		index := -1
		for i, ls := range g {
			if d, _ := lineStringDistanceFrom(ls, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case geo.Ring:
		return lineStringDistanceFrom(geo.LineString(g), p)
	case geo.Polygon:
		return polygonDistanceFrom(g, p)
	case geo.MultiPolygon:
		dist := math.Inf(1)
		index := -1
		for i, poly := range g {
			if d, _ := polygonDistanceFrom(poly, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case geo.Collection:
		dist := math.Inf(1)
		index := -1
		for i, ge := range g {
			if d, _ := DistanceFromWithIndex(ge, p); d < dist {
				dist = d
				index = i
			}
		}

		return dist, index
	case geo.Bound:
		return DistanceFromWithIndex(g.ToRing(), p)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

func multiPointDistanceFrom(mp geo.MultiPoint, p geo.Point) (float64, int) {
	dist := math.Inf(1)
	index := -1

	for i := range mp {
		if d := DistanceSquared(mp[i], p); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}

func lineStringDistanceFrom(ls geo.LineString, p geo.Point) (float64, int) {
	dist := math.Inf(1)
	index := -1

	for i := 0; i < len(ls)-1; i++ {
		if d := segmentDistanceFromSquared(ls[i], ls[i+1], p); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}

func polygonDistanceFrom(p geo.Polygon, point geo.Point) (float64, int) {
	if len(p) == 0 {
		return math.Inf(1), -1
	}

	dist, index := lineStringDistanceFrom(geo.LineString(p[0]), point)
	for i := 1; i < len(p); i++ {
		d, i := lineStringDistanceFrom(geo.LineString(p[i]), point)
		if d < dist {
			dist = d
			index = i
		}
	}

	return dist, index
}

func segmentDistanceFromSquared(p1, p2, point geo.Point) float64 {
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

	return dx*dx + dy*dy
}
