package topojson

import (
	"math"

	"github.com/paulmach/orb"
)

type quantize struct {
	Transform *Transform

	dx, dy, kx, ky float64
}

func newQuantize(dx, dy, kx, ky float64) *quantize {
	return &quantize{
		dx: dx,
		dy: dy,
		kx: kx,
		ky: ky,

		Transform: &Transform{
			Scale:     [2]float64{1 / kx, 1 / ky},
			Translate: [2]float64{-dx, -dy},
		},
	}
}

func (q *quantize) quantizePoint(p orb.Point) orb.Point {
	x := round((p[0] + q.dx) * q.kx)
	y := round((p[1] + q.dy) * q.ky)
	return orb.Point{x, y}
}

func (q *quantize) quantizeMultiPoint(in orb.MultiPoint, skipEqual bool) orb.MultiPoint {
	out := orb.MultiPoint{}

	var last []float64

	for _, p := range in {
		pt := q.quantizePoint(p)
		if !pointEquals([]float64{pt[0], pt[1]}, last) || !skipEqual {
			out = append(out, pt)
			last = []float64{pt[0], pt[1]}
		}
	}

	if len(out) < 2 {
		out = append(out, out[0])
	}

	return out
}

func (q *quantize) quantizeLine(in orb.LineString, skipEqual bool) orb.LineString {
	out := orb.LineString{}

	var last []float64

	for _, p := range in {
		pt := q.quantizePoint(p)
		if !pointEquals([]float64{pt[0], pt[1]}, last) || !skipEqual {
			out = append(out, pt)
			last = []float64{pt[0], pt[1]}
		}
	}

	if len(out) < 2 {
		out = append(out, out[0])
	}

	return out
}

func (q *quantize) quantizeRing(in orb.Ring, skipEqual bool) orb.Ring {
	out := orb.Ring{}

	var last []float64

	for _, p := range in {
		pt := q.quantizePoint(p)
		if !pointEquals([]float64{pt[0], pt[1]}, last) || !skipEqual {
			out = append(out, pt)
			last = []float64{pt[0], pt[1]}
		}
	}

	if len(out) < 2 {
		out = append(out, out[0])
	}

	return out
}

func (q *quantize) quantizeMultiLine(in orb.MultiLineString, skipEqual bool) orb.MultiLineString {
	out := make(orb.MultiLineString, len(in))
	for i, line := range in {
		line = q.quantizeLine(line, skipEqual)
		for len(line) < 4 {
			line = append(line, line[0])
		}
		out[i] = line
	}
	return out
}

func (q *quantize) quantizePolygon(in orb.Polygon, skipEqual bool) orb.Polygon {
	out := make(orb.Polygon, len(in))
	for i, ring := range in {
		out[i] = q.quantizeRing(ring, skipEqual)
	}
	return out
}

func (q *quantize) quantizeMultiPolygon(in orb.MultiPolygon, skipEqual bool) orb.MultiPolygon {
	out := make(orb.MultiPolygon, len(in))
	for i, ring := range in {
		out[i] = q.quantizePolygon(ring, skipEqual)
	}
	return out
}

func round(v float64) float64 {
	if v < 0 {
		return math.Ceil(v - 0.5)
	} else {
		return math.Floor(v + 0.5)
	}
}
