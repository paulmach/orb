package planar

import (
	"math"
	"testing"
)

func TestNewLineStringPreallocate(t *testing.T) {
	ls := NewLineStringPreallocate(10, 1000)
	if l := len(ls); l != 10 {
		t.Errorf("length not set correctly, got %d", l)
	}

	if c := cap(ls); c != 1000 {
		t.Errorf("capactity not set corrctly, got %d", c)
	}
}

func TestNewLineStringFromXYData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	ls := NewLineStringFromXYData(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := ls[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("first point incorrect, got %v", point)
	}

	if point := ls[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("first point incorrect, got %v", point)
	}
}

func TestNewLineStringFromYXData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	ls := NewLineStringFromYXData(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := ls[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("first point incorrect, got %v", point)
	}

	if point := ls[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("first point incorrect, got %v", point)
	}
}

func TestNewLineStringFromXYSlice(t *testing.T) {
	data := [][]float64{
		{1, 2, -1},
		nil,
		{3, 4},
	}

	ls := NewLineStringFromXYSlice(data)
	if l := len(ls); l != 2 {
		t.Errorf("should take full length of data, expected %d, got %d", 2, l)
	}

	if point := ls[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("first point incorrect, got %v", point)
	}

	if point := ls[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("first point incorrect, got %v", point)
	}
}

func TestNewLineStringFromYXSlice(t *testing.T) {
	data := [][]float64{
		{1, 2},
		{3, 4, -1},
	}

	ls := NewLineStringFromYXSlice(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := ls[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("first point incorrect, got %v", point)
	}

	if point := ls[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("first point incorrect, got %v", point)
	}
}

func TestLineStringDistance(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
	)

	if d := ls.Distance(); d != 7 {
		t.Errorf("distance got: %f, expected 7.0", d)
	}
}

func TestLineStringDistanceFrom(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
		NewPoint(4, 0),
	)

	cases := []struct {
		name   string
		point  Point
		result float64
	}{
		{
			point:  NewPoint(4.5, 1.5),
			result: 0.5,
		},
		{
			point:  NewPoint(0.4, 1.5),
			result: 0.4,
		},
		{
			point:  NewPoint(-0.3, 1.5),
			result: 0.3,
		},
		{
			point:  NewPoint(0.3, 2.8),
			result: 0.2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := ls.DistanceFrom(tc.point)
			if math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestLineStringDistanceFromSquared(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 3),
		NewPoint(4, 3),
		NewPoint(4, 0),
	)

	cases := []struct {
		name   string
		point  Point
		result float64
	}{
		{
			point:  NewPoint(4.5, 1.5),
			result: 0.25,
		},
		{
			point:  NewPoint(0.4, 1.5),
			result: 0.16,
		},
		{
			point:  NewPoint(-0.3, 1.5),
			result: 0.09,
		},
		{
			point:  NewPoint(0.3, 2.8),
			result: 0.04,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := ls.DistanceFromSquared(tc.point)
			if math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestLineStringProject(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(6, 8),
		NewPoint(12, 0),
	)

	cases := []struct {
		name   string
		point  Point
		result float64
	}{
		{
			name:   "middle of first line",
			point:  NewPoint(3, 4),
			result: 0.25,
		},
		{
			name:   "equal to first point",
			point:  NewPoint(0, 0),
			result: 0.0,
		},
		{
			name:   "equal to last point",
			point:  NewPoint(12, 0),
			result: 1.0,
		},
		{
			name:   "closest point on the line string",
			point:  NewPoint(-1, -1),
			result: 0.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ls.Project(tc.point)
			if result != tc.result {
				t.Errorf("project incorrect: %v != %v", result, tc.result)
			}
		})
	}
}

func TestLineStringInterpolate(t *testing.T) {
	ls := NewLineString()
	ls = append(ls, NewPoint(0, 0))
	ls = append(ls, NewPoint(1, 1))
	ls = append(ls, NewPoint(2, 2))
	ls = append(ls, NewPoint(3, 3))

	cases := []struct {
		percent float64
		result  Point
	}{
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

	for i, tc := range cases {
		r := ls.Interpolate(tc.percent)
		if r != tc.result {
			t.Errorf("incorrect result for %d: %v != %v", i, r, tc.result)
		}
	}
}

func TestLineStringCentroid(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(5, 0),
		NewPoint(5, 4),
		NewPoint(8, 4))

	expected := NewPoint(13.0/3.0, 5.0/3.0)
	if c := ls.Centroid(); !c.Equal(expected) {
		t.Errorf("incorrect result: %v != %v", c, expected)
	}
}

func TestLineStringReverse(t *testing.T) {
	t.Run("1 point line", func(t *testing.T) {
		ls := append(NewLineString(), NewPoint(1, 2))
		rs := ls.Reverse()

		if !rs.Equal(ls) {
			t.Errorf("1 point lines should be equal if reversed")
		}
	})

	cases := []struct {
		name   string
		input  LineString
		output LineString
	}{
		{
			name:   "2 point line",
			input:  append(NewLineString(), NewPoint(1, 2), NewPoint(3, 4)),
			output: append(NewLineString(), NewPoint(3, 4), NewPoint(1, 2)),
		},
		{
			name:   "3 point line",
			input:  append(NewLineString(), NewPoint(1, 2), NewPoint(3, 4), NewPoint(5, 6)),
			output: append(NewLineString(), NewPoint(5, 6), NewPoint(3, 4), NewPoint(1, 2)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reversed := tc.input.Reverse()

			if !reversed.Equal(tc.output) {
				t.Errorf("line should be reversed: %v", reversed)
			}

			if tc.input.Equal(reversed) {
				t.Errorf("should create new line string object")
			}
		})
	}
}

func TestLineStringWKT(t *testing.T) {
	ls := NewLineString()

	answer := "EMPTY"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	ls = append(ls, NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	ls = append(ls, NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}
}
