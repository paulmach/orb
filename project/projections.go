package project

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
)

const earthRadius = 6378137.0
const earthRadiusPi = 6378137.0 * math.Pi

// Transformation functions that define how projections work.
type (
	// A Projection can transform between both planar and geo spaces.
	Projection struct {
		// ToXY is a function that projects from WGS84 (lon/lat) to the projection.
		ToXY orb.Projection

		// ToWGS84 is a function that projects from the projection to WGS84 (lon/lat).
		ToWGS84 orb.Projection
	}
)

// Mercator performs the Spherical Pseudo-Mercator projection used by most web maps.
var Mercator = struct {
	ToWGS84 orb.Projection
}{
	ToWGS84: func(p orb.Point) orb.Point {
		return orb.Point{
			180.0 * p[0] / earthRadiusPi,
			180.0 / math.Pi * (2*math.Atan(math.Exp(p[1]/earthRadius)) - math.Pi/2.0),
		}
	},
}

// WGS84 is what common uses lon/lat projection.
var WGS84 = struct {
	// ToMercator projections from WGS to Mercator, used by most web maps
	ToMercator orb.Projection
}{
	ToMercator: func(g orb.Point) orb.Point {
		y := math.Log(math.Tan((90.0+g[1])*math.Pi/360.0)) * earthRadius
		return orb.Point{
			earthRadiusPi / 180.0 * g[0],
			math.Max(-earthRadiusPi, math.Min(y, earthRadiusPi)),
		}
	},
}

// MercatorScaleFactor returns the mercator scaling factor for a given degree latitude.
func MercatorScaleFactor(g orb.Point) float64 {
	if g[1] < -90.0 || g[1] > 90.0 {
		panic(fmt.Sprintf("orb: latitude out of range, given %f", g[1]))
	}

	return 1.0 / math.Cos(g[1]/180.0*math.Pi)
}
