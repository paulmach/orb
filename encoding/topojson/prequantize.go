package topojson

import (
	"github.com/paulmach/orb"
)

func (t *Topology) preQuantize() {
	if t.opts.PreQuantize == 0 {
		return
	}
	if t.opts.PostQuantize == 0 {
		t.opts.PostQuantize = t.opts.PreQuantize
	}

	q0 := t.opts.PreQuantize
	q1 := t.opts.PostQuantize

	x0 := t.BBox[0]
	y0 := t.BBox[1]
	x1 := t.BBox[2]
	y1 := t.BBox[3]

	kx := float64(1)
	if x1-x0 != 0 {
		kx = (q1 - 1) / (x1 - x0) * q0 / q1
	}

	ky := float64(1)
	if y1-y0 != 0 {
		ky = (q1 - 1) / (y1 - y0) * q0 / q1
	}

	q := newQuantize(-x0, -y0, kx, ky)

	for _, f := range t.input {
		t.preQuantizeGeometry(q, &f.Geometry)
	}

	t.Transform = q.Transform
}

func (t *Topology) preQuantizeGeometry(q *quantize, g *orb.Geometry) {
	switch v := (*g).(type) {
	case orb.Collection:
		for _, g := range v {
			t.preQuantizeGeometry(q, &g)
		}
	case orb.Point:
		*g = q.quantizePoint(v)
	case orb.MultiPoint:
		*g = q.quantizeMultiPoint(v, false)
	case orb.LineString:
		*g = q.quantizeLine(v, true)
	case orb.MultiLineString:
		*g = q.quantizeMultiLine(v, true)
	case orb.Polygon:
		*g = q.quantizePolygon(v, true)
	case orb.MultiPolygon:
		*g = q.quantizeMultiPolygon(v, true)
	}
}
