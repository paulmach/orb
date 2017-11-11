package wrap

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func TestNexts(t *testing.T) {
	for i, next := range nexts[orb.CW] {
		if next == -1 {
			continue
		}

		if i != nexts[orb.CCW][next] {
			t.Errorf("incorrect %d: %d != %d", i, i, nexts[orb.CCW][next])
		}
	}
}

func TestInternalAroundBound(t *testing.T) {
	cases := []struct {
		name        string
		box         orb.Bound
		input       orb.Ring
		output      orb.Ring
		orientation orb.Orientation
		expected    orb.Orientation
	}{
		{
			name:   "simple ccw",
			box:    orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:  orb.Ring{{-2, -2}, {2, 2}},
			output: orb.Ring{{-2, -2}, {2, 2}, {0, 1}, {-1, 1}, {-1, 0}, {-2, -2}}, expected: orb.CCW,
		},
		{
			name:     "simple cw",
			box:      orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:    orb.Ring{{-2, -2}, {2, 2}},
			output:   orb.Ring{{-2, -2}, {2, 2}, {1, 0}, {1, -1}, {0, -1}, {-2, -2}},
			expected: orb.CW,
		},
		{
			name:     "wrap around whole box ccw",
			box:      orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:    orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:   orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-2, 0.5}},
			expected: orb.CCW,
		},
		{
			name:     "wrap around whole box cw",
			box:      orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:    orb.Ring{{-2, -0.5}, {0, -0.5}, {0, 0.5}, {-2, 0.5}},
			output:   orb.Ring{{-2, -0.5}, {0, -0.5}, {0, 0.5}, {-2, 0.5}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-2, -0.5}},
			expected: orb.CW,
		},
		{
			name:        "already cw with endpoints in same section",
			box:         orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:       orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:      orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-2, 0.5}},
			orientation: orb.CW,
			expected:    orb.CW,
		},
		{
			name:        "cw but want ccw with endpoints in same section",
			box:         orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:       orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:      orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-2, 0.5}},
			orientation: orb.CW,
			expected:    orb.CCW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := aroundBound(tc.box, tc.input, tc.expected)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("does not match")
				t.Logf("%v", out)
				t.Logf("%v", tc.output)
			}
		})
	}
}
