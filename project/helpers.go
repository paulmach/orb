package project

import "github.com/paulmach/orb"

// ToPlanar projects a geometry from geo -> planar
func ToPlanar(g orb.Geometry, proj *Projection) orb.Geometry {
	if g == nil {
		return nil
	}

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
	if g == nil {
		return nil
	}

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
	for i := range mp {
		mp[i] = proj.ToPlanar(mp[i])
	}

	return mp
}

// MultiPointToGeo is a helper to project an entire multi point.
func MultiPointToGeo(mp orb.MultiPoint, proj *Projection) orb.MultiPoint {
	for i := range mp {
		mp[i] = proj.ToGeo(mp[i])
	}

	return mp
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
	for i := range mls {
		mls[i] = LineStringToPlanar(mls[i], proj)
	}

	return mls
}

// MultiLineStringToGeo is a helper to project an entire multi linestring.
func MultiLineStringToGeo(mls orb.MultiLineString, proj *Projection) orb.MultiLineString {
	for i := range mls {
		mls[i] = LineStringToGeo(mls[i], proj)
	}

	return mls
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
	for i := range p {
		p[i] = RingToPlanar(p[i], proj)
	}

	return p
}

// PolygonToGeo is a helper to project an entire line string.
func PolygonToGeo(p orb.Polygon, proj *Projection) orb.Polygon {
	for i := range p {
		p[i] = RingToGeo(p[i], proj)
	}

	return p
}

// MultiPolygonToPlanar is a helper to project an entire multi polygon.
func MultiPolygonToPlanar(mp orb.MultiPolygon, proj *Projection) orb.MultiPolygon {
	for i := range mp {
		mp[i] = PolygonToPlanar(mp[i], proj)
	}

	return mp
}

// MultiPolygonToGeo is a helper to project an entire multi linestring.
func MultiPolygonToGeo(mp orb.MultiPolygon, proj *Projection) orb.MultiPolygon {
	for i := range mp {
		mp[i] = PolygonToGeo(mp[i], proj)
	}

	return mp
}

// CollectionToPlanar is a helper to project a rectangle.
func CollectionToPlanar(c orb.Collection, proj *Projection) orb.Collection {
	for i := range c {
		c[i] = ToPlanar(c[i], proj)
	}

	return c
}

// CollectionToGeo is a helper to project a rectangle.
func CollectionToGeo(c orb.Collection, proj *Projection) orb.Collection {
	for i := range c {
		c[i] = ToGeo(c[i], proj)
	}

	return c
}

// BoundToPlanar is a helper to project a rectangle.
func BoundToPlanar(bound orb.Bound, proj *Projection) orb.Bound {
	min := proj.ToPlanar(bound.Min)
	return orb.Bound{Min: min, Max: min}.Extend(proj.ToPlanar(bound.Max))
}

// BoundToGeo is a helper to project a rectangle.
func BoundToGeo(bound orb.Bound, proj *Projection) orb.Bound {
	min := proj.ToGeo(bound.Min)
	return orb.Bound{Min: min, Max: min}.Extend(proj.ToGeo(bound.Max))
}
