package planar

import "testing"

func TestMultiPointCentroid(t *testing.T) {
	mp := NewMultiPoint()
	mp = append(mp,
		Point{0, 0},
		Point{1, 1.5},
		Point{2, 0},
	)

	centroid := mp.Centroid()
	expected := Point{1, 0.5}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}
}

func TestMultiPointDistanceFrom(t *testing.T) {
	mp := append(NewMultiPoint(),
		Point{0, 0},
		Point{1, 1},
		Point{2, 2},
	)

	fromPoint := Point{3, 2}

	if distance := mp.DistanceFrom(fromPoint); distance != 1 {
		t.Errorf("distance incorrect: %v != %v", distance, 1)
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
		t.Errorf("should incorrectpoint of new multi point: %v", mp[0])
	}

	if !mp[2].Equal(Point{-122.40206146, 37.77962363}) {
		t.Errorf("incorrect first point of new multi point: %v", mp[2])
	}
}

func TestMultiPointBound(t *testing.T) {
	mp := append(NewMultiPoint(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
		NewPoint(1, 8),
	)

	answer := NewRect(-1, 1, 0, 10)
	if b := mp.Bound(); !b.Equal(answer) {
		t.Errorf("bound, %v != %v", b, answer)
	}

	mp = NewMultiPoint()
	if !mp.Bound().IsZero() {
		t.Error("expect empty multi point to have zero bound")
	}
}

func TestMultiPointEqual(t *testing.T) {
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
		t.Error("sets should be equal")
	}

	p2[1] = NewPoint(1, 0)
	if p1.Equal(p2) {
		t.Error("sets should not be equal")
	}

	p1[1] = NewPoint(1, 0)
	p1 = append(p1, NewPoint(0, 0))
	if p2.Equal(p1) {
		t.Error("sets should not be equal")
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
		t.Error("clone sets should be equal")
	}
}

func TestMultiPointWKT(t *testing.T) {
	mp := NewMultiPoint()

	answer := "EMPTY"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(1, 2))
	answer = "MULTIPOINT(1 2)"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	mp = append(mp, NewPoint(3, 4))
	answer = "MULTIPOINT(1 2,3 4)"
	if s := mp.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}
}
