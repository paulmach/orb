package planar

import (
	"math"
	"testing"
)

func TestNewPathPreallocate(t *testing.T) {
	p := NewPathPreallocate(10, 1000)
	if l := len(p); l != 10 {
		t.Errorf("path, length not set correctly, got %d", l)
	}

	if c := cap(p); c != 1000 {
		t.Errorf("path, capactity not set corrctly, got %d", c)
	}
}

func TestNewPathFromXYData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	p := NewPathFromXYData(data)
	if l := len(p); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromYXData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	p := NewPathFromYXData(data)
	if l := len(p); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromXYSlice(t *testing.T) {
	data := [][]float64{
		{1, 2, -1},
		nil,
		{3, 4},
	}

	p := NewPathFromXYSlice(data)
	if l := len(p); l != 2 {
		t.Errorf("path, should take full length of data, expected %d, got %d", 2, l)
	}

	if point := p[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromYXSlice(t *testing.T) {
	data := [][]float64{
		{1, 2},
		{3, 4, -1},
	}

	p := NewPathFromYXSlice(data)
	if l := len(p); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestPathDistance(t *testing.T) {
	p := append(NewPath(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
	)

	if d := p.Distance(); d != 7 {
		t.Errorf("path, distance got: %f, expected 7.0", d)
	}
}

func TestPathDistanceFrom(t *testing.T) {
	var answer float64

	p := append(NewPath(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
		NewPoint(4, 0),
	)

	answer = 0.5
	if d := p.DistanceFrom(NewPoint(4.5, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.4
	if d := p.DistanceFrom(NewPoint(0.4, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.3
	if d := p.DistanceFrom(NewPoint(-0.3, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.2
	if d := p.DistanceFrom(NewPoint(0.3, 2.8)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}
}

func TestPathDistanceFromSquared(t *testing.T) {
	var answer float64

	p := append(NewPath(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
		NewPoint(4, 0),
	)

	answer = 0.25
	if d := p.DistanceFromSquared(NewPoint(4.5, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.16
	if d := p.DistanceFromSquared(NewPoint(0.4, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.09
	if d := p.DistanceFromSquared(NewPoint(-0.3, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.04
	if d := p.DistanceFromSquared(NewPoint(0.3, 2.8)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}
}

func TestPathMeasure(t *testing.T) {
	p := append(NewPath(),
		NewPoint(0, 0),
		NewPoint(6, 8),
		NewPoint(12, 0),
	)

	result := p.Measure(NewPoint(3, 4))
	expected := 5.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// coincident with start point
	result = p.Measure(NewPoint(0, 0))
	expected = 0.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// coincident with end point
	result = p.Measure(NewPoint(12, 0))
	expected = 20.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// closest point on path
	result = p.Measure(NewPoint(-1, -1))
	expected = 0.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}
}

func TestPathInterpolate(t *testing.T) {
	p := NewPath()
	p = append(p, NewPoint(0, 0))
	p = append(p, NewPoint(1, 1))
	p = append(p, NewPoint(2, 2))
	p = append(p, NewPoint(3, 3))

	type testCase struct {
		Percent float64
		Result  Point
	}

	tests := []testCase{
		{-0.1, NewPoint(0, 0)},
		{0.0, NewPoint(0, 0)},
		{0.1, NewPoint(0.3, 0.3)},
		{0.2, NewPoint(0.6, 0.6)},
		{0.25, NewPoint(0.75, 0.75)},
		{0.5, NewPoint(1.5, 1.5)},
		{0.75, NewPoint(2.25, 2.25)},
		{1.0, NewPoint(3, 3)},
		{1.1, NewPoint(3, 3)},
	}

	for i, test := range tests {
		r := p.Interpolate(test.Percent)
		if r != test.Result {
			t.Errorf("incorrect result for %d: got %v, expected %v", i, r, test.Result)
		}
	}
}

func TestPathWKT(t *testing.T) {
	p := NewPath()

	answer := "EMPTY"
	if s := p.WKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p = append(p, NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := p.WKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p = append(p, NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := p.WKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}
}

func TestPathString(t *testing.T) {
	p := NewPath()

	answer := "EMPTY"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p = append(p, NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p = append(p, NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}
}
