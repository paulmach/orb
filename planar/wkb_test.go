package planar

import (
	"testing"

	"github.com/paulmach/orb"
)

var testPathWKB = []byte{1, 2, 0, 0, 0, 6, 0, 0, 0, 205, 228, 155, 109, 110, 114, 87, 192, 174, 158, 147, 222, 55, 50, 64, 64, 134, 56, 214, 197, 109, 114, 87, 192, 238, 235, 192, 57, 35, 50, 64, 64, 173, 47, 18, 218, 114, 114, 87, 192, 25, 4, 86, 14, 45, 50, 64, 64, 10, 75, 60, 160, 108, 114, 87, 192, 224, 161, 40, 208, 39, 50, 64, 64, 149, 159, 84, 251, 116, 114, 87, 192, 96, 147, 53, 234, 33, 50, 64, 64, 195, 158, 118, 248, 107, 114, 87, 192, 89, 139, 79, 1, 48, 50, 64, 64}

func TestPointScan(t *testing.T) {
	p := NewPoint(0, 0)

	if err := p.Scan(123); err != orb.ErrUnsupportedDataType {
		t.Errorf("incorrect error, got %v", err)
	}

	type testCase struct {
		x, y float64
		data []byte
	}

	tests := []testCase{
		{ // little endian
			x: -122.4546440212, y: 37.7382859071,
			data: []byte{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		},
		{ // big endian
			x: -122.4546440212, y: 37.7382859071,
			data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39, 17, 11, 205},
		},
		{ // mysql srid+wkb
			x: -122.671129, y: 38.177484,
			data: []byte{215, 15, 0, 0, 1, 1, 0, 0, 0, 107, 153, 12, 199, 243, 170, 94, 192, 25, 200, 179, 203, 183, 22, 67, 64},
		},
		{
			x: -93.787988, y: 32.392335,
			data: []byte{1, 1, 0, 0, 0, 253, 104, 56, 101, 110, 114, 87, 192, 192, 9, 133, 8, 56, 50, 64, 64},
		},
	}

	for i, test := range tests {
		err := p.Scan(test.data)
		if err != nil {
			t.Errorf("test %d had error %v", i, err)
		}

		if test.x != p[0] {
			t.Errorf("test %d incorrect x, got %v", i, p[0])
		}

		if test.y != p[1] {
			t.Errorf("test %d incorrect y, got %v", i, p[1])
		}
	}

	// error conditions
	err := p.Scan([]byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}

	err = p.Scan([]byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64})
	if err != orb.ErrNotWKB {
		t.Errorf("incorrect error, got %v", err)
	}

	err = p.Scan([]byte{0, 2, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}
}

func TestLineScan(t *testing.T) {
	l := Line{}

	if err := l.Scan(123); err != orb.ErrUnsupportedDataType {
		t.Errorf("incorrect error, got %v", err)
	}

	type testCase struct {
		line Line
		data []byte
	}

	tests := []testCase{
		{
			line: NewLine(NewPoint(-123.016508, 38.040608), NewPoint(-122.670176, 38.548019)),
			data: []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		},
		{
			line: NewLine(NewPoint(-123.016508, 38.040608), NewPoint(-122.670176, 38.548019)),
			data: []byte{215, 15, 0, 0, 1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		},
		{
			line: NewLine(NewPoint(-72.796408, -45.407131), NewPoint(-72.688541, -45.384987)),
			data: []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 117, 145, 66, 89, 248, 50, 82, 192, 9, 24, 93, 222, 28, 180, 70, 192, 33, 61, 69, 14, 17, 44, 82, 192, 77, 49, 7, 65, 71, 177, 70, 192},
		},
	}

	for i, test := range tests {
		err := l.Scan(test.data)
		if err != nil {
			t.Errorf("test %d had error %v", i, err)
		}

		if !l.Equal(test.line) {
			t.Errorf("test %d incorrect line, got %v", i, l)
		}
	}

	// error conditions
	err := l.Scan([]byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}

	err = l.Scan([]byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}

	err = l.Scan([]byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}
}

func TestPathScan(t *testing.T) {
	path := NewPath()

	if err := path.Scan(123); err != orb.ErrUnsupportedDataType {
		t.Errorf("incorrect error, got %v", err)
	}

	type testCase struct {
		ps   PointSet
		data []byte
	}

	ps := append(NewPointSet(),
		NewPoint(-93.78799, 32.39233),
		NewPoint(-93.78795, 32.3917),
		NewPoint(-93.78826, 32.392),
		NewPoint(-93.78788, 32.39184),
		NewPoint(-93.78839, 32.39166),
		NewPoint(-93.78784, 32.39209),
	)

	tests := []testCase{
		{
			ps:   ps,
			data: testPathWKB,
		},
		{
			ps:   ps,
			data: append([]byte{215, 15, 0, 0}, testPathWKB...),
		},
	}

	for i, test := range tests {
		err := path.Scan(test.data)
		if err != nil {
			t.Errorf("test %d had error %v", i, err)
		}

		if !ps.Equal(test.ps) {
			t.Errorf("test %d incorrect point set, got %v", i, ps)
		}
	}

	// error conditions
	err := path.Scan([]byte{0, 0, 0, 0, 1})
	if err != orb.ErrNotWKB {
		t.Errorf("incorrect error, got %v", err)
	}

	err = path.Scan([]byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}

	err = path.Scan([]byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64})
	if err != orb.ErrIncorrectGeometry {
		t.Errorf("incorrect error, got %v", err)
	}
}

func TestWKBPolygon(t *testing.T) {
	// raw WKB polygon data
	data := []byte{1, 3, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 222, 90, 38, 195, 241, 110, 73, 64, 229, 179, 60, 15, 238, 190, 22, 64, 94, 189, 138, 140, 14, 110, 73, 64, 24, 11, 67, 228, 244, 213, 22, 64, 29, 119, 74, 7, 235, 109, 73, 64, 190, 22, 244, 222, 24, 178, 22, 64, 222, 90, 38, 195, 241, 110, 73, 64, 229, 179, 60, 15, 238, 190, 22, 64}

	// pointset
	ps := NewPointSet()

	expected := "MULTIPOINT(50.866753 5.686455,50.859819 5.708942,50.858735 5.673923,50.866753 5.686455)"
	err := ps.Scan(data)
	if err != nil {
		t.Errorf("should not return error, got %v", err)
	}

	if ps.String() != expected {
		t.Errorf("incorrect point set, got %v", ps)
	}

	// path
	path := NewPath()

	expected = "LINESTRING(50.866753 5.686455,50.859819 5.708942,50.858735 5.673923,50.866753 5.686455)"
	err = path.Scan(data)
	if err != nil {
		t.Errorf("should not return error, got %v", err)
	}

	if path.String() != expected {
		t.Errorf("incorrect path, got %v", path)
	}
}
