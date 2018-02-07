package mvt

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/simplify"
)

func TestLayerSimplify(t *testing.T) {
	// should remove feature that are empty.
	l := &Layer{
		Features: []*geojson.Feature{
			geojson.NewFeature(orb.LineString(nil)),
			geojson.NewFeature(orb.LineString{{0, 0}, {1, 1}}),
		},
	}

	simplifier := simplify.DouglasPeucker(10)
	l.Simplify(simplifier)

	if len(l.Features) != 1 {
		t.Errorf("should remove empty feature")
	}

	if v := l.Features[0].Geometry.GeoJSONType(); v != "LineString" {
		t.Errorf("incorrect type: %v", v)
	}
}

func TestLayerRemoveEmpty(t *testing.T) {
	// should remove empty features
	l := &Layer{
		Features: []*geojson.Feature{
			geojson.NewFeature(orb.Ring{{0, 0}, {1, 1}, {0, 1}, {0, 0}}),
			geojson.NewFeature(orb.LineString{{0, 0}, {5, 5}, {0, 0}}),
		},
	}

	l.RemoveEmpty(2, 0.5)
	if len(l.Features) != 2 {
		t.Errorf("should not remove things above the limit")
	}

	// remove the area
	l.RemoveEmpty(2, 15)

	if len(l.Features) != 1 {
		t.Errorf("should remove empty feature")
	}

	if v := l.Features[0].Geometry.GeoJSONType(); v != "LineString" {
		t.Errorf("incorrect type: %v", v)
	}

	// remove the line
	l.RemoveEmpty(15, 2)

	if len(l.Features) != 0 {
		t.Errorf("should remove empty feature")
	}
}
