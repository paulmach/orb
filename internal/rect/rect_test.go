package rect

import "testing"

func TestNew(t *testing.T) {
	rect := New(5, 0, 3, 0)
	if rect.SW != NewPoint(0, 0) {
		t.Errorf("incorrect sw: expected %v, got %v", NewPoint(0, 0), rect.SW)
	}

	if rect.NE != NewPoint(5, 3) {
		t.Errorf("incorrect ne: expected %v, got %v", NewPoint(5, 3), rect.NE)
	}

	rect = FromPoints(NewPoint(0, 3), NewPoint(4, 0))
	if rect.SW != NewPoint(0, 0) {
		t.Errorf("incorrect sw: expected %v, got %v", NewPoint(0, 0), rect.SW)
	}

	if rect.NE != NewPoint(4, 3) {
		t.Errorf("incorrect ne: expected %v, got %v", NewPoint(4, 3), rect.NE)
	}

	rect1 := New(1, 2, 3, 4)
	rect2 := FromPoints(NewPoint(1, 3), NewPoint(2, 4))
	if !rect1.Equal(rect2) {
		t.Errorf("expected %v == %v", rect1, rect2)
	}
}

func TestRectExtend(t *testing.T) {
	rect := New(3, 0, 5, 0)

	if r := rect.Extend(NewPoint(2, 1)); !r.Equal(rect) {
		t.Errorf("extend expected %v, got %v", rect, r)
	}

	answer := New(6, 0, 5, -1)
	if r := rect.Extend(NewPoint(6, -1)); !r.Equal(answer) {
		t.Errorf("extend expected %v, got %v", answer, r)
	}
}

func TestRectUnion(t *testing.T) {
	r1 := New(0, 1, 0, 1)
	r2 := New(0, 2, 0, 2)

	expected := New(0, 2, 0, 2)
	if r := r1.Union(r2); !r.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, r)
	}

	if r := r2.Union(r1); !r.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestRectContains(t *testing.T) {
	var p Point
	rect := New(2, -2, 1, -1)

	p = NewPoint(0, 0)
	if !rect.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, rect)
	}

	p = NewPoint(-1, 0)
	if !rect.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, rect)
	}

	p = NewPoint(2, 1)
	if !rect.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, rect)
	}

	p = NewPoint(0, 3)
	if rect.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, rect)
	}

	p = NewPoint(0, -3)
	if rect.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, rect)
	}

	p = NewPoint(3, 0)
	if rect.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, rect)
	}

	p = NewPoint(-3, 0)
	if rect.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, rect)
	}
}

func TestRectIntersects(t *testing.T) {
	var tester Rect
	rect := New(0, 1, 2, 3)

	tester = New(5, 6, 7, 8)
	if rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, rect)
	}

	tester = New(-6, -5, 7, 8)
	if rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, rect)
	}

	tester = New(0, 0.5, 7, 8)
	if rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, rect)
	}

	tester = New(0, 0.5, 1, 4)
	if !rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, rect)
	}

	tester = New(-1, 2, 1, 4)
	if !rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, rect)
	}

	tester = New(0.3, 0.6, 2.3, 2.6)
	if !rect.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, rect)
	}

	a := New(7, 8, 6, 7)
	b := New(6.1, 8.1, 6.1, 8.1)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}

	a = New(1, 4, 2, 3)
	b = New(2, 3, 1, 4)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}
}

func TestRectCenter(t *testing.T) {
	var p Point
	var r Rect

	r = New(0, 1, 2, 3)
	p = NewPoint(0.5, 2.5)
	if c := r.Center(); c != p {
		t.Errorf("center expected %v, got %v", p, c)
	}

	r = New(0, 0, 2, 2)
	p = NewPoint(0, 2)
	if c := r.Center(); c != p {
		t.Errorf("center expected %v, got %v", p, c)
	}
}

func TestRectIsEmpty(t *testing.T) {
	rect := New(1, 2, 3, 4)
	if rect.IsEmpty() {
		t.Error("IsEmpty expected false, got true")
	}

	rect = New(1, 1, 2, 2)
	if rect.IsEmpty() {
		t.Error("IsEmpty expected false, got true")
	}

	// horizontal bar
	rect = New(1, 1, 2, 3)
	if rect.IsEmpty() {
		t.Error("IsEmpty expected false, got true")
	}

	// vertical bar
	rect = New(1, 2, 2, 2)
	if rect.IsEmpty() {
		t.Error("IsEmpty expected false, got true")
	}

	// negative/malformed area
	rect = New(1, 1, 2, 2)
	rect.NE[0] = 0
	if !rect.IsEmpty() {
		t.Error("IsEmpty expected true, got false")
	}

	// negative/malformed area
	rect = New(1, 1, 2, 2)
	rect.SW[1] = 3
	if !rect.IsEmpty() {
		t.Error("IsEmpty expected true, got false")
	}
}

func TestRectIsZero(t *testing.T) {
	rect := New(1, 1, 2, 2)
	if rect.IsZero() {
		t.Error("IsZero expected false, got true")
	}

	rect = New(0, 0, 0, 0)
	if !rect.IsZero() {
		t.Error("IsZero expected true, got false")
	}

	var r Rect
	if !r.IsZero() {
		t.Error("IsZero expected true, got false")
	}
}

func TestWKT(t *testing.T) {
	rect := New(1, 2, 3, 4)

	answer := "POLYGON((1 3, 1 4, 2 4, 2 3, 1 3))"
	if s := rect.WKT(); s != answer {
		t.Errorf("wkt expected %s, got %s", answer, s)
	}
}
