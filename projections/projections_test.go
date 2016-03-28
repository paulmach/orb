package projections

import (
	"math"
	"testing"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
)

func TestMercator(t *testing.T) {
	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		p := Mercator.Project(g)
		g = Mercator.Inverse(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("Mercator, latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lng()-city[1]) > mercator.Epsilon {
			t.Errorf("Mercator, longitude miss match: %f != %f", g.Lng(), city[1])
		}
	}
}

func TestMercatorScaleFactor(t *testing.T) {
	expected := 1.154701
	if f := MercatorScaleFactor(geo.NewPoint(0, 30.0)); math.Abs(expected-f) > mercator.Epsilon {
		t.Errorf("TestMercatorScaleFactor, wrong, expected %f, got %f", expected, f)
	}

	expected = 1.414214
	if f := MercatorScaleFactor(geo.NewPoint(0, 45.0)); math.Abs(expected-f) > mercator.Epsilon {
		t.Errorf("TestMercatorScaleFactor, wrong, expected %f, got %f", expected, f)
	}

	expected = 2.0
	if f := MercatorScaleFactor(geo.NewPoint(0, 60.0)); math.Abs(expected-f) > mercator.Epsilon {
		t.Errorf("TestMercatorScaleFactor, wrong, expected %f, got %f", expected, f)
	}

	expected = 5.758770
	if f := MercatorScaleFactor(geo.NewPoint(0, 80.0)); math.Abs(expected-f) > mercator.Epsilon {
		t.Errorf("TestMercatorScaleFactor, wrong, expected %f, got %f", expected, f)
	}
}

func TestTransverseMercator(t *testing.T) {
	tested := 0

	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		if math.Abs(g.Lng()) > 10 {
			continue
		}

		p := TransverseMercator.Project(g)
		g = TransverseMercator.Inverse(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("TransverseMercator, latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lng()-city[1]) > mercator.Epsilon {
			t.Errorf("TransverseMercator, longitude miss match: %f != %f", g.Lng(), city[1])
		}

		tested++
	}

	if tested == 0 {
		t.Error("TransverseMercator, no points tested")
	}
}

func TestTransverseMercatorScaling(t *testing.T) {

	// points on the 0 longitude should have the same
	// projected distance as geo distance
	g1 := geo.NewPoint(0, 15)
	g2 := geo.NewPoint(0, 30)

	geoDistance := g1.DistanceFrom(g2)

	p1 := TransverseMercator.Project(g1)
	p2 := TransverseMercator.Project(g2)
	projectedDistance := p1.DistanceFrom(p2)

	if math.Abs(geoDistance-projectedDistance) > mercator.Epsilon {
		t.Errorf("TransverseMercatorScaling: values mismatch: %f != %f", geoDistance, projectedDistance)
	}
}

func TestBuildTransverseMercator(t *testing.T) {
	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		offset := math.Floor(g.Lng()/10.0) * 10.0
		projector := BuildTransverseMercator(offset)

		p := projector.Project(g)
		g = projector.Inverse(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lng()-city[1]) > mercator.Epsilon {
			t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", g.Lng(), city[1])
		}
	}

	// test anti-meridian from right
	projector := BuildTransverseMercator(-178.0)

	test := geo.NewPoint(-175.0, 30)

	g := test
	p := projector.Project(g)
	g = projector.Inverse(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lng()-test.Lng()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", g.Lng(), test.Lat())
	}

	test = geo.NewPoint(179.0, 30)

	g = test
	p = projector.Project(g)
	g = projector.Inverse(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lng()-test.Lng()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", g.Lng(), test.Lat())
	}

	// test anti-meridian from left
	projector = BuildTransverseMercator(178.0)

	test = geo.NewPoint(175.0, 30)

	g = test
	p = projector.Project(g)
	g = projector.Inverse(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lng()-test.Lng()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", g.Lng(), test.Lat())
	}

	test = geo.NewPoint(-179.0, 30)

	g = test
	p = projector.Project(g)
	g = projector.Inverse(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lng()-test.Lng()) > mercator.Epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", g.Lng(), test.Lat())
	}
}

func TestScalarMercator(t *testing.T) {

	x, y := ScalarMercator.Project(geo.NewPoint(0, 0))
	g := ScalarMercator.Inverse(x, y)

	if g.Lat() != 0.0 {
		t.Errorf("Scalar Mercator, latitude should be 0: %f", g.Lat())
	}

	if g.Lng() != 0.0 {
		t.Errorf("Scalar Mercator, longitude should be 0: %f", g.Lng())
	}

	// specific case
	if x, y := ScalarMercator.Project(geo.NewPoint(-87.65005229999997, 41.850033), 20); x != 268988 || y != 389836 {
		t.Errorf("Scalar Mercator, projection incorrect, got %d %d", x, y)
	}

	ScalarMercator.Level = 28
	if x, y := ScalarMercator.Project(geo.NewPoint(-87.65005229999997, 41.850033)); x != 68861112 || y != 99798110 {
		t.Errorf("Scalar Mercator, projection incorrect, got %d %d", x, y)
	}

	// default level
	ScalarMercator.Level = 31
	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		x, y := ScalarMercator.Project(g, 31)
		g = ScalarMercator.Inverse(x, y)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("Scalar Mercator, latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lng()-city[1]) > mercator.Epsilon {
			t.Errorf("Scalar Mercator, longitude miss match: %f != %f", g.Lng(), city[1])
		}
	}

	// test polar regions
	if _, y := ScalarMercator.Project(geo.NewPoint(0, 89.9)); y != (1<<ScalarMercator.Level)-1 {
		t.Errorf("Scalar Mercator, top of the world error, got %d", y)
	}

	if _, y := ScalarMercator.Project(geo.NewPoint(0, -89.9)); y != 0 {
		t.Errorf("Scalar Mercator, bottom of the world error, got %d", y)
	}
}
