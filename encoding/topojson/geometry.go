package topojson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/paulmach/orb/geojson"
)

type Geometry struct {
	ID         string                 `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	BBox       []float64              `json:"bbox,omitempty"`

	Point           []float64
	MultiPoint      [][]float64
	LineString      []int
	MultiLineString [][]int
	Polygon         [][]int
	MultiPolygon    [][][]int
	Geometries      []*Geometry
}

// MarshalJSON converts the geometry object into the correct JSON.
// This fulfills the json.Marshaler interface.
func (g *Geometry) MarshalJSON() ([]byte, error) {
	// defining a struct here lets us define the order of the JSON elements.
	type geometry struct {
		ID          string                 `json:"id,omitempty"`
		Type        string                 `json:"type"`
		Properties  map[string]interface{} `json:"properties"`
		BBox        []float64              `json:"bbox,omitempty"`
		Coordinates interface{}            `json:"coordinates,omitempty"`
		Arcs        interface{}            `json:"arcs,omitempty"`
		Geometries  interface{}            `json:"geometries,omitempty"`
	}

	geo := &geometry{
		ID:         g.ID,
		Type:       g.Type,
		Properties: g.Properties,
		BBox:       g.BBox,
	}

	switch g.Type {
	case geojson.TypePoint:
		geo.Coordinates = g.Point
	case geojson.TypeMultiPoint:
		geo.Coordinates = g.MultiPoint
	case geojson.TypeLineString:
		geo.Arcs = g.LineString
	case geojson.TypeMultiLineString:
		geo.Arcs = g.MultiLineString
	case geojson.TypePolygon:
		geo.Arcs = g.Polygon
	case geojson.TypeMultiPolygon:
		geo.Arcs = g.MultiPolygon
	default:
		geo.Geometries = g.Geometries
	}

	return json.Marshal(geo)
}

// UnmarshalJSON decodes the data into a TopoJSON geometry.
// This fulfills the json.Unmarshaler interface.
func (g *Geometry) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return decodeGeometry(g, object)
}

func decodeGeometry(g *Geometry, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = string(s)
	} else {
		return errors.New("type property not string")
	}

	if p, ok := object["properties"].(map[string]interface{}); ok {
		g.Properties = p
	}

	if id, ok := object["id"].(string); ok {
		g.ID = id
	}

	var err error
	switch g.Type {
	case geojson.TypePoint:
		g.Point, err = decodePosition(object["coordinates"])
	case geojson.TypeMultiPoint:
		g.MultiPoint, err = decodePositionSet(object["coordinates"])
	case geojson.TypeLineString:
		g.LineString, err = decodeArcs(object["arcs"])
	case geojson.TypeMultiLineString:
		g.MultiLineString, err = decodeArcsSet(object["arcs"])
	case geojson.TypePolygon:
		g.Polygon, err = decodeArcsSet(object["arcs"])
	case geojson.TypeMultiPolygon:
		g.MultiPolygon, err = decodePolygonArcs(object["arcs"])
	default:
		g.Geometries, err = decodeGeometries(object["geometries"])
	}

	return err
}

func decodePosition(data interface{}) ([]float64, error) {
	coords, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid position, got %v", data)
	}

	result := make([]float64, 0, len(coords))
	for _, coord := range coords {
		if f, ok := coord.(float64); ok {
			result = append(result, f)
		} else {
			return nil, fmt.Errorf("not a valid coordinate, got %v", coord)
		}
	}

	return result, nil
}

func decodePositionSet(data interface{}) ([][]float64, error) {
	points, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %v", data)
	}

	result := make([][]float64, 0, len(points))
	for _, point := range points {
		if p, err := decodePosition(point); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodeArcs(data interface{}) ([]int, error) {
	arcs, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of arcs, got %v", data)
	}

	result := make([]int, 0, len(arcs))
	for _, arc := range arcs {
		if i, ok := arc.(int); ok {
			result = append(result, i)
		} else if i, ok := arc.(float64); ok {
			result = append(result, int(i))
		} else {
			return nil, fmt.Errorf("not a valid arc index, got %#v", arc)
		}
	}

	return result, nil
}

func decodeArcsSet(data interface{}) ([][]int, error) {
	sets, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of arcs, got %v", data)
	}

	result := make([][]int, 0, len(sets))
	for _, arcs := range sets {
		if s, err := decodeArcs(arcs); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePolygonArcs(data interface{}) ([][][]int, error) {
	rings, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of rings, got %v", data)
	}

	result := make([][][]int, 0, len(rings))
	for _, sets := range rings {
		if s, err := decodeArcsSet(sets); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodeGeometries(data interface{}) ([]*Geometry, error) {
	if vs, ok := data.([]interface{}); ok {
		geometries := make([]*Geometry, 0, len(vs))
		for _, v := range vs {
			g := &Geometry{}

			vmap, ok := v.(map[string]interface{})
			if !ok {
				break
			}

			err := decodeGeometry(g, vmap)
			if err != nil {
				return nil, err
			}

			geometries = append(geometries, g)
		}

		if len(geometries) == len(vs) {
			return geometries, nil
		}
	}

	return nil, fmt.Errorf("not a valid set of geometries, got %v", data)
}
