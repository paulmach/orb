package clip

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  orb.Bound
		input  orb.Ring
		output orb.Ring
	}{
		{
			name:  "bound to the top",
			bound: orb.NewBound(-1, 3, 3, 4),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: orb.Ring{},
		},
		{
			name:  "bound in lower left",
			bound: orb.NewBound(-1, 0, -1, 0),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: orb.Ring{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Ring(tc.bound, tc.input)

			if !result.Equal(tc.output) {
				t.Errorf("not equal")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestBound(t *testing.T) {
	cases := []struct {
		name string
		b1   orb.Bound
		b2   orb.Bound
		rs   orb.Bound
	}{
		{
			name: "normal intersection",
			b1:   orb.NewBound(0, 3, 1, 4),
			b2:   orb.NewBound(1, 4, 2, 5),
			rs:   orb.NewBound(1, 3, 2, 4),
		},
		{
			name: "1 contains 2",
			b1:   orb.NewBound(0, 3, 1, 4),
			b2:   orb.NewBound(1, 2, 2, 3),
			rs:   orb.NewBound(1, 2, 2, 3),
		},
		{
			name: "no overlap",
			b1:   orb.NewBound(0, 3, 1, 4),
			b2:   orb.NewBound(4, 5, 5, 6),
			rs:   orb.NewBound(1, 0, 1, 0), // empty
		},
		{
			name: "same bound",
			b1:   orb.NewBound(0, 3, 1, 4),
			b2:   orb.NewBound(0, 3, 1, 4),
			rs:   orb.NewBound(0, 3, 1, 4),
		},
		{
			name: "1 is empty",
			b1:   orb.NewBound(1, 0, 1, 0),
			b2:   orb.NewBound(0, 3, 1, 4),
			rs:   orb.NewBound(0, 3, 1, 4),
		},
		{
			name: "both are empty",
			b1:   orb.NewBound(1, 0, 1, 0),
			b2:   orb.NewBound(1, 0, 1, 0),
			rs:   orb.NewBound(1, 0, 1, 0),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r1 := Bound(tc.b1, tc.b2)
			r2 := Bound(tc.b1, tc.b2)

			if tc.rs.IsEmpty() && (!r1.IsEmpty() || !r2.IsEmpty()) {
				t.Errorf("should be empty")
				t.Logf("%v", r1)
				t.Logf("%v", r2)
			}

			if !tc.rs.IsEmpty() {
				if !r1.Equal(tc.rs) {
					t.Errorf("r1 not equal")
					t.Logf("%v", r1)
					t.Logf("%v", tc.rs)
				}
				if !r2.Equal(tc.rs) {
					t.Errorf("r2 not equal")
					t.Logf("%v", r2)
					t.Logf("%v", tc.rs)
				}
			}
		})
	}
}
