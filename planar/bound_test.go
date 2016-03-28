package planar

import "testing"

func TestBoundPad(t *testing.T) {
	var bound, tester Bound

	bound = NewBound(0, 1, 2, 3)
	tester = NewBound(-0.5, 1.5, 1.5, 3.5)
	if bound = bound.Pad(0.5); !bound.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, bound)
	}

	bound = NewBound(0, 1, 2, 3)
	tester = NewBound(0.1, 0.9, 2.1, 2.9)
	if bound = bound.Pad(-0.1); !bound.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, bound)
	}
}
