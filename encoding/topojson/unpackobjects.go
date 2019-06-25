package topojson

import (
	"github.com/paulmach/orb/geojson"
)

type arcEntry struct {
	Start int
	End   int
}

func (t *Topology) unpackObjects() {
	for _, o := range t.objects {
		obj := t.unpackObject(o)
		t.Objects[obj.ID] = obj
	}
	t.objects = nil
	t.deletedArcs = nil
	t.shiftArcs = nil
	t.arcIndexes = nil
}

func (t *Topology) unpackObject(o *topologyObject) *Geometry {
	obj := &Geometry{
		ID:         o.ID,
		Type:       o.Type,
		Properties: o.Properties,
		BBox:       o.BBox,
	}

	switch o.Type {
	default:
		for _, geom := range o.Geometries {
			obj.Geometries = append(obj.Geometries, t.unpackObject(geom))
		}
	case geojson.TypeLineString:
		obj.LineString = t.lookupArc(o.Arc)
	case geojson.TypeMultiLineString:
		obj.MultiLineString = t.lookupArcs(o.Arcs)
	case geojson.TypePolygon:
		obj.Polygon = t.lookupArcs(o.Arcs)
	case geojson.TypeMultiPolygon:
		obj.MultiPolygon = t.lookupMultiArcs(o.MultiArcs)
	case geojson.TypePoint:
		obj.Point = o.Point
	case geojson.TypeMultiPoint:
		obj.MultiPoint = o.MultiPoint
	}

	return obj
}

func (t *Topology) lookupArc(a *arc) []int {
	result := make([]int, 0)

	for a != nil {
		if a.Start < a.End {
			index := t.arcIndexes[arcEntry{a.Start, a.End}]
			if !t.deletedArcs[index] {
				result = append(result, index-t.shiftArcs[index])
			}
		} else {
			index := t.arcIndexes[arcEntry{a.End, a.Start}]
			if !t.deletedArcs[index] {
				result = append(result, ^(index - t.shiftArcs[index]))
			}
		}
		a = a.Next
	}

	return result
}

func (t *Topology) lookupArcs(a []*arc) [][]int {
	result := make([][]int, 0)
	for _, arc := range a {
		result = append(result, t.lookupArc(arc))
	}
	return result
}

func (t *Topology) lookupMultiArcs(a [][]*arc) [][][]int {
	result := make([][][]int, 0)
	for _, s := range a {
		result = append(result, t.lookupArcs(s))
	}
	return result
}
