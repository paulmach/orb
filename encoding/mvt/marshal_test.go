package mvt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt/vectortile"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
)

func TestMarshalMarshalGzipped_Full(t *testing.T) {
	tile := maptile.New(8956, 12223, 15)
	ls := orb.LineString{
		{-81.60346275, 41.50998572},
		{-81.6033669, 41.50991259},
		{-81.60355599, 41.50976036},
		{-81.6040648, 41.50936811},
		{-81.60404411, 41.50935405},
	}
	expected := ls.Clone()

	f := geojson.NewFeature(ls)
	f.Properties = geojson.Properties{
		"source":       "openstreetmap.org",
		"kind":         "path",
		"name":         "Uptown Alley",
		"landuse_kind": "retail",
		"sort_rank":    float64(354),
		"kind_detail":  "pedestrian",
		"min_zoom":     float64(13),
		"id":           float64(246698394),
	}

	fc := geojson.NewFeatureCollection()
	fc.Append(f)

	layers := Layers{NewLayer("roads", fc)}

	// project to the tile coords
	layers.ProjectToTile(tile)

	// marshal
	encoded, err := MarshalGzipped(layers)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// unmarshal
	decoded, err := UnmarshalGzipped(encoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// project back
	decoded.ProjectToWGS84(tile)

	// compare the results
	result := decoded[0].Features[0]
	compareProperties(t, result.Properties, f.Properties)

	// compare geometry
	xe, ye := tileEpsilon(tile)
	compareOrbGeometry(t, result.Geometry, expected, xe, ye)
}

func TestMarshalUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		tile maptile.Tile
	}{
		{
			name: "15-8956-12223",
			tile: maptile.New(8956, 12223, 15),
		},
		{
			name: "16-17896-24449",
			tile: maptile.New(17896, 24449, 16),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := loadGeoJSON(t, tc.tile)
			layers := NewLayers(loadGeoJSON(t, tc.tile))
			layers.ProjectToTile(tc.tile)
			data, err := Marshal(layers)
			if err != nil {
				t.Errorf("marshal error: %v", err)
			}

			layers, err = Unmarshal(data)
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}

			layers.ProjectToWGS84(tc.tile)
			result := layers.ToFeatureCollections()

			xEpsilon, yEpsilon := tileEpsilon(tc.tile)
			for key := range expected {
				for i := range expected[key].Features {
					r := result[key].Features[i]
					e := expected[key].Features[i]

					// t.Logf("checking layer %s: feature %d", key, i)
					compareProperties(t, r.Properties, e.Properties)
					compareOrbGeometry(t, r.Geometry, e.Geometry, xEpsilon, yEpsilon)
				}
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		tile maptile.Tile
	}{
		{
			name: "15-8956-12223",
			tile: maptile.New(8956, 12223, 15),
		},
		{
			name: "16-17896-24449",
			tile: maptile.New(17896, 24449, 16),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := loadGeoJSON(t, tc.tile)
			layers, err := UnmarshalGzipped(loadMVT(t, tc.tile))
			if err != nil {
				t.Fatalf("error unmarshalling: %v", err)
			}

			layers.ProjectToWGS84(tc.tile)
			result := layers.ToFeatureCollections()

			xEpsilon, yEpsilon := tileEpsilon(tc.tile)
			for key := range expected {
				for i := range expected[key].Features {
					r := result[key].Features[i]
					e := expected[key].Features[i]

					t.Logf("checking layer %s: feature %d", key, i)
					compareProperties(t, r.Properties, e.Properties)
					compareOrbGeometry(t, r.Geometry, e.Geometry, xEpsilon, yEpsilon)
				}
			}
		})
	}
}

func tileEpsilon(tile maptile.Tile) (float64, float64) {
	b := tile.Bound()
	xEpsilon := (b.Max[0] - b.Min[0]) / 4096 * 2 // allowed error
	yEpsilon := (b.Max[1] - b.Min[1]) / 4096 * 2

	return xEpsilon, yEpsilon
}

