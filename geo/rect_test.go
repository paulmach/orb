package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb/internal/mercator"
)

func TestRectAroundPoint(t *testing.T) {
	p := Point{
		5.42553,
		50.0359,
	}
	rect := NewRectAroundPoint(p, 1000000)
	if rect.Center().Lat() != p.Lat() {
		t.Errorf("rect, should have correct center lat point")
	}

	if rect.Center().Lon() != p.Lon() {
		t.Errorf("rect, should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if !rect.Contains(Point{3.412, 58.3838}) {
		t.Errorf("rect, should have point included in rect")
	}

	rect = NewRectAroundPoint(p, 10000.0)
	if rect.Center().Lat() != p.Lat() {
		t.Errorf("rect, should have correct center lat point")
	}

	if rect.Center().Lon() != p.Lon() {
		t.Errorf("rect, should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if rect.Contains(Point{3.412, 58.3838}) {
		t.Errorf("rect, should not have point included in rect")
	}
}

func TestNewRectFromMapTile(t *testing.T) {
	rect := NewRectFromMapTile(7, 8, 9)

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the rect
	lon, lat := mercator.ScalarInverse(7<<factor+1, 8<<factor+1, level)
	if !rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("rect, should contain point")
	}

	lon, lat = mercator.ScalarInverse(7<<factor-1, 8<<factor-1, level)
	if rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("rect, should not contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor-1, 9<<factor-1, level)
	if !rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("rect, should contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor+1, 9<<factor+1, level)
	if rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("rect, should not contain point")
	}

	rect = NewRectFromMapTile(7, 8, 35)
}

func TestRectPad(t *testing.T) {
	tests := []Rect{
		NewRectFromPoints(NewPoint(-122.559, 37.887), NewPoint(-122.521, 37.911)),
		NewRectFromPoints(NewPoint(-122.559, 15), NewPoint(-122.521, 15)),
		NewRectFromPoints(NewPoint(20, -15), NewPoint(20, -15)),
	}

	for i, b1 := range tests {
		b2 := b1.Pad(100)

		if math.Abs(b1.Height()+200-b2.Height()) > 1.0 {
			t.Errorf("rect, geoPad height incorrected for %d, expected %v, got %v", i, b1.Height()+200, b2.Height())
		}

		if math.Abs(b1.Width()+200-b2.Width()) > 1.0 {
			t.Errorf("rect, geoPad width incorrected for %d, expected %v, got %v", i, b1.Width()+200, b2.Width())
		}
	}
}
