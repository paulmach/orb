package topojson

import (
	"fmt"

	orb "github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func NewTestFeature(id string, geom orb.Geometry) *geojson.Feature {
	feature := geojson.NewFeature(geom)
	feature.ID = id
	return feature
}

func GetFeature(topo *Topology, id string) *topologyObject {
	for _, o := range topo.objects {
		if o.ID == id {
			return o
		}
	}
	panic(fmt.Sprintf("No such object: %s", id))
}
