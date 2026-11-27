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
	cases := []struct {
		name     string
		point    orb.Point
		bearing  float64
		distance float64
		expected orb.Point
	}{
		{
			name:     "simple",
			point:    orb.Point{-1.8444, 53.1506},
			bearing:  127.373,
			distance: 85194.89,
			expected: orb.Point{-0.841153, 52.68179432},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := PointAtBearingAndDistance(tc.point, tc.bearing, tc.distance)

			if d := DistanceHaversine(actual, tc.expected); d > 1 {
				t.Errorf("expected %v, got %v (%vm away)", tc.expected, actual, d)
			}
		})
	}

	t.Run("midpoint", func(t *testing.T) {
		a := orb.Point{-1.8444, 53.1506}
		b := orb.Point{0.1406, 52.2047}
		bearing := Bearing(a, b)
		distance := DistanceHaversine(a, b)

		p1 := PointAtBearingAndDistance(a, bearing, distance/2)
		p2 := Midpoint(a, b)

		if d := DistanceHaversine(p1, p2); d > epsilon {
			t.Errorf("expected %v to be within %vm of %v", p1, epsilon, p2)
		}
	})
}

func TestMidpointAntimeridian(t *testing.T) {
	cases := []struct {
		name string
		a, b orb.Point
	}{
		{"180_to_neg90", orb.Point{180, 0}, orb.Point{-90, 0}},
		{"179.9_45_to_neg90_0", orb.Point{179.9, 45}, orb.Point{-90, 0}},
		{"neg179_to_179", orb.Point{-179, 10}, orb.Point{179, 10}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := Midpoint(tc.a, tc.b)
			if m[0] < -180 || m[0] > 180 {
				t.Fatalf("midpoint longitude out of [-180, 180]: got %v", m)
			}
			da := DistanceHaversine(tc.a, m)
			db := DistanceHaversine(m, tc.b)
			if math.Abs(da-db) > 1.0 {
				t.Errorf("midpoint %v not equidistant: d(a,m)=%f, d(m,b)=%f", m, da, db)
			}
		})
	}
}

func TestPointAtBearingAndDistanceAntimeridian(t *testing.T) {
	east := PointAtBearingAndDistance(orb.Point{179.5, 0}, 90, 100000)
	if east[0] < -180 || east[0] > 180 {
		t.Fatalf("eastbound crossing produced unnormalized lon: %v", east)
	}
	if east[0] >= -179 || east[0] < -180 {
		t.Errorf("expected eastbound lon in [-180, -179), got %v", east)
	}

	west := PointAtBearingAndDistance(orb.Point{-179.5, 0}, -90, 100000)
	if west[0] < -180 || west[0] > 180 {
		t.Fatalf("westbound crossing produced unnormalized lon: %v", west)
	}
	if west[0] <= 179 || west[0] > 180 {
		t.Errorf("expected westbound lon in (179, 180], got %v", west)
	}
}

func TestPointAtDistanceAlongLineWithSinglePoint(t *testing.T) {
	cases := []struct {
		name            string
		line            orb.LineString
		distance        float64
		expectedPoint   orb.Point
		expectedBearing float64
	}{
		{
			name: "with single point",
			line: orb.LineString{
				orb.Point{-1.8444, 53.1506},
			},
			distance:        9000,
			expectedPoint:   orb.Point{-1.8444, 53.1506},
			expectedBearing: 0,
		},
		{
			name: "with minimal points",
			line: orb.LineString{
				orb.Point{-1.8444, 53.1506},
				orb.Point{0.1406, 52.2047},
			},
			distance:      85194.89,
			expectedPoint: orb.Point{-0.841153, 52.68179432},
			expectedBearing: Bearing(
				orb.Point{-1.8444, 53.1506},
				orb.Point{0.1406, 52.2047},
			),
		},
		{
			name: "with single point",
			line: orb.LineString{
				orb.Point{-1.8444, 53.1506},
				orb.Point{-0.8411, 52.6817},
				orb.Point{0.1406, 52.2047},
			},
			distance:      90000,
			expectedPoint: orb.Point{-0.78526, 52.65506},
			expectedBearing: Bearing(
				orb.Point{-0.8411, 52.6817},
				orb.Point{0.1406, 52.2047},
			),
		},
		{
			name: "past end of line",
			line: orb.LineString{
				orb.Point{-1.8444, 53.1506},
				orb.Point{-0.8411, 52.6817},
				orb.Point{0.1406, 52.2047},
			},
			distance:      200000,
			expectedPoint: orb.Point{0.1406, 52.2047},
			expectedBearing: Bearing(
				orb.Point{-0.8411, 52.6817},
				orb.Point{0.1406, 52.2047},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualPoint, actualBearing := PointAtDistanceAlongLine(tc.line, tc.distance)

			if d := DistanceHaversine(actualPoint, tc.expectedPoint); d > 1 {
				t.Errorf("point %v != %v", actualPoint, tc.expectedPoint)
			}

			if d := math.Abs(actualBearing - tc.expectedBearing); d > 1 {
				t.Errorf("bearing %v != %v", tc.expectedBearing, actualBearing)
			}
		})
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
