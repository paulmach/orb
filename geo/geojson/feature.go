package geojson

import (
	"encoding/json"
	"errors"

	"github.com/paulmach/orb/geo"
)

// A Feature corresponds to GeoJSON feature object
type Feature struct {
	ID         interface{}  `json:"id,omitempty"`
	Type       string       `json:"type"`
	Geometry   geo.Geometry `json:"geometry"`
	Properties Properties   `json:"properties"`
}

// NewFeature creates and initializes a GeoJSON feature given the required attributes.
func NewFeature(geometry geo.Geometry) *Feature {
	return &Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(map[string]interface{}),
	}
}

// MarshalJSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
func (f Feature) MarshalJSON() ([]byte, error) {
	jf := &jsonFeature{
		ID:         f.ID,
		Type:       "Feature",
		Properties: f.Properties,
	}

	if len(jf.Properties) == 0 {
		jf.Properties = nil
	}

	if f.Geometry != nil {
		var (
			coords []byte
			err    error
		)

		if ring, ok := f.Geometry.(geo.Ring); ok {
			coords, err = json.Marshal(geo.Polygon{ring})
		} else {
			coords, err = json.Marshal(f.Geometry)
		}
		if err != nil {
			return nil, err
		}

		jf.Geometry = &jsonGeometry{
			Type:        f.Geometry.GeoJSONType(),
			Coordinates: coords,
		}
	}

	return json.Marshal(jf)
}

// UnmarshalFeature decodes the data into a GeoJSON feature.
// Alternately one can call json.Unmarshal(f) directly for the same result.
func UnmarshalFeature(data []byte) (*Feature, error) {
	f := &Feature{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalJSON handles the correct unmarshalling of the data
// into the geo.Geometry types.
func (f *Feature) UnmarshalJSON(data []byte) error {
	jf := &jsonFeature{}
	err := json.Unmarshal(data, &jf)
	if err != nil {
		return err
	}

	*f = Feature{
		ID:         jf.ID,
		Type:       jf.Type,
		Properties: jf.Properties,
	}

	switch jf.Geometry.Type {
	case "Point":
		p := geo.Point{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &p)
		f.Geometry = p
	case "MultiPoint":
		mp := geo.MultiPoint{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &mp)
		f.Geometry = mp
	case "LineString":
		ls := geo.LineString{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &ls)
		f.Geometry = ls
	case "MultiLineString":
		mls := geo.MultiLineString{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &mls)
		f.Geometry = mls
	case "Polygon":
		p := geo.Polygon{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &p)
		f.Geometry = p
	case "MultiPolygon":
		mp := geo.MultiPolygon{}
		err = json.Unmarshal(jf.Geometry.Coordinates, &mp)
		f.Geometry = mp
	default:
		return errors.New("geojson: invalid geometry")
	}

	return err
}

type jsonFeature struct {
	ID         interface{}   `json:"id,omitempty"`
	Type       string        `json:"type"`
	Geometry   *jsonGeometry `json:"geometry"`
	Properties Properties    `json:"properties"`
}

type jsonGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}
