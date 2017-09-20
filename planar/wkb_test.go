package planar

import (
	"testing"

	"github.com/paulmach/orb/internal/wkb"
)

func TestPointScan(t *testing.T) {
	p := NewPoint(0, 0)
	for i, test := range wkb.PointTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if test.X != p[0] {
			t.Errorf("test %d, incorrect x: %v != %v", i, p[0], test.X)
		}

		if test.Y != p[1] {
			t.Errorf("test %d, incorrect y: %v != %v", i, p[1], test.Y)
		}
	}
}

func TestBoundScan(t *testing.T) {
	b := Bound{}
	for i, test := range wkb.SegmentTestCases {
		err := b.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !b.Equal(pointsToLineString(test.Points).Bound()) {
			t.Errorf("test %d, incorrect rectangle: %v", i, b)
		}
	}
}

func TestLineStringScan(t *testing.T) {
	p := LineString{}
	for i, test := range wkb.LineStringTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToLineString(test.Points)) {
			t.Errorf("test %d, incorrect line string: %v", i, p)
		}
	}
}

func TestPolygonScan(t *testing.T) {
	p := Polygon{}
	for i, test := range wkb.PolygonTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToPolygon(test.Points)) {
			t.Errorf("test %d, incorrect polygon: %v", i, p)
		}
	}
}

func TestMultiPointScan(t *testing.T) {
	p := MultiPoint{}
	for i, test := range wkb.MultiPointTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToMultiPoint(test.Points)) {
			t.Errorf("test %d, incorrect multi point: %v", i, p)
		}
	}
}

func TestMultiLineStringScan(t *testing.T) {
	p := MultiLineString{}
	for i, test := range wkb.MultiLineStringTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToMultiLineString(test.Points)) {
			t.Errorf("test %d, incorrect multi linestring: %v", i, p)
		}
	}
}
func TestMultiPolygonScan(t *testing.T) {
	p := MultiPolygon{}
	for i, test := range wkb.MultiPolygonTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToMultiPolygon(test.Points)) {
			t.Errorf("test %d, incorrect multi polygon: \n%v\n%v", i, p, pointsToMultiPolygon(test.Points))
		}
	}
}

func pointsToLineString(points [][2]float64) LineString {
	ls := NewLineString()
	for _, point := range points {
		ls = append(ls, Point(point))
	}
	return ls
}

func pointsToPolygon(points [][][2]float64) Polygon {
	p := NewPolygon()
	for _, ring := range points {
		p = append(p, Ring(pointsToLineString(ring)))
	}

	return p
}

func pointsToMultiPoint(points [][2]float64) MultiPoint {
	p := NewMultiPoint()
	for _, point := range points {
		p = append(p, Point(point))
	}
	return p
}

func pointsToMultiLineString(points [][][2]float64) MultiLineString {
	mls := NewMultiLineString()
	for _, line := range points {
		mls = append(mls, pointsToLineString(line))
	}
	return mls
}

func pointsToMultiPolygon(points [][][][2]float64) MultiPolygon {
	mp := NewMultiPolygon()
	for _, poly := range points {
		mp = append(mp, pointsToPolygon(poly))
	}
	return mp
}
