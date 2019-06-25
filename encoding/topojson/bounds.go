package topojson

import (
	"math"

	"github.com/paulmach/orb"
)

func (t *Topology) bounds() {
	t.BBox = []float64{
		math.MaxFloat64,
		math.MaxFloat64,
		-math.MaxFloat64,
		-math.MaxFloat64,
	}

	for _, f := range t.input {
		t.boundGeometry(f.Geometry)
	}

}

func (t *Topology) boundGeometry(g orb.Geometry) {
	switch c := g.(type) {
	case orb.Point:
		t.Bound(c.Bound())
	case orb.MultiPoint:
		t.Bound(c.Bound())
	case orb.LineString:
		t.Bound(c.Bound())
	case orb.MultiLineString:
		t.Bound(c.Bound())
	case orb.Polygon:
		t.Bound(c.Bound())
	case orb.MultiPolygon:
		t.Bound(c.Bound())
	case orb.Collection:
		t.Bound(c.Bound())
		// for _, geo := range c {
		// 	t.boundGeometry(geo)
		// }
	}
}

func (t *Topology) Bound(b orb.Bound) {
	xx := []float64{b.Min[0], b.Max[0]}
	yy := []float64{b.Min[1], b.Max[1]}
	for _, x := range xx {
		if x < t.BBox[0] {
			t.BBox[0] = x
		}
		if x > t.BBox[2] {
			t.BBox[2] = x
		}
	}
	for _, y := range yy {
		if y < t.BBox[1] {
			t.BBox[1] = y
		}
		if y > t.BBox[3] {
			t.BBox[3] = y
		}
	}
}

func (t *Topology) boundPoint(p []float64) {
	x := p[0]
	y := p[1]

	if x < t.BBox[0] {
		t.BBox[0] = x
	}
	if x > t.BBox[2] {
		t.BBox[2] = x
	}
	if y < t.BBox[1] {
		t.BBox[1] = y
	}
	if y > t.BBox[3] {
		t.BBox[3] = y
	}
}

func (t *Topology) boundPoints(l [][]float64) {
	for _, p := range l {
		t.boundPoint(p)
	}
}

func (t *Topology) boundMultiPoints(ml [][][]float64) {
	for _, l := range ml {
		for _, p := range l {
			t.boundPoint(p)
		}
	}
}
