package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb/internal/mercator"
)

func center(r Rect) Point {
	return Point{
		(r[0][0] + r[1][0]) / 2.0,
		(r[0][1] + r[1][1]) / 2.0,
	}
}

func TestRectAroundPoint(t *testing.T) {

	p := Point{
		5.42553,
		50.0359,
	}
	rect := NewRectAroundPoint(p, 1000000)
	if center(rect).Lat() != p.Lat() {
		t.Errorf("should have correct center lat point")
	}

	if center(rect).Lon() != p.Lon() {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if !rect.Contains(Point{3.412, 58.3838}) {
		t.Errorf("should have point included in rect")
	}

	rect = NewRectAroundPoint(p, 10000.0)
	if center(rect).Lat() != p.Lat() {
		t.Errorf("should have correct center lat point")
	}

	if center(rect).Lon() != p.Lon() {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if rect.Contains(Point{3.412, 58.3838}) {
		t.Errorf("should not have point included in rect")
	}
}

func TestNewRectFromMapTile(t *testing.T) {
	rect := NewRectFromMapTile(7, 8, 9)

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the rect
	lon, lat := mercator.ScalarInverse(7<<factor+1, 8<<factor+1, level)
	if !rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("should contain point")
	}

	lon, lat = mercator.ScalarInverse(7<<factor-1, 8<<factor-1, level)
	if rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("should not contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor-1, 9<<factor-1, level)
	if !rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("should contain point")
	}

	lon, lat = mercator.ScalarInverse(8<<factor+1, 9<<factor+1, level)
	if rect.Contains(NewPoint(lon, lat)) {
		t.Errorf("should not contain point")
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

func TestRectExtend(t *testing.T) {
	rect := NewRect(3, 0, 5, 0)

	if r := rect.Extend(NewPoint(2, 1)); !r.Equal(rect) {
		t.Errorf("extend incorrect: %v != %v", r, rect)
	}

	answer := NewRect(6, 0, 5, -1)
	if r := rect.Extend(NewPoint(6, -1)); !r.Equal(answer) {
		t.Errorf("extend incorrect: %v != %v", r, answer)
	}
}

func TestRectUnion(t *testing.T) {
	r1 := NewRect(0, 1, 0, 1)
	r2 := NewRect(0, 2, 0, 2)

	expected := NewRect(0, 2, 0, 2)
	if r := r1.Union(r2); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}

	if r := r2.Union(r1); !r.Equal(expected) {
		t.Errorf("union incorrect: %v != %v", r, expected)
	}
}

func TestRectContains(t *testing.T) {
	rect := NewRect(2, -2, 1, -1)

	cases := []struct {
		name   string
		point  Point
		result bool
	}{
		{
			name:   "middle",
			point:  NewPoint(0, 0),
			result: true,
		},
		{
			name:   "left border",
			point:  NewPoint(-1, 0),
			result: true,
		},
		{
			name:   "ne corner",
			point:  NewPoint(2, 1),
			result: true,
		},
		{
			name:   "above",
			point:  NewPoint(0, 3),
			result: false,
		},
		{
			name:   "below",
			point:  NewPoint(0, -3),
			result: false,
		},
		{
			name:   "left",
			point:  NewPoint(-3, 0),
			result: false,
		},
		{
			name:   "right",
			point:  NewPoint(3, 0),
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := rect.Contains(tc.point)
			if v != tc.result {
				t.Errorf("incorrect contains: %v != %v", v, tc.result)
			}
		})
	}
}

func TestRectIntersects(t *testing.T) {
	rect := NewRect(0, 1, 2, 3)

	cases := []struct {
		name   string
		rect   Rect
		result bool
	}{
		{
			name:   "outside, top right",
			rect:   NewRect(5, 6, 7, 8),
			result: false,
		},
		{
			name:   "outside, top left",
			rect:   NewRect(-6, -5, 7, 8),
			result: false,
		},
		{
			name:   "outside, above",
			rect:   NewRect(0, 0.5, 7, 8),
			result: false,
		},
		{
			name:   "over the middle",
			rect:   NewRect(0, 0.5, 1, 4),
			result: true,
		},
		{
			name:   "over the left",
			rect:   NewRect(-1, 2, 1, 4),
			result: true,
		},
		{
			name:   "completely inside",
			rect:   NewRect(0.3, 0.6, 2.3, 2.6),
			result: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := rect.Intersects(tc.rect)
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}

	a := NewRect(7, 8, 6, 7)
	b := NewRect(6.1, 8.1, 6.1, 8.1)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}

	a = NewRect(1, 4, 2, 3)
	b = NewRect(2, 3, 1, 4)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}
}

func TestRectIsEmpty(t *testing.T) {
	cases := []struct {
		name   string
		rect   Rect
		result bool
	}{
		{
			name:   "regular rect",
			rect:   NewRect(1, 2, 3, 4),
			result: false,
		},
		{
			name:   "single point",
			rect:   NewRect(1, 1, 2, 2),
			result: false,
		},
		{
			name:   "horizontal bar",
			rect:   NewRect(1, 1, 2, 3),
			result: false,
		},
		{
			name:   "vertical bar",
			rect:   NewRect(1, 2, 2, 2),
			result: false,
		},
		{
			name:   "vertical bar",
			rect:   NewRect(1, 2, 2, 2),
			result: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.rect.IsEmpty()
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}

		})
	}

	// negative/malformed area
	rect := NewRect(1, 1, 2, 2)
	rect[1][0] = 0
	if !rect.IsEmpty() {
		t.Error("expected true, got false")
	}

	// negative/malformed area
	rect = NewRect(1, 1, 2, 2)
	rect[0][1] = 3
	if !rect.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestRectIsZero(t *testing.T) {
	rect := NewRect(1, 1, 2, 2)
	if rect.IsZero() {
		t.Error("expected false, got true")
	}

	rect = NewRect(0, 0, 0, 0)
	if !rect.IsZero() {
		t.Error("expected true, got false")
	}

	var r Rect
	if !r.IsZero() {
		t.Error("expected true, got false")
	}
}

func TestWKT(t *testing.T) {
	rect := NewRect(1, 2, 3, 4)

	answer := "POLYGON((1 3,1 4,2 4,2 3,1 3))"
	if s := rect.WKT(); s != answer {
		t.Errorf("wkt expected %s, got %s", answer, s)
	}
}
