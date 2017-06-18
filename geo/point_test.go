package geo_test

import (
	"math"
	"testing"

	. "github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/tile"
)

var epsilon = 1e-6

func TestNewPoint(t *testing.T) {
	p := NewPoint(1, 2)
	if p.Lon() != 1 {
		t.Errorf("incorrect lon: %v != 1", p.Lon())
	}

	if p.Lat() != 2 {
		t.Errorf("incorrect lat: %v != 2", p.Lat())
	}
}

func TestPointQuadkey(t *testing.T) {
	p := Point{
		-87.65005229999997,
		41.850033,
	}

	if k := p.Quadkey(15); k != 212521785 {
		t.Errorf("incorrect quadkey: %v != 212521785", k)
	}

	// default level
	level := uint64(30)
	for _, city := range mercator.Cities {
		p := Point{
			city[1],
			city[0],
		}
		key := p.Quadkey(level)

		p = tile.FromQuadkey(key, level).Center()

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > epsilon {
			t.Errorf("longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}

func TestPointDistanceFrom(t *testing.T) {
	p1 := NewPoint(-1.8444, 53.1506)
	p2 := NewPoint(0.1406, 52.2047)

	if d := p1.DistanceFrom(p2, true); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	if d := p1.DistanceFrom(p2, false); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = NewPoint(0.5, 30)
	p2 = NewPoint(-0.5, 30)

	dFast := p1.DistanceFrom(p2, false)
	dHav := p1.DistanceFrom(p2, true)

	p1 = NewPoint(179.5, 30)
	p2 = NewPoint(-179.5, 30)

	if d := p1.DistanceFrom(p2, false); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	if d := p1.DistanceFrom(p2, true); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestPointBearingTo(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(0, 1)

	if d := p1.BearingTo(p2); d != 0 {
		t.Errorf("expected 0, got %f", d)
	}

	if d := p2.BearingTo(p1); d != 180 {
		t.Errorf("expected 180, got %f", d)
	}

	p1 = NewPoint(0, 0)
	p2 = NewPoint(1, 0)

	if d := p1.BearingTo(p2); d != 90 {
		t.Errorf("expected 90, got %f", d)
	}

	if d := p2.BearingTo(p1); d != -90 {
		t.Errorf("expected -90, got %f", d)
	}

	p1 = NewPoint(-1.8444, 53.1506)
	p2 = NewPoint(0.1406, 52.2047)

	if d := p1.BearingTo(p2); math.Abs(127.373351-d) > epsilon {
		t.Errorf("point, bearingTo got %f", d)
	}
}

func TestPointMidpoint(t *testing.T) {
	answer := NewPoint(-0.841153, 52.68179432)
	m := NewPoint(-1.8444, 53.1506).Midpoint(NewPoint(0.1406, 52.2047))

	if d := m.DistanceFrom(answer); d > 1 {
		t.Errorf("expected %v, got %v", answer, m)
	}
}

func TestPointEqual(t *testing.T) {
	p1 := NewPoint(1, 0)
	p2 := NewPoint(1, 0)

	p3 := NewPoint(2, 3)
	p4 := NewPoint(2, 4)

	if !p1.Equal(p2) {
		t.Errorf("expected: %v == %v", p1, p2)
	}

	if p2.Equal(p3) {
		t.Errorf("expected: %v != %v", p2, p3)
	}

	if p3.Equal(p4) {
		t.Errorf("expected: %v != %v", p3, p4)
	}
}

func TestPointWKT(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.WKT(); s != answer {
		t.Errorf("expected %s, got %s", answer, s)
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.String(); s != answer {
		t.Errorf("expected %s, got %s", answer, s)
	}
}
