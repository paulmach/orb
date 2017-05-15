package geo

import (
	"errors"
	"math"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

// A Bound represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Bound [2]Point

// NewBound creates a new bound given the parameters.
func NewBound(west, east, south, north float64) Bound {
	return Bound{
		Point{math.Min(west, east), math.Min(south, north)},
		Point{math.Max(west, east), math.Max(south, north)},
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{corner, corner}.Extend(oppositeCorner)
}

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point in meters.
func NewBoundAroundPoint(center Point, distance float64) Bound {
	radDist := distance / orb.EarthRadius
	radLat := deg2rad(center.Lat())
	radLon := deg2rad(center.Lon())
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
		Point{rad2deg(minLon), rad2deg(minLat)},
		Point{rad2deg(maxLon), rad2deg(maxLat)},
	}
}

// NewBoundFromMapTile creates a bound given an online map tile index.
// Panics if x or y is out of range for zoom level.
func NewBoundFromMapTile(x, y, z uint64) (Bound, error) {
	maxIndex := uint64(1) << z
	if x < 0 || x >= maxIndex {
		return Bound{}, errors.New("geo: x index out of range for this zoom")
	}
	if y < 0 || y >= maxIndex {
		return Bound{}, errors.New("geo: y index out of range for this zoom")
	}

	shift := 31 - z
	if z > 31 {
		shift = 0
	}

	lon1, lat1 := mercator.ScalarInverse(x<<shift, y<<shift, 31)
	lon2, lat2 := mercator.ScalarInverse((x+1)<<shift, (y+1)<<shift, 31)

	return Bound{
		Point{math.Min(lon1, lon2), math.Min(lat1, lat2)},
		Point{math.Max(lon1, lon2), math.Max(lat1, lat2)},
	}, nil
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
	lonMin, lonMax := -180.0, 180.0
	even := true

	for _, b := range hash {
		// TODO: index step could probably be done better
		i := strings.Index("0123456789bcdefghjkmnpqrstuvwxyz", string(b))
		for j := 0x10; j != 0; j >>= 1 {
			if even {
				mid := (lonMin + lonMax) / 2.0
				if i&j == 0 {
					lonMax = mid
				} else {
					lonMin = mid
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

	return lonMin, lonMax, latMin, latMax
}

func geoHashInt2ranges(hash int64, bits int) (float64, float64, float64, float64) {
	latMin, latMax := -90.0, 90.0
	lonMin, lonMax := -180.0, 180.0

	var i int64
	i = 1 << uint(bits)

	for i != 0 {
		i >>= 1

		mid := (lonMin + lonMax) / 2.0
		if hash&i == 0 {
			lonMax = mid
		} else {
			lonMin = mid
		}

		i >>= 1
		mid = (latMin + latMax) / 2.0
		if hash&i == 0 {
			latMax = mid
		} else {
			latMin = mid
		}
	}

	return lonMin, lonMax, latMin, latMax
}

// GeoJSONType returns the GeoJSON type for the object.
func (b Bound) GeoJSONType() string {
	return "Polygon"
}

// ToPolygon converts the bound into a Polygon object.
func (b Bound) ToPolygon() Polygon {
	return Polygon{b.ToRing()}
}

// ToRing converts the bound into a loop defined
// by the boundary of the box.
func (b Bound) ToRing() Ring {
	r := make(Ring, 5)
	r[0] = b[0]
	r[1] = Point{b[0][0], b[1][1]}
	r[2] = b[1]
	r[3] = Point{b[1][0], b[0][1]}
	r[4] = b[0]

	return r
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b[0][0] = math.Min(b[0].Lon(), point.Lon())
	b[1][0] = math.Max(b[1].Lon(), point.Lon())

	b[0][1] = math.Min(b[0].Lat(), point.Lat())
	b[1][1] = math.Max(b[1].Lat(), point.Lat())

	return b
}

// Union extends this bound to contain the union of this and the given bound.
func (b Bound) Union(other Bound) Bound {
	b = b.Extend(other[0])
	b = b.Extend(other[1])
	b = b.Extend(Point{other[0][0], other[1][1]})
	b = b.Extend(Point{other[1][0], other[0][1]})

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	if point.Lat() < b[0].Lat() || b[1].Lat() < point.Lat() {
		return false
	}

	if point.Lon() < b[0].Lon() || b[1].Lon() < point.Lon() {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	if (b[1][0] < bound[0][0]) ||
		(b[0][0] > bound[1][0]) ||
		(b[1][1] < bound[0][1]) ||
		(b[0][1] > bound[1][1]) {
		return false
	}

	return true
}

// Pad expands the bound in all directions by the given amount of meters.
func (b Bound) Pad(meters float64) Bound {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(b[1].Lat()))
	dx = math.Max(dx, dy/math.Cos(deg2rad(b[0].Lat())))

	b[0][0] -= dx
	b[0][1] -= dy

	b[1][0] += dx
	b[1][1] += dy

	return b
}

// Height returns the approximate height in meters.
func (b Bound) Height() float64 {
	return 111131.75 * (b[1][1] - b[0][1])
}

// Width returns the approximate width in meters
// of the center of the bound.
func (b Bound) Width(haversine ...bool) float64 {
	c := (b[0][1] + b[1][1]) / 2.0

	s1 := Point{b[0][0], c}
	s2 := Point{b[1][0], c}

	return s1.DistanceFrom(s2, yesHaversine(haversine))
}

// North returns the top of the bound.
func (b Bound) North() float64 {
	return b[1][1]
}

// South returns the bottom of the bound.
func (b Bound) South() float64 {
	return b[0][1]
}

// East returns the right of the bound.
func (b Bound) East() float64 {
	return b[1][0]
}

// West returns the left of the bound.
func (b Bound) West() float64 {
	return b[0][0]
}

// SouthWest returns the lower left point of the bound.
func (b Bound) SouthWest() Point {
	return NewPoint(b[0][0], b[0][1])
}

// NorthEast return the upper right point of the bound.
func (b Bound) NorthEast() Point {
	return NewPoint(b[1][0], b[1][1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b[0].Lon() > b[1].Lon() || b[0].Lat() > b[1].Lat()
}

// IsZero return true if the bound just includes just null island.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}

// String returns the string respentation of the bound in WKT format.
func (b Bound) String() string {
	return b.WKT()
}

// WKT returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) WKT() string {
	return b.ToPolygon().WKT()
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b[0] == c[0] && b[1] == c[1]
}
