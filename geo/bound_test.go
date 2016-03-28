package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb/internal/mercator"
)

func TestBoundAroundPoint(t *testing.T) {
	p := Point{
		5.42553,
		50.0359,
	}
	bound := NewBoundAroundPoint(p, 1000000)
	if bound.Center().Lat() != p.Lat() {
		t.Errorf("bound, should have correct center lat point")
	}

	if bound.Center().Lng() != p.Lng() {
		t.Errorf("bound, should have correct center lng point")
	}

	//Given point is 968.9 km away from center
	if !bound.Contains(Point{3.412, 58.3838}) {
		t.Errorf("bound, should have point included in bound")
	}

	bound = NewBoundAroundPoint(p, 10000.0)
	if bound.Center().Lat() != p.Lat() {
		t.Errorf("bound, should have correct center lat point")
	}

	if bound.Center().Lng() != p.Lng() {
		t.Errorf("bound, should have correct center lng point")
	}

	//Given point is 968.9 km away from center
	if bound.Contains(Point{3.412, 58.3838}) {
		t.Errorf("bound, should not have point included in bound")
	}
}

func TestNewBoundFromMapTile(t *testing.T) {
	bound := NewBoundFromMapTile(7, 8, 9)

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the bounds
	lng, lat := mercator.ScalarInverse(7<<factor+1, 8<<factor+1, level)
	if !bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should contain point")
	}

	lng, lat = mercator.ScalarInverse(7<<factor-1, 8<<factor-1, level)
	if bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should not contain point")
	}

	lng, lat = mercator.ScalarInverse(8<<factor-1, 9<<factor-1, level)
	if !bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should contain point")
	}

	lng, lat = mercator.ScalarInverse(8<<factor+1, 9<<factor+1, level)
	if bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should not contain point")
	}

	bound = NewBoundFromMapTile(7, 8, 35)
}

func TestBoundPad(t *testing.T) {
	tests := []Bound{
		NewBoundFromPoints(NewPoint(-122.559, 37.887), NewPoint(-122.521, 37.911)),
		NewBoundFromPoints(NewPoint(-122.559, 15), NewPoint(-122.521, 15)),
		NewBoundFromPoints(NewPoint(20, -15), NewPoint(20, -15)),
	}

	for i, b1 := range tests {
		b2 := b1.Pad(100)

		if math.Abs(b1.Height()+200-b2.Height()) > 1.0 {
			t.Errorf("bound, geoPad height incorrected for %d, expected %v, got %v", i, b1.Height()+200, b2.Height())
		}

		if math.Abs(b1.Width()+200-b2.Width()) > 1.0 {
			t.Errorf("bound, geoPad width incorrected for %d, expected %v, got %v", i, b1.Width()+200, b2.Width())
		}
	}
}
