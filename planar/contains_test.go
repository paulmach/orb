package planar

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestRingContains(t *testing.T) {
	ring := orb.Ring{
		{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5},
		{2, 1}, {3, 1}, {3, 0}, {0, 0},
	}

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	cases := []struct {
		name   string
		point  orb.Point
		result bool
	}{
		{
			name:   "in base",
			point:  orb.Point{1.5, 0.25},
			result: true,
		},
		{
			name:   "in right tower",
			point:  orb.Point{0.5, 0.75},
			result: true,
		},
		{
			name:   "in middle",
			point:  orb.Point{1.5, 0.75},
			result: false,
		},
		{
			name:   "in left tower",
			point:  orb.Point{2.5, 0.75},
			result: true,
		},
		{
			name:   "in tp middle",
			point:  orb.Point{1.5, 1.0},
			result: false,
		},
		{
			name:   "above",
			point:  orb.Point{2.5, 1.75},
			result: false,
		},
		{
			name:   "below",
			point:  orb.Point{2.5, -1.75},
			result: false,
		},
		{
			name:   "left",
			point:  orb.Point{-2.5, -0.75},
			result: false,
		},
		{
			name:   "right",
			point:  orb.Point{3.5, 0.75},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring.Reverse()
			val := RingContains(ring, tc.point)

			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}

			// should not care about orientation
			ring.Reverse()
			val = RingContains(ring, tc.point)
			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}
		})
	}

	// points should all be in
	for i, p := range ring {
		if !RingContains(ring, p) {
			t.Errorf("point index %d: should be inside", i)
		}
	}

	// on all the segments should be in.
	for i := 1; i < len(ring); i++ {
		c := interpolate(ring[i], ring[i-1], 0.5)
		if !RingContains(ring, c) {
			t.Errorf("index %d centroid: should be inside", i)
		}
	}

	// collinear with segments but outside
	for i := 1; i < len(ring); i++ {
		p := interpolate(ring[i], ring[i-1], 5)
		if RingContains(ring, p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}

		p = interpolate(ring[i], ring[i-1], -5)
		if RingContains(ring, p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}
	}
}

func TestPolygonContains(t *testing.T) {
	// should exclude holes
	p := orb.Polygon{
		{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
	}

	if !PolygonContains(p, orb.Point{1.5, 1.5}) {
		t.Errorf("should contain point")
	}

	// ring oriented same as outer ring
	p = append(p, orb.Ring{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}})
	if PolygonContains(p, orb.Point{1.5, 1.5}) {
		t.Errorf("should not contain point in hole")
	}

	p[1].Reverse() // oriented correctly as opposite of outer
	if PolygonContains(p, orb.Point{1.5, 1.5}) {
		t.Errorf("should not contain point in hole")
	}
}

func TestMultiPolygonContains(t *testing.T) {
	// should exclude holes
	mp := orb.MultiPolygon{
		{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
	}

	if !MultiPolygonContains(mp, orb.Point{0.5, 0.5}) {
		t.Errorf("should contain point")
	}

	if MultiPolygonContains(mp, orb.Point{1.5, 1.5}) {
		t.Errorf("should not contain point")
	}

	mp = append(mp, orb.Polygon{{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}}})

	if !MultiPolygonContains(mp, orb.Point{2.5, 0.5}) {
		t.Errorf("should contain point")
	}

	if MultiPolygonContains(mp, orb.Point{1.5, 0.5}) {
		t.Errorf("should not contain point")
	}
}

func TestRingContainsRing(t *testing.T) {
	outer := orb.Ring{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	inner := orb.Ring{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}

	if !RingContainsRing(outer, inner) {
		t.Errorf("inner ring should be within outer ring")
	}

	larger := orb.Ring{{-1, -1}, {11, -1}, {11, 11}, {-1, 11}, {-1, -1}}
	if RingContainsRing(outer, larger) {
		t.Errorf("larger ring should not be within smaller ring")
	}

	partial := orb.Ring{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}
	if RingContainsRing(outer, partial) {
		t.Errorf("partial ring should not be within outer ring")
	}
}

func TestPolygonContainsPolygon(t *testing.T) {
	outer := orb.Polygon{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	}
	inner := orb.Polygon{
		{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}},
	}

	if !PolygonContainsPolygon(outer, inner) {
		t.Errorf("inner polygon should be within outer polygon")
	}

	outerWithHole := orb.Polygon{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
		{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}},
	}
	if PolygonContainsPolygon(outerWithHole, inner) {
		t.Errorf("inner polygon should not be within hole")
	}

	outerSmaller := orb.Polygon{
		{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}},
	}
	if PolygonContainsPolygon(outerSmaller, inner) {
		t.Errorf("inner polygon should not be within smaller outer")
	}
}

func TestPolygonContainsLineString(t *testing.T) {
	poly := orb.Polygon{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	}

	ls := orb.LineString{{2, 2}, {4, 4}, {6, 6}}
	if !PolygonContainsLineString(poly, ls) {
		t.Errorf("linestring should be within polygon")
	}

	lsOutside := orb.LineString{{2, 2}, {12, 12}}
	if PolygonContainsLineString(poly, lsOutside) {
		t.Errorf("linestring outside should not be within polygon")
	}

	lsOnBoundary := orb.LineString{{0, 0}, {5, 0}}
	if !PolygonContainsLineString(poly, lsOnBoundary) {
		t.Errorf("linestring on boundary should be within polygon")
	}
}

