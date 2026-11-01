package planar

import (
	"math"

	"github.com/paulmach/orb"
)

// RingContains returns true if the point is inside the ring.
// Points on the boundary are considered in.
func RingContains(r orb.Ring, point orb.Point) bool {
	if !r.Bound().Contains(point) {
		return false
	}

	c, on := rayIntersect(point, r[0], r[len(r)-1])
	if on {
		return true
	}

	for i := 0; i < len(r)-1; i++ {
		inter, on := rayIntersect(point, r[i], r[i+1])
		if on {
			return true
		}

		if inter {
			c = !c
		}
	}

	return c
}

func boundContainsBound(outer, inner orb.Bound) bool {
	return outer.Contains(inner.Min) && outer.Contains(inner.Max)
}

func segmentIntersects(a, b, c, d orb.Point) bool {
	denom := (d[1]-c[1])*(b[0]-a[0]) - (d[0]-c[0])*(b[1]-a[1])
	if denom == 0 {
		return false
	}

	ua := ((d[0]-c[0])*(a[1]-c[1]) - (d[1]-c[1])*(a[0]-c[0])) / denom
	ub := ((b[0]-a[0])*(a[1]-c[1]) - (b[1]-a[1])*(a[0]-c[0])) / denom

	return ua >= 0 && ua <= 1 && ub >= 0 && ub <= 1
}

// PolygonContains checks if the point is within the polygon.
// Points on the boundary are considered in.
func PolygonContains(p orb.Polygon, point orb.Point) bool {
	if !RingContains(p[0], point) {
		return false
	}

	for i := 1; i < len(p); i++ {
		if RingContains(p[i], point) {
			return false
		}
	}

	return true
}

// MultiPolygonContains checks if the point is within the multi-polygon.
// Points on the boundary are considered in.
func MultiPolygonContains(mp orb.MultiPolygon, point orb.Point) bool {
	for _, p := range mp {
		if PolygonContains(p, point) {
			return true
		}
	}

	return false
}

// RingContainsRing checks if the inner ring is completely within the outer ring.
// Points on the boundary are considered within.
func RingContainsRing(outer, inner orb.Ring) bool {
	if !boundContainsBound(outer.Bound(), inner.Bound()) {
		return false
	}

	for _, p := range inner {
		if !RingContains(outer, p) {
			return false
		}
	}

	return true
}

// PolygonContainsPolygon checks if the inner polygon is completely within the outer polygon.
// This includes checking the outer rings and ensuring the inner polygon
// is not in any holes of the outer polygon.
// Points on the boundary are considered within.
func PolygonContainsPolygon(outer, inner orb.Polygon) bool {
	if !boundContainsBound(outer.Bound(), inner.Bound()) {
		return false
	}

	if !RingContainsRing(outer[0], inner[0]) {
		return false
	}

	if len(outer) > 1 {
		center := inner[0].Bound().Center()
		for i := 1; i < len(outer); i++ {
			if RingContains(outer[i], center) {
				return false
			}
		}
	}

	return true
}

// PolygonContainsLineString checks if the linestring is completely within the polygon.
// Points on the boundary are considered within.
func PolygonContainsLineString(p orb.Polygon, ls orb.LineString) bool {
	if !boundContainsBound(p.Bound(), ls.Bound()) {
		return false
	}

	for _, pt := range ls {
		if !PolygonContains(p, pt) {
			return false
		}
	}

	return true
}

// MultiPolygonContainsPolygon checks if the polygon is completely within the multi-polygon.
func MultiPolygonContainsPolygon(mp orb.MultiPolygon, poly orb.Polygon) bool {
	for _, p := range mp {
		if PolygonContainsPolygon(p, poly) {
			return true
		}
	}

	return false
}

