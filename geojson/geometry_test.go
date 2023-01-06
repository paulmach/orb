package geojson

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/paulmach/orb"
)

func TestGeometry(t *testing.T) {
	for _, g := range orb.AllGeometries {
		NewGeometry(g)
	}
}

func TestGeometryMarshal(t *testing.T) {
	cases := []struct {
		name    string
		geom    orb.Geometry
		include string
	}{
		{
			name:    "point",
			geom:    orb.Point{},
			include: `"type":"Point"`,
		},
		{
			name:    "multi point",
			geom:    orb.MultiPoint{},
			include: `"type":"MultiPoint"`,
		},
		{
			name:    "linestring",
			geom:    orb.LineString{},
			include: `"type":"LineString"`,
		},
		{
			name:    "multi linestring",
			geom:    orb.MultiLineString{},
			include: `"type":"MultiLineString"`,
		},
		{
			name:    "polygon",
			geom:    orb.Polygon{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "multi polygon",
			geom:    orb.MultiPolygon{},
			include: `"type":"MultiPolygon"`,
		},
		{
			name:    "ring",
			geom:    orb.Ring{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "bound",
			geom:    orb.Bound{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "collection",
			geom:    orb.Collection{orb.LineString{}},
			include: `"type":"GeometryCollection"`,
		},
		{
			name:    "collection2",
			geom:    orb.Collection{orb.Point{}, orb.Point{}},
			include: `"geometries":[`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NewGeometry(tc.geom).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}

			g := &Geometry{Coordinates: tc.geom}
			data, err = g.MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}
		})
	}
}

func TestGeometryUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		geom orb.Geometry
	}{
		{
			name: "point",
			geom: orb.Point{1, 2},
		},
		{
			name: "multi point",
			geom: orb.MultiPoint{{1, 2}, {3, 4}},
		},
		{
			name: "linestring",
			geom: orb.LineString{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name: "multi linestring",
			geom: orb.MultiLineString{},
		},
		{
			name: "polygon",
			geom: orb.Polygon{},
		},
		{
			name: "multi polygon",
			geom: orb.MultiPolygon{},
		},
		{
			name: "collection",
			geom: orb.Collection{orb.LineString{{1, 2}, {3, 4}}, orb.Point{5, 6}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NewGeometry(tc.geom).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			// unmarshal
			g, err := UnmarshalGeometry(data)
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}

			if g.Type != tc.geom.GeoJSONType() {
				t.Errorf("incorrenct type: %v != %v", g.Type, tc.geom.GeoJSONType())
			}

			if !orb.Equal(g.Geometry(), tc.geom) {
				t.Errorf("incorrect geometry")
				t.Logf("%[1]T, %[1]v", g.Geometry())
				t.Log(tc.geom)
			}
		})
	}

	// invalid type
	_, err := UnmarshalGeometry([]byte(`{
		"type": "arc",
		"coordinates": [[0, 0]]
	}`))
	if err == nil {
		t.Errorf("should return error for invalid type")
	}

	if !strings.Contains(err.Error(), "invalid geometry") {
		t.Errorf("incorrect error: %v", err)
	}

	// invalid json
	_, err = UnmarshalGeometry([]byte(`{"type": "arc",`)) // truncated
	if err == nil {
		t.Errorf("should return error for invalid json")
	}

	g := &Geometry{}
	err = g.UnmarshalJSON([]byte(`{"type": "arc",`)) // truncated
	if err == nil {
		t.Errorf("should return error for invalid json")
	}

	// invalid type (null)
	_, err = UnmarshalGeometry([]byte(`null`))
	if err == nil {
		t.Errorf("should return error for invalid type")
	}

	if !strings.Contains(err.Error(), "invalid geometry") {
		t.Errorf("incorrect error: %v", err)
	}
}

