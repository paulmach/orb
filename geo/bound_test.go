package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func center(r orb.Bound) orb.Point {
	return orb.Point{
		(r[0][0] + r[1][0]) / 2.0,
		(r[0][1] + r[1][1]) / 2.0,
	}
}

func TestBoundAroundPoint(t *testing.T) {
	p := orb.Point{
		5.42553,
		50.0359,
	}

	bound := NewBoundAroundPoint(p, 1000000)
	if center(bound)[1] != p[1] {
		t.Errorf("should have correct center lat point")
	}

	if center(bound)[0] != p[0] {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if !bound.Contains(orb.Point{3.412, 58.3838}) {
		t.Errorf("should have point included in bound")
	}

	bound = NewBoundAroundPoint(p, 10000.0)
	if center(bound)[1] != p[1] {
		t.Errorf("should have correct center lat point")
	}

	if center(bound)[0] != p[0] {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if bound.Contains(orb.Point{3.412, 58.3838}) {
		t.Errorf("should not have point included in bound")
	}
}

func TestBoundPad(t *testing.T) {
	cases := []struct {
		name  string
		bound orb.Bound
	}{
		{
			name:  "test bound",
			bound: orb.NewBoundFromPoints(orb.NewPoint(-122.559, 37.887), orb.NewPoint(-122.521, 37.911)),
		},
		{
			name:  "no width",
			bound: orb.NewBoundFromPoints(orb.NewPoint(-122.559, 15), orb.NewPoint(-122.521, 15)),
		},
		{
			name:  "no area",
			bound: orb.NewBoundFromPoints(orb.NewPoint(20, -15), orb.NewPoint(20, -15)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b2 := BoundPad(tc.bound, 100)

			v1 := BoundHeight(tc.bound) + 200
			v2 := BoundHeight(b2)
			if math.Abs(v1-v2) > 1.0 {
				t.Errorf("height incorrected: %f != %f", v1, v2)
			}

			v1 = BoundWidth(tc.bound) + 200
			v2 = BoundWidth(b2)
			if math.Abs(v1-v2) > 1.0 {
				t.Errorf("height incorrected: %f != %f", v1, v2)
			}
		})
	}
}
