package topojson

import (
	geo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/reducers"
)

func (t *Topology) simplify() {
	t.deletedArcs = make(map[int]bool)
	t.shiftArcs = make(map[int]int)

	if t.opts.Simplify == 0 {
		for i := range t.Arcs {
			t.deletedArcs[i] = false
			t.shiftArcs[i] = 0
		}
		return
	}

	newArcs := make([][][]float64, 0)
	for i, arc := range t.Arcs {
		path := geo.NewPathFromYXSlice(arc)
		path = reducers.VisvalingamThreshold(path, t.opts.Simplify)
		points := path.Points()
		newArc := make([][]float64, len(points))
		for j, p := range points {
			newArc[j] = []float64{p[1], p[0]}
		}

		if i == 0 {
			t.shiftArcs[i] = 0
		} else {
			t.shiftArcs[i] = t.shiftArcs[i-1]
		}

		remove := len(newArc) <= 2 && pointEquals(newArc[0], newArc[1])
		if remove {
			// Zero-length arc, remove it!
			t.deletedArcs[i] = true
			t.shiftArcs[i] += 1
		} else {
			newArcs = append(newArcs, newArc)
		}
	}
	t.Arcs = newArcs
}
