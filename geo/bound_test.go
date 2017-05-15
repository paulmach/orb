package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb/internal/mercator"
)

func center(r Bound) Point {
	return Point{
		(r[0][0] + r[1][0]) / 2.0,
		(r[0][1] + r[1][1]) / 2.0,
	}
}

func TestBoundAroundPoint(t *testing.T) {

	p := Point{
		5.42553,
		50.0359,
	}
	bound := NewBoundAroundPoint(p, 1000000)
	if center(bound).Lat() != p.Lat() {
		t.Errorf("should have correct center lat point")
	}

	if center(bound).Lon() != p.Lon() {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if !bound.Contains(Point{3.412, 58.3838}) {
		t.Errorf("should have point included in bound")
	}

	bound = NewBoundAroundPoint(p, 10000.0)
	if center(bound).Lat() != p.Lat() {
		t.Errorf("should have correct center lat point")
	}

	if center(bound).Lon() != p.Lon() {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if bound.Contains(Point{3.412, 58.3838}) {
		t.Errorf("should not have point included in bound")
	}
}

func TestNewBoundFromMapTile(t *testing.T) {
	bound, _ := NewBoundFromMapTile(7, 8, 9)

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the bound
	lon, lat := mercator.ScalarInverse(7<<factor+1, 8<<factor+1, level)
	if !bound.Contains(NewPoint(lon, lat)) {
		t.Errorf("should contain point")
	}

	lon, lat = mercator.ScalarInverse(7<<factor-1, 8<<factor-1, level)
	if bound.Contains(NewPoint(lon, lat)) {
		t.Errorf("should not contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor-1, 9<<factor-1, level)
	if !bound.Contains(NewPoint(lon, lat)) {
		t.Errorf("should contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor+1, 9<<factor+1, level)
	if bound.Contains(NewPoint(lon, lat)) {
		t.Errorf("should not contain point")
	}
}

func TestBoundPad(t *testing.T) {
	tests := []Bound{
		NewBoundFromPoints(NewPoint(-122.559, 37.887), NewPoint(-122.521, 37.911)),
		NewBoundFromPoints(NewPoint(-122.559, 15), NewPoint(-122.521, 15)),
		NewBoundFromPoints(NewPoint(20, -15), NewPoint(20, -15)),
	}

	for i, b1 := range tests {
		b2 := b1.Pad(100)

		if math.Abs(b1.Height()+200-b2.Height()) > 1.0 {
			t.Errorf("geoPad height incorrected for %d, expected %v, got %v", i, b1.Height()+200, b2.Height())
		}

		if math.Abs(b1.Width()+200-b2.Width()) > 1.0 {
			t.Errorf("geoPad width incorrected for %d, expected %v, got %v", i, b1.Width()+200, b2.Width())
		}
	}
}

func TestBoundExtend(t *testing.T) {
	bound := NewBound(3, 0, 5, 0)

	if r := bound.Extend(NewPoint(2, 1)); !r.Equal(bound) {
		t.Errorf("extend incorrect: %v != %v", r, bound)
	}

	answer := NewBound(6, 0, 5, -1)
	if r := bound.Extend(NewPoint(6, -1)); !r.Equal(answer) {
		t.Errorf("extend incorrect: %v != %v", r, answer)
	}
}

func TestBoundUnion(t *testing.T) {
	r1 := NewBound(0, 1, 0, 1)
	r2 := NewBound(0, 2, 0, 2)

	expected := NewBound(0, 2, 0, 2)
	if r := r1.Union(r2); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}

	if r := r2.Union(r1); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}
}

func TestBoundContains(t *testing.T) {
	bound := NewBound(2, -2, 1, -1)

	cases := []struct {
		name   string
		point  Point
		result bool
	}{
		{
			name:   "middle",
			point:  NewPoint(0, 0),
			result: true,
		},
		{
			name:   "left border",
			point:  NewPoint(-1, 0),
			result: true,
		},
		{
			name:   "ne corner",
			point:  NewPoint(2, 1),
			result: true,
		},
		{
			name:   "above",
			point:  NewPoint(0, 3),
			result: false,
		},
		{
			name:   "below",
			point:  NewPoint(0, -3),
			result: false,
		},
		{
			name:   "left",
			point:  NewPoint(-3, 0),
			result: false,
		},
		{
			name:   "right",
			point:  NewPoint(3, 0),
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := bound.Contains(tc.point)
			if v != tc.result {
				t.Errorf("incorrect contains: %v != %v", v, tc.result)
			}
		})
	}
}

func TestBoundIntersects(t *testing.T) {
	bound := NewBound(0, 1, 2, 3)

	cases := []struct {
		name   string
		bound  Bound
		result bool
	}{
		{
			name:   "outside, top right",
			bound:  NewBound(5, 6, 7, 8),
			result: false,
		},
		{
			name:   "outside, top left",
			bound:  NewBound(-6, -5, 7, 8),
			result: false,
		},
		{
			name:   "outside, above",
			bound:  NewBound(0, 0.5, 7, 8),
			result: false,
		},
		{
			name:   "over the middle",
			bound:  NewBound(0, 0.5, 1, 4),
			result: true,
		},
		{
			name:   "over the left",
			bound:  NewBound(-1, 2, 1, 4),
			result: true,
		},
		{
			name:   "completely inside",
			bound:  NewBound(0.3, 0.6, 2.3, 2.6),
			result: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := bound.Intersects(tc.bound)
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}

	a := NewBound(7, 8, 6, 7)
	b := NewBound(6.1, 8.1, 6.1, 8.1)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}

	a = NewBound(1, 4, 2, 3)
	b = NewBound(2, 3, 1, 4)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}
}

func TestBoundIsEmpty(t *testing.T) {
	cases := []struct {
		name   string
		bound  Bound
		result bool
	}{
		{
			name:   "regular bound",
			bound:  NewBound(1, 2, 3, 4),
			result: false,
		},
		{
			name:   "single point",
			bound:  NewBound(1, 1, 2, 2),
			result: false,
		},
		{
			name:   "horizontal bar",
			bound:  NewBound(1, 1, 2, 3),
			result: false,
		},
		{
			name:   "vertical bar",
			bound:  NewBound(1, 2, 2, 2),
			result: false,
		},
		{
			name:   "vertical bar",
			bound:  NewBound(1, 2, 2, 2),
			result: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.bound.IsEmpty()
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}

		})
	}

	// negative/malformed area
	bound := NewBound(1, 1, 2, 2)
	bound[1][0] = 0
	if !bound.IsEmpty() {
		t.Error("expected true, got false")
	}

	// negative/malformed area
	bound = NewBound(1, 1, 2, 2)
	bound[0][1] = 3
	if !bound.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestBoundIsZero(t *testing.T) {
	bound := NewBound(1, 1, 2, 2)
	if bound.IsZero() {
		t.Error("expected false, got true")
	}

	bound = NewBound(0, 0, 0, 0)
	if !bound.IsZero() {
		t.Error("expected true, got false")
	}

	var r Bound
	if !r.IsZero() {
		t.Error("expected true, got false")
	}
}

func TestWKT(t *testing.T) {
	bound := NewBound(1, 2, 3, 4)

	answer := "POLYGON((1 3,1 4,2 4,2 3,1 3))"
	if s := bound.WKT(); s != answer {
		t.Errorf("wkt expected %s, got %s", answer, s)
	}
}
