package project

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

		p := Mercator.ToPlanar(g)
		g = Mercator.ToGeo(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g.Lon(), city[1])
		}
	}
}

func TestMercatorScaleFactor(t *testing.T) {
	cases := []struct {
		name   string
		point  geo.Point
		factor float64
	}{
		{
			name:   "30 deg",
			point:  geo.NewPoint(0, 30.0),
			factor: 1.154701,
		},
		{
			name:   "45 deg",
			point:  geo.NewPoint(0, 45.0),
			factor: 1.414214,
		},
		{
			name:   "60 deg",
			point:  geo.NewPoint(0, 60.0),
			factor: 2,
		},
		{
			name:   "80 deg",
			point:  geo.NewPoint(0, 80.0),
			factor: 5.758770,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if f := MercatorScaleFactor(tc.point); math.Abs(tc.factor-f) > mercator.Epsilon {
				t.Errorf("incorrect factor: %v != %v", f, tc.factor)
			}
		})
	}
}

func TestTransverseMercator(t *testing.T) {
	tested := 0

	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		if math.Abs(g.Lon()) > 10 {
			continue
		}

		p := TransverseMercator.ToPlanar(g)
		g = TransverseMercator.ToGeo(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g.Lon(), city[1])
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

	p1 := TransverseMercator.ToPlanar(g1)
	p2 := TransverseMercator.ToPlanar(g2)
	projectedDistance := p1.DistanceFrom(p2)

	if math.Abs(geoDistance-projectedDistance) > mercator.Epsilon {
		t.Errorf("incorrect scale: %f != %f", geoDistance, projectedDistance)
	}
}

func TestBuildTransverseMercator(t *testing.T) {
	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		offset := math.Floor(g.Lon()/10.0) * 10.0
		projector := BuildTransverseMercator(offset)

		p := projector.ToPlanar(g)
		g = projector.ToGeo(p)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g.Lon(), city[1])
		}
	}

	// test anti-meridian from right
	projector := BuildTransverseMercator(-178.0)

	test := geo.NewPoint(-175.0, 30)

	g := test
	p := projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lon()-test.Lon()) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g.Lon(), test.Lat())
	}

	test = geo.NewPoint(179.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lon()-test.Lon()) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g.Lon(), test.Lat())
	}

	// test anti-meridian from left
	projector = BuildTransverseMercator(178.0)

	test = geo.NewPoint(175.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lon()-test.Lon()) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g.Lon(), test.Lat())
	}

	test = geo.NewPoint(-179.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g.Lat()-test.Lat()) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g.Lat(), test.Lat())
	}

	if math.Abs(g.Lon()-test.Lon()) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g.Lon(), test.Lat())
	}
}

func TestScalarMercator(t *testing.T) {

	x, y := ScalarMercator.ToPlanar(geo.NewPoint(0, 0))
	g := ScalarMercator.ToGeo(x, y)

	if g.Lat() != 0.0 {
		t.Errorf("latitude should be 0: %f", g.Lat())
	}

	if g.Lon() != 0.0 {
		t.Errorf("longitude should be 0: %f", g.Lon())
	}

	// specific case
	if x, y := ScalarMercator.ToPlanar(geo.NewPoint(-87.65005229999997, 41.850033), 20); x != 268988 || y != 389836 {
		t.Errorf("projection incorrect: %d %d", x, y)
	}

	ScalarMercator.Level = 28
	if x, y := ScalarMercator.ToPlanar(geo.NewPoint(-87.65005229999997, 41.850033)); x != 68861112 || y != 99798110 {
		t.Errorf("projection incorrect: %d %d", x, y)
	}

	// default level
	ScalarMercator.Level = 31
	for _, city := range mercator.Cities {
		g := geo.Point{
			city[1],
			city[0],
		}

		x, y := ScalarMercator.ToPlanar(g, 31)
		g = ScalarMercator.ToGeo(x, y)

		if math.Abs(g.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g.Lat(), city[0])
		}

		if math.Abs(g.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g.Lon(), city[1])
		}
	}

	// test polar regions
	if _, y := ScalarMercator.ToPlanar(geo.NewPoint(0, 89.9)); y != (1<<ScalarMercator.Level)-1 {
		t.Errorf("top of the world error: %d", y)
	}

	if _, y := ScalarMercator.ToPlanar(geo.NewPoint(0, -89.9)); y != 0 {
		t.Errorf("bottom of the world error: %d", y)
	}
}
