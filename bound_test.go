package orb

import (
	"testing"
)

func center(r Bound) Point {
	return Point{
		(r[0][0] + r[1][0]) / 2.0,
		(r[0][1] + r[1][1]) / 2.0,
	}
}

func TestBoundExtend(t *testing.T) {
	bound := NewBound(0, 3, 0, 5)

	if r := bound.Extend(NewPoint(2, 1)); !r.Equal(bound) {
		t.Errorf("extend incorrect: %v != %v", r, bound)
	}

	answer := NewBound(0, 6, -1, 5)
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
	bound := NewBound(-2, 2, -1, 1)

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

func TestBoundToRing(t *testing.T) {
	bound := NewBound(1, 2, 1, 2)

	if bound.ToRing().Orientation() != CCW {
		t.Errorf("orientation should be ccw")
	}
}
