package orb

import (
	"testing"
)

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

func TestPathBound(t *testing.T) {
	mp := append(NewMultiPoint(),
		NewPoint(0.5, .2),
		NewPoint(-1, 0),
		NewPoint(1, 10),
		NewPoint(1, 8),
	)

	expected := Bound{Min: Point{-1, 0}, Max: Point{1, 10}}
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
