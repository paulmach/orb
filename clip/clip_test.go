package clip

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func TestInternalLine(t *testing.T) {
	cases := []struct {
		name   string
		bound  orb.Bound
		input  orb.LineString
		output orb.MultiLineString
	}{
		{
			name:  "clip line",
			bound: orb.NewBound(0, 30, 0, 30),
			input: orb.LineString{
				{-10, 10}, {10, 10}, {10, -10}, {20, -10}, {20, 10}, {40, 10},
				{40, 20}, {20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20}, {-10, 20},
			},
			output: orb.MultiLineString{
				{{0, 10}, {10, 10}, {10, 0}},
				{{20, 0}, {20, 10}, {30, 10}},
				{{30, 20}, {20, 20}, {20, 30}},
				{{10, 30}, {10, 20}, {5, 20}, {0, 20}},
			},
		},
		{
			name:  "clips line crossign many times",
			bound: orb.NewBound(0, 20, 0, 20),
			input: orb.LineString{
				{10, -10}, {10, 30}, {20, 30}, {20, -10},
			},
			output: orb.MultiLineString{
				{{10, 0}, {10, 20}},
				{{20, 20}, {20, 0}},
			},
		},
		{
			name:  "no changes if all inside",
			bound: orb.NewBound(0, 20, 0, 20),
			input: orb.LineString{
				{1, 1}, {2, 2}, {3, 3},
			},
			output: orb.MultiLineString{
				{{1, 1}, {2, 2}, {3, 3}},
			},
		},
		{
			name:  "empty if nothing in bound",
			bound: orb.NewBound(0, 2, 0, 2),
			input: orb.LineString{
				{10, 10}, {20, 20}, {30, 30},
			},
			output: nil,
		},
		{
			name:  "floating point example",
			bound: orb.NewBound(-91.93359375, -91.7578125, 42.29356419217009, 42.42345651793831),
			input: orb.LineString{
				{-86.66015624999999, 42.22851735620852}, {-81.474609375, 38.51378825951165},
				{-85.517578125, 37.125286284966776}, {-85.8251953125, 38.95940879245423},
				{-90.087890625, 39.53793974517628}, {-91.93359375, 42.32606244456202},
				{-86.66015624999999, 42.22851735620852},
			},
			output: orb.MultiLineString{
				{
					{-91.91208030440808, 42.29356419217009},
					{-91.93359375, 42.32606244456202},
					{-91.7578125, 42.3228109416169},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := line(tc.bound, tc.input)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestInternalRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  orb.Bound
		input  orb.Ring
		output orb.Ring
	}{
		{
			name:  "clips polygon",
			bound: orb.NewBound(0, 30, 0, 30),
			input: orb.Ring{
				{-10, 10}, {0, 10}, {10, 10}, {10, 5}, {10, -5},
				{10, -10}, {20, -10}, {20, 10}, {40, 10}, {40, 20},
				{20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20},
				{-10, 20}},
			// note: we allow duplicate points if polygon endpoints are
			// on the box boundary.
			output: orb.Ring{
				{0, 10}, {0, 10}, {10, 10}, {10, 5}, {10, 0},
				{20, 0}, {20, 10}, {30, 10}, {30, 20}, {20, 20},
				{20, 30}, {10, 30}, {10, 20}, {5, 20}, {0, 20},
			},
		},
		{
			name:  "completely inside bound",
			bound: orb.NewBound(0, 10, 0, 10),
			input: orb.Ring{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
			output: orb.Ring{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
		},
		{
			name:  "completely around bound",
			bound: orb.NewBound(1, 2, 1, 2),
			input: orb.Ring{
				{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0},
			},
			output: orb.Ring{{1, 2}, {1, 1}, {2, 1}, {2, 2}, {1, 2}},
		},
		{
			name:  "completely around touching corners",
			bound: orb.NewBound(1, 3, 1, 3),
			input: orb.Ring{
				{0, 2}, {2, 0}, {4, 2}, {2, 4}, {0, 2},
			},
			output: orb.Ring{{1, 1}, {1, 1}, {3, 1}, {3, 1}, {3, 3}, {3, 3}, {1, 3}, {1, 3}, {1, 1}},
		},
		{
			name:  "around but cut corners",
			bound: orb.NewBound(0.5, 3.5, 0.5, 3.5),
			input: orb.Ring{
				{0, 2}, {2, 4}, {4, 2}, {2, 0}, {0, 2},
			},
			output: orb.Ring{{0.5, 2.5}, {1.5, 3.5}, {2.5, 3.5}, {3.5, 2.5}, {3.5, 1.5}, {2.5, 0.5}, {1.5, 0.5}, {0.5, 1.5}, {0.5, 2.5}},
		},
		{
			name:  "unclosed ring",
			bound: orb.NewBound(1, 4, 1, 4),
			input: orb.Ring{
				{2, 0}, {3, 0}, {3, 5}, {2, 5},
			},
			output: orb.Ring{{3, 1}, {3, 4}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Ring(tc.bound, tc.input)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestInternalRing_CompletelyOutside(t *testing.T) {
	cases := []struct {
		name   string
		bound  orb.Bound
		input  orb.Ring
		output orb.Ring
	}{
		{
			name:  "bound in lower left",
			bound: orb.NewBound(-1, 0, -1, 0),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound in lower right",
			bound: orb.NewBound(3, 4, -1, 0),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound in upper right",
			bound: orb.NewBound(3, 4, 3, 4),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound in upper left",
			bound: orb.NewBound(-1, 0, 3, 4),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound to the left",
			bound: orb.NewBound(-1, 0, -1, 3),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound to the right",
			bound: orb.NewBound(3, 4, -1, 3),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound to the top",
			bound: orb.NewBound(-1, 3, 3, 4),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
		{
			name:  "bound to the bottom",
			bound: orb.NewBound(-1, 3, -1, 0),
			input: orb.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ring(tc.bound, tc.input)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v %+v", result == nil, result)
				t.Logf("%v %+v", tc.output == nil, tc.output)
			}
		})
	}
}
