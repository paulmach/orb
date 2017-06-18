package tilecover

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/tile"
)

// ForBound creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func ForBound(b geo.Bound, z uint64) tile.Tiles {
	return tile.NewBound(b[0], b[1], z).Covering(z)
}

// ForPoint creates a tile cover for the point, i.e. just the tile
// containing the point.
func ForPoint(ll geo.Point, z uint64) tile.Tiles {
	return tile.Tiles{tile.New(ll, z)}
}
