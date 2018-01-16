package mvt

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt/vectortile"
)

func TestGeometry_Point(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.Point
	}{
		{
			name:   "basic point",
			input:  []uint32{9, 50, 34},
			output: orb.Point{25, 17},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_POINT, tc.input, tc.output)
		})
	}
}

func TestGeometry_MultiPoint(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.MultiPoint
	}{
		{
			name:   "basic multi point",
			input:  []uint32{17, 10, 14, 3, 9},
			output: orb.MultiPoint{{5, 7}, {3, 2}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_POINT, tc.input, tc.output)
		})
	}
}

func TestGeometry_LineString(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.LineString
	}{
		{
			name:   "basic line string",
			input:  []uint32{9, 4, 4, 18, 0, 16, 16, 0},
			output: orb.LineString{{2, 2}, {2, 10}, {10, 10}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_LINESTRING, tc.input, tc.output)
		})
	}
}

func TestGeometry_MultiLineString(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.MultiLineString
	}{
		{
			name:  "basic multi line string",
			input: []uint32{9, 4, 4, 18, 0, 16, 16, 0, 9, 17, 17, 10, 4, 8},
			output: orb.MultiLineString{
				{{2, 2}, {2, 10}, {10, 10}},
				{{1, 1}, {3, 5}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_LINESTRING, tc.input, tc.output)
		})
	}
}

func TestGeometry_Polygon(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.Polygon
	}{
		{
			name:  "basic polygon",
			input: []uint32{9, 6, 12, 18, 10, 12, 24, 44, 15},
			output: orb.Polygon{
				{{3, 6}, {8, 12}, {20, 34}, {3, 6}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_POLYGON, tc.input, tc.output)
		})
	}
}

func TestGeometry_MultiPolygon(t *testing.T) {
	cases := []struct {
		name   string
		input  []uint32
		output orb.MultiPolygon
	}{
		{
			name: "multi polygon",
			input: []uint32{9, 0, 0, 26, 20, 0, 0, 20, 19, 0, 15, 9, 22, 2, 26,
				18, 0, 0, 18, 17, 0, 15, 9, 4, 13, 26, 0, 8, 8, 0, 0, 7, 15},
			output: orb.MultiPolygon{
				{
					{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
				},
				{
					{{11, 11}, {20, 11}, {20, 20}, {11, 20}, {11, 11}},
					{{13, 13}, {13, 17}, {17, 17}, {17, 13}, {13, 13}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareGeometry(t, vectortile.Tile_POLYGON, tc.input, tc.output)
		})
	}
}

func TestKeyValueEncoder_JSON(t *testing.T) {
	kve := newKeyValueEncoder()
	i, err := kve.Value([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	value := decodeValue(kve.Values[i])
	if value != "[1,2,3]" {
		t.Errorf("should encode non standard types as json")
	}
}

func compareGeometry(
	t testing.TB,
	geomType vectortile.Tile_GeomType,
	input []uint32,
	expected orb.Geometry,
) {
	t.Helper()

	// test encoding
	gt, encoded, err := encodeGeometry(expected)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	if gt != geomType {
		t.Errorf("type mismatch: %v != %v", gt, geomType)
	}

	if !reflect.DeepEqual(encoded, input) {
		t.Logf("%v", encoded)
		t.Logf("%v", input)
		t.Errorf("different encoding")
	}

	result, err := decodeGeometry(geomType, input)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if result.GeoJSONType() != expected.GeoJSONType() {
		t.Errorf("types different: %s != %s", result.GeoJSONType(), expected.GeoJSONType())
	}

	if !orb.Equal(result, expected) {
		t.Logf("%v", result)
		t.Logf("%v", expected)
		t.Errorf("geometry not equal")
	}
}