// RingIntersectsRing checks if two rings intersect.
// Returns true if they touch or overlap.
func RingIntersectsRing(r1, r2 orb.Ring) bool {
	if !r1.Bound().Intersects(r2.Bound()) {
		return false
	}

	for i := 0; i < len(r1)-1; i++ {
		for j := 0; j < len(r2)-1; j++ {
			if segmentIntersects(r1[i], r1[i+1], r2[j], r2[j+1]) {
				return true
			}
		}
	}

	if RingContains(r1, r2[0]) || RingContains(r2, r1[0]) {
		return true
	}

	return false
}

// PolygonIntersectsPolygon checks if two polygons intersect.
// Returns true if they touch or overlap.
func PolygonIntersectsPolygon(p1, p2 orb.Polygon) bool {
	if !p1.Bound().Intersects(p2.Bound()) {
		return false
	}

	for i := 0; i < len(p1); i++ {
		for j := 0; j < len(p2); j++ {
			if RingIntersectsRing(p1[i], p2[j]) {
				return true
			}
		}
	}

	return false
}

// LineStringIntersectsPolygon checks if a linestring intersects a polygon.
// Returns true if any part of the linestring touches or is within the polygon.
func LineStringIntersectsPolygon(ls orb.LineString, p orb.Polygon) bool {
	if !ls.Bound().Intersects(p.Bound()) {
		return false
	}

	for i := 0; i < len(ls)-1; i++ {
		for _, ring := range p {
			for j := 0; j < len(ring)-1; j++ {
				if segmentIntersects(ls[i], ls[i+1], ring[j], ring[j+1]) {
					return true
				}
			}
		}
	}

	for _, pt := range ls {
		if PolygonContains(p, pt) {
			return true
		}
	}

	return false
}

// RingCovers checks if the ring covers the point.
// Points on the boundary are considered covered.
func RingCovers(r orb.Ring, point orb.Point) bool {
	return RingContains(r, point)
}

// PolygonCovers checks if the polygon covers the point.
// Points on the boundary are considered covered.
func PolygonCovers(p orb.Polygon, point orb.Point) bool {
	return PolygonContains(p, point)
}

// PolygonCoversLineString checks if the polygon covers the linestring.
// Points on the boundary are considered covered.
func PolygonCoversLineString(p orb.Polygon, ls orb.LineString) bool {
	return PolygonContainsLineString(p, ls)
}

// PolygonCoversPolygon checks if the outer polygon covers the inner polygon.
// Points on the boundary are considered covered.
func PolygonCoversPolygon(outer, inner orb.Polygon) bool {
	return PolygonContainsPolygon(outer, inner)
}

// RingWithin returns true if the ring is completely within the other ring.
// Points on the boundary are considered within.
func RingWithin(inner, outer orb.Ring) bool {
	return RingContainsRing(outer, inner)
}

// PolygonWithin returns true if the polygon is completely within the other polygon.
// Points on the boundary are considered within.
func PolygonWithin(inner, outer orb.Polygon) bool {
	return PolygonContainsPolygon(outer, inner)
}

// LineStringWithinPolygon returns true if the linestring is completely within the polygon.
// Points on the boundary are considered within.
func LineStringWithinPolygon(ls orb.LineString, p orb.Polygon) bool {
	return PolygonContainsLineString(p, ls)
}

// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func rayIntersect(p, s, e orb.Point) (intersects, on bool) {
	if s[0] > e[0] {
		s, e = e, s
	}

	switch p[0] {
	case s[0]:
		if p[1] == s[1] {
			// p == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= p[1] && p[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= p[1] && p[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		p[0] = math.Nextafter(p[0], math.Inf(1))
	case e[0]:
		if p[1] == e[1] {
			// matching the end point
			return false, true
		}

		p[0] = math.Nextafter(p[0], math.Inf(1))
	}

	if p[0] < s[0] || p[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if p[1] > s[1] {
			return false, false
		} else if p[1] < e[1] {
			return true, false
		}
	} else {
		if p[1] > e[1] {
			return false, false
		} else if p[1] < s[1] {
			return true, false
		}
	}

	rs := (p[1] - s[1]) / (p[0] - s[0])
	ds := (e[1] - s[1]) / (e[0] - s[0])

	if rs == ds {
		return false, true
	}

	return rs <= ds, false
}
