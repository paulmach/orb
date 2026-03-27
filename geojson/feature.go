package geojson

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// A FeatureOf corresponds to GeoJSON feature object but allows for a generic type for the properties.
// This allows users to unmarshal into a struct instead of a map if they choose.
//
// The code assumes type of P is a struct, or map as the GeoJSON spec requires it
// marshal into the a json object.
type FeatureOf[P any] struct {
	ID         any          `json:"id,omitempty"`
	Type       string       `json:"type"`
	BBox       BBox         `json:"bbox,omitempty"`
	Geometry   orb.Geometry `json:"geometry"`
	Properties P            `json:"properties"`

	// ExtraMembers can be used to encoded/decode extra key/members in
	// the base of the feature object. Note that keys of "id", "type", "bbox"
	// "geometry" and "properties" will not work as those are reserved by the
	// GeoJSON spec.
	ExtraMembers Properties `json:"-"`
}

// A Feature corresponds to GeoJSON feature object.
type Feature = FeatureOf[Properties]

// NewFeature creates and initializes a GeoJSON feature given the required attributes.
func NewFeature(geometry orb.Geometry) *Feature {
	return &Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(map[string]any),
	}
}

// Point implements the orb.Pointer interface so that Features can be used
// with quadtrees. The point returned is the center of the Bound of the geometry.
// To represent the geometry with another point you must create a wrapper type.
func (f *FeatureOf[P]) Point() orb.Point {
	return f.Geometry.Bound().Center()
}

var _ orb.Pointer = &FeatureOf[any]{}

// MarshalJSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
// Items in the ExtraMembers map will be included in the base of the
// feature object.
func (f FeatureOf[P]) MarshalJSON() ([]byte, error) {
	jProperties, err := f.jsonProperties()
	if err != nil {
		return nil, err
	}
	return marshalJSON(f.newFeatureDoc(jProperties))
}

// MarshalBSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
// Items in the ExtraMembers map will be included in the base of the
// feature object.
func (f FeatureOf[P]) MarshalBSON() ([]byte, error) {
	properties, err := f.bsonProperties()
	if err != nil {
		return nil, err
	}
	return bson.Marshal(f.newFeatureDoc(properties))
}

func (f FeatureOf[P]) jsonProperties() (json.RawMessage, error) {
	jProperties, err := json.Marshal(f.Properties)
	if err != nil {
		return nil, err
	}

	if len(jProperties) <= 2 { // empty
		// we assume it's an object so an empty {} is 2 bytes
		// in that case the properties should be nil according to the geojson spec
		jProperties = nil
	}

	return jProperties, nil
}

func (f FeatureOf[P]) bsonProperties() (any, error) {
	t, value, err := bson.MarshalValue(f.Properties)
	if err != nil {
		return nil, err
	}

	if t == bson.TypeEmbeddedDocument && bytes.Equal(value, []byte{5, 0, 0, 0, 0}) {
		return nil, nil
	}

	return bson.RawValue{Type: t, Value: value}, nil
}

func (f FeatureOf[P]) newFeatureDoc(properties any) any {
	if len(f.ExtraMembers) == 0 {
		return &featureDoc[any]{
			ID:         f.ID,
			Type:       "Feature",
			BBox:       f.BBox,
			Geometry:   NewGeometry(f.Geometry),
			Properties: properties,
		}
	}

	var tmp map[string]any
	if f.ExtraMembers != nil {
		tmp = f.ExtraMembers.Clone()
	} else {
		tmp = make(map[string]any, 3)
	}

	delete(tmp, "id")
	if f.ID != nil {
		tmp["id"] = f.ID
	}
	tmp["type"] = "Feature"

	delete(tmp, "bbox")
	if f.BBox != nil {
		tmp["bbox"] = f.BBox
	}

	tmp["geometry"] = NewGeometry(f.Geometry)
	tmp["properties"] = properties

	return tmp
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
func (f *FeatureOf[P]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`null`)) {
		*f = FeatureOf[P]{}
		return nil
	}

	tmp := make(map[string]nocopyRawMessage, 4)

	err := unmarshalJSON(data, &tmp)
	if err != nil {
		return err
	}

	*f = FeatureOf[P]{}
	for key, value := range tmp {
		switch key {
		case "id":
			err := unmarshalJSON(value, &f.ID)
			if err != nil {
				return err
			}
		case "type":
			err := unmarshalJSON(value, &f.Type)
			if err != nil {
				return err
			}
		case "bbox":
			err := unmarshalJSON(value, &f.BBox)
			if err != nil {
				return err
			}
		case "geometry":
			g := &Geometry{}
			err := unmarshalJSON(value, &g)
			if err != nil {
				return err
			}

			if g != nil {
				f.Geometry = g.Geometry()
			}
		case "properties":
			err := unmarshalJSON(value, &f.Properties)
			if err != nil {
				return err
			}
		default:
			if f.ExtraMembers == nil {
				f.ExtraMembers = Properties{}
			}

			var val any
			err := unmarshalJSON(value, &val)
			if err != nil {
				return err
			}
			f.ExtraMembers[key] = val
		}
	}

	if f.Type != "Feature" {
		return fmt.Errorf("geojson: not a feature: type=%s", f.Type)
	}

	return nil
}

// UnmarshalBSON will unmarshal a BSON document created with bson.Marshal.
func (f *FeatureOf[P]) UnmarshalBSON(data []byte) error {
	tmp := make(map[string]bson.RawValue, 4)

	err := bson.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	*f = FeatureOf[P]{}
	for key, value := range tmp {
		switch key {
		case "id":
			err := value.Unmarshal(&f.ID)
			if err != nil {
				return err
			}
		case "type":
			f.Type, _ = bson.RawValue(value).StringValueOK()
		case "bbox":
			err := value.Unmarshal(&f.BBox)
			if err != nil {
				return err
			}
		case "geometry":
			g := &Geometry{}
			err := value.Unmarshal(&g)
			if err != nil {
				return err
			}

			if g != nil {
				f.Geometry = g.Geometry()
			}
		case "properties":
			err := value.Unmarshal(&f.Properties)
			if err != nil {
				return err
			}
		default:
			if f.ExtraMembers == nil {
				f.ExtraMembers = Properties{}
			}

			var val any
			err := value.Unmarshal(&val)
			if err != nil {
				return err
			}
			f.ExtraMembers[key] = val
		}
	}

	if f.Type != "Feature" {
		return fmt.Errorf("geojson: not a feature: type=%s", f.Type)
	}

	return nil
}

type featureDoc[P any] struct {
	ID         any       `json:"id,omitempty" bson:"id"`
	Type       string    `json:"type" bson:"type"`
	BBox       BBox      `json:"bbox,omitempty" bson:"bbox,omitempty"`
	Geometry   *Geometry `json:"geometry" bson:"geometry"`
	Properties P         `json:"properties" bson:"properties"`
}
