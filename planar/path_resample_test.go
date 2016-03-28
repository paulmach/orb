package planar

import "testing"

func TestPathResample(t *testing.T) {
	p := NewPath()
	p.Resample(10) // should not panic

	p = append(p, NewPoint(0, 0))
	p.Resample(10) // should not panic

	p = append(p, NewPoint(1.5, 1.5))
	p = append(p, NewPoint(2, 2))

	// resample to 0?
	result := p.Clone().Resample(0)
	if len(result) != 0 {
		t.Error("path, resample down to zero should be empty line")
	}

	// resample to 1
	result = p.Clone().Resample(1)
	answer := NewPath()
	answer = append(answer, NewPoint(0, 0))

	if !result.Equal(answer) {
		t.Error("path, resample down to 1 should be first point")
	}

	result = p.Clone().Resample(2)
	answer = append(NewPath(), NewPoint(0, 0), NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("path, resample downsampling")
	}

	result = p.Clone().Resample(5)
	answer = NewPath()
	answer = append(answer, NewPoint(0, 0), NewPoint(0.5, 0.5))
	answer = append(answer, NewPoint(1, 1), NewPoint(1.5, 1.5), NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("path, resample upsampling")
		t.Error(result)
		t.Error(answer)
	}

	// round off error case, triggered on my laptop
	p1 := append(NewPath(), NewPoint(-88.145243, 42.321059), NewPoint(-88.145232, 42.325902))
	p1 = p1.Resample(109)
	if len(p1) != 109 {
		t.Errorf("path, resample incorrect length, expected 109, got %d", len(p1))
	}

	// duplicate points
	p = append(NewPath(),
		NewPoint(1, 0),
		NewPoint(1, 0),
		NewPoint(1, 0),
	)

	p = p.Resample(10)
	if l := len(p); l != 10 {
		t.Errorf("path, resample length incorrect, got %d", l)
	}

	for i := 0; i < len(p); i++ {
		if !p[i].Equal(NewPoint(1, 0)) {
			t.Errorf("path, resample not correct point, got %v", p[i])
		}
	}
}

func TestPathResampleWithInterval(t *testing.T) {
	p := append(NewPath(),
		NewPoint(0, 0),
		NewPoint(0, 10),
	)

	p = p.ResampleWithInterval(5.0)
	if l := len(p); l != 3 {
		t.Errorf("incorrect resample, got %v", l)
	}

	if v := p[1]; !v.Equal(NewPoint(0, 5.0)) {
		t.Errorf("incorrect point, got %v", v)
	}
}

func TestPathResampleEdgeCases(t *testing.T) {
	p := append(NewPath(),
		NewPoint(0, 0),
	)

	_, ret := p.resampleEdgeCases(10)
	if !ret {
		t.Errorf("should return true")
	}

	// duplicate points
	p = append(p, NewPoint(0, 0))
	p, ret = p.resampleEdgeCases(10)
	if !ret {
		t.Errorf("should return true")
	}

	if l := len(p); l != 10 {
		t.Errorf("should reset to suggested points, got %v", l)
	}

	p, _ = p.resampleEdgeCases(5)
	if l := len(p); l != 5 {
		t.Errorf("should shorten if necessary, got %v", l)
	}
}
