package planar

import (
	"testing"

	"github.com/paulmach/orb/geo"
)

func TestCentroidArea_MultiPoint(t *testing.T) {
	mp := geo.MultiPoint{{0, 0}, {1, 1.5}, {2, 0}}

	centroid, area := CentroidArea(mp)
	expected := geo.NewPoint(1, 0.5)
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 0 {
		t.Errorf("area should be 0: %f", area)
	}
}

func TestCentroid_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   geo.Ring
		result geo.Point
	}{
		{
			name:   "triangle, cw",
			ring:   geo.Ring{{0, 0}, {1, 3}, {2, 0}, {0, 0}},
			result: geo.Point{1, 1},
		},
		{
			name:   "triangle, ccw",
			ring:   geo.Ring{{0, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: geo.Point{1, 1},
		},
		{
			name:   "square, cw",
			ring:   geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "triangle, ccw",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "redudent points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: geo.Point{1, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e8
				tc.ring[i][1] -= 1e8
			}

			tc.result[0] += 1e8
			tc.result[1] -= 1e8

			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestArea_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   geo.Ring
		result float64
	}{
		{
			name:   "simple box, ccw",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			ring:   geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: -1,
		},
		{
			name:   "even number of points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0.4, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "4 points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			result: 0.5,
		},
		{
			name:   "6 points",
			ring:   geo.Ring{{1, 1}, {2, 1}, {2, 1.5}, {2, 2}, {1, 2}, {1, 1}},
			result: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, val := CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e15
				tc.ring[i][1] -= 1e15
			}

			_, val = CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that are rendant last point is implicit
			tc.ring = tc.ring[:len(tc.ring)-1]
			_, val = CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

func TestCentroid_RingAdv(t *testing.T) {
	ring := geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5}, {2, 1}, {3, 1}, {3, 0}, {0, 0}}

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	expected := geo.NewPoint(1.5, 0.45)
	if c, _ := CentroidArea(ring); !c.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", c, expected)
	}
}

func TestCentroidArea_Polygon(t *testing.T) {
	r1 := geo.Ring{{0, 0}, {4, 0}, {4, 3}, {0, 3}, {0, 0}}
	r1.Reverse()

	r2 := geo.Ring{{2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 1}}
	poly := geo.Polygon{r1, r2}

	centroid, area := CentroidArea(poly)
	if !centroid.Equal(geo.NewPoint(21.5/11.0, 1.5)) {
		t.Errorf("%v", 21.5/11.0)
		t.Errorf("incorrect centroid: %v", centroid)
	}

	if area != 11 {
		t.Errorf("incorrect area: %v != 11", area)
	}
}

func TestCentroidArea_Bound(t *testing.T) {
	r := geo.NewBound(0, 1, 2, 3)
	centroid, area := CentroidArea(r)

	expected := geo.NewPoint(0.5, 2.5)
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 1 {
		t.Errorf("incorrect area: %f != 1", area)
	}

	r = geo.NewBound(0, 0, 2, 2)
	centroid, area = CentroidArea(r)

	expected = geo.NewPoint(0, 2)
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 0 {
		t.Errorf("area should be zero: %f", area)
	}
}
