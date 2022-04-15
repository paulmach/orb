package geojson

import (
	"fmt"

	"github.com/paulmach/orb"
)

// A Feature corresponds to GeoJSON feature object
type Feature struct {
	ID         interface{}  `json:"id,omitempty"`
	Type       string       `json:"type"`
	BBox       BBox         `json:"bbox,omitempty"`
	Geometry   orb.Geometry `json:"geometry"`
	Properties Properties   `json:"properties"`
}

// NewFeature creates and initializes a GeoJSON feature given the required attributes.
func NewFeature(geometry orb.Geometry) *Feature {
	return &Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(map[string]interface{}),
	}
}

// Point implements the orb.Pointer interface so that Features can be used
// with quadtrees. The point returned is the center of the Bound of the geometry.
// To represent the geometry with another point you must create a wrapper type.
func (f *Feature) Point() orb.Point {
	return f.Geometry.Bound().Center()
}

var _ orb.Pointer = &Feature{}

// MarshalJSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
func (f Feature) MarshalJSON() ([]byte, error) {
	jf := &jsonFeature{
		ID:         f.ID,
		Type:       "Feature",
		Properties: f.Properties,
		BBox:       f.BBox,
		Geometry:   NewGeometry(f.Geometry),
	}

	if len(jf.Properties) == 0 {
		jf.Properties = nil
	}

	return marshalJSON(jf)
}

// UnmarshalFeature decodes the data into a GeoJSON feature.
// Alternately one can call json.Unmarshal(f) directly for the same result.
func UnmarshalFeature(data []byte) (*Feature, error) {
	f := &Feature{}
	err := f.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalJSON handles the correct unmarshalling of the data
// into the orb.Geometry types.
func (f *Feature) UnmarshalJSON(data []byte) error {
	jf := &jsonFeature{}
	err := unmarshalJSON(data, &jf)
	if err != nil {
		return err
	}

	if jf.Type != "Feature" {
		return fmt.Errorf("geojson: not a feature: type=%s", jf.Type)
	}

	var g orb.Geometry
	if jf.Geometry != nil {
		if jf.Geometry.Coordinates == nil && jf.Geometry.Geometries == nil {
			return ErrInvalidGeometry
		}
		g = jf.Geometry.Geometry()
	}

	*f = Feature{
		ID:         jf.ID,
		Type:       jf.Type,
		Properties: jf.Properties,
		BBox:       jf.BBox,
		Geometry:   g,
	}

	return nil
}

type jsonFeature struct {
	ID         interface{} `json:"id,omitempty"`
	Type       string      `json:"type"`
	BBox       BBox        `json:"bbox,omitempty"`
	Geometry   *Geometry   `json:"geometry"`
	Properties Properties  `json:"properties"`
}
