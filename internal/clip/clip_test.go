package clip

import (
	"reflect"
	"testing"
)

func TestLine(t *testing.T) {
	cases := []struct {
		name   string
		bound  Bound
		input  *lineString
		output *multiLineString
	}{
		{
			name:  "clip line",
			bound: Bound{0, 30, 0, 30},
			input: &lineString{
				{-10, 10}, {10, 10}, {10, -10}, {20, -10}, {20, 10}, {40, 10},
				{40, 20}, {20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20}, {-10, 20},
			},
			output: &multiLineString{
				&lineString{{0, 10}, {10, 10}, {10, 0}},
				&lineString{{20, 0}, {20, 10}, {30, 10}},
				&lineString{{30, 20}, {20, 20}, {20, 30}},
				&lineString{{10, 30}, {10, 20}, {5, 20}, {0, 20}},
			},
		},
		{
			name:  "clips line crossign many times",
			bound: Bound{0, 20, 0, 20},
			input: &lineString{
				{10, -10}, {10, 30}, {20, 30}, {20, -10},
			},
			output: &multiLineString{
				&lineString{{10, 0}, {10, 20}},
				&lineString{{20, 20}, {20, 0}},
			},
		},
		{
			name:  "no changes if all inside",
			bound: Bound{0, 20, 0, 20},
			input: &lineString{
				{1, 1}, {2, 2}, {3, 3},
			},
			output: &multiLineString{
				&lineString{{1, 1}, {2, 2}, {3, 3}},
			},
		},
		{
			name:  "empty if nothing in bound",
			bound: Bound{0, 2, 0, 2},
			input: &lineString{
				{10, 10}, {20, 20}, {30, 30},
			},
			output: &multiLineString{},
		},
		{
			name:  "floating point example",
			bound: Bound{-91.93359375, -91.7578125, 42.29356419217009, 42.42345651793831},
			input: &lineString{
				{-86.66015624999999, 42.22851735620852}, {-81.474609375, 38.51378825951165},
				{-85.517578125, 37.125286284966776}, {-85.8251953125, 38.95940879245423},
				{-90.087890625, 39.53793974517628}, {-91.93359375, 42.32606244456202},
				{-86.66015624999999, 42.22851735620852},
			},
			output: &multiLineString{
				&lineString{
					{-91.91208030440808, 42.29356419217009},
					{-91.93359375, 42.32606244456202},
					{-91.7578125, 42.3228109416169},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := &multiLineString{}
			Line(tc.bound, tc.input, result)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  Bound
		input  *lineString
		output *lineString
	}{
		{
			name:  "clips polygon",
			bound: Bound{0, 30, 0, 30},
			input: &lineString{
				{-10, 10}, {0, 10}, {10, 10}, {10, 5}, {10, -5},
				{10, -10}, {20, -10}, {20, 10}, {40, 10}, {40, 20},
				{20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20},
				{-10, 20}},
			output: &lineString{
				{0, 10}, {0, 10}, {10, 10}, {10, 5}, {10, 0},
				{20, 0}, {20, 10}, {30, 10}, {30, 20}, {20, 20},
				{20, 30}, {10, 30}, {10, 20}, {5, 20}, {0, 20},
			},
		},
		{
			name:  "completely inside bound",
			bound: Bound{0, 10, 0, 10},
			input: &lineString{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
			output: &lineString{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
		},
		{
			name:  "completely outside bound",
			bound: Bound{0, 2, 0, 2},
			input: &lineString{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
			output: &lineString{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := &lineString{}
			Ring(tc.bound, tc.input, result)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

type multiLineString []*lineString

func (mls *multiLineString) Append(i int, x, y float64) {
	if i >= len(*mls) {
		ls := &lineString{}
		ls.Append(x, y)
		*mls = append(*mls, ls)
	} else {
		(*mls)[i].Append(x, y)
	}
}
