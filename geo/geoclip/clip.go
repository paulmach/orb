package geoclip

import (
	"math"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/clip"
)

// Clip will clip the geometry to the bounding box using the
// correct functions for the type.
func Clip(b geo.Bound, g geo.Geometry) geo.Geometry {
	if !b.Intersects(g.Bound()) {
		return nil
	}

	switch g := g.(type) {
	case geo.Point:
		return g // Intersect check above
	case geo.MultiPoint:
		return MultiPoint(b, g)
	case geo.LineString:
		mls := LineString(b, g)
		if len(mls) == 1 {
			return mls[0]
		}
	case geo.MultiLineString:
		return MultiLineString(b, g)
	case geo.Ring:
		return Ring(b, g)
	case geo.Polygon:
		return Polygon(b, g)
	case geo.MultiPolygon:
		return MultiPolygon(b, g)
	case geo.Collection:
		return Collection(b, g)
	case geo.Bound:
		b = Bound(b, g)
		if b.IsEmpty() {
			return nil
		}

		return b
	}

	panic("geometry type not supported")
}

// MultiPoint returns a new set with the points outside the bound removed.
func MultiPoint(b geo.Bound, mp geo.MultiPoint) geo.MultiPoint {
	var result geo.MultiPoint
	for _, p := range mp {
		if b.Contains(p) {
			result = append(result, p)
		}
	}

	return result
}

// LineString clips the linestring to the bounding box.
func LineString(b geo.Bound, ls geo.LineString) geo.MultiLineString {
	result := &multiLineString{}
	clip.Line(mapBound(b), &lineString{ls: ls}, result)

	if len(result.mls) == 0 {
		return nil
	}

	return result.mls
}

// MultiLineString clips the linestrings to the bounding box
// and returns a linestring union.
func MultiLineString(b geo.Bound, mls geo.MultiLineString) geo.MultiLineString {
	result := geo.MultiLineString{}
	for _, ls := range mls {
		result = append(result, LineString(b, ls)...)
	}
	return result
}

// Ring clips the ring to the bounding box and returns another ring.
func Ring(b geo.Bound, r geo.Ring) geo.Ring {
	result := &lineString{}
	clip.Ring(mapBound(b), &lineString{ls: geo.LineString(r)}, result)

	if len(result.ls) == 0 {
		return nil
	}

	return geo.Ring(result.ls)
}

// Polygon clips the polygon to the bounding box excluding the inner rings
// if they do not intersect the bounding box.
func Polygon(b geo.Bound, p geo.Polygon) geo.Polygon {
	r := Ring(b, p[0])
	if r == nil {
		return nil
	}

	result := geo.Polygon{r}
	for i := 1; i < len(p); i++ {
		r := Ring(b, p[i])
		if r != nil {
			result = append(result, r)
		}
	}

	return result
}

// MultiPolygon clips the multi polygon to the bounding box excluding
// any polygons if they don't intersect the bounding box.
func MultiPolygon(b geo.Bound, mp geo.MultiPolygon) geo.MultiPolygon {
	var result geo.MultiPolygon
	for _, polygon := range mp {
		p := Polygon(b, polygon)
		if p != nil {
			result = append(result, p)
		}
	}

	return result
}

// Collection clips each element in the collection to the bounding box.
// It will exclude elements if they don't intersect the bounding box.
func Collection(b geo.Bound, c geo.Collection) geo.Collection {
	var result geo.Collection
	for _, g := range c {
		clipped := Clip(b, g)
		if clipped != nil {
			result = append(result, clipped)
		}
	}

	return result
}

// Bound intersects the two bounds. May result in an
// empty/degenerate bound.
func Bound(b, bound geo.Bound) geo.Bound {
	if b.IsEmpty() && bound.IsEmpty() {
		return bound
	}

	if b.IsEmpty() {
		return bound
	} else if bound.IsEmpty() {
		return b
	}

	return geo.Bound{
		geo.Point{
			math.Max(b[0][0], bound[0][0]),
			math.Max(b[0][1], bound[0][1]),
		},
		geo.Point{
			math.Min(b[1][0], bound[1][0]),
			math.Min(b[1][1], bound[1][1]),
		},
	}
}

type lineString struct {
	ls geo.LineString
}

func (ls *lineString) Len() int {
	return len(ls.ls)
}

func (ls *lineString) Get(i int) (x, y float64) {
	return ls.ls[i][0], ls.ls[i][0]
}

func (ls *lineString) Append(x, y float64) {
	ls.ls = append(ls.ls, geo.NewPoint(x, y))
}

func (ls *lineString) Clear() {
	ls.ls = ls.ls[:0]
}

type multiLineString struct {
	mls geo.MultiLineString
}

func (mls *multiLineString) Append(i int, x, y float64) {
	if i >= len(mls.mls) {
		mls.mls = append(mls.mls, geo.LineString{{x, y}})
	} else {
		mls.mls[i] = append(mls.mls[i], geo.Point{x, y})
	}
}

func mapBound(b geo.Bound) clip.Bound {
	return clip.Bound{
		Left:   b.West(),
		Right:  b.East(),
		Bottom: b.South(),
		Top:    b.North(),
	}
}
