package planar

import "testing"

func TestNewRect(t *testing.T) {
	rect := NewRect(5, 0, 3, 0)
	if rect[0] != NewPoint(0, 0) {
		t.Errorf("incorrect sw: %v != %v", rect[0], NewPoint(0, 0))
	}

	if rect[1] != NewPoint(5, 3) {
		t.Errorf("incorrect ne: %v != %v", rect[1], NewPoint(5, 3))
	}

	rect = NewRectFromPoints(NewPoint(0, 3), NewPoint(4, 0))
	if rect[0] != NewPoint(0, 0) {
		t.Errorf("incorrect sw: %v != %v", rect[0], NewPoint(0, 0))
	}

	if rect[1] != NewPoint(4, 3) {
		t.Errorf("incorrect ne: %v != %v", rect[1], NewPoint(4, 3))
	}

	rect1 := NewRect(1, 2, 3, 4)
	rect2 := NewRectFromPoints(NewPoint(1, 3), NewPoint(2, 4))
	if !rect1.Equal(rect2) {
		t.Errorf("incorrect rect: %v != %v", rect1, rect2)
	}
}

func TestRectPad(t *testing.T) {
	var rect, tester Rect

	rect = NewRect(0, 1, 2, 3)
	tester = NewRect(-0.5, 1.5, 1.5, 3.5)
	if rect = rect.Pad(0.5); !rect.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, rect)
	}

	rect = NewRect(0, 1, 2, 3)
	tester = NewRect(0.1, 0.9, 2.1, 2.9)
	if rect = rect.Pad(-0.1); !rect.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, rect)
	}
}

func TestRectExtend(t *testing.T) {
	rect := NewRect(3, 0, 5, 0)

	if r := rect.Extend(NewPoint(2, 1)); !r.Equal(rect) {
		t.Errorf("extend incorrect: %v != %v", r, rect)
	}

	answer := NewRect(6, 0, 5, -1)
	if r := rect.Extend(NewPoint(6, -1)); !r.Equal(answer) {
		t.Errorf("extend incorrect: %v != %v", r, answer)
	}
}

func TestRectUnion(t *testing.T) {
	r1 := NewRect(0, 1, 0, 1)
	r2 := NewRect(0, 2, 0, 2)

	expected := NewRect(0, 2, 0, 2)
	if r := r1.Union(r2); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}

	if r := r2.Union(r1); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}
}

func TestRectContains(t *testing.T) {
	rect := NewRect(2, -2, 1, -1)

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
			v := rect.Contains(tc.point)
			if v != tc.result {
				t.Errorf("incorrect contains: %v != %v", v, tc.result)
			}
		})
	}
}

func TestRectIntersects(t *testing.T) {
	rect := NewRect(0, 1, 2, 3)

	cases := []struct {
		name   string
		rect   Rect
		result bool
	}{
		{
			name:   "outside, top right",
			rect:   NewRect(5, 6, 7, 8),
			result: false,
		},
		{
			name:   "outside, top left",
			rect:   NewRect(-6, -5, 7, 8),
			result: false,
		},
		{
			name:   "outside, above",
			rect:   NewRect(0, 0.5, 7, 8),
			result: false,
		},
		{
			name:   "over the middle",
			rect:   NewRect(0, 0.5, 1, 4),
			result: true,
		},
		{
			name:   "over the left",
			rect:   NewRect(-1, 2, 1, 4),
			result: true,
		},
		{
			name:   "completely inside",
			rect:   NewRect(0.3, 0.6, 2.3, 2.6),
			result: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := rect.Intersects(tc.rect)
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}

	a := NewRect(7, 8, 6, 7)
	b := NewRect(6.1, 8.1, 6.1, 8.1)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}

	a = NewRect(1, 4, 2, 3)
	b = NewRect(2, 3, 1, 4)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}
}

func TestRectCentroid(t *testing.T) {
	var p Point
	var r Rect

	r = NewRect(0, 1, 2, 3)
	p = NewPoint(0.5, 2.5)
	if c := r.Centroid(); c != p {
		t.Errorf("incorrect centroid: %v != %v", c, p)
	}

	r = NewRect(0, 0, 2, 2)
	p = NewPoint(0, 2)
	if c := r.Centroid(); c != p {
		t.Errorf("incorrect centroid: %v != %v", c, p)
	}
}

func TestRectIsEmpty(t *testing.T) {
	cases := []struct {
		name   string
		rect   Rect
		result bool
	}{
		{
			name:   "regular rect",
			rect:   NewRect(1, 2, 3, 4),
			result: false,
		},
		{
			name:   "single point",
			rect:   NewRect(1, 1, 2, 2),
			result: false,
		},
		{
			name:   "horizontal bar",
			rect:   NewRect(1, 1, 2, 3),
			result: false,
		},
		{
			name:   "vertical bar",
			rect:   NewRect(1, 2, 2, 2),
			result: false,
		},
		{
			name:   "vertical bar",
			rect:   NewRect(1, 2, 2, 2),
			result: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.rect.IsEmpty()
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}

		})
	}

	// negative/malformed area
	rect := NewRect(1, 1, 2, 2)
	rect[1][0] = 0
	if !rect.IsEmpty() {
		t.Error("expected true, got false")
	}

	// negative/malformed area
	rect = NewRect(1, 1, 2, 2)
	rect[0][1] = 3
	if !rect.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestRectIsZero(t *testing.T) {
	rect := NewRect(1, 1, 2, 2)
	if rect.IsZero() {
		t.Error("expected false, got true")
	}

	rect = NewRect(0, 0, 0, 0)
	if !rect.IsZero() {
		t.Error("expected true, got false")
	}

	var r Rect
	if !r.IsZero() {
		t.Error("expected true, got false")
	}
}

func TestWKT(t *testing.T) {
	rect := NewRect(1, 2, 3, 4)

	answer := "POLYGON((1 3,1 4,2 4,2 3,1 3))"
	if s := rect.WKT(); s != answer {
		t.Errorf("wkt expected %s, got %s", answer, s)
	}
}
