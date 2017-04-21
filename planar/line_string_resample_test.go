package planar

import "testing"

func TestLineStringResample(t *testing.T) {
	ls := NewLineString()
	ls.Resample(10) // should not panic

	ls = append(ls, NewPoint(0, 0))
	ls.Resample(10) // should not panic

	ls = append(ls, NewPoint(1.5, 1.5))
	ls = append(ls, NewPoint(2, 2))

	// resample to 0?
	result := ls.Clone().Resample(0)
	if len(result) != 0 {
		t.Error("down to zero should be empty line")
	}

	// resample to 1
	result = ls.Clone().Resample(1)
	answer := NewLineString()
	answer = append(answer, NewPoint(0, 0))

	if !result.Equal(answer) {
		t.Error("down to 1 should be first point")
	}

	result = ls.Clone().Resample(2)
	answer = append(NewLineString(), NewPoint(0, 0), NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("resample downsampling")
	}

	result = ls.Clone().Resample(5)
	answer = NewLineString()
	answer = append(answer, NewPoint(0, 0), NewPoint(0.5, 0.5))
	answer = append(answer, NewPoint(1, 1), NewPoint(1.5, 1.5), NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("resample upsampling")
		t.Log(result)
		t.Log(answer)
	}

	// round off error case, triggered on my laptop
	p1 := append(NewLineString(), NewPoint(-88.145243, 42.321059), NewPoint(-88.145232, 42.325902))
	p1 = p1.Resample(109)
	if len(p1) != 109 {
		t.Errorf("incorrect length: %v != 109", len(p1))
	}

	// duplicate points
	ls = append(NewLineString(),
		NewPoint(1, 0),
		NewPoint(1, 0),
		NewPoint(1, 0),
	)

	ls = ls.Resample(10)
	if l := len(ls); l != 10 {
		t.Errorf("length incorrect: %d != 10", l)
	}

	expected := NewPoint(1, 0)
	for i := 0; i < len(ls); i++ {
		if !ls[i].Equal(expected) {
			t.Errorf("incorrect point: %v != %v", ls[i], expected)
		}
	}
}

func TestLineStringResampleWithInterval(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 10),
	)

	ls = ls.ResampleWithInterval(5.0)
	if l := len(ls); l != 3 {
		t.Errorf("incorrect length: %v != 3", l)
	}

	expected := NewPoint(0, 5.0)
	if v := ls[1]; !v.Equal(expected) {
		t.Errorf("incorrect point: %v != %v", v, expected)
	}
}

func TestLineStringResampleEdgeCases(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
	)

	_, ret := ls.resampleEdgeCases(10)
	if !ret {
		t.Errorf("should return true")
	}

	// duplicate points
	ls = append(ls, NewPoint(0, 0))
	ls, ret = ls.resampleEdgeCases(10)
	if !ret {
		t.Errorf("should return true")
	}

	if l := len(ls); l != 10 {
		t.Errorf("should reset to suggested points: %v != 10", l)
	}

	ls, _ = ls.resampleEdgeCases(5)
	if l := len(ls); l != 5 {
		t.Errorf("should shorten if necessary: %v != 5", l)
	}
}
