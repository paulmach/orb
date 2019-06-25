package topojson

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func (t *Topology) extract() {
	t.objects = make([]*topologyObject, 0, len(t.input))

	for i, g := range t.input {
		feature := t.extractFeature(g)
		if len(feature.ID) == 0 {
			// if multiple features exist without ids only one will be retained, so provide a synthetic id
			feature.ID = fmt.Sprintf("feature_%d", i)
		}
		t.objects = append(t.objects, feature)
	}
	t.input = nil // no longer needed
}

func (t *Topology) extractFeature(f *geojson.Feature) *topologyObject {
	g := f.Geometry
	o := t.extractGeometry(geojson.NewGeometry(g))

	// TODO
	// idProp := "id"
	// if t.opts != nil && t.opts.IDProperty != "" {
	// 	idProp = t.opts.IDProperty
	// }

	if f.ID != nil {
		o.ID = fmt.Sprint(f.ID)
	}
	o.Properties = f.Properties
	o.BBox = f.BBox
	return o
}

func (t *Topology) extractGeometry(g *geojson.Geometry) *topologyObject {
	o := &topologyObject{
		Type: g.Type,
	}

	// TODO
	// if g.Coordinates != nil {
	// 	o.BBox = []float64{
	// 		g.Coordinates.Bound().Min[0],
	// 		g.Coordinates.Bound().Min[1],
	// 		g.Coordinates.Bound().Max[0],
	// 		g.Coordinates.Bound().Max[1],
	// 	}
	// }

	switch g.Type {
	default:
		for _, geom := range g.Geometries {
			o.Geometries = append(o.Geometries, t.extractGeometry(geom))
		}
	case geojson.TypeLineString:
		o.Arc = t.extractLine(g.Coordinates.(orb.LineString))
	case geojson.TypeMultiLineString:
		o.Arcs = make([]*arc, len(g.Coordinates.(orb.MultiLineString)))
		for i, l := range g.Coordinates.(orb.MultiLineString) {
			o.Arcs[i] = t.extractLine(l)
		}
	case geojson.TypePolygon:
		o.Arcs = make([]*arc, len(g.Coordinates.(orb.Polygon)))
		for i, r := range g.Coordinates.(orb.Polygon) {
			o.Arcs[i] = t.extractRing(r)
		}
	case geojson.TypeMultiPolygon:
		o.MultiArcs = make([][]*arc, len(g.Coordinates.(orb.MultiPolygon)))
		for i, p := range g.Coordinates.(orb.MultiPolygon) {
			arcs := make([]*arc, len(p))
			for j, r := range p {
				arcs[j] = t.extractRing(r)
			}
			o.MultiArcs[i] = arcs
		}
	case geojson.TypePoint:
		o.Point = []float64{g.Coordinates.(orb.Point)[0], g.Coordinates.(orb.Point)[1]}
	case geojson.TypeMultiPoint:
		for _, v := range g.Coordinates.(orb.MultiPoint) {
			o.MultiPoint = append(o.MultiPoint, []float64{v[0], v[1]})
		}
	}

	return o
}

func (t *Topology) extractLine(line orb.LineString) *arc {
	n := len(line)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, []float64{line[i][0], line[i][1]})
	}

	index := len(t.coordinates) - 1
	arc := &arc{Start: index - n + 1, End: index}
	t.lines = append(t.lines, arc)

	return arc
}

func (t *Topology) extractRing(ring orb.Ring) *arc {
	n := len(ring)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, []float64{ring[i][0], ring[i][1]})
	}

	index := len(t.coordinates) - 1
	arc := &arc{Start: index - n + 1, End: index}
	t.rings = append(t.rings, arc)

	return arc
}
