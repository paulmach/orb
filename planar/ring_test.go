package planar

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestRingCentroid(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result Point
	}{
		{
			name:   "triangle, cw",
			points: []Point{{0, 0}, {1, 3}, {2, 0}, {0, 0}},
			result: Point{1, 1},
		},
		{
			name:   "triangle, ccw",
			points: []Point{{0, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: Point{1, 1},
		},
		{
			name:   "square, cw",
			points: []Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: Point{0.5, 0.5},
		},
		{
			name:   "triangle, ccw",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: Point{0.5, 0.5},
		},
		{
			name:   "redudent points",
			points: []Point{{0, 0}, {1, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: Point{1, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			if c := ring.Centroid(); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.points {
				tc.points[i][0] += 1e8
				tc.points[i][1] -= 1e8
			}

			tc.result[0] += 1e8
			tc.result[1] -= 1e8

			if c := ring.Centroid(); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestRingCentroidAdv(t *testing.T) {
	ring := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(0, 1),
		NewPoint(1, 1),
		NewPoint(1, 0.5),
		NewPoint(2, 0.5),
		NewPoint(2, 1),
		NewPoint(3, 1),
		NewPoint(3, 0),
		NewPoint(0, 0),
	)

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	expected := NewPoint(1.5, 0.45)
	if c := ring.Centroid(); !c.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", c, expected)
	}
}

func TestPolygonContains(t *testing.T) {
	ring := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(0, 1),
		NewPoint(1, 1),
		NewPoint(1, 0.5),
		NewPoint(2, 0.5),
		NewPoint(2, 1),
		NewPoint(3, 1),
		NewPoint(3, 0),
		NewPoint(0, 0),
	)

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	cases := []struct {
		name   string
		point  Point
		result bool
	}{
		{
			name:   "in base",
			point:  Point{1.5, 0.25},
			result: true,
		},
		{
			name:   "in right tower",
			point:  Point{0.5, 0.75},
			result: true,
		},
		{
			name:   "in middle",
			point:  Point{1.5, 0.75},
			result: false,
		},
		{
			name:   "in left tower",
			point:  Point{2.5, 0.75},
			result: true,
		},
		{
			name:   "in tp middle",
			point:  Point{1.5, 1.0},
			result: false,
		},
		{
			name:   "above",
			point:  Point{2.5, 1.75},
			result: false,
		},
		{
			name:   "below",
			point:  Point{2.5, -1.75},
			result: false,
		},
		{
			name:   "left",
			point:  Point{-2.5, -0.75},
			result: false,
		},
		{
			name:   "right",
			point:  Point{3.5, 0.75},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := ring.Contains(tc.point)

			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}
		})
	}

	// points should all be in
	for i, p := range ring {
		if !ring.Contains(p) {
			t.Errorf("point index %d: should be inside", i)
		}
	}

	// on all the segments should be in.
	for i := 1; i < len(ring); i++ {
		c := newSegment(ring[i], ring[i-1]).Centroid()
		if !ring.Contains(c) {
			t.Errorf("index %d centroid: should be inside", i)
		}
	}

	// colinear with segments but outside
	for i := 1; i < len(ring); i++ {
		p := newSegment(ring[i], ring[i-1]).Interpolate(5)
		if ring.Contains(p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}

		p = newSegment(ring[i], ring[i-1]).Interpolate(-5)
		if ring.Contains(p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}
	}
}

func TestRingSignedArea(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result float64
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			points: []Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: -1,
		},
		{
			name:   "even number of points",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0.4, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "4 points",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			result: 0.5,
		},
		{
			name:   "6 points",
			points: []Point{{1, 1}, {2, 1}, {2, 1.5}, {2, 2}, {1, 2}, {1, 1}},
			result: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			val := ring.SignedArea()
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.points {
				tc.points[i][0] += 1e15
				tc.points[i][1] -= 1e15
			}

			val = ring.SignedArea()
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that are rendant last point is implicit
			ring = ring[:len(ring)-1]
			val = ring.SignedArea()
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

func TestRingOrientation(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result orb.Orientation
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: orb.CCW,
		},
		{
			name:   "simple box, cw",
			points: []Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: orb.CW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			val := ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong winding: %v != %v", val, tc.result)
			}
		})
	}
}
