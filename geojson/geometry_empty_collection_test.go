package geojson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGeometryCollectionMarshal(t *testing.T) {
	decoded, err := UnmarshalGeometry([]byte(`{"type":"GeometryCollection","geometries":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	empty := &Geometry{Geometries: []*Geometry{}}
	cases := []struct {
		name string
		geom *Geometry
		want string
	}{
		{"constructor", NewGeometry(orb.Collection{}), `{"type":"GeometryCollection","geometries":[]}`},
		{"coordinates", &Geometry{Coordinates: orb.Collection{}}, `{"type":"GeometryCollection","geometries":[]}`},
		{"decoded", decoded, `{"type":"GeometryCollection","geometries":[]}`},
		{"explicit empty", empty, `{"type":"GeometryCollection","geometries":[]}`},
		{"explicit nonempty", &Geometry{Geometries: []*Geometry{NewGeometry(orb.Point{1, 2})}}, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}]}`},
		{"nested empty", &Geometry{Geometries: []*Geometry{empty}}, `{"type":"GeometryCollection","geometries":[{"type":"GeometryCollection","geometries":[]}]}`},
		{"point coordinates win", &Geometry{Type: "GeometryCollection", Coordinates: orb.Point{1, 2}, Geometries: []*Geometry{}}, `{"type":"Point","coordinates":[1,2]}`},
		{"collection coordinates win", &Geometry{Coordinates: orb.Collection{orb.Point{1, 2}}, Geometries: []*Geometry{}}, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}]}`},
	}
	codecs := []struct {
		name    string
		marshal func(*Geometry) ([]byte, error)
	}{
		{"JSON", func(g *Geometry) ([]byte, error) { return json.Marshal(g) }},
		{"BSON document", func(g *Geometry) ([]byte, error) { return geometryBSONJSON(bson.Marshal(g)) }},
		{"BSON value", func(g *Geometry) ([]byte, error) {
			kind, data, err := g.MarshalBSONValue()
			if err != nil {
				return nil, err
			}
			if bson.Type(kind) != bson.TypeEmbeddedDocument {
				return nil, fmt.Errorf("geometry value type = %v, want embedded document", bson.Type(kind))
			}
			return geometryBSONJSON(data, nil)
		}},
	}
	for _, codec := range codecs {
		for _, tc := range cases {
			t.Run(codec.name+"/"+tc.name, func(t *testing.T) {
				data, err := codec.marshal(tc.geom)
				if err != nil {
					t.Fatal(err)
				}
				var got, want any
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("marshal = %s, want %s", data, tc.want)
				}
				// Check the wire shape before decoding nested collections into orb.
				if _, err := UnmarshalGeometry(data); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func geometryBSONJSON(data []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	var document bson.D
	if err := bson.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	return bson.MarshalExtJSON(document, false, false)
}

func TestGeometryMarshalZeroValue(t *testing.T) {
	g := &Geometry{}
	data, err := json.Marshal(g)
	if err != nil || string(data) != "null" {
		t.Fatalf("zero JSON = %s, %v; want null", data, err)
	}
	kind, data, err := g.MarshalBSONValue()
	if err != nil || bson.Type(kind) != bson.TypeNull || data != nil {
		t.Fatalf("zero BSON value = %v, %v, %v; want null", bson.Type(kind), data, err)
	}
}

func TestFeatureEmptyCollectionRoundTrip(t *testing.T) {
	codecs := []struct {
		name      string
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{
		{"JSON", json.Marshal, json.Unmarshal},
		{"BSON", bson.Marshal, bson.Unmarshal},
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(NewFeature(orb.Collection{}))
			if err != nil {
				t.Fatal(err)
			}
			var got Feature
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			collection, ok := got.Geometry.(orb.Collection)
			if !ok || collection == nil || len(collection) != 0 {
				t.Fatalf("feature geometry = %#v; want an empty collection", got.Geometry)
			}
		})
	}
}
