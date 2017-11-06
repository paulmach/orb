package geo

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

// A Point is a Lon/Lat 2d point.
type Point [2]float64

// NewPoint creates a new point.
func NewPoint(lon, lat float64) Point {
	return Point{lon, lat}
}

// GeoJSONType returns the GeoJSON type for the object.
func (p Point) GeoJSONType() string {
	return "Point"
}

// Dimensions returns 0 because a point is a 0d object.
func (p Point) Dimensions() int {
	return 0
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
func (p Point) Quadkey(level uint32) uint64 {
	x, y := mercator.ToPlanar(p.Lon(), p.Lat(), level)

	var i, result uint64
	for i = 0; i < uint64(level); i++ {
		result |= (uint64(x) & (1 << i)) << i
		result |= (uint64(y) & (1 << i)) << (i + 1)
	}

	return result
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
