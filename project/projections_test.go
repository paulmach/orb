package project

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/planar"
)

func TestMercator(t *testing.T) {
	for _, city := range mercator.Cities {
		g := orb.Point{
			city[1],
			city[0],
		}

		p := Mercator.ToPlanar(g)
		g = Mercator.ToGeo(p)

		if math.Abs(g[1]-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g[1], city[0])
		}

		if math.Abs(g[0]-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g[0], city[1])
		}
	}
}

func TestMercatorScaleFactor(t *testing.T) {
	cases := []struct {
		name   string
		point  orb.Point
		factor float64
	}{
		{
			name:   "30 deg",
			point:  orb.NewPoint(0, 30.0),
			factor: 1.154701,
		},
		{
			name:   "45 deg",
			point:  orb.NewPoint(0, 45.0),
			factor: 1.414214,
		},
		{
			name:   "60 deg",
			point:  orb.NewPoint(0, 60.0),
			factor: 2,
		},
		{
			name:   "80 deg",
			point:  orb.NewPoint(0, 80.0),
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
		g := orb.Point{
			city[1],
			city[0],
		}

		if math.Abs(g[0]) > 10 {
			continue
		}

		p := TransverseMercator.ToPlanar(g)
		g = TransverseMercator.ToGeo(p)

		if math.Abs(g[1]-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g[1], city[0])
		}

		if math.Abs(g[0]-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g[0], city[1])
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
	g1 := orb.NewPoint(0, 15)
	g2 := orb.NewPoint(0, 30)

	geoDistance := geo.Distance(g1, g2)

	p1 := TransverseMercator.ToPlanar(g1)
	p2 := TransverseMercator.ToPlanar(g2)
	projectedDistance := planar.Distance(p1, p2)

	if math.Abs(geoDistance-projectedDistance) > mercator.Epsilon {
		t.Errorf("incorrect scale: %f != %f", geoDistance, projectedDistance)
	}
}

func TestBuildTransverseMercator(t *testing.T) {
	for _, city := range mercator.Cities {
		g := orb.Point{
			city[1],
			city[0],
		}

		offset := math.Floor(g[0]/10.0) * 10.0
		projector := BuildTransverseMercator(offset)

		p := projector.ToPlanar(g)
		g = projector.ToGeo(p)

		if math.Abs(g[1]-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", g[1], city[0])
		}

		if math.Abs(g[0]-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", g[0], city[1])
		}
	}

	// test anti-meridian from right
	projector := BuildTransverseMercator(-178.0)

	test := orb.NewPoint(-175.0, 30)

	g := test
	p := projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g[1]-test[1]) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g[1], test[1])
	}

	if math.Abs(g[0]-test[0]) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g[0], test[1])
	}

	test = orb.NewPoint(179.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g[1]-test[1]) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g[1], test[1])
	}

	if math.Abs(g[0]-test[0]) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g[0], test[1])
	}

	// test anti-meridian from left
	projector = BuildTransverseMercator(178.0)

	test = orb.NewPoint(175.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g[1]-test[1]) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g[1], test[1])
	}

	if math.Abs(g[0]-test[0]) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g[0], test[1])
	}

	test = orb.NewPoint(-179.0, 30)

	g = test
	p = projector.ToPlanar(g)
	g = projector.ToGeo(p)

	if math.Abs(g[1]-test[1]) > mercator.Epsilon {
		t.Errorf("latitude miss match: %f != %f", g[1], test[1])
	}

	if math.Abs(g[0]-test[0]) > mercator.Epsilon {
		t.Errorf("longitude miss match: %f != %f", g[0], test[1])
	}
}
