package wkb

import (
	"fmt"
	"strings"

	"github.com/paulmach/orb"
)

/*
This purpose of this file is to house the wkt functions. These functions are
use to take a tagola.Geometry and convert it to a wkt string. It will, also,
contain functions to parse a wkt string into a wkb.Geometry.
*/

func wkt(geom orb.Geometry) string {
	switch g := geom.(type) {
	case orb.Point:
		return fmt.Sprintf("%v %v", g.X(), g.Y())

	case orb.MultiPoint:
		var points []string
		for _, p := range g {
			points = append(points, wkt(p))
		}
		return "(" + strings.Join(points, ",") + ")"

	case orb.LineString:
		var points []string
		for _, p := range g {
			points = append(points, wkt(p))
		}
		return "(" + strings.Join(points, ",") + ")"

	case orb.MultiLineString:
		var lines []string
		for _, l := range g {
			lines = append(lines, wkt(l))
		}
		return "(" + strings.Join(lines, ",") + ")"

	case orb.Polygon:
		var lines []string
		for _, r := range g {
			lines = append(lines, wkt(orb.LineString(r)))
		}
		return "(" + strings.Join(lines, ",") + ")"

	case orb.MultiPolygon:
		var polygons []string
		for _, p := range g {
			polygons = append(polygons, wkt(p))
		}
		return "(" + strings.Join(polygons, ",") + ")"

	}

	panic("unknown geometry")
}

//WKT returns a WKT representation of the Geometry if possible.
// the Error will be non-nil if geometry is unknown.
func WKT(geom orb.Geometry) string {
	switch g := geom.(type) {
	default:
		return ""
	case orb.Point:
		return "POINT (" + wkt(g) + ")"
	// case tegola.Point3: TODO
	// 	// POINT M ( 10 10 10 )
	// 	if g == nil {
	// 		return "POINT M EMPTY"
	// 	}
	// 	return "POINT M (" + wkt(g) + ")"
	case orb.MultiPoint:
		if g == nil {
			return "MULTIPOINT EMPTY"
		}
		return "MULTIPOINT " + wkt(g)
	case orb.LineString:
		if g == nil {
			return "LINESTRING EMPTY"
		}
		return "LINESTRING " + wkt(g)
	case orb.MultiLineString:
		if g == nil {
			return "MULTILINE EMPTY"
		}
		return "MULTILINE " + wkt(g)
	case orb.Polygon:
		if g == nil {
			return "POLYGON EMPTY"
		}
		return "POLYGON " + wkt(g)
	case orb.MultiPolygon:
		if g == nil {
			return "MULTIPOLYGON EMPTY"
		}
		return "MULTIPOLYGON " + wkt(g)
	case orb.Collection:
		if g == nil {
			return "GEOMETRYCOLLECTION EMPTY"

		}
		var geometries []string
		for _, c := range g {
			s := WKT(c)
			geometries = append(geometries, s)
		}
		return "GEOMETRYCOLLECTION (" + strings.Join(geometries, ",") + ")"
	}
}
