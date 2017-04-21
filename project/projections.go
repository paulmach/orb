package project

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/planar"
)

// Transformation functions that define how projections work.
type (
	// A Projection can transform between both planar and geo spaces.
	Projection struct {
		// Forward is a function that projects from geo to planar.
		Forward func(geo.Point) planar.Point

		// Reverse is a function that projects from planar to geo.
		Reverse func(planar.Point) geo.Point
	}
)

const mercatorPole = 20037508.34

// Mercator projection, performs EPSG:3857, sometimes also described as EPSG:900913.
var Mercator = Projection{
	Forward: func(g geo.Point) planar.Point {
		y := math.Log(math.Tan((90.0+g.Lat())*math.Pi/360.0)) / math.Pi * mercatorPole
		return planar.Point{
			mercatorPole / 180.0 * g.Lon(),
			math.Max(-mercatorPole, math.Min(y, mercatorPole)),
		}
	},
	Reverse: func(p planar.Point) geo.Point {
		return geo.Point{
			p.X() * 180.0 / mercatorPole,
			180.0 / math.Pi * (2*math.Atan(math.Exp((p.Y()/mercatorPole)*math.Pi)) - math.Pi/2.0),
		}
	},
}

// MercatorScaleFactor returns the mercator scaling factor for a given degree latitude.
func MercatorScaleFactor(g geo.Point) float64 {
	if g.Lat() < -90.0 || g.Lat() > 90.0 {
		panic(fmt.Sprintf("geo: latitude out of range, given %f", g.Lat()))
	}

	return 1.0 / math.Cos(g.Lat()/180.0*math.Pi)
}

// BuildTransverseMercator builds a transverse Mercator projection
// that automatically recenters the longitude around the provided centerLon.
// Works correctly around the anti-meridian.
// http://en.wikipedia.org/wiki/Transverse_Mercator_projection
func BuildTransverseMercator(centerLon float64) Projection {
	return Projection{
		Forward: func(g geo.Point) planar.Point {
			lon := g.Lon() - centerLon
			if lon < 180 {
				lon += 360.0
			}

			if lon > 180 {
				lon -= 360.0
			}

			g[0] = lon
			return TransverseMercator.Forward(g)
		},
		Reverse: func(p planar.Point) geo.Point {
			g := TransverseMercator.Reverse(p)

			lon := g.Lon() + centerLon
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
var TransverseMercator = Projection{
	Forward: func(g geo.Point) planar.Point {
		radLat := deg2rad(g.Lat())
		radLon := deg2rad(g.Lon())

		sincos := math.Sin(radLon) * math.Cos(radLat)
		return planar.Point{
			0.5 * math.Log((1+sincos)/(1-sincos)) * orb.EarthRadius,
			math.Atan(math.Tan(radLat)/math.Cos(radLon)) * orb.EarthRadius,
		}
	},
	Reverse: func(p planar.Point) geo.Point {
		x := p.X() / orb.EarthRadius
		y := p.Y() / orb.EarthRadius

		lon := math.Atan(math.Sinh(x) / math.Cos(y))
		lat := math.Asin(math.Sin(y) / math.Cosh(x))

		return geo.Point{
			rad2deg(lon),
			rad2deg(lat),
		}
	},
}

// ScalarMercator converts from lon/lat float64 to x,y uint64.
// This is the same as Google's world coordinates.
var ScalarMercator struct {
	Level   uint64
	Forward func(g geo.Point, level ...uint64) (x, y uint64)
	Reverse func(x, y uint64, level ...uint64) geo.Point
}

func init() {
	ScalarMercator.Level = 31
	ScalarMercator.Forward = func(g geo.Point, level ...uint64) (x, y uint64) {
		l := ScalarMercator.Level
		if len(level) != 0 {
			l = level[0]
		}
		return mercator.ScalarProject(g.Lon(), g.Lat(), l)
	}
	ScalarMercator.Reverse = func(x, y uint64, level ...uint64) geo.Point {
		l := ScalarMercator.Level
		if len(level) != 0 {
			l = level[0]
		}

		lon, lat := mercator.ScalarInverse(x, y, l)
		return geo.Point{lon, lat}
	}
}
