package wrap

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func TestAroundBound(t *testing.T) {
	for _, g := range orb.AllGeometries {
		AroundBound(orb.Bound{}, g, orb.CW)
	}
}

func TestRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  orb.Bound
		input  orb.Ring
		output orb.Ring
		orient orb.Orientation
	}{
		{
			name:   "wrap around whole box ccw",
			bound:  orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:  orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output: orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-2, 0.5}},
			orient: orb.CCW,
		},
		{
			name:   "just close the ring",
			bound:  orb.Bound{Min: orb.Point{-1, -1}, Max: orb.Point{1, 1}},
			input:  orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}},
			output: orb.Ring{{-2, 0.5}, {0, 0.5}, {0, -0.5}, {-2, -0.5}, {-2, 0.5}},
			orient: orb.CW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Ring(tc.bound, tc.input, tc.orient)
			if err != nil {
				t.Fatalf("should not get error: %v", err)
			}

			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("does not match")
				t.Logf("%v", out)
				t.Logf("%v", tc.output)
			}
		})
	}
}
