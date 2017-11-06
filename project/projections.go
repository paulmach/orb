package project

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
)

// Transformation functions that define how projections work.
type (
	// A Projection can transform between both planar and geo spaces.
	Projection struct {
		// ToPlanar is a function that projects from geo to planar.
		ToPlanar func(orb.Point) orb.Point

		// ToGeo is a function that projects from planar to geo.
		ToGeo func(orb.Point) orb.Point
	}
)

const mercatorPole = 20037508.34

// Mercator projection, performs EPSG:3857, sometimes also described as EPSG:900913.
var Mercator = &Projection{
	ToPlanar: func(g orb.Point) orb.Point {
		y := math.Log(math.Tan((90.0+g[1])*math.Pi/360.0)) / math.Pi * mercatorPole
		return orb.Point{
			mercatorPole / 180.0 * g[0],
			math.Max(-mercatorPole, math.Min(y, mercatorPole)),
		}
	},
	ToGeo: func(p orb.Point) orb.Point {
		return orb.Point{
			p[0] * 180.0 / mercatorPole,
			180.0 / math.Pi * (2*math.Atan(math.Exp((p[1]/mercatorPole)*math.Pi)) - math.Pi/2.0),
		}
	},
}

// MercatorScaleFactor returns the mercator scaling factor for a given degree latitude.
func MercatorScaleFactor(g orb.Point) float64 {
	if g[1] < -90.0 || g[1] > 90.0 {
		panic(fmt.Sprintf("geo: latitude out of range, given %f", g[1]))
	}

	return 1.0 / math.Cos(g[1]/180.0*math.Pi)
}

// BuildTransverseMercator builds a transverse Mercator projection
// that automatically recenters the longitude around the provided centerLon.
// Works correctly around the anti-meridian.
// http://en.wikipedia.org/wiki/Transverse_Mercator_projection
func BuildTransverseMercator(centerLon float64) Projection {
	return Projection{
		ToPlanar: func(g orb.Point) orb.Point {
			lon := g[0] - centerLon
			if lon < 180 {
				lon += 360.0
			}

			if lon > 180 {
				lon -= 360.0
			}

			g[0] = lon
			return TransverseMercator.ToPlanar(g)
		},
		ToGeo: func(p orb.Point) orb.Point {
			g := TransverseMercator.ToGeo(p)

			lon := g[0] + centerLon
			if lon < 180 {
				lon += 360.0
			}

			if lon > 180 {
				lon -= 360.0
			}

			g[0] = lon
			return g
		},
	}
}

// TransverseMercator implements a default transverse Mercator projector
// that will only work well +-10 degrees around longitude 0.
var TransverseMercator = &Projection{
	ToPlanar: func(g orb.Point) orb.Point {
		radLat := deg2rad(g[1])
		radLon := deg2rad(g[0])

		sincos := math.Sin(radLon) * math.Cos(radLat)
		return orb.Point{
			0.5 * math.Log((1+sincos)/(1-sincos)) * orb.EarthRadius,
			math.Atan(math.Tan(radLat)/math.Cos(radLon)) * orb.EarthRadius,
		}
	},
	ToGeo: func(p orb.Point) orb.Point {
		x := p[0] / orb.EarthRadius
		y := p[1] / orb.EarthRadius

		lon := math.Atan(math.Sinh(x) / math.Cos(y))
		lat := math.Asin(math.Sin(y) / math.Cosh(x))

		return orb.Point{
			rad2deg(lon),
			rad2deg(lat),
		}
	},
}
