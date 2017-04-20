package planar

import "testing"

func TestRectPad(t *testing.T) {
	var rect, tester Rect

	rect = NewRect(0, 1, 2, 3)
	tester = NewRect(-0.5, 1.5, 1.5, 3.5)
	if rect = rect.Pad(0.5); !rect.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, rect)
	}

	rect = NewRect(0, 1, 2, 3)
	tester = NewRect(0.1, 0.9, 2.1, 2.9)
	if rect = rect.Pad(-0.1); !rect.Equal(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, rect)
	}
}
