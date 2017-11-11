package wrap

import (
	"fmt"

	"github.com/paulmach/orb"
)

// AroundBound takes a ring and if invalid (i.e. endpoints don't match) will
// connect the endpoints around the boundary of the bound in the direction provided.
// Will append to the original geometry.
func AroundBound(b orb.Bound, g orb.Geometry, o orb.Orientation) (orb.Geometry, error) {
	if g == nil {
		return nil, nil
	}

	switch g := g.(type) {
	case orb.Point, orb.MultiPoint:
		return g, nil
	case orb.LineString, orb.MultiLineString:
		return g, nil
	case orb.Bound:
		return g, nil
	case orb.Ring:
		return Ring(b, g, o)
	case orb.Polygon:
		return Polygon(b, g, o)
	case orb.MultiPolygon:
		return MultiPolygon(b, g, o)
	case orb.Collection:
		return Collection(b, g, o)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// Ring will connect the ring round the bound in the direction provided.
func Ring(b orb.Bound, r orb.Ring, o orb.Orientation) (orb.Ring, error) {
	result, err := aroundBound(b, r, o)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Polygon will connect the polygon rings around the bound assuming the outer
// ring is in the direction provided and the inner rings are the opposite.
func Polygon(b orb.Bound, p orb.Polygon, o orb.Orientation) (orb.Polygon, error) {
	if len(p) == 0 {
		return p, nil
	}

	r, err := Ring(b, p[0], o)
	if err != nil {
		return nil, err
	}

	result := orb.Polygon{r}
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
func MultiPolygon(b orb.Bound, mp orb.MultiPolygon, o orb.Orientation) (orb.MultiPolygon, error) {
	if len(mp) == 0 {
		return mp, nil
	}

	result := make(orb.MultiPolygon, 0, len(mp))
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
func Collection(b orb.Bound, c orb.Collection, o orb.Orientation) (orb.Collection, error) {
	if len(c) == 0 {
		return c, nil
	}

	result := make(orb.Collection, 0, len(c))
	for _, g := range c {
		ng, err := AroundBound(b, g, o)
		if err != nil {
			return nil, err
		}

		result = append(result, ng)
	}

	return result, nil
}
