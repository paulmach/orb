package tilecover

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// ForBound creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func ForBound(b orb.Bound, z maptile.Zoom) maptile.Tiles {
	lo := maptile.At(b[0], z)
	hi := maptile.At(b[1], z)

	result := make(maptile.Tiles, 0, (hi.X-lo.X+1)*(lo.Y-hi.Y+1))

	for x := lo.X; x <= hi.X; x++ {
		for y := hi.Y; y <= lo.Y; y++ {
			result = append(result, maptile.Tile{X: x, Y: y, Z: z})
		}
	}

	return result
}

// ForPoint creates a tile cover for the point, i.e. just the tile
// containing the point.
func ForPoint(ll orb.Point, z maptile.Zoom) maptile.Tiles {
	return maptile.Tiles{maptile.At(ll, z)}
}
