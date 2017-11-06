package project

import "github.com/paulmach/orb"

// ToPlanar projects a geometry from geo -> planar
func ToPlanar(g orb.Geometry, proj *Projection) orb.Geometry {
	switch g := g.(type) {
	case orb.Point:
		return proj.ToPlanar(g)
	case orb.MultiPoint:
		return MultiPointToPlanar(g, proj)
	case orb.LineString:
		return LineStringToPlanar(g, proj)
	case orb.MultiLineString:
		return MultiLineStringToPlanar(g, proj)
	case orb.Ring:
		return RingToPlanar(g, proj)
	case orb.Polygon:
		return PolygonToPlanar(g, proj)
	case orb.MultiPolygon:
		return MultiPolygonToPlanar(g, proj)
	case orb.Collection:
		return CollectionToPlanar(g, proj)
	case orb.Bound:
		return BoundToPlanar(g, proj)
	}

	panic("geometry type not supported")
}

// ToGeo projects a geometry from planar -> geo
func ToGeo(g orb.Geometry, proj *Projection) orb.Geometry {
	switch g := g.(type) {
	case orb.Point:
		return proj.ToGeo(g)
	case orb.MultiPoint:
		return MultiPointToGeo(g, proj)
	case orb.LineString:
		return LineStringToGeo(g, proj)
	case orb.MultiLineString:
		return MultiLineStringToGeo(g, proj)
	case orb.Ring:
		return RingToGeo(g, proj)
	case orb.Polygon:
		return PolygonToGeo(g, proj)
	case orb.MultiPolygon:
		return MultiPolygonToGeo(g, proj)
	case orb.Collection:
		return CollectionToGeo(g, proj)
	case orb.Bound:
		return BoundToGeo(g, proj)
	}

	panic("geometry type not supported")
}

// MultiPointToPlanar is a helper to project an entire multi point.
func MultiPointToPlanar(mp orb.MultiPoint, proj *Projection) orb.MultiPoint {
	n := make(orb.MultiPoint, len(mp))
	for i := range mp {
		n[i] = proj.ToPlanar(mp[i])
	}

	return n
}

// MultiPointToGeo is a helper to project an entire multi point.
func MultiPointToGeo(mp orb.MultiPoint, proj *Projection) orb.MultiPoint {
	n := make(orb.MultiPoint, len(mp))
	for i := range mp {
		n[i] = proj.ToGeo(mp[i])
	}

	return n
}

// LineStringToPlanar is a helper to project an entire line string.
func LineStringToPlanar(ls orb.LineString, proj *Projection) orb.LineString {
	return orb.LineString(MultiPointToPlanar(orb.MultiPoint(ls), proj))
}

// LineStringToGeo is a helper to project an entire line string.
func LineStringToGeo(ls orb.LineString, proj *Projection) orb.LineString {
	return orb.LineString(MultiPointToGeo(orb.MultiPoint(ls), proj))
}

// MultiLineStringToPlanar is a helper to project an entire multi linestring.
func MultiLineStringToPlanar(mls orb.MultiLineString, proj *Projection) orb.MultiLineString {
	n := make(orb.MultiLineString, len(mls))
	for i := range mls {
		n[i] = LineStringToPlanar(mls[i], proj)
	}

	return n
}

// MultiLineStringToGeo is a helper to project an entire multi linestring.
func MultiLineStringToGeo(mls orb.MultiLineString, proj *Projection) orb.MultiLineString {
	n := make(orb.MultiLineString, len(mls))
	for i := range mls {
		n[i] = LineStringToGeo(mls[i], proj)
	}

	return n
}

// RingToPlanar is a helper to project an entire ring.
func RingToPlanar(r orb.Ring, proj *Projection) orb.Ring {
	return orb.Ring(LineStringToPlanar(orb.LineString(r), proj))
}

// RingToGeo is a helper to project an entire ring.
func RingToGeo(r orb.Ring, proj *Projection) orb.Ring {
	return orb.Ring(LineStringToGeo(orb.LineString(r), proj))
}

// PolygonToPlanar is a helper to project an entire polygon.
func PolygonToPlanar(p orb.Polygon, proj *Projection) orb.Polygon {
	n := make(orb.Polygon, len(p))
	for i := range p {
		n[i] = RingToPlanar(p[i], proj)
	}

	return n
}

// PolygonToGeo is a helper to project an entire line string.
func PolygonToGeo(p orb.Polygon, proj *Projection) orb.Polygon {
	n := make(orb.Polygon, len(p))
	for i := range p {
		n[i] = RingToGeo(p[i], proj)
	}

	return n
}

// MultiPolygonToPlanar is a helper to project an entire multi polygon.
func MultiPolygonToPlanar(mp orb.MultiPolygon, proj *Projection) orb.MultiPolygon {
	n := make(orb.MultiPolygon, len(mp))
	for i := range mp {
		n[i] = PolygonToPlanar(mp[i], proj)
	}

	return n
}

// MultiPolygonToGeo is a helper to project an entire multi linestring.
func MultiPolygonToGeo(mp orb.MultiPolygon, proj *Projection) orb.MultiPolygon {
	n := make(orb.MultiPolygon, len(mp))
	for i := range mp {
		n[i] = PolygonToGeo(mp[i], proj)
	}

	return n
}

// CollectionToPlanar is a helper to project a rectangle.
func CollectionToPlanar(c orb.Collection, proj *Projection) orb.Collection {
	n := make(orb.Collection, len(c))
	for i := range c {
		n[i] = ToPlanar(c[i], proj)
	}

	return n
}

// CollectionToGeo is a helper to project a rectangle.
func CollectionToGeo(c orb.Collection, proj *Projection) orb.Collection {
	n := make(orb.Collection, len(c))
	for i := range c {
		n[i] = ToGeo(c[i], proj)
	}

	return n
}

// BoundToPlanar is a helper to project a rectangle.
func BoundToPlanar(bound orb.Bound, proj *Projection) orb.Bound {
	return orb.NewBoundFromPoints(
		proj.ToPlanar(bound[0]),
		proj.ToPlanar(bound[1]),
	)
}

// BoundToGeo is a helper to project a rectangle.
func BoundToGeo(bound orb.Bound, proj *Projection) orb.Bound {
	return orb.NewBoundFromPoints(
		proj.ToGeo(bound[0]),
		proj.ToGeo(bound[1]),
	)
}
