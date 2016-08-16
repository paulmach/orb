package geo

import (
	"math"
	"strings"
	"testing"

	"github.com/paulmach/orb/internal/mercator"
)

var citiesGeoHash = [][3]interface{}{
	{57.09700, 9.85000, "u4phb4hw"},
	{49.03000, -122.32000, "c29nbt9k3q"},
	{39.23500, -76.17490, "dqcz4we0k"},
	{-34.7666, 138.53670, "r1fd0qzmg"},
}

func TestNewPoint(t *testing.T) {
	p := NewPoint(1, 2)
	if p.Lon() != 1 {
		t.Errorf("point, expected 1, got %f", p.Lon())
	}

	if p.Lat() != 2 {
		t.Errorf("point, expected 2, got %f", p.Lat())
	}
}

func TestPointQuadkey(t *testing.T) {
	p := Point{
		-87.65005229999997,
		41.850033,
	}

	if k := p.Quadkey(15); k != 212521785 {
		t.Errorf("point quadkey, incorrect got %d", k)
	}

	// default level
	level := 30
	for _, city := range mercator.Cities {
		p := Point{
			city[1],
			city[0],
		}
		key := p.Quadkey(level)

		p = NewPointFromQuadkey(key, level)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("point quadkey, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > epsilon {
			t.Errorf("point quadkey, longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}

func TestPointQuadkeyString(t *testing.T) {
	p := Point{
		-87.65005229999997,
		41.850033,
	}

	if k := p.QuadkeyString(15); k != "030222231030321" {
		t.Errorf("point quadkey string, incorrect got %s", k)
	}

	// default level
	level := 30
	for _, city := range mercator.Cities {
		p := Point{
			city[1],
			city[0],
		}

		key := p.QuadkeyString(level)

		p = NewPointFromQuadkeyString(key)
		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("point quadkey, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > epsilon {
			t.Errorf("point quadkey, longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}

func TestNewPointFromGeoHash(t *testing.T) {
	for _, c := range citiesGeoHash {
		p := NewPointFromGeoHash(c[2].(string))
		if d := p.DistanceFrom(NewPoint(c[1].(float64), c[0].(float64))); d > 10 {
			t.Errorf("point, new from geohash expected distance %f", d)
		}
	}
}

func TestNewPointFromGeoHashInt64(t *testing.T) {
	for _, c := range citiesGeoHash {
		var hash int64
		for _, r := range c[2].(string) {
			hash <<= 5
			hash |= int64(strings.Index("0123456789bcdefghjkmnpqrstuvwxyz", string(r)))
		}

		p := NewPointFromGeoHashInt64(hash, 5*len(c[2].(string)))
		if d := p.DistanceFrom(NewPoint(c[1].(float64), c[0].(float64))); d > 10 {
			t.Errorf("point, new from geohash expected distance %f", d)
		}
	}
}
func TestPointDistanceFrom(t *testing.T) {
	p1 := NewPoint(-1.8444, 53.1506)
	p2 := NewPoint(0.1406, 52.2047)

	if d := p1.DistanceFrom(p2, true); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	if d := p1.DistanceFrom(p2, false); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	p1 = NewPoint(0.5, 30)
	p2 = NewPoint(-0.5, 30)

	dFast := p1.DistanceFrom(p2, false)
	dHav := p1.DistanceFrom(p2, true)

	p1 = NewPoint(179.5, 30)
	p2 = NewPoint(-179.5, 30)

	if d := p1.DistanceFrom(p2, false); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	if d := p1.DistanceFrom(p2, true); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}
}

func TestPointBearingTo(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(0, 1)

	if d := p1.BearingTo(p2); d != 0 {
		t.Errorf("point, bearingTo expected 0, got %f", d)
	}

	if d := p2.BearingTo(p1); d != 180 {
		t.Errorf("point, bearingTo expected 180, got %f", d)
	}

	p1 = NewPoint(0, 0)
	p2 = NewPoint(1, 0)

	if d := p1.BearingTo(p2); d != 90 {
		t.Errorf("point, bearingTo expected 90, got %f", d)
	}

	if d := p2.BearingTo(p1); d != -90 {
		t.Errorf("point, bearingTo expected -90, got %f", d)
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
		t.Errorf("line, midpoint expected %v, got %v", answer, m)
	}
}

func TestPointGeoHash(t *testing.T) {
	for _, c := range citiesGeoHash {
		hash := NewPoint(c[1].(float64), c[0].(float64)).GeoHash(12)
		if !strings.HasPrefix(hash, c[2].(string)) {
			t.Errorf("point, geohash expected %s, got %s", c[2].(string), hash)
		}
	}

	for _, c := range citiesGeoHash {
		hash := NewPoint(c[1].(float64), c[0].(float64)).GeoHash(len(c[2].(string)))
		if hash != c[2].(string) {
			t.Errorf("point, geohash expected %s, got %s", c[2].(string), hash)
		}
	}
}

func TestPointEqual(t *testing.T) {
	p1 := NewPoint(1, 0)
	p2 := NewPoint(1, 0)

	p3 := NewPoint(2, 3)
	p4 := NewPoint(2, 4)

	if !p1.Equal(p2) {
		t.Errorf("point, equals expect %v == %v", p1, p2)
	}

	if p2.Equal(p3) {
		t.Errorf("point, equals expect %v != %v", p2, p3)
	}

	if p3.Equal(p4) {
		t.Errorf("point, equals expect %v != %v", p3, p4)
	}
}

func TestPointGeoJSON(t *testing.T) {
	p := NewPoint(1, 2.5)

	f := p.GeoJSON()
	if !f.Geometry.IsPoint() {
		t.Errorf("point, should be point geometry")
	}
}

func TestPointWKT(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.WKT(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.String(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}
