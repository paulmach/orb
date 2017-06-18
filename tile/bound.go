package tile

import "github.com/paulmach/orb/geo"

// Interval represents a closed interval of values.
// Used to devine a bound.
type Interval struct {
	Lo, Hi uint64
}

// Contains checks if the interval contains (inclusive) the value.
func (i Interval) Contains(v uint64) bool {
	return !(v < i.Lo || i.Hi < v)
}

// ContainsInterval check is the interval full contains the provided interval.
func (i Interval) ContainsInterval(i2 Interval) bool {
	return !(i2.Hi < i.Lo || i.Hi < i2.Lo)
}

// Bound represents a rectangle of tiles tiles.
type Bound struct {
	X, Y Interval
	Z    uint64
}

func newInterval(tx, tz, z uint64) Interval {
	if z > tz {
		return Interval{
			Lo: tx << (z - tz),
			Hi: ((tx + 1) << (z - tz)) - 1,
		}
	}

	return Interval{
		Lo: tx >> (tz - z),
		Hi: tx >> (tz - z),
	}
}

// Bound converts the tile into a bound at the given zoom.
func (t Tile) Bound(z uint64) Bound {
	return Bound{
		X: newInterval(t.X, t.Z, z),
		Y: newInterval(t.Y, t.Z, z),
		Z: z,
	}
}

// NewBound creates a "coverting" bound for the area defined by
// the given points.
func NewBound(lo, hi geo.Point, z uint64) Bound {
	lot := New(lo, z)
	hit := New(hi, z)

	return Bound{
		X: Interval{lot.X, hit.X},
		Y: Interval{hit.Y, lot.Y},
		Z: z,
	}
}

// Contains evaluates if the tile is within the bound.
func (b Bound) Contains(t Tile) bool {
	if t.Z < b.Z {
		return b.X.ContainsInterval(newInterval(t.X, t.Z, b.Z)) &&
			b.Y.ContainsInterval(newInterval(t.Y, t.Z, b.Z))
	}

	offset := t.Z - b.Z
	return b.X.Contains(t.X>>offset) && b.Y.Contains(t.Y>>offset)
}

// Covering returns the set of zoom z tiles that cover the bound.
func (b Bound) Covering(z uint64) Tiles {
	var lx, hx, ly, hy uint64
	if z > b.Z {
		lx = b.X.Lo << (z - b.Z)
		hx = ((b.X.Hi + 1) << (b.Z - z)) - 1
		ly = b.Y.Lo << (z - b.Z)
		hy = ((b.Y.Hi + 1) << (b.Z - z)) - 1
	} else {
		lx = b.X.Lo >> (b.Z - z)
		hx = b.X.Hi >> (b.Z - z)
		ly = b.Y.Lo >> (b.Z - z)
		hy = b.Y.Hi >> (b.Z - z)
	}

	result := make(Tiles, 0, (hx-lx+1)*(hy-ly+1))
	for x := lx; x <= hx; x++ {
		for y := ly; y <= hy; y++ {
			result = append(result, Tile{X: x, Y: y, Z: z})
		}
	}

	return result
}