func TestGeometryUnmarshal_errors(t *testing.T) {
	cases := []struct {
		name string
		data string
	}{
		{
			name: "point",
			data: `{"type":"Point","coordinates":1}`,
		},
		{
			name: "multi point",
			data: `{"type":"MultiPoint","coordinates":2}`,
		},
		{
			name: "linestring",
			data: `{"type":"LineString","coordinates":3}`,
		},
		{
			name: "multi linestring",
			data: `{"type":"MultiLineString","coordinates":4}`,
		},
		{
			name: "polygon",
			data: `{"type":"Polygon","coordinates":10.2}`,
		},
		{
			name: "multi polygon",
			data: `{"type":"MultiPolygon","coordinates":{}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := UnmarshalGeometry([]byte(tc.data))
			if err == nil {
				t.Errorf("expected error, got nothing")
			}
		})
	}
}

func TestHelperTypes(t *testing.T) {
	// This test makes sure the marshal-unmarshal loop does the same thing.
	// The code and types here are complicated to avoid duplicate code.
	cases := []struct {
		name   string
		geom   orb.Geometry
		helper interface{}
		output interface{}
	}{
		{
			name:   "point",
			geom:   orb.Point{1, 2},
			helper: Point(orb.Point{1, 2}),
			output: &Point{},
		},
		{
			name:   "multi point",
			geom:   orb.MultiPoint{{1, 2}, {3, 4}},
			helper: MultiPoint(orb.MultiPoint{{1, 2}, {3, 4}}),
			output: &MultiPoint{},
		},
		{
			name:   "linestring",
			geom:   orb.LineString{{1, 2}, {3, 4}},
			helper: LineString(orb.LineString{{1, 2}, {3, 4}}),
			output: &LineString{},
		},
		{
			name:   "multi linestring",
			geom:   orb.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			helper: MultiLineString(orb.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}),
			output: &MultiLineString{},
		},
		{
			name:   "polygon",
			geom:   orb.Polygon{{{1, 2}, {3, 4}}},
			helper: Polygon(orb.Polygon{{{1, 2}, {3, 4}}}),
			output: &Polygon{},
		},
		{
			name:   "multi polygon",
			geom:   orb.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}}},
			helper: MultiPolygon(orb.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}}}),
			output: &MultiPolygon{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// check marshalling
			data, err := tc.helper.(json.Marshaler).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			geoData, err := NewGeometry(tc.geom).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !reflect.DeepEqual(data, geoData) {
				t.Errorf("should marshal the same")
				t.Log(string(data))
				t.Log(string(geoData))
			}

			// check unmarshalling
			err = tc.output.(json.Unmarshaler).UnmarshalJSON(data)
			if err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			geo := &Geometry{}
			err = geo.UnmarshalJSON(data)
			if err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !orb.Equal(tc.output.(geom).Geometry(), geo.Coordinates) {
				t.Errorf("should unmarshal the same")
				t.Log(tc.output)
				t.Log(geo.Coordinates)
			}

			// invalid json should return error
			err = tc.output.(json.Unmarshaler).UnmarshalJSON([]byte(`{invalid}`))
			if err == nil {
				t.Errorf("should return error for invalid json")
			}

			// not the correct type should return error.
			// non of they types directly supported are geometry collections.
			data, err = NewGeometry(orb.Collection{orb.Point{}}).MarshalJSON()
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}

			err = tc.output.(json.Unmarshaler).UnmarshalJSON(data)
			if err == nil {
				t.Fatalf("should return error for invalid json")
			}
		})
	}
}

type geom interface {
	Geometry() orb.Geometry
}

func BenchmarkGeometryMarshalJSON(b *testing.B) {
	ls := orb.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, orb.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := g.MarshalJSON()
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkGeometryUnmarshalJSON(b *testing.B) {
	ls := orb.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, orb.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}
	data, err := g.MarshalJSON()
	if err != nil {
		b.Fatalf("marshal error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.UnmarshalJSON(data)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
