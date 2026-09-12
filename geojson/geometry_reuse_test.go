package geojson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGeometryUnmarshalReuse(t *testing.T) {
	codecs := []struct {
		name      string
		marshal   func(interface{}) ([]byte, error)
		unmarshal func([]byte, interface{}) error
	}{
		{"JSON", json.Marshal, json.Unmarshal},
		{"BSON", bson.Marshal, bson.Unmarshal},
	}
	point := orb.Point{1, 2}
	collection := orb.Collection{orb.Point{3, 4}}
	for _, codec := range codecs {
		for _, initial := range []orb.Geometry{point, collection} {
			t.Run(codec.name+"/"+initial.GeoJSONType(), func(t *testing.T) {
				var expected orb.Geometry = point
				if _, ok := initial.(orb.Point); ok {
					expected = collection
				}
				data, err := codec.marshal(NewGeometry(expected))
				if err != nil {
					t.Fatal(err)
				}
				reused := NewGeometry(initial)
				if err := codec.unmarshal(data, reused); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(reused, NewGeometry(expected)) {
					t.Fatalf("reused geometry = %#v, want %#v", reused, NewGeometry(expected))
				}
				if !reflect.DeepEqual(reused.Geometry(), expected) {
					t.Errorf("Geometry() = %#v, want %#v", reused.Geometry(), expected)
				}
			})
		}
	}
}
