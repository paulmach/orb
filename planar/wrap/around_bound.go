package wrap

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/clip"
	"github.com/paulmach/orb/planar"
)

// AroundBound takes a ring and if invalid (i.e. endpoints don't match) will
// connect the endpoints around the boundary of the bound in the direction provided.
// Will append to the original geometry.
func AroundBound(b planar.Bound, g planar.Geometry, o orb.Orientation) (planar.Geometry, error) {
	if o != orb.CCW && o != orb.CW {
		panic("invalid orientation")
	}

	switch g := g.(type) {
	case planar.Point, planar.MultiPoint:
		return g, nil
	case planar.LineString, planar.MultiLineString:
		return g, nil
	case planar.Bound:
		return g, nil
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
func Ring(b planar.Bound, r planar.Ring, o orb.Orientation) (planar.Ring, error) {
	result, err := clip.AroundBound(
		mapBound(b),
		&lineString{ls: planar.LineString(r)},
		o,
		func(in clip.LineString) orb.Orientation {
			ls := in.(*lineString)
			return planar.Ring(ls.ls).Orientation()
		},
	)

	if err != nil {
		return nil, err
	}

	return planar.Ring(result.(*lineString).ls), nil
}

// Polygon will connect the polygon rings around the bound assuming the outer
// ring is in the direction provided and the inner rings are the opposite.
func Polygon(b planar.Bound, p planar.Polygon, o orb.Orientation) (planar.Polygon, error) {
	r, err := Ring(b, p[0], o)
	if err != nil {
		return nil, err
	}

	result := planar.Polygon{r}
	if len(p) <= 1 {
		return result, nil
	}

	for i := 1; i < len(p); i++ {
		r, err := Ring(b, p[i], -1*o)
		if err != nil {
			return nil, err
		}
		if r != nil {
			result = append(result, r)
		}
	}

	return result, nil
}

// MultiPolygon will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
func MultiPolygon(b planar.Bound, mp planar.MultiPolygon, o orb.Orientation) (planar.MultiPolygon, error) {
	if len(mp) == 0 {
		return mp, nil
	}

	result := make(planar.MultiPolygon, 0, len(mp))
	for _, polygon := range mp {
		p, err := Polygon(b, polygon, o)
		if err != nil {
			return nil, err
		}

		if p != nil {
			result = append(result, p)
		}
	}

	return result, nil
}

// Collection will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
// It will noop non-2d geometry.
func Collection(b planar.Bound, c planar.Collection, o orb.Orientation) (planar.Collection, error) {
	if len(c) == 0 {
		return c, nil
	}

	result := make(planar.Collection, 0, len(c))
	for _, g := range c {
		ng, err := AroundBound(b, g, o)
		if err != nil {
			return nil, err
		}

		result = append(result, ng)
	}

	return result, nil
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
