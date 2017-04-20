package project

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// PathForward is a helper to project an entire path.
func PathForward(path geo.Path, f Forward) planar.Path {
	n := planar.NewPathPreallocate(len(path), len(path))
	for i := range path {
		n[i] = f(path[i])
	}

	return n
}

// PathReverse is a helper to project an entire path.
func PathReverse(path planar.Path, r Reverse) geo.Path {
	n := geo.NewPathPreallocate(len(path), len(path))
	for i := range path {
		n[i] = r(path[i])
	}

	return n
}

// RectForward is a helper to project a bound.
func RectForward(bound geo.Rect, f Forward) planar.Rect {
	return planar.NewRectFromPoints(
		f(geo.Point(bound.SW)),
		f(geo.Point(bound.NE)),
	)

}

// RectReverse is a helper to project a bound.
func RectReverse(bound planar.Rect, r Reverse) geo.Rect {
	return geo.NewRectFromPoints(
		r(planar.Point(bound.SW)),
		r(planar.Point(bound.NE)),
	)
}
