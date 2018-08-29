package geojson

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/paulmach/orb"
)

func TestNewFeatureCollection(t *testing.T) {
	fc := NewFeatureCollection()

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}
}

func TestUnmarshalFeatureCollection(t *testing.T) {
	rawJSON := `
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      },
	      { "type": "Feature",
	        "geometry": {
	          "type": "LineString",
	          "coordinates": [
	            [102.0, 0.0], [103.0, 1.0], [104.0, 0.0], [105.0, 1.0]
	            ]
	          },
	        "properties": {
	          "prop0": "value0",
	          "prop1": 0.0
	        }
	      },
	      { "type": "Feature",
	         "geometry": {
	           "type": "Polygon",
	           "coordinates": [
	             [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0],
	               [100.0, 1.0], [100.0, 0.0] ]
	             ]
	         },
	         "properties": {
	           "prop0": "value0",
	           "prop1": {"this": "that"}
	         }
	       }
	     ]
	  }`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature collection without issue, err %v", err)
	}

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}

	if len(fc.Features) != 3 {
		t.Errorf("should have 3 features but got %d", len(fc.Features))
	}

	f := fc.Features[0]
	if gt := f.Geometry.GeoJSONType(); gt != "Point" {
		t.Errorf("incorrect feature type: %v != %v", gt, "Point")
	}

	f = fc.Features[1]
	if gt := f.Geometry.GeoJSONType(); gt != "LineString" {
		t.Errorf("incorrect feature type: %v != %v", gt, "LineString")
	}

	f = fc.Features[2]
	if gt := f.Geometry.GeoJSONType(); gt != "Polygon" {
		t.Errorf("incorrect feature type: %v != %v", gt, "Polygon")
	}

	// check unmarshal/marshal loop
	var expected interface{}
	err = json.Unmarshal([]byte(rawJSON), &expected)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	data, err := json.MarshalIndent(fc, "", " ")
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	var raw interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(raw, expected) {
		t.Errorf("invalid marshalling: \n%v", string(data))
	}

	// not a feature collection
	data, _ = NewFeature(orb.Point{}).MarshalJSON()
	_, err = UnmarshalFeatureCollection(data)
	if err == nil {
		t.Error("should return error if not a feature collection")
	}

	if !strings.Contains(err.Error(), "not a feature collection") {
		t.Errorf("incorrect error: %v", err)
	}

	// invalid json
	_, err = UnmarshalFeatureCollection([]byte(`{"type": "FeatureCollection",`)) // truncated
	if err == nil {
		t.Errorf("should return error for invalid json")
	}
}

func TestFeatureCollectionMarshalJSON(t *testing.T) {
	fc := NewFeatureCollection()
	blob, err := fc.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}

func TestFeatureCollectionMarshal(t *testing.T) {
	fc := NewFeatureCollection()
	fc.Features = nil
	blob, err := json.Marshal(fc)

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}

func TestFeatureCollectionMarshal_BBox(t *testing.T) {
	fc := NewFeatureCollection()

	// nil bbox
	fc.BBox = nil
	blob, err := json.Marshal(fc)
	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if bytes.Contains(blob, []byte(`"bbox"`)) {
		t.Errorf("should not contain bbox attribute if empty")
	}

	// with a bbox
	fc.BBox = []float64{1, 2, 3, 4}
	blob, err = json.Marshal(fc)
	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"bbox":[1,2,3,4]`)) {
		t.Errorf("did not marshal bbox correctly: %v", string(blob))
	}
}

func TestFeatureCollectionMarshalValue(t *testing.T) {
	fc := NewFeatureCollection()
	fc.Features = nil
	blob, err := json.Marshal(*fc)

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}