func TestMultiPolygonContainsPolygon(t *testing.T) {
	mp := orb.MultiPolygon{
		{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}},
		{{{10, 0}, {15, 0}, {15, 5}, {10, 5}, {10, 0}}},
	}

	poly := orb.Polygon{{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}}
	if !MultiPolygonContainsPolygon(mp, poly) {
		t.Errorf("polygon should be within multi-polygon")
	}

	polyOutside := orb.Polygon{{{7, 7}, {8, 7}, {8, 8}, {7, 8}, {7, 7}}}
	if MultiPolygonContainsPolygon(mp, polyOutside) {
		t.Errorf("polygon outside should not be within multi-polygon")
	}
}

func TestRingIntersectsRing(t *testing.T) {
	r1 := orb.Ring{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	r2 := orb.Ring{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}

	if !RingIntersectsRing(r1, r2) {
		t.Errorf("rings should intersect")
	}

	r3 := orb.Ring{{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 20}}
	if RingIntersectsRing(r1, r3) {
		t.Errorf("non-overlapping rings should not intersect")
	}
}

func TestPolygonIntersectsPolygon(t *testing.T) {
	p1 := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}
	p2 := orb.Polygon{{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}}

	if !PolygonIntersectsPolygon(p1, p2) {
		t.Errorf("polygons should intersect")
	}

	p3 := orb.Polygon{{{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 30}}}
	if PolygonIntersectsPolygon(p1, p3) {
		t.Errorf("non-overlapping polygons should not intersect")
	}
}

func TestLineStringIntersectsPolygon(t *testing.T) {
	p := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}

	ls := orb.LineString{{-5, 5}, {5, 5}}
	if !LineStringIntersectsPolygon(ls, p) {
		t.Errorf("linestring crossing polygon should intersect")
	}

	lsInside := orb.LineString{{2, 2}, {4, 4}, {6, 6}}
	if !LineStringIntersectsPolygon(lsInside, p) {
		t.Errorf("linestring inside polygon should intersect")
	}

	lsOutside := orb.LineString{{20, 20}, {30, 30}}
	if LineStringIntersectsPolygon(lsOutside, p) {
		t.Errorf("linestring outside should not intersect")
	}
}

func TestRingCovers(t *testing.T) {
	r := orb.Ring{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}

	if !RingCovers(r, orb.Point{5, 5}) {
		t.Errorf("point inside should be covered")
	}

	if !RingCovers(r, orb.Point{0, 0}) {
		t.Errorf("point on boundary should be covered")
	}

	if RingCovers(r, orb.Point{15, 15}) {
		t.Errorf("point outside should not be covered")
	}
}

func TestPolygonCovers(t *testing.T) {
	p := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}

	if !PolygonCovers(p, orb.Point{5, 5}) {
		t.Errorf("point inside should be covered")
	}

	if PolygonCovers(p, orb.Point{15, 15}) {
		t.Errorf("point outside should not be covered")
	}
}

func TestPolygonCoversLineString(t *testing.T) {
	p := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}

	ls := orb.LineString{{2, 2}, {4, 4}}
	if !PolygonCoversLineString(p, ls) {
		t.Errorf("linestring inside should be covered")
	}

	lsOutside := orb.LineString{{2, 2}, {12, 12}}
	if PolygonCoversLineString(p, lsOutside) {
		t.Errorf("linestring outside should not be covered")
	}
}

func TestPolygonCoversPolygon(t *testing.T) {
	outer := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}
	inner := orb.Polygon{{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}}

	if !PolygonCoversPolygon(outer, inner) {
		t.Errorf("inner polygon should be covered by outer")
	}

	larger := orb.Polygon{{{0, 0}, {15, 0}, {15, 15}, {0, 15}, {0, 0}}}
	if PolygonCoversPolygon(outer, larger) {
		t.Errorf("larger polygon should not be covered by smaller")
	}
}

func TestRingWithin(t *testing.T) {
	outer := orb.Ring{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	inner := orb.Ring{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 8}}

	if !RingWithin(inner, outer) {
		t.Errorf("inner should be within outer")
	}

	if RingWithin(outer, inner) {
		t.Errorf("outer should not be within inner")
	}
}

func TestPolygonWithin(t *testing.T) {
	outer := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}
	inner := orb.Polygon{{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}}

	if !PolygonWithin(inner, outer) {
		t.Errorf("inner should be within outer")
	}

	if PolygonWithin(outer, inner) {
		t.Errorf("outer should not be within inner")
	}
}

func TestLineStringWithinPolygon(t *testing.T) {
	p := orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}

	ls := orb.LineString{{2, 2}, {4, 4}, {6, 6}}
	if !LineStringWithinPolygon(ls, p) {
		t.Errorf("linestring should be within polygon")
	}

	lsOutside := orb.LineString{{-1, 5}, {5, 5}}
	if LineStringWithinPolygon(lsOutside, p) {
		t.Errorf("linestring outside should not be within polygon")
	}
}

func interpolate(a, b orb.Point, percent float64) orb.Point {
	return orb.Point{
		a[0] + percent*(b[0]-a[0]),
		a[1] + percent*(b[1]-a[1]),
	}
}
