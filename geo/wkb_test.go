package geo

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
				t.Errorf("test %d, incorrect error, got %v", i, err)
			}
			continue
		}

		if test.X != p[0] {
			t.Errorf("test %d, incorrect x: %v", i, p[0])
		}

		if test.Y != p[1] {
			t.Errorf("test %d, incorrect y: %v", i, p[1])
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
			t.Errorf("test %d, incorrect bound: %v", i, b)
		}
	}
}

func TestMultiPointScan(t *testing.T) {
	mp := MultiPoint{}
	for i, test := range wkb.MultiPointTestCases {
		err := mp.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !mp.Equal(pointsToMultiPoint(test.Points)) {
			t.Errorf("test %d, incorrect points: %v", i, mp)
		}
	}
}

func TestLineStringScan(t *testing.T) {
	ls := LineString{}
	for i, test := range wkb.LineStringTestCases {
		err := ls.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error: %v", i, err)
			}
			continue
		}

		if !ls.Equal(pointsToLineString(test.Points)) {
			t.Errorf("test %d, incorrect line string: %v", i, ls)
		}
	}
}

func pointsToMultiPoint(points [][2]float64) MultiPoint {
	mp := NewMultiPoint()
	for _, p := range points {
		mp = append(mp, Point(p))
	}
	return mp
}

func pointsToLineString(points [][2]float64) LineString {
	ls := NewLineString()
	for _, point := range points {
		ls = append(ls, Point(point))
	}
	return ls
}
