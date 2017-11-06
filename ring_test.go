package orb

import (
	"testing"
)

func TestRingOrientation(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result Orientation
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: CCW,
		},
		{
			name:   "simple box, cw",
			points: []Point{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: CW,
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
