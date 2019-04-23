package geo

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

var epsilon = 1e-6

func TestDistance(t *testing.T) {
	p1 := orb.Point{-1.8444, 53.1506}
	p2 := orb.Point{0.1406, 52.2047}

	if d := Distance(p1, p2); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = orb.Point{0.5, 30}
	p2 = orb.Point{-0.5, 30}

	dFast := Distance(p1, p2)

	p1 = orb.Point{179.5, 30}
	p2 = orb.Point{-179.5, 30}

	if d := Distance(p1, p2); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestDistanceHaversine(t *testing.T) {
	p1 := orb.Point{-1.8444, 53.1506}
	p2 := orb.Point{0.1406, 52.2047}

	if d := DistanceHaversine(p1, p2); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
	p1 = orb.Point{0.5, 30}
	p2 = orb.Point{-0.5, 30}

	dHav := DistanceHaversine(p1, p2)

	p1 = orb.Point{179.5, 30}
	p2 = orb.Point{-179.5, 30}

	if d := DistanceHaversine(p1, p2); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestBearing(t *testing.T) {
	p1 := orb.Point{0, 0}
	p2 := orb.Point{0, 1}

	if d := Bearing(p1, p2); d != 0 {
		t.Errorf("expected 0, got %f", d)
	}

	if d := Bearing(p2, p1); d != 180 {
		t.Errorf("expected 180, got %f", d)
	}

	p1 = orb.Point{0, 0}
	p2 = orb.Point{1, 0}

	if d := Bearing(p1, p2); d != 90 {
		t.Errorf("expected 90, got %f", d)
	}

	if d := Bearing(p2, p1); d != -90 {
		t.Errorf("expected -90, got %f", d)
	}

	p1 = orb.Point{-1.8444, 53.1506}
	p2 = orb.Point{0.1406, 52.2047}

	if d := Bearing(p1, p2); math.Abs(127.373351-d) > epsilon {
		t.Errorf("point, bearingTo got %f", d)
	}
}

func TestMidpoint(t *testing.T) {
	answer := orb.Point{-0.841153, 52.68179432}
	m := Midpoint(orb.Point{-1.8444, 53.1506}, orb.Point{0.1406, 52.2047})

	if d := Distance(m, answer); d > 1 {
		t.Errorf("expected %v, got %v", answer, m)
	}
}

func TestIntermediatePoint(t *testing.T) {
	p1 := orb.Point{-77.035382, 38.898269}
	p2 := orb.Point{-77.011250, 38.889789}

	data := [3][3]float64{
		{0.25, -77.02934845975295, 38.896149465676906},
		{0.50, -77.02331527966795, 38.894029620874510},
		{0.75, -77.01728245974898, 38.891909465634875},
	}
	for _, d := range data {
		f := d[0]
		answer := orb.Point{d[1], d[2]}

		p := IntermediatePoint(p1, p2, f)
		if d := Distance(p, answer); d > 1 {
			t.Errorf("expected %v for f = %v; got %v", answer, f, p)
		}
	}
}
