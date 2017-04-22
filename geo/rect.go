package geo

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

// A Rect represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Rect [2]Point

// NewRect creates a new rect given the parameters.
func NewRect(west, east, south, north float64) Rect {
	return Rect{
		Point{math.Min(west, east), math.Min(south, north)},
		Point{math.Max(west, east), math.Max(south, north)},
	}
}

// NewRectFromPoints creates a new rect given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewRectFromPoints(corner, oppositeCorner Point) Rect {
	return Rect{corner, corner}.Extend(oppositeCorner)
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
		Point{rad2deg(minLon), rad2deg(minLat)},
		Point{rad2deg(maxLon), rad2deg(maxLat)},
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
		Point{math.Min(lon1, lon2), math.Min(lat1, lat2)},
		Point{math.Max(lon1, lon2), math.Max(lat1, lat2)},
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

// ToLineString converts the Rect into a loop defined
// by the Rect boundary.
func (r Rect) ToLineString() LineString {
	return append(NewLineStringPreallocate(0, 5),
		r[0],
		NewPoint(r[0][0], r[1][1]),
		r[1],
		NewPoint(r[1][0], r[0][1]),
		r[0],
	)
}

// Extend grows the rect to include the new point.
func (r Rect) Extend(point Point) Rect {
	// already included, no big deal
	if r.Contains(point) {
		return r
	}

	r[0][0] = math.Min(r[0].Lon(), point.Lon())
	r[1][0] = math.Max(r[1].Lon(), point.Lon())

	r[0][1] = math.Min(r[0].Lat(), point.Lat())
	r[1][1] = math.Max(r[1].Lat(), point.Lat())

	return r
}

// Union extends this rect to contain the union of this and the given rect.
func (r Rect) Union(other Rect) Rect {
	r = r.Extend(other[0])
	r = r.Extend(other[1])
	r = r.Extend(Point{other[0][0], other[1][1]})
	r = r.Extend(Point{other[1][0], other[0][1]})

	return r
}

// Contains determines if the point is within the rect.
// Points on the boundary are considered within.
func (r Rect) Contains(point Point) bool {
	if point.Lat() < r[0].Lat() || r[1].Lat() < point.Lat() {
		return false
	}

	if point.Lon() < r[0].Lon() || r[1].Lon() < point.Lon() {
		return false
	}

	return true
}

// Intersects determines if two rectangles intersect.
// Returns true if they are touching.
func (r Rect) Intersects(rect Rect) bool {
	if (r[1][0] < rect[0][0]) ||
		(r[0][0] > rect[1][0]) ||
		(r[1][1] < rect[0][1]) ||
		(r[0][1] > rect[1][1]) {
		return false
	}

	return true
}

// Pad expands the rect in all directions by the given amount of meters.
func (r Rect) Pad(meters float64) Rect {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(r[1].Lat()))
	dx = math.Max(dx, dy/math.Cos(deg2rad(r[0].Lat())))

	r[0][0] -= dx
	r[0][1] -= dy

	r[1][0] += dx
	r[1][1] += dy

	return r
}

// Height returns the approximate height in meters.
func (r Rect) Height() float64 {
	return 111131.75 * (r[1][1] - r[0][1])
}

// Width returns the approximate width in meters
// of the center of the rectangle.
func (r Rect) Width(haversine ...bool) float64 {
	c := (r[0][1] + r[1][1]) / 2.0

	a := Point{r[0][0], c}
	b := Point{r[1][0], c}

	return a.DistanceFrom(b, yesHaversine(haversine))
}

// North returns the top of the rect.
func (r Rect) North() float64 {
	return r[1][1]
}

// South returns the bottom of the rect.
func (r Rect) South() float64 {
	return r[0][1]
}

// East returns the right of the rect.
func (r Rect) East() float64 {
	return r[1][0]
}

// West returns the left of the rect.
func (r Rect) West() float64 {
	return r[0][0]
}

// SouthWest returns the lower left point of the rect.
func (r Rect) SouthWest() Point {
	return NewPoint(r[0][0], r[0][1])
}

// NorthEast return the upper right point of the rect.
func (r Rect) NorthEast() Point {
	return NewPoint(r[1][0], r[1][1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (r Rect) IsEmpty() bool {
	return r[0].Lon() > r[1].Lon() || r[0].Lat() > r[1].Lat()
}

// IsZero return true if the rect just includes just null island.
func (r Rect) IsZero() bool {
	return r == Rect{}
}

// String returns the string respentation of the rect in WKT format.
func (r Rect) String() string {
	return r.WKT()
}

// WKT returns the string respentation of the rect in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (r Rect) WKT() string {
	if r.IsEmpty() {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "POLYGON(")
	wktPoints(buff, r.ToLineString())
	fmt.Fprintf(buff, ")")

	return buff.String()
}

// Pointer is a helper for using rectangle in a nullable context.
func (r Rect) Pointer() *Rect {
	return &r
}

// Equal returns if two rectangles are equal.
func (r Rect) Equal(c Rect) bool {
	return r[0] == c[0] && r[1] == c[1]
}
