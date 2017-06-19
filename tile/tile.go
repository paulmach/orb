package tile

import (
	"math"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

// Tiles is a set of tiles, later we can add methods to this.
type Tiles []Tile

// Tile is an x, y, z web mercator tile.
type Tile struct {
	X, Y, Z uint32
}

// New creates a tile for the point at the given zoom.
func New(ll geo.Point, z uint32) Tile {
	t := Tile{Z: z}
	t.X, t.Y = project.ScalarMercator.ToPlanar(ll, z)

	return t
}

// FromQuadkey creates the tile from the quadkey.
func FromQuadkey(k uint64, z uint32) Tile {
	t := Tile{Z: z}

	for i := uint32(0); i < z; i++ {
		t.X |= uint32((k & (1 << (2 * i))) >> i)
		t.Y |= uint32((k & (1 << (2*i + 1))) >> (i + 1))
	}

	return t
}

// Valid returns if the tile's x/y are within the range for the tile's zoom.
func (t Tile) Valid() bool {
	maxIndex := uint32(1) << t.Z
	return t.X < maxIndex && t.Z < maxIndex
}

// GeoBound returns the geo bound for the tile.
func (t Tile) GeoBound() geo.Bound {
	lon1, lat1 := mercator.ScalarInverse(t.X, t.Y, t.Z)
	lon2, lat2 := mercator.ScalarInverse(t.X+1, t.Y+1, t.Z)

	return geo.Bound{
		geo.Point{lon1, lat2},
		geo.Point{lon2, lat1},
	}
}

// Center returns the center of the tile.
func (t Tile) Center() geo.Point {
	return t.GeoBound().Center()
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

// Fraction returns the precise tile fraction at the given zoom.
func (t Tile) Fraction(ll geo.Point, z uint32) planar.Point {
	var p planar.Point

	factor := uint32(1 << z)
	maxtiles := float64(factor)

	lng := ll[0]/360.0 + 0.5
	p[0] = lng * maxtiles

	// bound it because we have a top of the world problem
	siny := math.Sin(ll[1] * math.Pi / 180.0)

	if siny < -0.9999 {
		p[1] = 0
	} else if siny > 0.9999 {
		p[1] = maxtiles
	} else {
		lat := 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		p[1] = lat * maxtiles
	}

	return p
}

// SharedParent returns the tile that contains both the tiles.
func (t Tile) SharedParent(tile Tile) Tile {
	// bring both tiles to the lowest zoom.
	if t.Z < tile.Z {
		tile = tile.toZoom(t.Z)
	} else {
		t = t.toZoom(tile.Z)
	}

	if t == tile {
		return t
	}

	// move from most significant to least until there isn't a match.
	// TODO: this can be improved using the go1.9 bits package.
	for i := t.Z; i > 0; i-- {
		if t.X&(1<<i) != tile.X&(1<<i) ||
			t.Y&(1<<i) != tile.Y&(1<<i) {
			return Tile{
				t.X >> (t.Z - i),
				t.Y >> (t.Z - i),
				i,
			}
		}
	}

	// if we reach here the tiles are the same, which was checked above.
	panic("unreachable")
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
	for i = 0; i < uint64(t.Z); i++ {
		result |= (uint64(t.X) & (1 << i)) << i
		result |= (uint64(t.Y) & (1 << i)) << (i + 1)
	}

	return result
}

// Range returns the min and max tile "range" to cover the tile
// at the given zoom.
func (t Tile) Range(z uint32) (min, max Tile) {
	if z < t.Z {
		t = t.toZoom(z)
		return t, t
	}

	offset := z - t.Z
	return Tile{
			X: t.X << offset,
			Y: t.Y << offset,
			Z: z,
		}, Tile{
			X: ((t.X + 1) << offset) - 1,
			Y: ((t.Y + 1) << offset) - 1,
			Z: z,
		}
}

func (t Tile) toZoom(z uint32) Tile {
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
