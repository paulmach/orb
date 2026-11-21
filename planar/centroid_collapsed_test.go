package planar

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func TestCentroidAreaMultiLineStringCollapsedComponents(t *testing.T) {
	for _, test := range []struct {
		name  string
		lines orb.MultiLineString
		want  orb.Point
	}{
		{"mixed", orb.MultiLineString{{{0, 0}, {2, 0}}, {{100, 100}, {100, 100}}}, orb.Point{1, 0}},
		{"collapsed first", orb.MultiLineString{{{100, 100}, {100, 100}}, {{0, 0}, {2, 0}}}, orb.Point{1, 0}},
		{"unequal lengths", orb.MultiLineString{{{0, 0}, {2, 0}}, {{0, 2}, {0, 6}}, {{100, 100}, {100, 100}}}, orb.Point{1.0 / 3, 8.0 / 3}},
		{"single point", orb.MultiLineString{{{100, 100}}, {{0, 0}, {2, 0}}}, orb.Point{1, 0}},
		{"all collapsed", orb.MultiLineString{{{10, 20}, {10, 20}}, {{30, 40}, {30, 40}}}, orb.Point{20, 30}},
		{"empty components", orb.MultiLineString{nil, {{0, 0}, {2, 0}}, {}, {{100, 100}, {100, 100}}}, orb.Point{1, 0}},
		{"ordinary", orb.MultiLineString{{{0, 0}, {2, 0}}, {{0, 2}, {0, 6}}}, orb.Point{1.0 / 3, 8.0 / 3}},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, shift := range []orb.Point{{0, 0}, {16, -8}} {
				lines := make(orb.MultiLineString, len(test.lines))
				for i, line := range test.lines {
					for _, p := range line {
						lines[i] = append(lines[i], orb.Point{p[0] + shift[0], p[1] + shift[1]})
					}
				}
				got, area := CentroidArea(lines)
				want := orb.Point{test.want[0] + shift[0], test.want[1] + shift[1]}
				if math.Abs(got[0]-want[0]) > 1e-12 || math.Abs(got[1]-want[1]) > 1e-12 {
					t.Errorf("shift %v: centroid=%v, want %v", shift, got, want)
				}
				if area != 0 {
					t.Errorf("area=%v, want 0", area)
				}
			}
		})
	}
}
