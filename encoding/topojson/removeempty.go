package topojson

import (
	geojson "github.com/paulmach/orb/geojson"
)

func (t *Topology) removeEmpty() {
	objs := make(map[string]*Geometry, len(t.Objects))
	for _, o := range t.Objects {
		obj := t.removeEmptyObjects(o)
		if obj != nil {
			objs[obj.ID] = obj
		}
	}
	t.Objects = objs
}

func (t *Topology) removeEmptyObjects(obj *Geometry) *Geometry {
	switch obj.Type {
	case "GeometryCollection":
		geoms := make([]*Geometry, 0, len(obj.Geometries))
		for _, g := range obj.Geometries {
			geom := t.removeEmptyObjects(g)
			if geom != nil {
				geoms = append(geoms, geom)
			}
		}
		if len(geoms) == 0 {
			return nil
		}
		obj.Geometries = geoms
	case geojson.TypeLineString:
		if len(obj.LineString) == 0 {
			return nil
		}
	case geojson.TypeMultiLineString:
		linestrings := make([][]int, 0, len(obj.MultiLineString))
		for _, ls := range obj.MultiLineString {
			if len(ls) > 0 {
				linestrings = append(linestrings, ls)
			}
		}
		if len(linestrings) == 0 {
			return nil
		}

		if len(linestrings) == 1 {
			obj.LineString = linestrings[0]
			obj.MultiLineString = nil
			obj.Type = geojson.TypeLineString
		} else {
			obj.MultiLineString = linestrings
		}
	case geojson.TypePolygon:
		rings := t.removeEmptyPolygon(obj.Polygon)
		if rings == nil {
			return nil
		}
		obj.Polygon = rings
	case geojson.TypeMultiPolygon:
		polygons := make([][][]int, 0, len(obj.MultiPolygon))
		for _, polygon := range obj.MultiPolygon {
			rings := t.removeEmptyPolygon(polygon)
			if rings != nil {
				polygons = append(polygons, rings)
			}
		}
		if len(polygons) == 0 {
			return nil
		}
		if len(polygons) == 1 {
			obj.Polygon = polygons[0]
			obj.MultiPolygon = nil
			obj.Type = geojson.TypePolygon
		} else {
			obj.MultiPolygon = polygons
		}
	}

	return obj
}

func (t *Topology) removeEmptyPolygon(polygon [][]int) [][]int {
	rings := make([][]int, 0, len(polygon))
	for i, ring := range polygon {
		if i == 0 && len(ring) == 0 {
			return nil // Outer ring empty
		}
		if len(ring) > 0 {
			rings = append(rings, ring)
		}
	}
	return rings
}
