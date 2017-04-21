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

func TestSegmentScan(t *testing.T) {
	s := Segment{}
	for i, test := range wkb.SegmentTestCases {
		err := s.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !s.Equal(pointsToSegment(test.Points)) {
			t.Errorf("test %d, incorrect line: %v", i, s)
		}
	}
}

func TestRectScan(t *testing.T) {
	r := Rect{}
	for i, test := range wkb.SegmentTestCases {
		err := r.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !r.Equal(pointsToSegment(test.Points).Bound()) {
			t.Errorf("test %d, incorrect rectangle: %v", i, r)
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

func pointsToSegment(p [][2]float64) Segment {
	return NewSegment(Point(p[0]), Point(p[1]))
}

func pointsToMultiPoint(points [][2]float64) MultiPoint {
	p := NewMultiPoint()
	for _, point := range points {
		p = append(p, Point(point))
	}
	return p
}

func pointsToLineString(points [][2]float64) LineString {
	ls := NewLineString()
	for _, point := range points {
		ls = append(ls, Point(point))
	}
	return ls
}
