package project

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// LineStringForward is a helper to project an entire line string.
func LineStringForward(path geo.LineString, f Forward) planar.LineString {
	n := planar.NewLineStringPreallocate(len(path), len(path))
	for i := range path {
		n[i] = f(path[i])
	}

	return n
}

// LineStringReverse is a helper to project an entire line string.
func LineStringReverse(path planar.LineString, r Reverse) geo.LineString {
	n := geo.NewLineStringPreallocate(len(path), len(path))
	for i := range path {
		n[i] = r(path[i])
	}

	return n
}

// RectForward is a helper to project a rectangle.
func RectForward(bound geo.Rect, f Forward) planar.Rect {
	return planar.NewRectFromPoints(
		f(geo.Point(bound.SW)),
		f(geo.Point(bound.NE)),
	)

}

// RectReverse is a helper to project a rectangle.
func RectReverse(bound planar.Rect, r Reverse) geo.Rect {
	return geo.NewRectFromPoints(
		r(planar.Point(bound.SW)),
		r(planar.Point(bound.NE)),
	)
}
