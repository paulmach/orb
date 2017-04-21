package geo

import (
	"math"
	"testing"
)

func TestMultiPointCentroid(t *testing.T) {
	mp := NewMultiPoint()
	mp = append(mp,
		Point{-188.1298828125, -33.97980872872456},
		Point{-186.1083984375, -38.54816542304658},
		Point{-194.8974609375, -46.10370875598026},
		Point{-192.1728515625, -47.8721439688873},
		Point{-179.7802734375, -37.30027528134431},
	)

	centroid := mp.Centroid()

	// NOTE: input of longitude is outside of the -180:180 range but output is within.
	expected := Point{172.08523311057562, -40.87523942007359}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}
}

func TestMultiPointDistanceFrom(t *testing.T) {
	mp := append(NewMultiPoint(),
		Point{-122.42558918, 37.76159786},
		Point{-122.40206146, 37.77962363},
		Point{-122.41486043, 37.78138826},
	)

	fromPoint := Point{-122.41941550000001, 37.7749295}

	if distance, _ := mp.DistanceFrom(fromPoint); math.Floor(distance) != 823 {
		t.Errorf("geo distance incorrect: %v", distance)
	}

	if _, index := mp.DistanceFrom(fromPoint); index != 2 {
		t.Errorf("incorrect closest index: %v", index)
	}
}

func TestNewMultiPoint(t *testing.T) {
	mp := append(NewMultiPoint(),
		Point{-122.42558918, 37.76159786},
		Point{-122.41486043, 37.78138826},
		Point{-122.40206146, 37.77962363},
	)

	if len(mp) != 3 {
		t.Errorf("incorrect length of new multi point: %v", len(mp))
	}
}

func TestNewMultiPointPreallocate(t *testing.T) {
	mp := append(NewMultiPoint(),
		Point{-122.42558918, 37.76159786},
		Point{-122.41486043, 37.78138826},
		Point{-122.40206146, 37.77962363},
	)

	if l := len(mp); l != 3 {
		t.Errorf("incorrect length of new multi point: %v", l)
	}

	if !mp[0].Equal(Point{-122.42558918, 37.76159786}) {
		t.Errorf("incorrect first point of new multi point: %v", mp[0])
	}

	if !mp[2].Equal(Point{-122.40206146, 37.77962363}) {
		t.Errorf("incorrect first point of new multi point: %v", mp[2])
	}
}

func TestPathRect(t *testing.T) {
	mp := append(NewMultiPoint(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
		NewPoint(1, 8),
	)

	expected := NewRect(-1, 1, 0, 10)
	if b := mp.Bound(); !b.Equal(expected) {
		t.Errorf("incorrect bound, %v != %v", b, expected)
	}

	mp = NewMultiPoint()
	if !mp.Bound().IsZero() {
		t.Error("expect empty multi point to have zero bounds")
	}
}

func TestMultiPointEquals(t *testing.T) {
	p1 := append(NewMultiPoint(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
	)

	p2 := append(NewMultiPoint(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
	)

	if !p1.Equal(p2) {
		t.Error("should be equal")
	}

	p2[1] = NewPoint(1, 0)
	if p1.Equal(p2) {
		t.Error("should not be equal")
	}

	p1[1] = NewPoint(1, 0)
	p1 = append(p1, NewPoint(0, 0))
	if p2.Equal(p1) {
		t.Error("should not be equal")
	}
}

func TestMultiPointClone(t *testing.T) {
	p1 := append(NewMultiPoint(),
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
		t.Error("clone should be equal")
	}
}

func TestMultiPointToGeoJSON(t *testing.T) {
	p := append(NewMultiPoint(),
		NewPoint(1, 2),
	)

	f := p.GeoJSON()
	if !f.Geometry.IsMultiPoint() {
		t.Errorf("should be linestring geometry")
	}
}

func TestMultiPointWKT(t *testing.T) {
	mp := NewMultiPoint()

	answer := "EMPTY"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(1, 2))
	answer = "MULTIPOINT(1 2)"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(3, 4))
	answer = "MULTIPOINT(1 2,3 4)"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}
}

func TestMultiPointString(t *testing.T) {
	mp := NewMultiPoint()

	answer := "EMPTY"
	if s := mp.String(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(1, 2))
	answer = "MULTIPOINT(1 2)"
	if s := mp.String(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(3, 4))
	answer = "MULTIPOINT(1 2,3 4)"
	if s := mp.String(); s != answer {
		t.Errorf("incorrect string: %v != %v", s, answer)
	}
}
