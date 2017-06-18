package geo

import (
	"fmt"
	"math"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

// A Point is a Lon/Lat 2d point.
type Point [2]float64

// NewPoint creates a new point.
func NewPoint(lon, lat float64) Point {
	return Point{lon, lat}
}

// NewPointFromQuadkey creates a new point from a quadkey.
// See http://msdn.microsoft.com/en-us/library/bb259689.aspx for more information
// about this coordinate system.
func NewPointFromQuadkey(key uint64, level int) Point {
	var x, y uint64

	var i uint
	for i = 0; i < uint(level); i++ {
		x |= (key & (1 << (2 * i))) >> i
		y |= (key & (1 << (2*i + 1))) >> (i + 1)
	}

	lon, lat := mercator.ScalarInverse(x, y, uint64(level))
	return Point{lon, lat}
}

// NewPointFromQuadkeyString creates a new point from a quadkey string.
func NewPointFromQuadkeyString(key string) Point {
	i, _ := strconv.ParseUint(key, 4, 64)
	return NewPointFromQuadkey(i, len(key))
}

// NewPointFromGeoHash creates a new point at the center of the geohash range.
func NewPointFromGeoHash(hash string) Point {
	west, east, south, north := geoHash2ranges(hash)
	return NewPoint((west+east)/2.0, (north+south)/2.0)
}

// NewPointFromGeoHashInt64 creates a new point at the center of the
// integer version of a geohash range. bits indicates the precision of the hash.
func NewPointFromGeoHashInt64(hash int64, bits int) Point {
	west, east, south, north := geoHashInt2ranges(hash, bits)
	return NewPoint((west+east)/2.0, (north+south)/2.0)
}

// GeoJSONType returns the GeoJSON type for the object.
func (p Point) GeoJSONType() string {
	return "Point"
}

// Bound returns a single point bound of the point.
func (p Point) Bound() Bound {
	return NewBound(p[0], p[0], p[1], p[1])
}

// DistanceFrom returns the geodesic distance in meters.
func (p Point) DistanceFrom(point Point, haversine ...bool) float64 {
	dLat := deg2rad(point.Lat() - p.Lat())
	dLon := deg2rad(point.Lon() - p.Lon())

	if yesHaversine(haversine) {
		// yes trig functions
		dLat2Sin := math.Sin(dLat / 2)
		dLon2Sin := math.Sin(dLon / 2)
		a := dLat2Sin*dLat2Sin + math.Cos(deg2rad(p.Lat()))*math.Cos(deg2rad(point.Lat()))*dLon2Sin*dLon2Sin

		return 2.0 * orb.EarthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	}

	dLon = math.Abs(dLon)
	if dLon > math.Pi {
		dLon = 2*math.Pi - dLon
	}

	// fast way using pythagorean theorem on an equirectangular projection
	x := dLon * math.Cos(deg2rad((p.Lat()+point.Lat())/2.0))
	return math.Sqrt(dLat*dLat+x*x) * orb.EarthRadius
}

// BearingTo computes the direction one must start traveling on earth
// to be heading to the given point.
func (p Point) BearingTo(point Point) float64 {
	dLon := deg2rad(point.Lon() - p.Lon())

	pLatRad := deg2rad(p.Lat())
	pointLatRad := deg2rad(point.Lat())

	y := math.Sin(dLon) * math.Cos(pointLatRad)
	x := math.Cos(pLatRad)*math.Sin(pointLatRad) - math.Sin(pLatRad)*math.Cos(pointLatRad)*math.Cos(dLon)

	return rad2deg(math.Atan2(y, x))
}

// Midpoint returns the half-way point along a great circle path between the two points.
func (p Point) Midpoint(p2 Point) Point {
	dLon := deg2rad(p2.Lon() - p.Lon())

	aLatRad := deg2rad(p.Lat())
	bLatRad := deg2rad(p2.Lat())

	x := math.Cos(bLatRad) * math.Cos(dLon)
	y := math.Cos(bLatRad) * math.Sin(dLon)

	r := Point{
		deg2rad(p.Lon()) + math.Atan2(y, math.Cos(aLatRad)+x),
		math.Atan2(math.Sin(aLatRad)+math.Sin(bLatRad), math.Sqrt((math.Cos(aLatRad)+x)*(math.Cos(aLatRad)+x)+y*y)),
	}

	// convert back to degrees
	r[0] = rad2deg(r[0])
	r[1] = rad2deg(r[1])

	return r
}

// Quadkey returns the quad key for the given point at the provided level.
// See http://msdn.microsoft.com/en-us/library/bb259689.aspx for more information
// about this coordinate system.
func (p Point) Quadkey(level int) uint64 {
	x, y := mercator.ScalarProject(p.Lon(), p.Lat(), uint64(level))

	var i uint
	var result uint64
	for i = 0; i < uint(level); i++ {
		result |= (x & (1 << i)) << i
		result |= (y & (1 << i)) << (i + 1)
	}

	return result
}

// QuadkeyString returns the quad key for the given point at the provided level in string form
// See http://msdn.microsoft.com/en-us/library/bb259689.aspx for more information
// about this coordinate system.
func (p Point) QuadkeyString(level int) string {
	s := strconv.FormatUint(p.Quadkey(level), 4)

	// for zero padding
	zeros := "000000000000000000000000000000"
	return zeros[:((level+1)-len(s))/2] + s
}

const base32 = "0123456789bcdefghjkmnpqrstuvwxyz"

// GeoHash returns the geohash string of a point representing a lon/lat location.
// The resulting hash will be `GeoHashPrecision` characters long, default is 12.
// Optionally one can include their required number of chars precision.
func (p Point) GeoHash(precision int) string {
	// 15 must be greater than GeoHashPrecision. If not, panic!!
	var result [15]byte

	hash := p.GeoHashInt64(5 * precision)
	for i := 1; i <= precision; i++ {
		result[precision-i] = base32[hash&0x1F]
		hash >>= 5
	}

	return string(result[:precision])
}

// GeoHashInt64 returns the integer version of the geohash
// down to the given number of bits.
// The main usecase for this function is to be able to do integer based ordering of points.
// In that case the number of bits should be the same for all encodings.
func (p Point) GeoHashInt64(bits int) (hash int64) {
	// This code was inspired by https://github.com/broady/gogeohash

	latMin, latMax := -90.0, 90.0
	lonMin, lonMax := -180.0, 180.0

	for i := 0; i < bits; i++ {
		hash <<= 1

		// interleave bits
		if i%2 == 0 {
			mid := (lonMin + lonMax) / 2.0
			if p[0] > mid {
				lonMin = mid
				hash |= 1
			} else {
				lonMax = mid
			}
		} else {
			mid := (latMin + latMax) / 2.0
			if p[1] > mid {
				latMin = mid
				hash |= 1
			} else {
				latMax = mid
			}
		}
	}

	return
}

// Equal checks if the point represents the same point or vector.
func (p Point) Equal(point Point) bool {
	return p[0] == point[0] && p[1] == point[1]
}

// Lat returns the latitude/vertical component of the point.
func (p Point) Lat() float64 {
	return p[1]
}

// Lon returns the longitude/horizontal component of the point.
func (p Point) Lon() float64 {
	return p[0]
}

// WKT returns the point in WKT format, eg. POINT(30.5 10.5)
func (p Point) WKT() string {
	return p.String()
}

// String returns a string representation of the point.
// The format is WKT, e.g. POINT(30.5 10.5)
func (p Point) String() string {
	return fmt.Sprintf("POINT(%g %g)", p[0], p[1])
}
