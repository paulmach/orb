package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	orb "github.com/paulmach/orb"
	geojson "github.com/paulmach/orb/geojson"
)

func TestBBox(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("bar", orb.LineString{
			orb.Point{-1, 0}, orb.Point{1, 0}, orb.Point{-2, 3},
		}),
	}

	topo := &Topology{input: in}
	topo.bounds()

	is.Equal(topo.BBox, []float64{-2, 0, 2, 3})
}
