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

func TestBoundScan(t *testing.T) {
	b := Bound{}
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

func TestPointSetScan(t *testing.T) {
	p := PointSet{}
	for i, test := range wkb.PointSetTestCases {
		err := p.Scan(test.Data)
		if err != nil {
			if err != test.Err {
				t.Errorf("test %d, incorrect error, got %v", i, err)
			}
			continue
		}

		if !p.Equal(pointsToPointSet(test.Points)) {
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

func pointsToPointSet(points [][2]float64) PointSet {
	p := NewPointSet()
	for _, point := range points {
		p = append(p, Point(point))
	}
	return p
}

func pointsToPath(points [][2]float64) Path {
	p := NewPath()
	for _, point := range points {
		p = append(p, Point(point))
	}
	return p
}
