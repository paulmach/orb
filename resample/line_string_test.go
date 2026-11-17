package resample

import (
	"math"
	"strconv"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

func TestResample(t *testing.T) {
	ls := orb.LineString{}
	Resample(ls, planar.Distance, 10) // should not panic

	ls = append(ls, orb.Point{0, 0})
	Resample(ls, planar.Distance, 10) // should not panic

	ls = append(ls, orb.Point{1.5, 1.5})
	ls = append(ls, orb.Point{2, 2})

	// resample to 0?
	result := Resample(ls, planar.Distance, 0)
	if len(result) != 0 {
		t.Error("down to zero should be empty line")
	}

	// resample to 1
	result = Resample(ls, planar.Distance, 1)
	answer := orb.LineString{{0, 0}}

	if !result.Equal(answer) {
		t.Error("down to 1 should be first point")
	}

	result = Resample(ls, planar.Distance, 2)
	answer = orb.LineString{{0, 0}, {2, 2}}
	if !result.Equal(answer) {
		t.Error("resample downsampling")
	}

	result = Resample(ls, planar.Distance, 5)
	answer = orb.LineString{{0, 0}, {0.5, 0.5}, {1, 1}, {1.5, 1.5}, {2, 2}}
	if !result.Equal(answer) {
		t.Error("resample upsampling")
		t.Log(result)
		t.Log(answer)
	}

	// round off error case, triggered on my laptop
	p1 := orb.LineString{{-88.145243, 42.321059}, {-88.145232, 42.325902}}
	p1 = Resample(p1, planar.Distance, 109)
	if len(p1) != 109 {
		t.Errorf("incorrect length: %v != 109", len(p1))
	}

	// duplicate points
	ls = orb.LineString{{1, 0}, {1, 0}, {1, 0}}
	ls = Resample(ls, planar.Distance, 10)
	if l := len(ls); l != 10 {
		t.Errorf("length incorrect: %d != 10", l)
	}

	expected := orb.Point{1, 0}
	for i := 0; i < len(ls); i++ {
		if !ls[i].Equal(expected) {
			t.Errorf("incorrect point: %v != %v", ls[i], expected)
		}
	}
}

func TestToInterval(t *testing.T) {
	ls := orb.LineString{{0, 0}, {0, 1}, {0, 10}}

	cases := []struct {
		name     string
		distance float64
		line     orb.LineString
		expected orb.LineString
	}{
		{
			name:     "same number of points",
			distance: 5.0,
			expected: orb.LineString{{0, 0}, {0, 5}, {0, 10}},
		},
		{
			name:     "dist less than 0",
			distance: -5.0,
			expected: nil,
		},
		{
			name:     "dist less than 0",
			distance: -5.0,
			expected: nil,
		},
		{
			name:     "return same if short line",
			distance: 5.0,
			line:     orb.LineString{{0, 0}},
			expected: orb.LineString{{0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := ls.Clone()
			if tc.line != nil {
				in = tc.line
			}

			ls := ToInterval(in, planar.Distance, tc.distance)
			if !ls.Equal(tc.expected) {
				t.Errorf("incorrect point: %v != %v", ls, tc.expected)
			}
		})
	}
}

func TestToIntervalEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		line orb.LineString
	}{
		{name: "nil"},
		{name: "empty", line: orb.LineString{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ToInterval(tc.line, planar.Distance, 5)
			if len(got) != 0 || (got == nil) != (tc.line == nil) {
				t.Errorf("expected unchanged empty line: got %#v, want %#v", got, tc.line)
			}
		})
	}
}

func TestLineStringResampleEdgeCases(t *testing.T) {
	ls := orb.LineString{{0, 0}}

	_, ret := resampleEdgeCases(ls, 10)
	if !ret {
		t.Errorf("should return true")
	}

	// duplicate points
	ls = append(ls, orb.Point{0, 0})
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

func TestResampleZeroDistance(t *testing.T) {
	for _, tc := range []struct {
		name string
		line orb.LineString
		df   orb.DistanceFunc
	}{
		{"antimeridian", orb.LineString{{-180, 0}, {180, 0}}, geo.Distance},
		{"reversed", orb.LineString{{180, 0}, {-180, 0}}, geo.Distance},
		{"latitude", orb.LineString{{-180, 45}, {180, 45}}, geo.Distance},
		{"multiple segments", orb.LineString{{-180, 0}, {180, 0}, {-180, 0}}, geo.Distance},
		{"duplicate vertices", orb.LineString{{-180, 0}, {-180, 0}, {180, 0}, {180, 0}}, geo.Distance},
		{"planar underflow", orb.LineString{{0, 0}, {1e-200, 0}}, planar.Distance},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if total, _ := precomputeDistances(tc.line, tc.df); total != 0 {
				t.Fatalf("regression requires zero total distance, got %v", total)
			}
			for _, count := range []int{0, 1, 2, 3, 8} {
				t.Run(strconv.Itoa(count), func(t *testing.T) {
					got := Resample(tc.line.Clone(), tc.df, count)
					if len(got) != count {
						t.Fatalf("got %d points, want %d", len(got), count)
					}
					for i, p := range got {
						want := tc.line[0]
						if count > 1 && i == count-1 {
							want = tc.line[len(tc.line)-1]
						}
						if !p.Equal(want) {
							t.Errorf("point %d = %v, want %v", i, p, want)
						}
						for _, coordinate := range p {
							if math.IsNaN(coordinate) || math.IsInf(coordinate, 0) {
								t.Errorf("point %d is not finite: %v", i, p)
							}
						}
					}
				})
			}
			// Zero total distance implies one sample at any positive interval.
			got := ToInterval(tc.line.Clone(), tc.df, 1)
			if !got.Equal(orb.LineString{tc.line[0]}) {
				t.Errorf("ToInterval = %v, want first point", got)
			}
		})
	}
}

func TestResamplePositiveDistance(t *testing.T) {
	for _, tc := range []struct {
		name string
		line orb.LineString
		df   orb.DistanceFunc
	}{
		{"planar", orb.LineString{{0, 0}, {0, 1}, {0, 4}}, planar.Distance},
		{"duplicate start", orb.LineString{{0, 0}, {0, 0}, {0, 4}}, planar.Distance},
		{"duplicate middle", orb.LineString{{0, 0}, {0, 2}, {0, 2}, {0, 4}}, planar.Distance},
		{"geo", orb.LineString{{0, 0}, {0, 1}, {0, 4}}, geo.Distance},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, count := range []int{1, 2, 3, 5, 9} {
				got := Resample(tc.line.Clone(), tc.df, count)
				if len(got) != count {
					t.Fatalf("got %d points, want %d", len(got), count)
				}
				for i, p := range got {
					want := 0.0
					if count > 1 {
						want = 4 * float64(i) / float64(count-1)
					}
					if p[0] != 0 || math.IsNaN(p[1]) || math.Abs(p[1]-want) > 1e-12 {
						t.Errorf("sample %d of %d = %v, want [0 %v]", i, count, p, want)
					}
				}
			}
		})
	}
}
