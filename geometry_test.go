package orb

import "testing"

func TestCollectionBound(t *testing.T) {
	// from the empty Point we get the zero bound.
	expected := Bound{}

	b2 := Collection(AllGeometries).Bound()
	if !b2.Equal(expected) {
		t.Errorf("wrong bound: %v != %v", b2, expected)
	}
}
