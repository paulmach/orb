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
			t.Errorf("test %d incorrect x, got %v", i, p[0])
		}

		if test.Y != p[1] {
			t.Errorf("test %d incorrect y, got %v", i, p[1])
		}
	}
}

func TestRectScan(t *testing.T) {
	b := Rect{}
	for i, test := range wkb.LineTestCases {
		err := b.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error, got %v", i, err)
			}
			continue
		}

		if !b.Equal(pointsToPath(test.Points).Bound()) {
			t.Errorf("test %d incorrect bound, got %v", i, b)
		}
	}
}

func TestMultiPointScan(t *testing.T) {
	p := MultiPoint{}
	for i, test := range wkb.MultiPointTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error, got %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToMultiPoint(test.Points)) {
			t.Errorf("test %d incorrect point set, got %v", i, p)
		}
	}
}

func TestPathScan(t *testing.T) {
	p := Path{}
	for i, test := range wkb.PathTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error, got %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToPath(test.Points)) {
			t.Errorf("test %d incorrect path, got %v", i, p)
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

func pointsToPath(points [][2]float64) Path {
	p := NewPath()
	for _, point := range points {
		p = append(p, Point(point))
	}
	return p
}
