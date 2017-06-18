package tile

import (
	"math"
	"testing"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
)

var epsilon = 1e-6

func TestTileQuadkey(t *testing.T) {
	// default level
	level := uint64(30)
	for _, city := range mercator.Cities {
		tile := New(geo.Point{city[1], city[0]}, level)
		key := tile.Quadkey()

		p := geo.NewPointFromQuadkey(key, int(level))

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > epsilon {
			t.Errorf("longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}
