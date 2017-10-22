package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func TestRingSignedArea(t *testing.T) {
	area := 12392.029
	cases := []struct {
		name   string
		points []Point
		result float64
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: area,
		},
		{
			name:   "simple box, cc",
			points: []Point{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: -area,
		},
		{
			name:   "even number of points",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0.0004, 0.001}, {0, 0.001}, {0, 0}},
			result: area,
		},
		{
			name:   "3 points",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}},
			result: area / 2.0,
		},
		{
			name:   "4 points",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0}},
			result: area / 2.0,
		},
		{
			name:   "6 points",
			points: []Point{{0.001, 0.001}, {0.002, 0.001}, {0.002, 0.0015}, {0.002, 0.002}, {0.001, 0.002}, {0.001, 0.001}},
			result: area,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			val := ring.SignedArea()
			if math.Abs(val-tc.result) > 1 {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// should work without redudant last point.
			if ring[0] == ring[len(ring)-1] {
				ring = ring[:len(ring)-1]
				val = ring.SignedArea()
				if math.Abs(val-tc.result) > 1 {
					t.Errorf("wrong area: %v != %v", val, tc.result)
				}
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
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: orb.CCW,
		},
		{
			name:   "simple box, cw",
			points: []Point{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: orb.CW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			val := ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong orientation: %v != %v", val, tc.result)
			}

			// should work without redudant last point.
			ring = ring[:len(ring)-1]
			val = ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong orientation: %v != %v", val, tc.result)
			}
		})
	}
}
