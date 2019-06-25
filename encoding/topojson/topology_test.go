package topojson

import (
	"testing"

	"github.com/paulmach/orb"

	"github.com/cheekybits/is"
	"github.com/paulmach/orb/geojson"
)

func TestTopology(t *testing.T) {
	is := is.New(t)

	poly := geojson.NewFeature(orb.Polygon{
		orb.Ring{
			{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0},
		},
	})
	poly.ID = "poly"

	fc := geojson.NewFeatureCollection()
	fc.Append(poly)

	topo := NewTopology(fc, nil)
	is.NotNil(topo)
	is.Equal(len(topo.Objects), 1)
	is.Equal(len(topo.Arcs), 1)
}

// Full roundtrip test
func TestFull(t *testing.T) {
	is := is.New(t)

	in := []byte(`{"type":"FeatureCollection","features":[{"type":"Feature", "id":"road", "properties":{"id":"road"},"geometry":{"type":"LineString","coordinates":[[4.707126617431641,50.88714360752515],[4.708242416381835,50.886683361842216],[4.708693027496338,50.886514152727386],[4.70914363861084,50.886321253586765],[4.709406495094299,50.88633140619302],[4.709567427635193,50.88636524819791],[4.709604978561401,50.88647354244835],[4.709545969963074,50.88664275171071],[4.708666205406189,50.88698116839158],[4.707368016242981,50.88743464288961]]}}]}`)
	out := []byte(`{"type":"Topology","bbox":[4.707126617431641,50.886321253586765,4.709604978561401,50.88743464288961],"objects":{"road":{"id":"road","properties":{"id":"road"},"type":"LineString","arcs":[0]}},"arcs":[[[0,7385],[4502,-4133],[1818,-1520],[1818,-1732],[1060,91],[650,304],[151,973],[-238,1519],[-3549,3039],[-5238,4073]]],"transform":{"scale":[2.478608990659808e-7,1.1135006529096161e-7],"translate":[4.707126617431641,50.886321253586765]}}`)

	fc, err := geojson.UnmarshalFeatureCollection(in)
	is.NoErr(err)

	topo := NewTopology(fc, &TopologyOptions{
		IDProperty:   "id",
		PreQuantize:  1000000,
		PostQuantize: 10000,
	})

	expected, err := UnmarshalTopology(out)
	is.NoErr(err)
	is.Equal(topo, expected)
}

// https://github.com/rubenv/topojson/issues/2
func TestMultiFeatures(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f1 := geojson.NewFeature(orb.Polygon{
		{{0, 0}, {1, 1}, {2, 0}, {1, -1}, {0, 0}},
		{{0.25, 0}, {0.75, -0.25}, {0.75, 0.25}, {0.25, 0}},
	})
	fc.Append(f1)
	f2 := geojson.NewFeature(orb.Polygon{
		{{0, 0}, {1, 1}, {2, 0}, {1, 2}, {0, 0}},
	})
	fc.Append(f2)

	topo := NewTopology(fc, nil)

	is.Equal(len(topo.Objects), len(fc.Features))
}

func TestCopyBounds(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f := NewTestFeature("one", orb.Polygon{
		orb.Ring{{0, 0}, {1, 1}, {2, 0}, {1, 2}, {0, 0}},
	})
	f.BBox = []float64{0, 0, 2, 2}
	fc.Append(f)

	topo := NewTopology(fc, nil)
	is.NotNil(topo)

	is.Equal(topo.Objects["one"].BBox, []float64{0, 0, 2, 2})
}
