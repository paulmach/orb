package topojson

import (
	"github.com/paulmach/orb"
)

func (t *Topology) postQuantize() {
	q0 := t.opts.PreQuantize
	q1 := t.opts.PostQuantize

	if q1 == 0 {
		return
	}

	var q *quantize

	if q0 != 0 {
		if q0 == q1 {
			return
		}

		k := q1 / q0

		q = newQuantize(0, 0, k, k)

		t.Transform.Scale[0] /= k
		t.Transform.Scale[1] /= k
	} else {
		x0 := t.BBox[0]
		y0 := t.BBox[1]
		x1 := t.BBox[2]
		y1 := t.BBox[3]

		kx := float64(1)
		if x1-x0 != 0 {
			kx = (q1 - 1) / (x1 - x0)
		}

		ky := float64(1)
		if y1-y0 != 0 {
			ky = (q1 - 1) / (y1 - y0)
		}

		q = newQuantize(-x0, -y0, kx, ky)
		t.Transform = q.Transform
	}

	for _, f := range t.input {
		t.postQuantizeGeometry(q, f.Geometry)
	}

	for i, arc := range t.Arcs {
		a := make(orb.LineString, len(arc))
		for i, v := range arc {
			a[i] = orb.Point{v[0], v[1]}
		}
		b := q.quantizeLine(a, true)
		c := make([][]float64, len(b))
		for i, v := range b {
			c[i] = []float64{v[0], v[1]}
		}
		t.Arcs[i] = c
	}
}

func (t *Topology) postQuantizeGeometry(q *quantize, g orb.Geometry) {
	switch v := g.(type) {
	default:
		for _, geom := range g.(orb.Collection) {
			t.postQuantizeGeometry(q, geom)
		}
	case orb.Point:
		v = q.quantizePoint(v)
	case orb.LineString:
		v = q.quantizeLine(v, false)
	}
}
