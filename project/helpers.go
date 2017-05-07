package project

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// ToPlanar projects a geometry from geo -> planar
func ToPlanar(g geo.Geometry, proj *Projection) planar.Geometry {
	switch g := g.(type) {
	case geo.Point:
		return proj.ToPlanar(g)
	case geo.MultiPoint:
		return MultiPointToPlanar(g, proj)
	case geo.LineString:
		return LineStringToPlanar(g, proj)
	case geo.MultiLineString:
		return MultiLineStringToPlanar(g, proj)
	case geo.Polygon:
		return PolygonToPlanar(g, proj)
	case geo.MultiPolygon:
		return MultiPolygonToPlanar(g, proj)
	case geo.Bound:
		return BoundToPlanar(g, proj)
	}

	panic("geometry type not supported")
}

// ToGeo projects a geometry from planar -> geo
func ToGeo(g planar.Geometry, proj *Projection) geo.Geometry {
	switch g := g.(type) {
	case planar.Point:
		return proj.ToGeo(g)
	case planar.MultiPoint:
		return MultiPointToGeo(g, proj)
	case planar.LineString:
		return LineStringToGeo(g, proj)
	case planar.MultiLineString:
		return MultiLineStringToGeo(g, proj)
	case planar.Polygon:
		return PolygonToGeo(g, proj)
	case planar.MultiPolygon:
		return MultiPolygonToGeo(g, proj)
	case planar.Bound:
		return BoundToGeo(g, proj)
	}

	panic("geometry type not supported")
}

// MultiPointToPlanar is a helper to project an entire multi point.
func MultiPointToPlanar(mp geo.MultiPoint, proj *Projection) planar.MultiPoint {
	n := make(planar.MultiPoint, len(mp))
	for i := range mp {
		n[i] = proj.ToPlanar(mp[i])
	}

	return n
}

// MultiPointToGeo is a helper to project an entire multi point.
func MultiPointToGeo(mp planar.MultiPoint, proj *Projection) geo.MultiPoint {
	n := make(geo.MultiPoint, len(mp))
	for i := range mp {
		n[i] = proj.ToGeo(mp[i])
	}

	return n
}

// LineStringToPlanar is a helper to project an entire line string.
func LineStringToPlanar(ls geo.LineString, proj *Projection) planar.LineString {
	n := make(planar.LineString, len(ls))
	for i := range ls {
		n[i] = proj.ToPlanar(ls[i])
	}

	return n
}

// LineStringToGeo is a helper to project an entire line string.
func LineStringToGeo(ls planar.LineString, proj *Projection) geo.LineString {
	n := make(geo.LineString, len(ls))
	for i := range ls {
		n[i] = proj.ToGeo(ls[i])
	}

	return n
}

// MultiLineStringToPlanar is a helper to project an entire multi linestring.
func MultiLineStringToPlanar(mls geo.MultiLineString, proj *Projection) planar.MultiLineString {
	n := make(planar.MultiLineString, len(mls))
	for i := range mls {
		n[i] = LineStringToPlanar(mls[i], proj)
	}

	return n
}

// MultiLineStringToGeo is a helper to project an entire multi linestring.
func MultiLineStringToGeo(mls planar.MultiLineString, proj *Projection) geo.MultiLineString {
	n := make(geo.MultiLineString, len(mls))
	for i := range mls {
		n[i] = LineStringToGeo(mls[i], proj)
	}

	return n
}

// PolygonToPlanar is a helper to project an entire polygon.
func PolygonToPlanar(p geo.Polygon, proj *Projection) planar.Polygon {
	n := make(planar.Polygon, len(p), len(p))
	for i := range p {
		n[i] = LineStringToPlanar(p[i], proj)
	}

	return n
}

// PolygonToGeo is a helper to project an entire line string.
func PolygonToGeo(p planar.Polygon, proj *Projection) geo.Polygon {
	n := make(geo.Polygon, len(p), len(p))
	for i := range p {
		n[i] = LineStringToGeo(p[i], proj)
	}

	return n
}

// MultiPolygonToPlanar is a helper to project an entire multi polygon.
func MultiPolygonToPlanar(mp geo.MultiPolygon, proj *Projection) planar.MultiPolygon {
	n := make(planar.MultiPolygon, len(mp))
	for i := range mp {
		n[i] = PolygonToPlanar(mp[i], proj)
	}

	return n
}

// MultiPolygonToGeo is a helper to project an entire multi linestring.
func MultiPolygonToGeo(mp planar.MultiPolygon, proj *Projection) geo.MultiPolygon {
	n := make(geo.MultiPolygon, len(mp))
	for i := range mp {
		n[i] = PolygonToGeo(mp[i], proj)
	}

	return n
}

// BoundToPlanar is a helper to project a rectangle.
func BoundToPlanar(bound geo.Bound, proj *Projection) planar.Bound {
	return planar.NewBoundFromPoints(
		proj.ToPlanar(bound[0]),
		proj.ToPlanar(bound[1]),
	)
}

// BoundToGeo is a helper to project a rectangle.
func BoundToGeo(bound planar.Bound, proj *Projection) geo.Bound {
	return geo.NewBoundFromPoints(
		proj.ToGeo(bound[0]),
		proj.ToGeo(bound[1]),
	)
}
