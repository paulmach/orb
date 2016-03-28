package geo

import (
	"math"
	"math/rand"
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

func TestNewPathFromEncoding(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		p := NewPath()
		for i := 0; i < 100; i++ {
			p = append(p, Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode(int(1.0 / epsilon))
		path := NewPathFromEncoding(encoded, int(1.0/epsilon))

		if len(path) != 100 {
			t.Fatalf("path, encodeDecode length mismatch: %d != 100", len(path))
		}

		for i := 0; i < 100; i++ {
			a := p[i]
			b := path[i]

			if e := math.Abs(a[0] - b[0]); e > epsilon {
				t.Errorf("path, encodeDecode X error too big: %f", e)
			}

			if e := math.Abs(a[1] - b[1]); e > epsilon {
				t.Errorf("path, encodeDecode Y error too big: %f", e)
			}
		}
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

func TestPathEncode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		p := NewPath()
		for i := 0; i < 100; i++ {
			p = append(p, Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode()
		for _, c := range encoded {
			if c < 63 || c > 127 {
				t.Errorf("path, encode result out of range: %d", c)
			}
		}
	}

	// empty path
	path := NewPath()
	if path.Encode() != "" {
		t.Error("path, encode empty path should be empty string")
	}
}

func TestPathGeoJSON(t *testing.T) {
	p := append(NewPath(), NewPoint(1, 2))

	f := p.GeoJSON()
	if !f.Geometry.IsLineString() {
		t.Errorf("path, should be linestring geometry")
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
