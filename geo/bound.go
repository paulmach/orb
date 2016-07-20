package geo

import (
	"math"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/bound"
	"github.com/paulmach/orb/internal/mercator"
)

// A Bound represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Bound struct {
	bound.Bound
}

// NewBound creates a new bound given the parameters.
func NewBound(west, east, south, north float64) Bound {
	return Bound{
		Bound: bound.New(west, east, south, north),
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{
		Bound: bound.FromPoints(bound.Point(corner), bound.Point(oppositeCorner)),
	}
}

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point in meters.
func NewBoundAroundPoint(center Point, distance float64) Bound {

	radDist := distance / orb.EarthRadius
	radLat := deg2rad(center.Lat())
	radLon := deg2rad(center.Lng())
	minLat := radLat - radDist
	maxLat := radLat + radDist

	var minLon, maxLon float64
	if minLat > minLatitude && maxLat < maxLatitude {
		deltaLon := math.Asin(math.Sin(radDist) / math.Cos(radLat))
		minLon = radLon - deltaLon
		if minLon < minLongitude {
			minLon += 2 * math.Pi
		}
		maxLon = radLon + deltaLon
		if maxLon > maxLongitude {
			maxLon -= 2 * math.Pi
		}
	} else {
		minLat = math.Max(minLat, minLatitude)
		maxLat = math.Min(maxLat, maxLatitude)
		minLon = minLongitude
		maxLon = maxLongitude
	}

	return Bound{
		Bound: bound.Bound{
			SW: bound.Point([2]float64{rad2deg(minLon), rad2deg(minLat)}),
			NE: bound.Point([2]float64{rad2deg(maxLon), rad2deg(maxLat)}),
		},
	}
}

// NewBoundFromMapTile creates a bound given an online map tile index.
// Panics if x or y is out of range for zoom level.
func NewBoundFromMapTile(x, y, z uint64) Bound {
	maxIndex := uint64(1) << z
	if x < 0 || y < 0 || x >= maxIndex || y >= maxIndex {
		panic("tile index out of range")
	}

	shift := 31 - z
	if z > 31 {
		shift = 0
	}

	lng1, lat1 := mercator.ScalarInverse(x<<shift, y<<shift, 31)
	lng2, lat2 := mercator.ScalarInverse((x+1)<<shift, (y+1)<<shift, 31)

	return Bound{
		Bound: bound.FromPoints(
			bound.Point([2]float64{math.Min(lng1, lng2), math.Min(lat1, lat2)}),
			bound.Point([2]float64{math.Max(lng1, lng2), math.Max(lat1, lat2)}),
		),
	}
}

// NewBoundFromGeoHash creates a new bound for the region defined by the GeoHash.
func NewBoundFromGeoHash(hash string) Bound {
	west, east, south, north := geoHash2ranges(hash)
	return NewBound(west, east, south, north)
}

// NewBoundFromGeoHashInt64 creates a new bound from the region defined by the GeoHesh.
// bits indicates the precision of the hash.
func NewBoundFromGeoHashInt64(hash int64, bits int) Bound {
	west, east, south, north := geoHashInt2ranges(hash, bits)
	return NewBound(west, east, south, north)
}

func geoHash2ranges(hash string) (float64, float64, float64, float64) {
	latMin, latMax := -90.0, 90.0
	lngMin, lngMax := -180.0, 180.0
	even := true

	for _, r := range hash {
		// TODO: index step could probably be done better
		i := strings.Index("0123456789bcdefghjkmnpqrstuvwxyz", string(r))
		for j := 0x10; j != 0; j >>= 1 {
			if even {
				mid := (lngMin + lngMax) / 2.0
				if i&j == 0 {
					lngMax = mid
				} else {
					lngMin = mid
				}
			} else {
				mid := (latMin + latMax) / 2.0
				if i&j == 0 {
					latMax = mid
				} else {
					latMin = mid
				}
			}
			even = !even
		}
	}

	return lngMin, lngMax, latMin, latMax
}

func geoHashInt2ranges(hash int64, bits int) (float64, float64, float64, float64) {
	latMin, latMax := -90.0, 90.0
	lngMin, lngMax := -180.0, 180.0

	var i int64
	i = 1 << uint(bits)

	for i != 0 {
		i >>= 1

		mid := (lngMin + lngMax) / 2.0
		if hash&i == 0 {
			lngMax = mid
		} else {
			lngMin = mid
		}

		i >>= 1
		mid = (latMin + latMax) / 2.0
		if hash&i == 0 {
			latMax = mid
		} else {
			latMin = mid
		}
	}

	return lngMin, lngMax, latMin, latMax
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	b.Bound = b.Bound.Extend(bound.Point(point))
	return b
}

// Union extends this bounds to contain the union of this and the given bounds.
func (b Bound) Union(other Bound) Bound {
	b.Bound = b.Bound.Union(other.Bound)
	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	return b.Bound.Contains(bound.Point(point))
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	return b.Bound.Intersects(bound.Bound)
}

// Pad expands the bound in all directions by the given amount of meters.
func (b Bound) Pad(meters float64) Bound {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(b.NE.Y()))
	dx = math.Max(dx, dy/math.Cos(deg2rad(b.SW.Y())))

	b.SW[0] -= dx
	b.SW[1] -= dy

	b.NE[0] += dx
	b.NE[1] += dy

	return b
}

// Height returns the approximate height in meters.
func (b Bound) Height() float64 {
	return 111131.75 * (b.NE[1] - b.SW[1])
}

// Width returns the approximate width in meters
// of the center of the bound.
func (b Bound) Width(haversine ...bool) float64 {
	c := b.Center()

	A := Point{b.SW[0], c[1]}
	B := Point{b.NE[0], c[1]}

	return A.DistanceFrom(B, yesHaversine(haversine))
}

// Center returns the center of the bound.
func (b Bound) Center() Point {
	return Point(b.Bound.Center())
}

// North returns the top of the bound.
func (b Bound) North() float64 {
	return b.NE[1]
}

// South returns the bottom of the bound.
func (b Bound) South() float64 {
	return b.SW[1]
}

// East returns the right of the bound.
func (b Bound) East() float64 {
	return b.NE[0]
}

// West returns the left of the bound.
func (b Bound) West() float64 {
	return b.SW[0]
}

// SouthWest returns the lower left point of the bound.
func (b Bound) SouthWest() Point {
	return NewPoint(b.SW[0], b.SW[1])
}

// NorthEast return the upper right point of the bound.
func (b Bound) NorthEast() Point {
	return NewPoint(b.NE[0], b.NE[1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b.Bound.IsEmpty()
}

// IsZero return true if the bound just includes just null island.
func (b Bound) IsZero() bool {
	return b.Bound.IsZero()
}

// String returns the string respentation of the bound in WKT format.
func (b Bound) String() string {
	return b.Bound.WKT()
}

// WKT returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) WKT() string {
	return b.Bound.WKT()
}

// Pointer is a helper for using bound in a nullable context.
func (b Bound) Pointer() *Bound {
	return &b
}

// MysqlIntersectsCondition returns a condition defining the intersection
// of the column and the bound. To be used in a MySQL query.
func (b Bound) MysqlIntersectsCondition(column string) string {
	return b.Bound.MysqlIntersectsCondition(column)
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b.SW == c.SW && b.NE == c.NE
}
