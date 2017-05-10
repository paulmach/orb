package geo

import (
	"math"
	"math/rand"
	"testing"
)

func TestNewLineStringFromEncoding(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		ls1 := NewLineString()
		for i := 0; i < 100; i++ {
			ls1 = append(ls1, Point{rand.Float64(), rand.Float64()})
		}

		encoded := ls1.Encode(int(1.0 / epsilon))
		ls2 := NewLineStringFromEncoding(encoded, int(1.0/epsilon))

		if len(ls2) != 100 {
			t.Fatalf("encodeDecode length mismatch: %d != 100", len(ls2))
		}

		for i := 0; i < 100; i++ {
			a := ls1[i]
			b := ls2[i]

			if e := math.Abs(a[0] - b[0]); e > epsilon {
				t.Errorf("x error too big: %v != %v", b[0], a[0])
			}

			if e := math.Abs(a[1] - b[1]); e > epsilon {
				t.Errorf("y error too big: %v != %v", b[1], a[1])
			}
		}
	}
}

func TestNewLineStringFromXYData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	ls := NewLineStringFromXYData(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data, %d != %d", l, len(data))
	}

	if point := ls[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("first point incorrect: %v", point)
	}

	if point := ls[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("last point incorrect: %v", point)
	}
}

func TestNewLineStringFromYXData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	ls := NewLineStringFromYXData(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data: %v != %v", l, len(data))
	}

	if point := ls[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("first point incorrect: %v", point)
	}

	if point := ls[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("last point incorrect: %v", point)
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
		t.Errorf("should take full length of data: %v != %v", l, 2)
	}

	if point := ls[0]; !point.Equal(Point{1, 2}) {
		t.Errorf("first point incorrect: %v", point)
	}

	if point := ls[1]; !point.Equal(Point{3, 4}) {
		t.Errorf("last point incorrect: %v", point)
	}
}

func TestNewLineStringFromYXSlice(t *testing.T) {
	data := [][]float64{
		{1, 2},
		{3, 4, -1},
	}

	ls := NewLineStringFromYXSlice(data)
	if l := len(ls); l != len(data) {
		t.Errorf("should take full length of data: %v != %v", l, len(data))
	}

	if point := ls[0]; !point.Equal(Point{2, 1}) {
		t.Errorf("first point incorrect: %v", point)
	}

	if point := ls[1]; !point.Equal(Point{4, 3}) {
		t.Errorf("last point incorrect: %v", point)
	}
}

func TestLineStringEncode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		ls := NewLineString()
		for i := 0; i < 100; i++ {
			ls = append(ls, Point{rand.Float64(), rand.Float64()})
		}

		encoded := ls.Encode()
		for _, c := range encoded {
			if c < 63 || c > 127 {
				t.Errorf("result out of range: %d", c)
			}
		}
	}

	// empty path
	path := NewLineString()
	if v := path.Encode(); v != "" {
		t.Errorf("empty path should be empty string: %v", v)
	}
}

func TestLineStringReverse(t *testing.T) {
	t.Run("1 point line", func(t *testing.T) {
		ls := append(NewLineString(), NewPoint(1, 2))
		rs := ls.Clone()
		rs.Reverse()

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
			reversed := tc.input
			reversed.Reverse()

			if !reversed.Equal(tc.output) {
				t.Errorf("line should be reversed: %v", reversed)
			}

			if !tc.input.Equal(reversed) {
				t.Errorf("should reverse inplace")
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
