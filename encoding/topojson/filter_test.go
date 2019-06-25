package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/orb"
	geojson "github.com/paulmach/orb/geojson"
)

func TestFilter(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	fc.Append(NewTestFeature("one", orb.LineString{
		orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{1, 1}, orb.Point{0, 1}, orb.Point{0, 0},
	}))
	fc.Append(NewTestFeature("two", orb.LineString{
		orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{2, 1}, orb.Point{1, 1}, orb.Point{1, 0},
	}))
	fc.Append(NewTestFeature("three", orb.LineString{
		orb.Point{1, 1}, orb.Point{2, 1}, orb.Point{2, 2}, orb.Point{1, 2}, orb.Point{1, 1},
	}))

	topo := NewTopology(fc, nil)
	is.NotNil(topo)

	al := len(topo.Arcs)
	is.True(al > 0)

	topo2 := topo.Filter([]string{"one", "two"})
	is.NotNil(topo2)

	al2 := len(topo2.Arcs)
	is.True(al > al2) // Arc has been eliminated

	expected := map[string][]orb.Point{
		"one": {{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
		"two": {{1, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 0}},
	}

	fc2 := topo2.ToGeoJSON()
	is.NotNil(fc2)

	for _, feat := range fc2.Features {
		exp, ok := expected[feat.ID.(string)]
		is.True(ok)
		is.Equal(feat.Geometry, exp)
	}
}
