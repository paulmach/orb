package project

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// ToPlanar projects a geometry from geo -> planar
func ToPlanar(g geo.Geometry, proj Projection) planar.Geometry {
	switch g := g.(type) {
	case geo.Point:
		return proj.ToPlanar(g)
	case geo.LineString:
		return LineStringToPlanar(g, proj)
	case geo.MultiPoint:
		return MultiPointToPlanar(g, proj)
	case geo.Polygon:
		return PolygonToPlanar(g, proj)
	}

	panic("geometry type not supported")
}

// ToGeo projects a geometry from planar -> geo
func ToGeo(g planar.Geometry, proj Projection) geo.Geometry {
	switch g := g.(type) {
	case planar.Point:
		return proj.ToGeo(g)
	case planar.LineString:
		return LineStringToGeo(g, proj)
	case planar.MultiPoint:
		return MultiPointToGeo(g, proj)
	case planar.Polygon:
		return PolygonToGeo(g, proj)
	}

	panic("geometry type not supported")
}

// MultiPointToPlanar is a helper to project an entire multi point.
func MultiPointToPlanar(mp geo.MultiPoint, proj Projection) planar.MultiPoint {
	n := planar.NewMultiPointPreallocate(len(mp), len(mp))
	for i := range mp {
		n[i] = proj.ToPlanar(mp[i])
	}

	return n
}

// MultiPointToGeo is a helper to project an entire multi point.
func MultiPointToGeo(mp planar.MultiPoint, proj Projection) geo.MultiPoint {
	n := geo.NewMultiPointPreallocate(len(mp), len(mp))
	for i := range mp {
		n[i] = proj.ToGeo(mp[i])
	}

	return n
}

// LineStringToPlanar is a helper to project an entire line string.
func LineStringToPlanar(ls geo.LineString, proj Projection) planar.LineString {
	n := planar.NewLineStringPreallocate(len(ls), len(ls))
	for i := range ls {
		n[i] = proj.ToPlanar(ls[i])
	}

	return n
}

// LineStringToGeo is a helper to project an entire line string.
func LineStringToGeo(ls planar.LineString, proj Projection) geo.LineString {
	n := geo.NewLineStringPreallocate(len(ls), len(ls))
	for i := range ls {
		n[i] = proj.ToGeo(ls[i])
	}

	return n
}

// PolygonToPlanar is a helper to project an entire polygon.
func PolygonToPlanar(p geo.Polygon, proj Projection) planar.Polygon {
	n := make(planar.Polygon, len(p), len(p))
	for i := range p {
		n[i] = LineStringToPlanar(p[i], proj)
	}

	return n
}

// PolygonToGeo is a helper to project an entire line string.
func PolygonToGeo(p planar.Polygon, proj Projection) geo.Polygon {
	n := make(geo.Polygon, len(p), len(p))
	for i := range p {
		n[i] = LineStringToGeo(p[i], proj)
	}

	return n
}

// RectToPlanar is a helper to project a rectangle.
func RectToPlanar(bound geo.Rect, proj Projection) planar.Rect {
	return planar.NewRectFromPoints(
		proj.ToPlanar(bound[0]),
		proj.ToPlanar(bound[1]),
	)
}

// RectToGeo is a helper to project a rectangle.
func RectToGeo(bound planar.Rect, proj Projection) geo.Rect {
	return geo.NewRectFromPoints(
		proj.ToGeo(bound[0]),
		proj.ToGeo(bound[1]),
	)
}
