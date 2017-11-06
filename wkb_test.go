package orb

import (
	"encoding/hex"
	"testing"
)

func TestPointScan(t *testing.T) {
	cases := []struct {
		Point Point
		X, Y  float64
		Data  interface{}
		Err   error
	}{
		{ // little endian
			Point: Point{-122.4546440212, 37.7382859071},
			Data:  []byte{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		},
		{ // big endian
			Point: Point{-122.4546440212, 37.7382859071},
			Data:  []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39, 17, 11, 205},
		},
		{ // mysql srid+wkb
			Point: Point{-122.671129, 38.177484},
			Data:  []byte{215, 15, 0, 0, 1, 1, 0, 0, 0, 107, 153, 12, 199, 243, 170, 94, 192, 25, 200, 179, 203, 183, 22, 67, 64},
		},
		{ // mysql srid+wkb, empty srid
			Point: Point{-122.671129, 38.177484},
			Data:  []byte{0, 0, 0, 0, 1, 1, 0, 0, 0, 107, 153, 12, 199, 243, 170, 94, 192, 25, 200, 179, 203, 183, 22, 67, 64},
		},
		{
			Point: Point{-93.787988, 32.392335},
			Data:  []byte{1, 1, 0, 0, 0, 253, 104, 56, 101, 110, 114, 87, 192, 192, 9, 133, 8, 56, 50, 64, 64},
		},
		{
			Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{0, 2, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
	}

	for i, tc := range cases {
		p := Point{}

		err := p.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if p != tc.Point {
			t.Errorf("test %d, incorrect point: %v != %v", i, p, tc.Point)
		}
	}
}

func TestBoundScan(t *testing.T) {
	cases := []struct {
		LineString LineString
		Data       interface{}
		Err        error
	}{
		{
			LineString: LineString{{-123.016508, 38.040608}, {-122.670176, 38.548019}},
			Data:       []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		},
		{
			LineString: LineString{{-123.016508, 38.040608}, {-122.670176, 38.548019}},
			Data:       []byte{215, 15, 0, 0, 1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		},
		{
			LineString: LineString{{-72.796408, -45.407131}, {-72.688541, -45.384987}},
			Data:       []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 117, 145, 66, 89, 248, 50, 82, 192, 9, 24, 93, 222, 28, 180, 70, 192, 33, 61, 69, 14, 17, 44, 82, 192, 77, 49, 7, 65, 71, 177, 70, 192},
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
			Err:  ErrIncorrectGeometry,
		},
	}

	for i, tc := range cases {
		b := Bound{}
		err := b.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !b.Equal(tc.LineString.Bound()) {
			t.Errorf("test %d, incorrect rectangle: %v", i, b)
		}
	}
}

func TestLineStringScan(t *testing.T) {
	testData := []byte{1, 2, 0, 0, 0, 6, 0, 0, 0, 205, 228, 155, 109, 110, 114, 87, 192, 174, 158, 147, 222, 55, 50, 64, 64, 134, 56, 214, 197, 109, 114, 87, 192, 238, 235, 192, 57, 35, 50, 64, 64, 173, 47, 18, 218, 114, 114, 87, 192, 25, 4, 86, 14, 45, 50, 64, 64, 10, 75, 60, 160, 108, 114, 87, 192, 224, 161, 40, 208, 39, 50, 64, 64, 149, 159, 84, 251, 116, 114, 87, 192, 96, 147, 53, 234, 33, 50, 64, 64, 195, 158, 118, 248, 107, 114, 87, 192, 89, 139, 79, 1, 48, 50, 64, 64}

	cases := []struct {
		LineString LineString
		Data       interface{}
		Err        error
	}{
		{
			LineString: LineString{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
			Data:       testData,
		},
		{
			LineString: LineString{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
			Data:       append([]byte{215, 15, 0, 0}, testData...),
		},
		{
			LineString: LineString{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
			Data:       append([]byte{0, 0, 0, 0}, testData...),
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
			Err:  ErrIncorrectGeometry,
		},
	}

	for i, tc := range cases {
		ls := LineString{}

		err := ls.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !ls.Equal(tc.LineString) {
			t.Errorf("test %d, incorrect line string: %v", i, ls)
		}
	}
}

func TestPolygonScan(t *testing.T) {
	cases := []struct {
		Polygon Polygon
		Data    interface{}
		Err     error
	}{
		{
			Polygon: Polygon{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}},
			Data:    mustDecode(`01030000000100000005000000000000000000000000000000000000000000000000000000000000000000F03F000000000000F03F000000000000F03F000000000000F03F000000000000000000000000000000000000000000000000`),
		},
		{
			Polygon: Polygon{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}, {{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
			Data:    mustDecode(`0103000000020000000500000000000000000000000000000000000000000000000000000000000000000000400000000000000040000000000000004000000000000000400000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000F03F0000000000000000000000000000F03F000000000000F03F0000000000000000000000000000F03F00000000000000000000000000000000`),
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1},
			Err:  ErrNotWKB,
		},
	}

	for i, tc := range cases {
		p := Polygon{}

		err := p.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(tc.Polygon) {
			t.Errorf("test %d, incorrect polygon: %v", i, p)
		}
	}
}

func TestMultiPointScan(t *testing.T) {
	cases := []struct {
		MultiPoint MultiPoint
		Data       interface{}
		Err        error
	}{
		{
			MultiPoint: MultiPoint{{1, 2}, {3, 4}, {5, 6}},
			Data:       []byte{0, 0, 0, 0, 1, 4, 0, 0, 0, 3, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 24, 64},
		},
		{
			MultiPoint: MultiPoint{{50.866753, 5.686455}, {50.859819, 5.708942}},
			Data:       []byte{1, 4, 0, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 0, 222, 90, 38, 195, 241, 110, 73, 64, 229, 179, 60, 15, 238, 190, 22, 64, 1, 1, 0, 0, 0, 94, 189, 138, 140, 14, 110, 73, 64, 24, 11, 67, 228, 244, 213, 22, 64},
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
			Err:  ErrIncorrectGeometry,
		},
	}

	for i, tc := range cases {
		mp := MultiPoint{}
		err := mp.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !mp.Equal(tc.MultiPoint) {
			t.Errorf("test %d, incorrect multi point: %v", i, mp)
		}
	}
}

func TestMultiLineStringScan(t *testing.T) {
	cases := []struct {
		MultiLineString MultiLineString
		Data            interface{}
		Err             error
	}{
		{
			Data:            mustDecode(`010500000001000000010200000003000000000000000000F03F00000000000000400000000000000840000000000000104000000000000014400000000000001840`),
			MultiLineString: MultiLineString{{{1, 2}, {3, 4}, {5, 6}}},
		},
		{
			Data:            mustDecode(`010500000002000000010200000003000000000000000000F03F0000000000000040000000000000084000000000000010400000000000001440000000000000184001020000000400000000000000000018400000000000001440000000000000104000000000000008400000000000000040000000000000F03F00000000000000000000000000000000`),
			MultiLineString: MultiLineString{{{1, 2}, {3, 4}, {5, 6}}, {{6, 5}, {4, 3}, {2, 1}, {0, 0}}},
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
			Err:  ErrIncorrectGeometry,
		},
	}

	for i, tc := range cases {
		mls := MultiLineString{}

		err := mls.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !mls.Equal(tc.MultiLineString) {
			t.Errorf("test %d, incorrect multi linestring: %v", i, mls)
		}
	}
}

func TestMultiPolygonScan(t *testing.T) {
	cases := []struct {
		MultiPolygon MultiPolygon
		Data         interface{}
		Err          error
	}{
		{
			Data:         mustDecode(`01060000000100000001030000000100000005000000000000000000000000000000000000000000000000000000000000000000F03F000000000000F03F000000000000F03F000000000000F03F000000000000000000000000000000000000000000000000`),
			MultiPolygon: MultiPolygon{{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}},
		},
		{
			Data:         mustDecode(`0106000000020000000103000000010000000500000000000000000000000000000000000000000000000000000000000000000000400000000000000040000000000000004000000000000000400000000000000000000000000000000000000000000000000103000000010000000500000000000000000000000000000000000000000000000000F03F0000000000000000000000000000F03F000000000000F03F0000000000000000000000000000F03F00000000000000000000000000000000`),
			MultiPolygon: MultiPolygon{{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}}, {{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}},
		},
		{
			Data: 123,
			Err:  ErrUnsupportedDataType,
		},
		{
			Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			Err:  ErrIncorrectGeometry,
		},
		{
			Data: []byte{3, 5, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			Err:  ErrNotWKB,
		},
		{
			Data: []byte{1, 4, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
			Err:  ErrIncorrectGeometry,
		},
	}

	for i, tc := range cases {
		mp := MultiPolygon{}

		err := mp.Scan(tc.Data)
		if err != nil {
			if err != tc.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !mp.Equal(tc.MultiPolygon) {
			t.Errorf("test %d, incorrect multi polygon: %v", i, mp)
		}
	}
}

func mustDecode(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return data
}
