package wrap

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/clip"
	"github.com/paulmach/orb/planar"
)

// AroundBound takes a ring and if invalid (i.e. endpoints don't match) will
// connect the endpoints around the boundary of the bound in the direction provided.
func AroundBound(b planar.Bound, g planar.Geometry, o orb.Orientation) planar.Geometry {
	if o != orb.CCW && o != orb.CW {
		panic("invalid orientation")
	}

	switch g := g.(type) {
	case planar.Point, planar.MultiPoint:
		return g
	case planar.LineString, planar.MultiLineString:
		return g
	case planar.Bound:
		return g
	case planar.Ring:
		return Ring(b, g, o)
	case planar.Polygon:
		return Polygon(b, g, o)
	case planar.MultiPolygon:
		return MultiPolygon(b, g, o)
	case planar.Collection:
		return Collection(b, g, o)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// Ring will connect the ring round the bound in the direction provided.
func Ring(b planar.Bound, r planar.Ring, o orb.Orientation) planar.Ring {
	result, _ := clip.AroundBound(
		mapBound(b),
		&lineString{ls: planar.LineString(r)},
		o,
		func(in clip.LineString) orb.Orientation {
			ls := in.(*lineString)
			return planar.Ring(ls.ls).Orientation()
		},
	)

	return planar.Ring(result.(*lineString).ls)
}

// Polygon will connect the polygon rings around the bound assuming the outer
// ring is in the direction provided and the inner rings are the opposite.
func Polygon(b planar.Bound, p planar.Polygon, o orb.Orientation) planar.Polygon {
	r := Ring(b, p[0], o)

	result := planar.Polygon{r}
	if len(p) <= 1 {
		return result
	}

	for i := 1; i < len(p); i++ {
		r := Ring(b, p[i], -1*o)
		if r != nil {
			result = append(result, r)
		}
	}

	return result
}

// MultiPolygon will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
func MultiPolygon(b planar.Bound, mp planar.MultiPolygon, o orb.Orientation) planar.MultiPolygon {
	if len(mp) == 0 {
		return mp
	}

	result := make(planar.MultiPolygon, 0, len(mp))
	for _, polygon := range mp {
		p := Polygon(b, polygon, o)
		if p != nil {
			result = append(result, p)
		}
	}

	return result
}

// Collection will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
// It will noop non-2d geometry.
func Collection(b planar.Bound, c planar.Collection, o orb.Orientation) planar.Collection {
	if len(c) == 0 {
		return c
	}

	result := make(planar.Collection, 0, len(c))
	for _, g := range c {
		result = append(result, AroundBound(b, g, o))
	}

	return result
}

type lineString struct {
	ls planar.LineString
}

func (ls *lineString) Len() int {
	return len(ls.ls)
}

func (ls *lineString) Get(i int) (x, y float64) {
	return ls.ls[i][0], ls.ls[i][1]
}

func (ls *lineString) Append(x, y float64) {
	ls.ls = append(ls.ls, planar.NewPoint(x, y))
}

func (ls *lineString) Clear() {
	ls.ls = ls.ls[:0]
}

func mapBound(b planar.Bound) clip.Bound {
	return clip.Bound{
		Left:   b.Left(),
		Right:  b.Right(),
		Bottom: b.Bottom(),
		Top:    b.Top(),
	}
}
