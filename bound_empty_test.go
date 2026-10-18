package orb

import "testing"

func TestEmptyBoundOperations(t *testing.T) {
	empty := MultiPoint{}.Bound()
	point := Point{10, 20}
	if got := empty.Extend(point); !got.Equal(point.Bound()) {
		t.Errorf("empty.Extend(%v) = %v, want %v", point, got, point.Bound())
	}

	other := Bound{Min: Point{-10, -20}, Max: Point{10, 20}}
	if empty.Intersects(other) || other.Intersects(empty) {
		t.Error("an empty bound must not intersect another bound")
	}
	if got := empty.Union(point.Bound()); !got.Equal(point.Bound()) {
		t.Errorf("empty union = %v, want %v", got, point.Bound())
	}
}
