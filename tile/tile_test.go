package tile

import (
	"math"
	"testing"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
)

func TestNew(t *testing.T) {
	tile := New(geo.NewPoint(0, 0), 28)
	if b := tile.Bound(); b.North() != 0 || b.West() != 0 {
		t.Errorf("incorrect tile bound: %v", b)
	}

	// specific case
	if tile := New(geo.NewPoint(-87.65005229999997, 41.850033), 20); tile.X != 268988 || tile.Y != 389836 {
		t.Errorf("projection incorrect: %v", tile)
	}

	if tile := New(geo.NewPoint(-87.65005229999997, 41.850033), 28); tile.X != 68861112 || tile.Y != 99798110 {
		t.Errorf("projection incorrect: %v", tile)
	}

	for _, city := range mercator.Cities {
		tile := New(geo.Point{city[1], city[0]}, 31)
		c := tile.Center()

		if math.Abs(c.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", c.Lat(), city[0])
		}

		if math.Abs(c.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", c.Lon(), city[1])
		}
	}

	// test polar regions
	if tile := New(geo.NewPoint(0, 89.9), 30); tile.Y != 0 {
		t.Errorf("top of the world error: %d != %d", tile.Y, 0)
	}

	if tile := New(geo.NewPoint(0, -89.9), 30); tile.Y != (1<<30)-1 {
		t.Errorf("bottom of the world error: %d != %d", tile.Y, (1<<30)-1)
	}
}

func TestTileQuadkey(t *testing.T) {
	// default level
	level := uint32(30)
	for _, city := range mercator.Cities {
		tile := New(geo.Point{city[1], city[0]}, level)
		p := tile.Center()

		if math.Abs(p.Lat()-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lon()-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", p.Lon(), city[1])
		}
	}
}

func TestTileBound(t *testing.T) {
	bound := Tile{7, 8, 9}.Bound()

	level := uint32(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint32(5)

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
