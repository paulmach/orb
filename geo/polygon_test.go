package geo

import "testing"

func TestPolygonWKT(t *testing.T) {
	r1 := append(NewLineString(),
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

	r2 := append(NewLineString(),
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
