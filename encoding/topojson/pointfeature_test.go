package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func TestPointFeature(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f := geojson.NewFeature(orb.Point{0, 0})
	f.ID = "point"
	fc.Append(f)

	topo := NewTopology(fc, nil)

	is.Equal([]float64{0, 0}, topo.Objects["point"].Point)
}

func TestMultiPointFeature(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f := geojson.NewFeature(orb.MultiPoint{orb.Point{0, 0}, orb.Point{1, 1}})
	f.ID = "multipoint"
	fc.Append(f)

	topo := NewTopology(fc, nil)

	is.Equal([][]float64{{0, 0}, {1, 1}}, topo.Objects["multipoint"].MultiPoint)
}
