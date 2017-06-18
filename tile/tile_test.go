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
		p := tile.Center()

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > epsilon {
			t.Errorf("longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}

func TestTileGeoBound(t *testing.T) {
	bound := Tile{7, 8, 9}.GeoBound()

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the bound
	p := Tile{7<<factor + 1, 8<<factor + 1, level}.Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = Tile{7<<factor - 1, 8<<factor - 1, level}.Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}

	p = Tile{8<<factor - 1, 9<<factor - 1, level}.Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = Tile{8<<factor + 1, 9<<factor + 1, level}.Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}
}
