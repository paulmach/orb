package geo

import (
	"math"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/internal/rect"
)

// A Rect represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Rect struct {
	rect.Rect
}

// NewRect creates a new rect given the parameters.
func NewRect(west, east, south, north float64) Rect {
	return Rect{
		Rect: rect.New(west, east, south, north),
	}
}

// NewRectFromPoints creates a new rect given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewRectFromPoints(corner, oppositeCorner Point) Rect {
	return Rect{
		Rect: rect.FromPoints(rect.Point(corner), rect.Point(oppositeCorner)),
	}
}

// NewRectAroundPoint creates a new rect given a center point,
// and a distance from the center point in meters.
func NewRectAroundPoint(center Point, distance float64) Rect {

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

	return Rect{
		Rect: rect.Rect{
			SW: rect.Point([2]float64{rad2deg(minLon), rad2deg(minLat)}),
			NE: rect.Point([2]float64{rad2deg(maxLon), rad2deg(maxLat)}),
		},
	}
}

// NewRectFromMapTile creates a rect given an online map tile index.
// Panics if x or y is out of range for zoom level.
func NewRectFromMapTile(x, y, z uint64) Rect {
	maxIndex := uint64(1) << z
	if x < 0 || y < 0 || x >= maxIndex || y >= maxIndex {
		panic("tile index out of range")
	}

	shift := 31 - z
	if z > 31 {
		shift = 0
	}

	lon1, lat1 := mercator.ScalarInverse(x<<shift, y<<shift, 31)
	lon2, lat2 := mercator.ScalarInverse((x+1)<<shift, (y+1)<<shift, 31)

	return Rect{
		Rect: rect.FromPoints(
			rect.Point([2]float64{math.Min(lon1, lon2), math.Min(lat1, lat2)}),
			rect.Point([2]float64{math.Max(lon1, lon2), math.Max(lat1, lat2)}),
		),
	}
}

// NewRectFromGeoHash creates a new rect for the region defined by the GeoHash.
func NewRectFromGeoHash(hash string) Rect {
	west, east, south, north := geoHash2ranges(hash)
	return NewRect(west, east, south, north)
}

// NewRectFromGeoHashInt64 creates a new rect from the region defined by the GeoHesh.
// bits indicates the precision of the hash.
func NewRectFromGeoHashInt64(hash int64, bits int) Rect {
	west, east, south, north := geoHashInt2ranges(hash, bits)
	return NewRect(west, east, south, north)
}

func geoHash2ranges(hash string) (float64, float64, float64, float64) {
	latMin, latMax := -90.0, 90.0
	lonMin, lonMax := -180.0, 180.0
	even := true

	for _, r := range hash {
		// TODO: index step could probably be done better
		i := strings.Index("0123456789bcdefghjkmnpqrstuvwxyz", string(r))
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

// Extend grows the rect to include the new point.
func (r Rect) Extend(point Point) Rect {
	r.Rect = r.Rect.Extend(rect.Point(point))
	return r
}

// Union extends this rect to contain the union of this and the given rect.
func (r Rect) Union(other Rect) Rect {
	r.Rect = r.Rect.Union(other.Rect)
	return r
}

// Contains determines if the point is within the rect.
// Points on the boundary are considered within.
func (r Rect) Contains(point Point) bool {
	return r.Rect.Contains(rect.Point(point))
}

// Intersects determines if two rectangles intersect.
// Returns true if they are touching.
func (r Rect) Intersects(rectangle Rect) bool {
	return r.Rect.Intersects(rectangle.Rect)
}

// Pad expands the rect in all directions by the given amount of meters.
func (r Rect) Pad(meters float64) Rect {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(r.NE.Y()))
	dx = math.Max(dx, dy/math.Cos(deg2rad(r.SW.Y())))

	r.SW[0] -= dx
	r.SW[1] -= dy

	r.NE[0] += dx
	r.NE[1] += dy

	return r
}

// Height returns the approximate height in meters.
func (r Rect) Height() float64 {
	return 111131.75 * (r.NE[1] - r.SW[1])
}

// Width returns the approximate width in meters
// of the center of the rectangle.
func (r Rect) Width(haversine ...bool) float64 {
	c := r.Center()

	a := Point{r.SW[0], c[1]}
	b := Point{r.NE[0], c[1]}

	return a.DistanceFrom(b, yesHaversine(haversine))
}

// Center returns the center of the rect.
func (r Rect) Center() Point {
	return Point(r.Rect.Center())
}

// North returns the top of the rect.
func (r Rect) North() float64 {
	return r.NE[1]
}

// South returns the bottom of the rect.
func (r Rect) South() float64 {
	return r.SW[1]
}

// East returns the right of the rect.
func (r Rect) East() float64 {
	return r.NE[0]
}

// West returns the left of the rect.
func (r Rect) West() float64 {
	return r.SW[0]
}

// SouthWest returns the lower left point of the rect.
func (r Rect) SouthWest() Point {
	return NewPoint(r.SW[0], r.SW[1])
}

// NorthEast return the upper right point of the rect.
func (r Rect) NorthEast() Point {
	return NewPoint(r.NE[0], r.NE[1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (r Rect) IsEmpty() bool {
	return r.Rect.IsEmpty()
}

// IsZero return true if the rect just includes just null island.
func (r Rect) IsZero() bool {
	return r.Rect.IsZero()
}

// String returns the string respentation of the rect in WKT format.
func (r Rect) String() string {
	return r.Rect.WKT()
}

// WKT returns the string respentation of the rect in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) WKT() string {
	return r.Rect.WKT()
}

// Pointer is a helper for using rectangle in a nullable context.
func (r Rect) Pointer() *Rect {
	return &r
}

// Equal returns if two rectangles are equal.
func (r Rect) Equal(c Rect) bool {
	return r.SW == c.SW && r.NE == c.NE
}
