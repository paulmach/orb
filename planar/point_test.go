package planar

import "testing"

func TestNewPoint(t *testing.T) {
	p := NewPoint(1, 2)
	if p.X() != 1 {
		t.Errorf("point, expected 1, got %f", p.X())
	}

	if p.Y() != 2 {
		t.Errorf("point, expected 2, got %f", p.Y())
	}
}

func TestPointDistanceFrom(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(3, 4)

	if d := p1.DistanceFrom(p2); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}

	if d := p2.DistanceFrom(p1); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}
}

func TestPointSquaredDistanceFrom(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(3, 4)

	if d := p1.DistanceFromSquared(p2); d != 25 {
		t.Errorf("point, squaredDistanceFrom expected 25, got %f", d)
	}

	if d := p2.DistanceFromSquared(p1); d != 25 {
		t.Errorf("point, squaredDistanceFrom expected 25, got %f", d)
	}
}

func TestPointMidpoint(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(10, 20)

	answer := NewPoint(5, 10)
	if m := p1.Midpoint(p2); !m.Equal(answer) {
		t.Errorf("point, midpoint expected %v, got %v", answer, m)
	}
}

func TestPointAdd(t *testing.T) {
	p := NewPoint(1, 2)
	v := NewVector(3, 4)

	answer := NewPoint(4, 6)
	p2 := p.Add(v)
	if !p2.Equal(answer) {
		t.Errorf("point, add expect %v == %v", p2, answer)
	}
}

func TestPointSub(t *testing.T) {
	p1 := NewPoint(3, 4)
	p2 := NewPoint(1, 3)

	answer := NewVector(2, 1)
	v := p1.Sub(p2)
	if !v.Equal(answer) {
		t.Errorf("point, subtract expect %v == %v", v, answer)
	}
}

func TestPointEqual(t *testing.T) {
	p1 := NewPoint(1, 0)
	p2 := NewPoint(1, 0)

	p3 := NewPoint(2, 3)
	p4 := NewPoint(2, 4)

	if !p1.Equal(p2) {
		t.Errorf("point, equals expect %v == %v", p1, p2)
	}

	if p2.Equal(p3) {
		t.Errorf("point, equals expect %v != %v", p2, p3)
	}

	if p3.Equal(p4) {
		t.Errorf("point, equals expect %v != %v", p3, p4)
	}
}

func TestPointWKT(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.WKT(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.String(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}
