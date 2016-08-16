package planar

import "testing"

func TestLineNew(t *testing.T) {
	a := NewPoint(1, 1)
	b := NewPoint(2, 2)

	l := NewLine(a, b)
	if !l.A().Equal(a) {
		t.Errorf("line, expected %v == %v", l.A(), a)
	}

	if !l.B().Equal(b) {
		t.Errorf("line, expected %v == %v", l.B(), b)
	}
}

func TestLineDistanceFrom(t *testing.T) {
	var answer float64
	l := NewLine(NewPoint(0, 0), NewPoint(0, 10))

	answer = 1
	if d := l.DistanceFrom(NewPoint(1, 5)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 0
	if d := l.DistanceFrom(NewPoint(0, 2)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 5
	if d := l.DistanceFrom(NewPoint(0, -5)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 3
	if d := l.DistanceFrom(NewPoint(0, 13)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	l = NewLine(NewPoint(0, 0), NewPoint(0, 0))
	answer = 5
	if d := l.DistanceFrom(NewPoint(3, 4)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}
}

func TestLineDistanceFromSquared(t *testing.T) {
	var answer float64
	l := NewLine(NewPoint(0, 0), NewPoint(0, 10))

	answer = 1
	if d := l.DistanceFromSquared(NewPoint(1, 5)); d != answer {
		t.Errorf("line, squaredDistanceFrom expected %f, got %f", answer, d)
	}

	answer = 0
	if d := l.DistanceFromSquared(NewPoint(0, 2)); d != answer {
		t.Errorf("line, squaredDistanceFrom expected %f, got %f", answer, d)
	}

	answer = 25
	if d := l.DistanceFromSquared(NewPoint(0, -5)); d != answer {
		t.Errorf("line, squaredDistanceFrom expected %f, got %f", answer, d)
	}

	answer = 9
	if d := l.DistanceFromSquared(NewPoint(0, 13)); d != answer {
		t.Errorf("line, squaredDistanceFrom expected %f, got %f", answer, d)
	}

	l = NewLine(NewPoint(0, 0), NewPoint(0, 0))
	answer = 25
	if d := l.DistanceFromSquared(NewPoint(3, 4)); d != answer {
		t.Errorf("line, squaredDistanceFrom expected %f, got %f", answer, d)
	}
}

func TestLineProject(t *testing.T) {
	l1 := NewLine(NewPoint(1, 2), NewPoint(3, 4))

	proj := l1.Project(NewPoint(1, 2))
	expected := 0.0
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	proj = l1.Project(NewPoint(3, 4))
	expected = 1.0
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	proj = l1.Project(NewPoint(2, 3))
	expected = 0.5
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	proj = l1.Project(NewPoint(5, 6))
	expected = 2.0
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	proj = l1.Project(NewPoint(-1, 0))
	expected = -1.0
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	// point off of line
	l2 := NewLine(NewPoint(1, 1), NewPoint(3, 3))
	proj = l2.Project(NewPoint(1, 2))
	expected = 0.25
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}

	// line of length 0
	l3 := NewLine(NewPoint(1, 1), NewPoint(1, 1))
	proj = l3.Project(NewPoint(1, 2))
	expected = 0.0
	if proj != expected {
		t.Errorf("line, project expected %v == %v", proj, expected)
	}
}

func TestLineMeasure(t *testing.T) {
	l1 := NewLine(NewPoint(0, 0), NewPoint(0, 4))

	measure := l1.Measure(NewPoint(1, 2))
	expected := 2.0
	if measure != expected {
		t.Errorf("line, measure expected %v == %v", measure, expected)
	}

	measure = l1.Measure(NewPoint(1, -2))
	expected = 0.0
	if measure != expected {
		t.Errorf("line, measure expected %v == %v", measure, expected)
	}

	measure = l1.Measure(NewPoint(1, 6))
	expected = 4.0
	if measure != expected {
		t.Errorf("line, measure expected %v == %v", measure, expected)
	}
}

func TestLineDistance(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(3, 4))
	if d := l.Distance(); d != 5 {
		t.Errorf("line, distance expected 5, got %f", d)
	}
}

func TestLineSquaredDistance(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(3, 4))
	if d := l.DistanceSquared(); d != 25 {
		t.Errorf("line, squaredDistance expected 25, got %f", d)
	}
}

func TestLineInterpolate(t *testing.T) {
	var answer Point
	l := NewLine(NewPoint(0, 0), NewPoint(5, 10))

	answer = NewPoint(1, 2)
	if p := l.Interpolate(.20); !p.Equal(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(4, 8)
	if p := l.Interpolate(.80); !p.Equal(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(-1, -2)
	if p := l.Interpolate(-.20); !p.Equal(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(6, 12)
	if p := l.Interpolate(1.20); !p.Equal(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}
}

func TestLineSide(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(0, 10))

	// colinear
	if o := l.Side(NewPoint(0, -5)); o != 0 {
		t.Errorf("point, expected to be colinear, got %d", o)
	}

	// right
	if o := l.Side(NewPoint(1, 5)); o != -1 {
		t.Errorf("point, expected to be on right, got %d", o)
	}

	// left
	if o := l.Side(NewPoint(-1, 5)); o != 1 {
		t.Errorf("point, expected to be on left, got %d", o)
	}
}

func TestLineIntersection(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(1, 1))

	if p, err := l.Intersection(NewLine(NewPoint(1, 0), NewPoint(2, 1))); err == nil {
		t.Errorf("line, no intersection expected, got %v", p)
	}

	if p, err := l.Intersection(NewLine(NewPoint(1, 0), NewPoint(3, 1))); err == nil {
		t.Errorf("line, no intersection expected, got %v", p)
	}

	if p, err := l.Intersection(NewLine(NewPoint(1, 1), NewPoint(2, 2))); err == nil {
		t.Errorf("line, no intersection expected, got %v", p)
	}

	answer := NewPoint(1, 1)
	if p, _ := l.Intersection(NewLine(NewPoint(1, 1), NewPoint(2, 3))); !p.Equal(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0, 0)
	if p, _ := l.Intersection(NewLine(NewPoint(1, 10), NewPoint(0, 0))); !p.Equal(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p, _ := l.Intersection(NewLine(NewPoint(0, 1), NewPoint(1, 0))); !p.Equal(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p, _ := l.Intersection(NewLine(NewPoint(0, 1), NewPoint(2, -1))); !p.Equal(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p, _ := l.Intersection(NewLine(NewPoint(0.5, 0.5), NewPoint(2, -1))); !p.Equal(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}
}

func TestLineIntersects(t *testing.T) {
	var answer bool
	l := NewLine(NewPoint(0, 0), NewPoint(1, 1))

	answer = false
	if p := l.Intersects(NewLine(NewPoint(1, 0), NewPoint(2, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	if p := l.Intersects(NewLine(NewPoint(1, 0), NewPoint(0, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	if p := l.Intersects(NewLine(NewPoint(1, 1), NewPoint(2, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	l2 := NewLine(NewPoint(0.5, 0.5), NewPoint(2, 2))
	if p := l.Intersects(l2); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	if p := l2.Intersects(l); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}
}

func TestLineMidpoint(t *testing.T) {
	var answer Point
	l := NewLine(NewPoint(0, 0), NewPoint(10, 20))

	answer = NewPoint(5, 10)
	if p := l.Midpoint(); !p.Equal(answer) {
		t.Errorf("line, midpoint expected %v, got %v", answer, p)
	}
}

func TestLineBound(t *testing.T) {
	var answer Bound
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	l := NewLine(a, b)

	answer = NewBound(1, 3, 2, 4)
	if b := l.Bound(); !b.Equal(answer) {
		t.Errorf("line, bounds expected %v, got %v", answer, b)
	}

	if b := l.Reverse().Bound(); !b.Equal(answer) {
		t.Errorf("line, bounds expected %v, got %v", answer, b)
	}
}

func TestLineReverse(t *testing.T) {
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	l := NewLine(a, b).Reverse()

	if !l.A().Equal(b) || !l.B().Equal(a) {
		t.Error("line, reverse did not work")
	}
}

func TestLineEqual(t *testing.T) {
	l1 := NewLine(NewPoint(1, 2), NewPoint(3, 4))
	l2 := NewLine(NewPoint(1, 2), NewPoint(3, 4))

	if !l1.Equal(l2) || !l2.Equal(l1) {
		t.Errorf("line, equals expcted %v == %v", l1, l2)
	}

	l3 := NewLine(NewPoint(3, 4), NewPoint(1, 2))
	if !l1.Equal(l3) || !l3.Equal(l1) {
		t.Errorf("line, equals expcted %v == %v", l1, l3)
	}
}

func TestLineWKT(t *testing.T) {
	l := NewLine(NewPoint(1, 2), NewPoint(3, 4))

	answer := "LINESTRING(1 2,3 4)"
	if s := l.WKT(); s != answer {
		t.Errorf("line, string expected %s, got %s", answer, s)
	}
}

func TestLineString(t *testing.T) {
	l := NewLine(NewPoint(1, 2), NewPoint(3, 4))

	answer := "LINESTRING(1 2,3 4)"
	if s := l.String(); s != answer {
		t.Errorf("line, string expected %s, got %s", answer, s)
	}
}
