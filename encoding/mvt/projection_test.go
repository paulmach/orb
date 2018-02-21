package mvt

import (
	"testing"

	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/project"
)

func TestNonPowerOfTwoProjection(t *testing.T) {
	tile := maptile.New(8956, 12223, 15)
	regProj := newProjection(tile, 4096)
	nonProj := nonPowerOfTwoProjection(tile, 4096)

	expected := loadGeoJSON(t, tile)
	layers := NewLayers(loadGeoJSON(t, tile))

	// loopy de loop of projections
	for _, l := range layers {
		for _, f := range l.Features {
			f.Geometry = project.Geometry(f.Geometry, regProj.ToTile)
		}
	}

	for _, l := range layers {
		for _, f := range l.Features {
			f.Geometry = project.Geometry(f.Geometry, nonProj.ToWGS84)
		}
	}

	for _, l := range layers {
		for _, f := range l.Features {
			f.Geometry = project.Geometry(f.Geometry, nonProj.ToTile)
		}
	}

	for _, l := range layers {
		for _, f := range l.Features {
			f.Geometry = project.Geometry(f.Geometry, regProj.ToWGS84)
		}
	}

	result := layers.ToFeatureCollections()

	xEpsilon, yEpsilon := tileEpsilon(tile)
	for key := range expected {
		for i := range expected[key].Features {
			r := result[key].Features[i]
			e := expected[key].Features[i]

			compareOrbGeometry(t, r.Geometry, e.Geometry, xEpsilon, yEpsilon)
		}
	}
}
