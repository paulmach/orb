package geo

import "testing"

func TestMultiPolygonWKT(t *testing.T) {
	p1 := Polygon{append(NewRing(),
		NewPoint(0, 0),
		NewPoint(1, 0),
		NewPoint(1, 1),
		NewPoint(0, 1),
		NewPoint(0, 0),
	)}

	mp := MultiPolygon{p1}
	expected := "MULTIPOLYGON(((0 0,1 0,1 1,0 1,0 0)))"
	if w := mp.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

	p2 := Polygon{append(NewRing(),
		NewPoint(0.4, 0.4),
		NewPoint(0.6, 0.4),
		NewPoint(0.6, 0.6),
		NewPoint(0.4, 0.6),
		NewPoint(0.4, 0.4),
	), append(NewRing(),
		NewPoint(0, 0),
		NewPoint(1, 1),
	)}

	mp = MultiPolygon{p1, p2}
	expected = "MULTIPOLYGON(((0 0,1 0,1 1,0 1,0 0)),((0.4 0.4,0.6 0.4,0.6 0.6,0.4 0.6,0.4 0.4),(0 0,1 1)))"
	if w := mp.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

}
