package planar

import (
	"math"

	"github.com/paulmach/orb"
)

const (
	// EarthRadius is the Earth's radius in kilometers
	EarthRadius = 6371
)

// Distance returns the distance between two points in 2d euclidean geometry.
func Distance(p1, p2 orb.Point) float64 {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

// DistanceSquared returns the square of the distance between two points in 2d euclidean geometry.
func DistanceSquared(p1, p2 orb.Point) float64 {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return d0*d0 + d1*d1
}

// HaversineDistance calculates the distance between two points on earth in kilometers
// Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func HaversineDistance(p1, p2 orb.Point) float64 {
	dLat := (p2[1] - p1[1]) * (math.Pi / 180.0)
	dLon := (p2[0] - p1[0]) * (math.Pi / 180.0)

	lat1 := p1[1] * (math.Pi / 180.0)
	lat2 := p2[1] * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}
