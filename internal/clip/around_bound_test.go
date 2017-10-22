package clip

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

func TestAroundBound(t *testing.T) {
	cases := []struct {
		name        string
		box         Bound
		input       *lineString
		output      *lineString
		orientation orb.Orientation
		expected    orb.Orientation
	}{
		{
			name:     "simple ccw",
			box:      Bound{-1, 1, -1, 1},
			input:    &lineString{{-2, -2}, {2, 2}},
			output:   &lineString{{-2, -2}, {2, 2}, {0, 1}, {-1, 1}, {-1, 0}, {-2, -2}},
			expected: orb.CCW,
		},
		{
			name:     "simple cw",
			box:      Bound{-1, 1, -1, 1},
			input:    &lineString{{-2, -2}, {2, 2}},
			output:   &lineString{{-2, -2}, {2, 2}, {1, 0}, {1, -1}, {0, -1}, {-2, -2}},
			expected: orb.CW,
		},
		{
			name:     "wrap around whole box ccw",
			box:      Bound{-1, 1, -1, 1},
			input:    &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:   &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-2, 0.5}},
			expected: orb.CCW,
		},
		{
			name:     "wrap around whole box cw",
			box:      Bound{-1, 1, -1, 1},
			input:    &lineString{{-2, -0.5}, {0, -0.5}, {0, 0.5}, {-2, 0.5}},
			output:   &lineString{{-2, -0.5}, {0, -0.5}, {0, 0.5}, {-2, 0.5}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-2, -0.5}},
			expected: orb.CW,
		},
		{
			name:        "already cw with endpoints in same section",
			box:         Bound{-1, 1, -1, 1},
			input:       &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:      &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-2, 0.5}},
			orientation: orb.CW,
			expected:    orb.CW,
		},
		{
			name:        "cw but want ccw with endpoints in same section",
			box:         Bound{-1, 1, -1, 1},
			input:       &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output:      &lineString{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-2, 0.5}},
			orientation: orb.CW,
			expected:    orb.CCW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := AroundBound(
				tc.box,
				tc.input,
				tc.expected,
				func(in LineString) orb.Orientation { return tc.orientation },
			)
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
