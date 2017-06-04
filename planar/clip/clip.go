package clip

import (
	"fmt"
	"math"

	"github.com/paulmach/orb/internal/clip"
	"github.com/paulmach/orb/planar"
)

// Clip will clip the geometry to the bounding box using the
// correct functions for the type.
func Clip(b planar.Bound, g planar.Geometry) planar.Geometry {
	if !b.Intersects(g.Bound()) {
		return nil
	}

	switch g := g.(type) {
	case planar.Point:
		return g // Intersect check above
	case planar.MultiPoint:
		mp := MultiPoint(b, g)
		if mp == nil {
			return nil
		}

		return mp
	case planar.LineString:
		mls := LineString(b, g)
		if len(mls) == 1 {
			return mls[0]
		}

		if mls == nil {
			return nil
		}
		return mls
	case planar.MultiLineString:
		mls := MultiLineString(b, g)
		if mls == nil {
			return nil
		}

		return mls
	case planar.Ring:
		r := Ring(b, g)
		if r == nil {
			return nil
		}

		return r
	case planar.Polygon:
		p := Polygon(b, g)
		if p == nil {
			return p
		}

		return p
	case planar.MultiPolygon:
		mp := MultiPolygon(b, g)
		if mp == nil {
			return nil
		}

		return mp
	case planar.Collection:
		c := Collection(b, g)
		if c == nil {
			return nil
		}

		return c
	case planar.Bound:
		b = Bound(b, g)
		if b.IsEmpty() {
			return nil
		}

		return b
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// MultiPoint returns a new set with the points outside the bound removed.
func MultiPoint(b planar.Bound, mp planar.MultiPoint) planar.MultiPoint {
	var result planar.MultiPoint
	for _, p := range mp {
		if b.Contains(p) {
			result = append(result, p)
		}
	}

	return result
}

// LineString clips the linestring to the bounding box.
func LineString(b planar.Bound, ls planar.LineString) planar.MultiLineString {
	result := &multiLineString{}
	clip.Line(mapBound(b), &lineString{ls: ls}, result)

	if len(result.mls) == 0 {
		return nil
	}

	return result.mls
}

// MultiLineString clips the linestrings to the bounding box
// and returns a linestring union.
func MultiLineString(b planar.Bound, mls planar.MultiLineString) planar.MultiLineString {
	var result planar.MultiLineString
	for _, ls := range mls {
		r := LineString(b, ls)
		if r != nil {
			result = append(result, r...)
		}
	}
	return result
}

// Ring clips the ring to the bounding box and returns another ring.
func Ring(b planar.Bound, r planar.Ring) planar.Ring {
	result := &lineString{}
	clip.Ring(mapBound(b), &lineString{ls: planar.LineString(r)}, result)

	if len(result.ls) == 0 {
		return nil
	}

	return planar.Ring(result.ls)
}

// Polygon clips the polygon to the bounding box excluding the inner rings
// if they do not intersect the bounding box.
func Polygon(b planar.Bound, p planar.Polygon) planar.Polygon {
	r := Ring(b, p[0])
	if r == nil {
		return nil
	}

	result := planar.Polygon{r}
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
func MultiPolygon(b planar.Bound, mp planar.MultiPolygon) planar.MultiPolygon {
	var result planar.MultiPolygon
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
func Collection(b planar.Bound, c planar.Collection) planar.Collection {
	var result planar.Collection
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
func Bound(b, bound planar.Bound) planar.Bound {
	if b.IsEmpty() && bound.IsEmpty() {
		return bound
	}

	if b.IsEmpty() {
		return bound
	} else if bound.IsEmpty() {
		return b
	}

	return planar.Bound{
		planar.Point{
			math.Max(b[0][0], bound[0][0]),
			math.Max(b[0][1], bound[0][1]),
		},
		planar.Point{
			math.Min(b[1][0], bound[1][0]),
			math.Min(b[1][1], bound[1][1]),
		},
	}
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

type multiLineString struct {
	mls planar.MultiLineString
}

func (mls *multiLineString) Append(i int, x, y float64) {
	if i >= len(mls.mls) {
		mls.mls = append(mls.mls, planar.LineString{{x, y}})
	} else {
		mls.mls[i] = append(mls.mls[i], planar.Point{x, y})
	}
}

func mapBound(b planar.Bound) clip.Bound {
	return clip.Bound{
		Left:   b.Left(),
		Right:  b.Right(),
		Bottom: b.Bottom(),
		Top:    b.Top(),
	}
}
