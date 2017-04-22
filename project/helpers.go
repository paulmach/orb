package project

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// ForwardLineString is a helper to project an entire line string.
func ForwardLineString(ls geo.LineString, f func(geo.Point) planar.Point) planar.LineString {
	n := planar.NewLineStringPreallocate(len(ls), len(ls))
	for i := range ls {
		n[i] = f(ls[i])
	}

	return n
}

// ReverseLineString is a helper to project an entire line string.
func ReverseLineString(ls planar.LineString, r func(planar.Point) geo.Point) geo.LineString {
	n := geo.NewLineStringPreallocate(len(ls), len(ls))
	for i := range ls {
		n[i] = r(ls[i])
	}

	return n
}

// ForwardRect is a helper to project a rectangle.
func ForwardRect(bound geo.Rect, f func(geo.Point) planar.Point) planar.Rect {
	return planar.NewRectFromPoints(f(bound[0]), f(bound[1]))
}

// ReverseRect is a helper to project a rectangle.
func ReverseRect(bound planar.Rect, r func(planar.Point) geo.Point) geo.Rect {
	return geo.NewRectFromPoints(r(bound[0]), r(bound[1]))
}
