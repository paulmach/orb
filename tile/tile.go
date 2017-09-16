package tile

import (
	"math"
	"math/bits"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/planar"
)

// Tiles is a set of tiles, later we can add methods to this.
type Tiles []Tile

// Tile is an x, y, z web mercator tile.
type Tile struct {
	X, Y, Z uint32
}

// New creates a tile for the point at the given zoom.
// Will create a valid tile for the zoom. Points outside
// the range lat [-85.0511, 85.0511] will be snapped to the
// max or min tile as appropriate.
func New(ll geo.Point, z uint32) Tile {
	f := Fraction(ll, z)
	t := Tile{
		X: uint32(f[0]),
		Y: uint32(f[1]),
		Z: z,
	}

	// things
	if t.Y >= 1<<z {
		t.Y = (1 << z) - 1
	}

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

// Bound returns the geo bound for the tile.
func (t Tile) Bound() geo.Bound {
	lon1, lat1 := mercator.ToGeo(t.X, t.Y, t.Z)
	lon2, lat2 := mercator.ToGeo(t.X+1, t.Y+1, t.Z)

	return geo.Bound{
		geo.Point{lon1, lat2},
		geo.Point{lon2, lat1},
	}
}

// Center returns the center of the tile.
func (t Tile) Center() geo.Point {
	return t.Bound().Center()
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
// Returns the range y range of [0, 2^zoom]. Will return 2^zoom if
// the point is below 85.0511 S.
func Fraction(ll geo.Point, z uint32) planar.Point {
	var p planar.Point

	factor := uint32(1 << z)
	maxtiles := float64(factor)

	for ll[0] >= 180 {
		ll[0] -= 360.0
	}

	for ll[0] < -180 {
		ll[0] += 360.0
	}

	lng := ll[0]/360.0 + 0.5
	p[0] = lng * maxtiles

	// bound it because we have a top of the world problem
	if ll[1] < -85.0511 {
		p[1] = maxtiles
	} else if ll[1] > 85.0511 {
		p[1] = 0
	} else {
		siny := math.Sin(ll[1] * math.Pi / 180.0)
		lat := 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		p[1] = lat * maxtiles
	}

	return p
}

// SharedParent returns the tile that contains both the tiles.
func (t Tile) SharedParent(tile Tile) Tile {
	// bring both tiles to the lowest zoom.
	if t.Z != tile.Z {
		if t.Z < tile.Z {
			tile = tile.toZoom(t.Z)
		} else {
			t = t.toZoom(tile.Z)
		}
	}

	if t == tile {
		return t
	}

	// go version < 1.9
	// bit package usage was about 10% faster
	//
	// TODO: use build flags to support older versions of go.
	//
	// move from most significant to least until there isn't a match.
	// for i := t.Z - 1; i >= 0; i-- {
	// 	if t.X&(1<<i) != tile.X&(1<<i) ||
	// 		t.Y&(1<<i) != tile.Y&(1<<i) {
	// 		return Tile{
	// 			t.X >> (i + 1),
	// 			t.Y >> (i + 1),
	// 			t.Z - (i + 1),
	// 		}
	// 	}
	// }
	//
	// if we reach here the tiles are the same, which was checked above.
	// panic("unreachable")

	// bits different for x and y
	xc := uint32(32 - bits.LeadingZeros32(t.X^tile.X))
	yc := uint32(32 - bits.LeadingZeros32(t.Y^tile.Y))

	// max of xc, yc
	maxc := xc
	if yc > maxc {
		maxc = yc

	}

	return Tile{
		X: t.X >> maxc,
		Y: t.Y >> maxc,
		Z: t.Z - maxc,
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
