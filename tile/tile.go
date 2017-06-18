package tile

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/project"
)

// Tiles is a set of tiles, later we can add methods to this.
type Tiles []Tile

// Tile is an x, y, z web mercator tile.
type Tile struct {
	X, Y, Z uint64
}

// New creates a tile for the point at the given zoom.
func New(ll geo.Point, z uint64) Tile {
	t := Tile{Z: z}
	t.X, t.Y = project.ScalarMercator.ToPlanar(ll, z)

	return t
}

// FromQuadkey creates the tile from the quadkey.
func FromQuadkey(k uint64, z uint64) Tile {
	t := Tile{Z: z}

	for i := uint64(0); i < z; i++ {
		t.X |= (k & (1 << (2 * i))) >> i
		t.Y |= (k & (1 << (2*i + 1))) >> (i + 1)
	}

	return t
}

// Contains returns if the given tile is fully contained (or equal to) the give tile.
func (t Tile) Contains(tile Tile) bool {
	if tile.Z < t.Z {
		return false
	}

	return t == tile.toZoom(t.Z)
}

// Parent returns the parent of the tile.
func (t Tile) Parent() Tile {
	if t.Z == 0 {
		return t
	}

	return Tile{
		X: t.X >> 1,
		Y: t.Y >> 1,
		Z: t.Z - 1,
	}
}

// Children returns the 4 children of the tile.
func (t Tile) Children() Tiles {
	return Tiles{
		Tile{t.X << 1, t.Y << 1, t.Z + 1},
		Tile{(t.X << 1) + 1, t.Y << 1, t.Z + 1},
		Tile{(t.X << 1) + 1, (t.Y << 1) + 1, t.Z + 1},
		Tile{t.X << 1, (t.Y << 1) + 1, t.Z + 1},
	}
}

// Siblings returns the 4 tiles that share this tile's parent.
func (t Tile) Siblings() Tiles {
	return t.Parent().Children()
}

// Quadkey returns the quad key for the tile.
func (t Tile) Quadkey() uint64 {
	var i, result uint64
	for i = 0; i < t.Z; i++ {
		result |= (t.X & (1 << i)) << i
		result |= (t.Y & (1 << i)) << (i + 1)
	}

	return result
}

func (t Tile) toZoom(z uint64) Tile {
	if z > t.Z {
		return Tile{
			X: t.X << (z - t.Z),
			Y: t.Y << (z - t.Z),
			Z: z,
		}
	}

	return Tile{
		X: t.X >> (t.Z - z),
		Y: t.Y >> (t.Z - z),
		Z: z,
	}
}
