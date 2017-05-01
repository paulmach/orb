package planar

import "testing"

func TestMultiLineStringWKT(t *testing.T) {
	ls1 := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(1, 0),
		NewPoint(1, 1),
		NewPoint(0, 1),
		NewPoint(0, 0),
	)

	mls := MultiLineString{ls1}
	expected := "MULTILINESTRING((0 0,1 0,1 1,0 1,0 0))"
	if w := mls.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

	ls2 := append(NewLineString(),
		NewPoint(0.4, 0.4),
		NewPoint(0.6, 0.4),
		NewPoint(0.6, 0.6),
		NewPoint(0.4, 0.6),
		NewPoint(0.4, 0.4),
	)

	mls = MultiLineString{ls1, ls2}
	expected = "MULTILINESTRING((0 0,1 0,1 1,0 1,0 0),(0.4 0.4,0.6 0.4,0.6 0.6,0.4 0.6,0.4 0.4))"
	if w := mls.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

}
