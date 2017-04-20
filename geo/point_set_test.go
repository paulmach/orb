package geo

import (
	"math"
	"testing"
)

func TestPointSetCentroid(t *testing.T) {
	ps := NewPointSet()
	ps = append(ps,
		Point{-188.1298828125, -33.97980872872456},
		Point{-186.1083984375, -38.54816542304658},
		Point{-194.8974609375, -46.10370875598026},
		Point{-192.1728515625, -47.8721439688873},
		Point{-179.7802734375, -37.30027528134431},
	)

	centroid := ps.Centroid()

	// NOTE: input of longitude is outside of the -180:180 ranage but output is within.
	expectedCenter := Point{172.08523311057562, -40.87523942007359}
	if !centroid.Equal(expectedCenter) {
		t.Errorf("should find centroid correctly, got %v", centroid)
	}
}

func TestPointSetDistanceFrom(t *testing.T) {
	ps := append(NewPointSet(),
		Point{-122.42558918, 37.76159786},
		Point{-122.40206146, 37.77962363},
		Point{-122.41486043, 37.78138826},
	)

	fromPoint := Point{-122.41941550000001, 37.7749295}

	if distance, _ := ps.DistanceFrom(fromPoint); math.Floor(distance) != 823 {
		t.Errorf("geo distance incorrect, got %v", distance)
	}

	if _, index := ps.DistanceFrom(fromPoint); index != 2 {
		t.Errorf("incorrect closest index, got %v", index)
	}
}

func TestNewPointSet(t *testing.T) {
	ps := append(NewPointSet(),
		Point{-122.42558918, 37.76159786},
		Point{-122.41486043, 37.78138826},
		Point{-122.40206146, 37.77962363},
	)

	if len(ps) != 3 {
		t.Errorf("should find correct length of new point set %v", len(ps))
	}
}

func TestNewPointSetPreallocate(t *testing.T) {
	ps := append(NewPointSet(),
		Point{-122.42558918, 37.76159786},
		Point{-122.41486043, 37.78138826},
		Point{-122.40206146, 37.77962363},
	)

	if l := len(ps); l != 3 {
		t.Errorf("should find correct length of new point set %v", l)
	}

	if !ps[0].Equal(Point{-122.42558918, 37.76159786}) {
		t.Errorf("should find correct first point of new point set %v", ps[0])
	}

	if !ps[2].Equal(Point{-122.40206146, 37.77962363}) {
		t.Errorf("should find correct first point of new point set %v", ps[2])
	}
}

func TestPathRect(t *testing.T) {
	ps := append(NewPointSet(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
		NewPoint(1, 8),
	)

	answer := NewRect(-1, 1, 0, 10)
	if b := ps.Bound(); !b.Equal(answer) {
		t.Errorf("bound, %v != %v", b, answer)
	}

	ps = NewPointSet()
	if !ps.Bound().IsZero() {
		t.Error("expect empty point set to have zero bounds")
	}
}

func TestPointSetEquals(t *testing.T) {
	p1 := append(NewPointSet(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
	)

	p2 := append(NewPointSet(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
	)

	if !p1.Equal(p2) {
		t.Error("equals paths should be equal")
	}

	p2[1] = NewPoint(1, 0)
	if p1.Equal(p2) {
		t.Error("equals paths should not be equal")
	}

	p1[1] = NewPoint(1, 0)
	p1 = append(p1, NewPoint(0, 0))
	if p2.Equal(p1) {
		t.Error("equals paths should not be equal")
	}
}

func TestPointSetClone(t *testing.T) {
	p1 := append(NewPointSet(),
		NewPoint(0, 0),
		NewPoint(0.5, .2),
		NewPoint(1, 0),
	)

	p2 := p1.Clone()
	p2 = append(p2, NewPoint(0, 0))
	if len(p1) == len(p2) {
		t.Errorf("clone length %d == %d", len(p1), len(p2))
	}

	if p2.Equal(p1) {
		t.Error("clone paths should be equal")
	}
}

func TestPointSetToGeoJSON(t *testing.T) {
	p := append(NewPointSet(),
		NewPoint(1, 2),
	)

	f := p.GeoJSON()
	if !f.Geometry.IsMultiPoint() {
		t.Errorf("pointset, should be linestring geometry")
	}
}

func TestPointSetWKT(t *testing.T) {
	ps := NewPointSet()

	answer := "EMPTY"
	if s := ps.WKT(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}

	ps = append(ps, NewPoint(1, 2))
	answer = "MULTIPOINT(1 2)"
	if s := ps.WKT(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}

	ps = append(ps, NewPoint(3, 4))
	answer = "MULTIPOINT(1 2,3 4)"
	if s := ps.WKT(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}
}

func TestPointSetString(t *testing.T) {
	ps := NewPointSet()

	answer := "EMPTY"
	if s := ps.String(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}

	ps = append(ps, NewPoint(1, 2))
	answer = "MULTIPOINT(1 2)"
	if s := ps.String(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}

	ps = append(ps, NewPoint(3, 4))
	answer = "MULTIPOINT(1 2,3 4)"
	if s := ps.String(); s != answer {
		t.Errorf("pointset, string expected %s, got %s", answer, s)
	}
}
