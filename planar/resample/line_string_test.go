package resample

import (
	"testing"

	"github.com/paulmach/orb/planar"
)

func TestLineStringResample(t *testing.T) {
	ls := planar.NewLineString()
	Resample(ls, 10) // should not panic

	ls = append(ls, planar.NewPoint(0, 0))
	Resample(ls, 10) // should not panic

	ls = append(ls, planar.NewPoint(1.5, 1.5))
	ls = append(ls, planar.NewPoint(2, 2))

	// resample to 0?
	result := Resample(ls, 0)
	if len(result) != 0 {
		t.Error("down to zero should be empty line")
	}

	// resample to 1
	result = Resample(ls, 1)
	answer := append(planar.NewLineString(), planar.NewPoint(0, 0))

	if !result.Equal(answer) {
		t.Error("down to 1 should be first point")
	}

	result = Resample(ls, 2)
	answer = append(planar.NewLineString(),
		planar.NewPoint(0, 0), planar.NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("resample downsampling")
	}

	result = Resample(ls, 5)
	answer = append(planar.NewLineString(),
		planar.NewPoint(0, 0), planar.NewPoint(0.5, 0.5),
		planar.NewPoint(1, 1), planar.NewPoint(1.5, 1.5),
		planar.NewPoint(2, 2))
	if !result.Equal(answer) {
		t.Error("resample upsampling")
		t.Log(result)
		t.Log(answer)
	}

	// round off error case, triggered on my laptop
	p1 := append(planar.NewLineString(),
		planar.NewPoint(-88.145243, 42.321059),
		planar.NewPoint(-88.145232, 42.325902))
	p1 = Resample(p1, 109)
	if len(p1) != 109 {
		t.Errorf("incorrect length: %v != 109", len(p1))
	}

	// duplicate points
	ls = append(planar.NewLineString(),
		planar.NewPoint(1, 0),
		planar.NewPoint(1, 0),
		planar.NewPoint(1, 0),
	)

	ls = Resample(ls, 10)
	if l := len(ls); l != 10 {
		t.Errorf("length incorrect: %d != 10", l)
	}

	expected := planar.NewPoint(1, 0)
	for i := 0; i < len(ls); i++ {
		if !ls[i].Equal(expected) {
			t.Errorf("incorrect point: %v != %v", ls[i], expected)
		}
	}
}

func TestLineStringResampleWithInterval(t *testing.T) {
	ls := append(planar.NewLineString(),
		planar.NewPoint(0, 0),
		planar.NewPoint(0, 10),
	)

	ls = ToInterval(ls, 5.0)
	if l := len(ls); l != 3 {
		t.Errorf("incorrect length: %v != 3", l)
	}

	expected := planar.NewPoint(0, 5.0)
	if v := ls[1]; !v.Equal(expected) {
		t.Errorf("incorrect point: %v != %v", v, expected)
	}
}

func TestLineStringResampleEdgeCases(t *testing.T) {
	ls := append(planar.NewLineString(),
		planar.NewPoint(0, 0),
	)

	_, ret := resampleEdgeCases(ls, 10)
	if !ret {
		t.Errorf("should return true")
	}

	// duplicate points
	ls = append(ls, planar.NewPoint(0, 0))
	ls, ret = resampleEdgeCases(ls, 10)
	if !ret {
		t.Errorf("should return true")
	}

	if l := len(ls); l != 10 {
		t.Errorf("should reset to suggested points: %v != 10", l)
	}

	ls, _ = resampleEdgeCases(ls, 5)
	if l := len(ls); l != 5 {
		t.Errorf("should shorten if necessary: %v != 5", l)
	}
}
