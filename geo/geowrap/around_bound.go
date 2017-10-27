package geowrap

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/clip"
)

// AroundBound takes a ring and if invalid (i.e. endpoints don't match) will
// connect the endpoints around the boundary of the bound in the direction provided.
// Will append to the original geometry.
func AroundBound(b geo.Bound, g geo.Geometry, o orb.Orientation) (geo.Geometry, error) {
	if o != orb.CCW && o != orb.CW {
		panic("invalid orientation")
	}

	switch g := g.(type) {
	case geo.Point, geo.MultiPoint:
		return g, nil
	case geo.LineString, geo.MultiLineString:
		return g, nil
	case geo.Bound:
		return g, nil
	case geo.Ring:
		return Ring(b, g, o)
	case geo.Polygon:
		return Polygon(b, g, o)
	case geo.MultiPolygon:
		return MultiPolygon(b, g, o)
	case geo.Collection:
		return Collection(b, g, o)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// Ring will connect the ring round the bound in the direction provided.
func Ring(b geo.Bound, r geo.Ring, o orb.Orientation) (geo.Ring, error) {
	result, err := clip.AroundBound(
		mapBound(b),
		&lineString{ls: geo.LineString(r)},
		o,
		func(in clip.LineString) orb.Orientation {
			ls := in.(*lineString)
			return geo.Ring(ls.ls).Orientation()
		},
	)

	if err != nil {
		return nil, err
	}

	return geo.Ring(result.(*lineString).ls), nil
}

// Polygon will connect the polygon rings around the bound assuming the outer
// ring is in the direction provided and the inner rings are the opposite.
func Polygon(b geo.Bound, p geo.Polygon, o orb.Orientation) (geo.Polygon, error) {
	r, err := Ring(b, p[0], o)
	if err != nil {
		return nil, err
	}

	result := geo.Polygon{r}
	if len(p) <= 1 {
		return result, nil
	}

	for i := 1; i < len(p); i++ {
		r, err := Ring(b, p[i], -1*o)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// MultiPolygon will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
func MultiPolygon(b geo.Bound, mp geo.MultiPolygon, o orb.Orientation) (geo.MultiPolygon, error) {
	if len(mp) == 0 {
		return mp, nil
	}

	result := make(geo.MultiPolygon, 0, len(mp))
	for _, polygon := range mp {
		p, err := Polygon(b, polygon, o)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

// Collection will connect the polygon rings around the bound assuming the outer
// rings are in the direction provided and the inner rings are the opposite.
// It will noop non-2d geometry.
func Collection(b geo.Bound, c geo.Collection, o orb.Orientation) (geo.Collection, error) {
	if len(c) == 0 {
		return c, nil
	}

	result := make(geo.Collection, 0, len(c))
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
	ls geo.LineString
}

func (ls *lineString) Len() int {
	return len(ls.ls)
}

func (ls *lineString) Get(i int) (x, y float64) {
	return ls.ls[i][0], ls.ls[i][1]
}

func (ls *lineString) Append(x, y float64) {
	ls.ls = append(ls.ls, geo.NewPoint(x, y))
}

func (ls *lineString) Clear() {
	ls.ls = ls.ls[:0]
}

func mapBound(b geo.Bound) clip.Bound {
	return clip.Bound{
		Left:   b.West(),
		Right:  b.East(),
		Bottom: b.South(),
		Top:    b.North(),
	}
}
