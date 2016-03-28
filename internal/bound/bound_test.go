package bound

import "testing"

func TestNew(t *testing.T) {
	bound := New(5, 0, 3, 0)
	if bound.SW != NewPoint(0, 0) {
		t.Errorf("incorrect sw: expected %v, got %v", NewPoint(0, 0), bound.SW)
	}

	if bound.NE != NewPoint(5, 3) {
		t.Errorf("incorrect ne: expected %v, got %v", NewPoint(5, 3), bound.NE)
	}

	bound = FromPoints(NewPoint(0, 3), NewPoint(4, 0))
	if bound.SW != NewPoint(0, 0) {
		t.Errorf("incorrect sw: expected %v, got %v", NewPoint(0, 0), bound.SW)
	}

	if bound.NE != NewPoint(4, 3) {
		t.Errorf("incorrect ne: expected %v, got %v", NewPoint(4, 3), bound.NE)
	}

	bound1 := New(1, 2, 3, 4)
	bound2 := FromPoints(NewPoint(1, 3), NewPoint(2, 4))
	if !bound1.Equal(bound2) {
		t.Errorf("expected %v == %v", bound1, bound2)
	}
}

func TestBoundExtend(t *testing.T) {
	bound := New(3, 0, 5, 0)

	if b := bound.Extend(NewPoint(2, 1)); !b.Equal(bound) {
		t.Errorf("extend expected %v, got %v", bound, b)
	}

	answer := New(6, 0, 5, -1)
	if b := bound.Extend(NewPoint(6, -1)); !b.Equal(answer) {
		t.Errorf("extend expected %v, got %v", answer, b)
	}
}

func TestBoundUnion(t *testing.T) {
	b1 := New(0, 1, 0, 1)
	b2 := New(0, 2, 0, 2)

	expected := New(0, 2, 0, 2)
	if b := b1.Union(b2); !b.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, b)
	}

	if b := b2.Union(b1); !b.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, b)
	}
}

func TestBoundContains(t *testing.T) {
	var p Point
	bound := New(2, -2, 1, -1)

	p = NewPoint(0, 0)
	if !bound.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(-1, 0)
	if !bound.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(2, 1)
	if !bound.Contains(p) {
		t.Errorf("contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(0, 3)
	if bound.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(0, -3)
	if bound.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(3, 0)
	if bound.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(-3, 0)
	if bound.Contains(p) {
		t.Errorf("contains expected %v, to not be within %v", p, bound)
	}
}

func TestBoundIntersects(t *testing.T) {
	var tester Bound
	bound := New(0, 1, 2, 3)

	tester = New(5, 6, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = New(-6, -5, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = New(0, 0.5, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = New(0, 0.5, 1, 4)
	if !bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, bound)
	}

	tester = New(-1, 2, 1, 4)
	if !bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, bound)
	}

	tester = New(0.3, 0.6, 2.3, 2.6)
	if !bound.Intersects(tester) {
		t.Errorf("intersects expected %v, to intersect %v", tester, bound)
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

func TestBoundCenter(t *testing.T) {
	var p Point
	var b Bound

	b = New(0, 1, 2, 3)
	p = NewPoint(0.5, 2.5)
	if c := b.Center(); c != p {
		t.Errorf("center expected %v, got %v", p, c)
	}

	b = New(0, 0, 2, 2)
	p = NewPoint(0, 2)
	if c := b.Center(); c != p {
		t.Errorf("center expected %v, got %v", p, c)
	}
}

func TestBoundEmpty(t *testing.T) {
	bound := New(1, 2, 3, 4)
	if bound.Empty() {
		t.Error("empty exported false, got true")
	}

	bound = New(1, 1, 2, 2)
	if !bound.Empty() {
		t.Error("empty exported true, got false")
	}

	// horizontal bar
	bound = New(1, 1, 2, 3)
	if !bound.Empty() {
		t.Error("empty exported true, got false")
	}

	// vertical bar
	bound = New(1, 2, 2, 2)
	if !bound.Empty() {
		t.Error("empty exported true, got false")
	}

	// negative/malformed area
	bound = New(1, 0, 2, 2)
	if !bound.Empty() {
		t.Error("empty exported true, got false")
	}

	// negative/malformed area
	bound = New(1, 1, 2, 1)
	if !bound.Empty() {
		t.Error("empty exported true, got false")
	}
}

func TestWKT(t *testing.T) {
	bound := New(1, 2, 3, 4)

	answer := "POLYGON((1 3, 1 4, 2 4, 2 3, 1 3))"
	if s := bound.WKT(); s != answer {
		t.Errorf("wkt expected %s, got %s", answer, s)
	}
}

func TestMysqlIntersectsCondition(t *testing.T) {
	b := New(1, 2, 3, 4)

	expected := "INTERSECTS(column, GEOMFROMTEXT('POLYGON((1 3, 1 4, 2 4, 2 3, 1 3))'))"
	if p := b.MysqlIntersectsCondition("column"); p != expected {
		t.Errorf("incorrect condition, got %v", p)
	}
}
