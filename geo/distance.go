package geo

import (
	"math"

	"github.com/paulmach/orb"
)

// Distance returns the distance between two points on the earth.
func Distance(p1, p2 orb.Point) float64 {
	dLat := deg2rad(p1[1] - p2[1])
	dLon := deg2rad(p1[0] - p2[0])

	dLon = math.Abs(dLon)
	if dLon > math.Pi {
		dLon = 2*math.Pi - dLon
	}

	// fast way using pythagorean theorem on an equirectangular projection
	x := dLon * math.Cos(deg2rad((p1[1]+p2[1])/2.0))
	return math.Sqrt(dLat*dLat+x*x) * orb.EarthRadius
}

// DistanceHaversine computes the distance on the earth using the
// more accurate haversine formula.
func DistanceHaversine(p1, p2 orb.Point) float64 {
	dLat := deg2rad(p1[1] - p2[1])
	dLon := deg2rad(p1[0] - p2[0])

	dLat2Sin := math.Sin(dLat / 2)
	dLon2Sin := math.Sin(dLon / 2)
	a := dLat2Sin*dLat2Sin + math.Cos(deg2rad(p2[1]))*math.Cos(deg2rad(p1[1]))*dLon2Sin*dLon2Sin

	return 2.0 * orb.EarthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// Bearing computes the direction one must start traveling on earth
// to be heading from, to the given points.
func Bearing(from, to orb.Point) float64 {
	dLon := deg2rad(to[0] - from[0])

	fromLatRad := deg2rad(from[1])
	toLatRad := deg2rad(to[1])

	y := math.Sin(dLon) * math.Cos(toLatRad)
	x := math.Cos(fromLatRad)*math.Sin(toLatRad) - math.Sin(fromLatRad)*math.Cos(toLatRad)*math.Cos(dLon)

	return rad2deg(math.Atan2(y, x))
}

// Midpoint returns the half-way point along a great circle path between the two points.
func Midpoint(p, p2 orb.Point) orb.Point {
	dLon := deg2rad(p2[0] - p[0])

	aLatRad := deg2rad(p[1])
	bLatRad := deg2rad(p2[1])

	x := math.Cos(bLatRad) * math.Cos(dLon)
	y := math.Cos(bLatRad) * math.Sin(dLon)

	r := orb.Point{
		deg2rad(p[0]) + math.Atan2(y, math.Cos(aLatRad)+x),
		math.Atan2(math.Sin(aLatRad)+math.Sin(bLatRad), math.Sqrt((math.Cos(aLatRad)+x)*(math.Cos(aLatRad)+x)+y*y)),
	}

	// convert back to degrees
	r[0] = rad2deg(r[0])
	r[1] = rad2deg(r[1])

	return r
}
