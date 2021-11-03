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

func TestPointAtBearingAndDistance(t *testing.T) {
	answer := orb.Point{-0.841153, 52.68179432}
	bearing := 127.373
	distance := 85194.89
	p := PointAtBearingAndDistance(orb.Point{-1.8444, 53.1506}, bearing, distance)

	if d := DistanceHaversine(p, answer); d > 1 {
		t.Errorf("expected %v, got %v (%vm away)", answer, p, d)
	}
}
func TestMidpointAgainstPointAtBearingAndDistance(t *testing.T) {
	a := orb.Point{-1.8444, 53.1506}
	b := orb.Point{0.1406, 52.2047}
	bearing := Bearing(a, b)
	distance := DistanceHaversine(a, b)
	acceptableTolerance := 1e-06 // unit is meters

	p1 := PointAtBearingAndDistance(a, bearing, distance/2)
	p2 := Midpoint(a, b)

	if d := DistanceHaversine(p1, p2); d > acceptableTolerance {
		t.Errorf("expected %v to be within %vm of %v", p1, acceptableTolerance, p2)
	}
}

func TestPointAtDistanceAlongLineWithEmptyLineString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("PointAtDistanceAlongLine did not panic")
		}
	}()

	line := orb.LineString{}
	PointAtDistanceAlongLine(line, 90000)
}

func TestPointAtDistanceAlongLineWithSinglePoint(t *testing.T) {
	expected := orb.Point{-1.8444, 53.1506}
	line := orb.LineString{
		expected,
	}
	actual := PointAtDistanceAlongLine(line, 90000)

	if actual != expected {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestPointAtDistanceAlongLineWithMinimalPoints(t *testing.T) {
	expected := orb.Point{-0.841153, 52.68179432}
	acceptableTolerance := 1.0 // unit is meters
	line := orb.LineString{
		orb.Point{-1.8444, 53.1506},
		orb.Point{0.1406, 52.2047},
	}
	actual := PointAtDistanceAlongLine(line, 85194.89)

	if d := DistanceHaversine(expected, actual); d > acceptableTolerance {
		t.Errorf("expected %v to be within %vm of %v (%vm away)", expected, acceptableTolerance, actual, d)
	}
}

func TestPointAtDistanceAlongLineWithMultiplePoints(t *testing.T) {
	expected := orb.Point{-0.78526, 52.65506}
	acceptableTolerance := 1.0 // unit is meters
	line := orb.LineString{
		orb.Point{-1.8444, 53.1506},
		orb.Point{-0.8411, 52.6817},
		orb.Point{0.1406, 52.2047},
	}
	actual := PointAtDistanceAlongLine(line, 90000)

	if d := DistanceHaversine(expected, actual); d > acceptableTolerance {
		t.Errorf("expected %v to be within %vm of %v (%vm away)", expected, acceptableTolerance, actual, d)
	}
}

func TestPointAtDistanceAlongLinePastEndOfLine(t *testing.T) {
	expected := orb.Point{0.1406, 52.2047}
	line := orb.LineString{
		orb.Point{-1.8444, 53.1506},
		orb.Point{-0.8411, 52.6817},
		expected,
	}
	actual := PointAtDistanceAlongLine(line, 200000)

	if actual != expected {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}
