package geo

import (
	"math"
	"testing"
)

func TestRingArea(t *testing.T) {
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
			val := ring.Area()
			if math.Abs(val-tc.result) > 1 {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

func TestRingOrientation(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result int
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			points: []Point{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: -1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring := Ring(tc.points)
			val := ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong orientation: %v != %v", val, tc.result)
			}
		})
	}
}
