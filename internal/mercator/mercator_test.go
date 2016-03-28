package mercator

import (
	"math"
	"testing"
)

func TestScalarMercator(t *testing.T) {

	x, y := ScalarProject(0, 0, 31)
	lat, lng := ScalarInverse(x, y, 31)

	if lat != 0.0 {
		t.Errorf("Scalar Mercator, latitude should be 0: %f", lat)
	}

	if lng != 0.0 {
		t.Errorf("Scalar Mercator, longitude should be 0: %f", lng)
	}

	// specific case
	if x, y := ScalarProject(-87.65005229999997, 41.850033, 20); x != 268988 || y != 389836 {
		t.Errorf("Scalar Mercator, projection incorrect, got %d %d", x, y)
	}

	if x, y := ScalarProject(-87.65005229999997, 41.850033, 28); x != 68861112 || y != 99798110 {
		t.Errorf("Scalar Mercator, projection incorrect, got %d %d", x, y)
	}

	// default level
	for _, city := range Cities {
		x, y := ScalarProject(city[1], city[0], 31)
		lng, lat = ScalarInverse(x, y, 31)

		if math.Abs(lat-city[0]) > Epsilon {
			t.Errorf("Scalar Mercator, latitude miss match: %f != %f", lat, city[0])
		}

		if math.Abs(lng-city[1]) > Epsilon {
			t.Errorf("Scalar Mercator, longitude miss match: %f != %f", lng, city[1])
		}
	}

	// test polar regions
	if _, y := ScalarProject(0, 89.9, 31); y != (1<<31)-1 {
		t.Errorf("Scalar Mercator, top of the world error, got %d", y)
	}

	if _, y := ScalarProject(0, -89.9, 31); y != 0 {
		t.Errorf("Scalar Mercator, bottom of the world error, got %d", y)
	}
}
