package topojson

import (
	"github.com/paulmach/orb"
	geojson "github.com/paulmach/orb/geojson"
)

func (t *Topology) ToGeoJSON() *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()

	for _, obj := range t.Objects {
		switch obj.Type {
		case "GeometryCollection":
			for _, geometry := range obj.Geometries {
				feat := geojson.NewFeature(t.toGeometry(geometry))
				feat.ID = geometry.ID
				feat.Properties = geometry.Properties
				feat.BBox = geometry.BBox
				fc.Append(feat)
			}
		default:
			feat := geojson.NewFeature(t.toGeometry(obj))
			feat.ID = obj.ID
			feat.Properties = obj.Properties
			feat.BBox = obj.BBox
			fc.Append(feat)
		}
	}

	return fc
}

func (t *Topology) toGeometry(g *Geometry) orb.Geometry {
	switch g.Type {
	case geojson.TypePoint:
		return t.packPoint(g.Point)
	case geojson.TypeMultiPoint:
		return t.packPoints(g.MultiPoint)
	case geojson.TypeLineString:
		return t.packLinestring(g.LineString)
	case geojson.TypeMultiLineString:
		return t.packMultiLinestring(g.MultiLineString)
	case geojson.TypePolygon:
		return t.packPolygon(g.Polygon)
	case geojson.TypeMultiPolygon:
		return t.packMultiPolygon(g.MultiPolygon)
	default:
		geometries := make([]orb.Geometry, len(g.Geometries))
		for i, geometry := range g.Geometries {
			geometries[i] = t.toGeometry(geometry)
		}
		return orb.Collection(geometries)
	}
	return nil
}

func (t *Topology) packPoint(in []float64) orb.Geometry {
	if t.Transform == nil {
		return orb.Point{in[0], in[1]}
	}

	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = v
		if i < 2 {
			out[i] = v*t.Transform.Scale[i] + t.Transform.Translate[i]
		}
	}

	return orb.Point{out[0], out[1]}
}

func (t *Topology) packPoints(in [][]float64) orb.Geometry {
	out := make(orb.Collection, len(in))
	for i, p := range in {
		out[i] = t.packPoint(p)
	}
	return out
}

func (t *Topology) packLinestring(ls []int) orb.Geometry {
	result := orb.LineString{}
	for _, a := range ls {
		reverse := false
		if a < 0 {
			a = ^a
			reverse = true
		}
		arc := t.Arcs[a]

		// Copy arc
		newArc := make([][]float64, len(arc))
		for i, point := range arc {
			newArc[i] = append([]float64{}, point...)
		}

		if t.Transform != nil {
			x := float64(0)
			y := float64(0)

			for k, p := range newArc {
				x += p[0]
				y += p[1]

				newArc[k][0] = x*t.Transform.Scale[0] + t.Transform.Translate[0]
				newArc[k][1] = y*t.Transform.Scale[1] + t.Transform.Translate[1]
			}
		}

		if reverse {
			for j := len(newArc) - 1; j >= 0; j-- {
				if len(result) > 0 && pointEquals([]float64{result[len(result)-1][0], result[len(result)-1][1]}, newArc[j]) {
					continue
				}
				result = append(result, orb.Point{newArc[j][0], newArc[j][1]})
			}
		} else {
			for j := 0; j < len(newArc); j++ {
				if len(result) > 0 && pointEquals([]float64{result[len(result)-1][0], result[len(result)-1][1]}, newArc[j]) {
					continue
				}
				result = append(result, orb.Point{newArc[j][0], newArc[j][1]})
			}
		}
	}
	return result
}

func (t *Topology) packMultiLinestring(ls [][]int) orb.Geometry {
	result := make(orb.MultiLineString, len(ls))
	for i, l := range ls {
		result[i] = t.packLinestring(l).(orb.LineString)
	}
	return result
}

func (t *Topology) packPolygon(ls [][]int) orb.Geometry {
	result := make(orb.Polygon, len(ls))
	for i, l := range ls {
		s := t.packLinestring(l).(orb.LineString)
		result[i] = make(orb.Ring, len(s))
		for j, l := range s {
			result[i][j] = l
		}
	}
	return result
}

func (t *Topology) packMultiPolygon(ls [][][]int) orb.Geometry {
	result := make(orb.MultiPolygon, len(ls))
	for i, l := range ls {
		result[i] = t.packPolygon(l).(orb.Polygon)
	}
	return result
}
