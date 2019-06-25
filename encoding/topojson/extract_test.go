package topojson

import (
	"reflect"
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/extract-test.js

// extract copies coordinates sequentially into a buffer
func TestCopiesCoordinates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("bar", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	expected := [][]float64{
		{0, 0}, {1, 0}, {2, 0},
		{0, 0}, {1, 0}, {2, 0},
	}

	topo := &Topology{input: in}
	topo.extract()
	is.Equal(len(topo.coordinates), len(expected))
	for k, v := range topo.coordinates {
		is.Equal(v, expected[k])
	}
}

// extract includes closing coordinates in polygons
func TestClosingCoordinates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
		}),
	}

	expected := [][]float64{
		{0, 0}, {1, 0}, {2, 0}, {0, 0},
	}

	topo := &Topology{input: in}
	topo.extract()
	is.Equal(len(topo.coordinates), len(expected))
	for k, v := range topo.coordinates {
		is.Equal(v, expected[k])
	}
}

// extract represents lines as contiguous slices of the coordinate buffer
func TestLineSlices(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("bar", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()

	foo := GetFeature(topo, "foo")
	is.Equal(foo.Type, geojson.TypeLineString)
	is.True(reflect.DeepEqual(foo.Arc, &arc{Start: 0, End: 2}))

	bar := GetFeature(topo, "bar")
	is.Equal(bar.Type, geojson.TypeLineString)
	is.True(reflect.DeepEqual(bar.Arc, &arc{Start: 3, End: 5}))
}

// extract exposes the constructed lines and rings in the order of construction
func TestExtractRingsOrder(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("line", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("multiline", orb.MultiLineString{
			orb.LineString{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
			},
		}),
		NewTestFeature("polygon", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()

	is.True(reflect.DeepEqual(topo.lines, []*arc{
		{Start: 0, End: 2},
		{Start: 3, End: 5},
	}))
	is.True(reflect.DeepEqual(topo.rings, []*arc{
		{Start: 6, End: 9},
	}))
}

// extract supports nested geometry collections
func TestExtractNested(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.Collection{
			orb.LineString{
				orb.Point{0, 0}, orb.Point{0, 1},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()

	foo := GetFeature(topo, "foo")
	is.Equal(foo.Type, "GeometryCollection")

	geometries := foo.Geometries
	is.Equal(len(geometries), 1)
	is.Equal(geometries[0].Type, geojson.TypeLineString)
	is.True(reflect.DeepEqual(geometries[0].Arc, &arc{Start: 0, End: 1}))
}
