package planar

import (
	"math"
	"testing"
)

func TestPolygonDistanceFrom(t *testing.T) {
	r1 := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(3, 0),
		NewPoint(3, 3),
		NewPoint(0, 3),
		NewPoint(0, 0),
	)

	r2 := append(NewRing(),
		NewPoint(1, 1),
		NewPoint(2, 1),
		NewPoint(2, 2),
		NewPoint(1, 2),
		NewPoint(1, 1),
	)

	poly := append(NewPolygon(), r1, r2)

	cases := []struct {
		name   string
		point  Point
		result float64
	}{
		{
			name:   "outside",
			point:  NewPoint(-1, 2),
			result: 1,
		},
		{
			name:   "inside",
			point:  NewPoint(0.4, 2),
			result: 0,
		},
		{
			name:   "in hole",
			point:  NewPoint(1.3, 1.4),
			result: 0.3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if d := poly.DistanceFrom(tc.point); math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestPolygonCentroidArea(t *testing.T) {
	r1 := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(4, 0),
		NewPoint(4, 3),
		NewPoint(0, 3),
		NewPoint(0, 0),
	).Reverse()

	r2 := append(NewRing(),
		NewPoint(2, 1),
		NewPoint(3, 1),
		NewPoint(3, 2),
		NewPoint(2, 2),
		NewPoint(2, 1),
	)

	poly := Polygon{r1, r2}

	centroid, area := poly.CentroidArea()
	if !centroid.Equal(NewPoint(21.5/11.0, 1.5)) {
		t.Errorf("%v", 21.5/11.0)
		t.Errorf("incorrect centroid: %v", centroid)
	}

	if area != 11 {
		t.Errorf("incorrect area: %v != 11", area)
	}
}

func TestPolygonArea(t *testing.T) {
	r1 := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(3, 0),
		NewPoint(3, 3),
		NewPoint(0, 3),
		NewPoint(0, 0),
	).Reverse()

	r2 := append(NewRing(),
		NewPoint(1, 1),
		NewPoint(2, 1),
		NewPoint(2, 2),
		NewPoint(1, 2),
		NewPoint(1, 1),
	)

	polygon := append(NewPolygon(), r1, r2)

	expected := 8.0
	if a := polygon.Area(); a != expected {
		t.Errorf("incorrect area: %v != %v", a, expected)
	}
}

func TestPolygonWKT(t *testing.T) {
	r1 := append(NewRing(),
		NewPoint(0, 0),
		NewPoint(1, 0),
		NewPoint(1, 1),
		NewPoint(0, 1),
		NewPoint(0, 0),
	)

	poly := Polygon{r1}
	expected := "POLYGON((0 0,1 0,1 1,0 1,0 0))"
	if w := poly.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

	r2 := append(NewRing(),
		NewPoint(0.4, 0.4),
		NewPoint(0.6, 0.4),
		NewPoint(0.6, 0.6),
		NewPoint(0.4, 0.6),
		NewPoint(0.4, 0.4),
	)

	poly = Polygon{r1, r2}
	expected = "POLYGON((0 0,1 0,1 1,0 1,0 0),(0.4 0.4,0.6 0.4,0.6 0.6,0.4 0.6,0.4 0.4))"
	if w := poly.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}
}