func compareProperties(t testing.TB, result, expected geojson.Properties) {
	t.Helper()

	// properties
	fr := map[string]interface{}(result)
	fe := map[string]interface{}(expected)

	for k, v := range fe {
		if _, ok := v.([]interface{}); ok {
			// arrays are not included
			delete(fr, k)
			delete(fe, k)
		}

		// https: //github.com/tilezen/mapbox-vector-tile/pull/97
		// mapzen error where a 1 is encoded a boolean true
		// just ignore all known cases of that.
		if k == "scale_rank" || k == "layer" {
			if v == 1.0 {
				delete(fr, k)
				delete(fe, k)
			}
		}
	}

	if !reflect.DeepEqual(fr, fe) {
		t.Errorf("properties not equal")
		if len(fr) != len(fe) {
			t.Errorf("properties length not equal: %v != %v", len(fr), len(fe))
		}

		for k := range fr {
			t.Logf("%s: %T %v -- %T %v", k, fr[k], fr[k], fe[k], fe[k])
		}
	}
}

func compareOrbGeometry(
	t testing.TB,
	result, expected orb.Geometry,
	xEpsilon, yEpsilon float64,
) {
	t.Helper()

	if result.GeoJSONType() != expected.GeoJSONType() {
		t.Errorf("different types: %v != %v", result.GeoJSONType(), expected.GeoJSONType())
		return
	}

	switch r := result.(type) {
	case orb.Point:
		comparePoints(t,
			[]orb.Point{r},
			[]orb.Point{expected.(orb.Point)},
			xEpsilon, yEpsilon,
		)
	case orb.MultiPoint:
		comparePoints(t,
			[]orb.Point(r),
			[]orb.Point(expected.(orb.MultiPoint)),
			xEpsilon, yEpsilon,
		)
	case orb.LineString:
		comparePoints(t,
			[]orb.Point(r),
			[]orb.Point(expected.(orb.LineString)),
			xEpsilon, yEpsilon,
		)
	case orb.MultiLineString:
		e := expected.(orb.MultiLineString)
		for i := range r {
			compareOrbGeometry(t, r[i], e[i], xEpsilon, yEpsilon)
		}
	case orb.Polygon:
		e := expected.(orb.Polygon)
		for i := range r {
			compareOrbGeometry(t, orb.LineString(r[i]), orb.LineString(e[i]), xEpsilon, yEpsilon)
		}
	case orb.MultiPolygon:
		e := expected.(orb.MultiPolygon)
		for i := range r {
			compareOrbGeometry(t, r[i], e[i], xEpsilon, yEpsilon)
		}
	default:
		t.Errorf("unsupported type: %T", result)
	}
}

func comparePoints(t testing.TB, e, r []orb.Point, xEpsilon, yEpsilon float64) {
	if len(r) != len(e) {
		t.Errorf("geometry length not equal: %v != %v", len(r), len(e))
	}

	for i := range e {
		xe := math.Abs(r[i][0] - e[i][0])
		ye := math.Abs(r[i][1] - e[i][1])

		if xe > xEpsilon {
			t.Errorf("%d x: %f != %f    %f", i, r[i][0], e[i][0], xe)
		}

		if ye > yEpsilon {
			t.Errorf("%d y: %f != %f    %f", i, r[i][1], e[i][1], ye)
		}
	}
}

func loadMVT(t testing.TB, tile maptile.Tile) []byte {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%d-%d-%d.mvt", tile.Z, tile.X, tile.Y))
	if err != nil {
		t.Fatalf("failed to load mvt file: %v", err)
	}

	return data
}

func loadGeoJSON(t testing.TB, tile maptile.Tile) map[string]*geojson.FeatureCollection {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%d-%d-%d.json", tile.Z, tile.X, tile.Y))
	if err != nil {
		t.Fatalf("failed to load mvt file: %v", err)
	}

	r := make(map[string]*geojson.FeatureCollection)
	err = json.Unmarshal(data, &r)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	return r
}

func BenchmarkMarshal(b *testing.B) {
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Marshal(layers)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	data, err := Marshal(layers)
	if err != nil {
		b.Fatalf("marshal error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Unmarshal(data)
	}
}

func BenchmarkDecode(b *testing.B) {
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	data, err := Marshal(layers)
	if err != nil {
		b.Fatalf("marshal error: %v", err)
	}

	vt := &vectortile.Tile{}
	err = vt.Unmarshal(data)
	if err != nil {
		b.Fatalf("unmarshal error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decode(vt)
	}
}

func BenchmarkProjectToTile(b *testing.B) {
	tile := maptile.New(17896, 24449, 16)
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		layers.ProjectToTile(tile)
	}
}

func BenchmarkProjectToGeo(b *testing.B) {
	tile := maptile.New(17896, 24449, 16)
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))

	layers.ProjectToTile(tile)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		layers.ProjectToWGS84(tile)
	}
}
