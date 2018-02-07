package orb

import (
	"testing"
)

func TestMultiPolygonBound(t *testing.T) {
	// should be union of polygons
	mp := MultiPolygon{
		{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}},
		{{{1, 1}, {1, 3}, {3, 3}, {3, 1}, {1, 1}}},
	}

	b := mp.Bound()
	if !b.Equal(Bound{Min: Point{0, 0}, Max: Point{3, 3}}) {
		t.Errorf("incorrect bound: %v", b)
	}
}
