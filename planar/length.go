package planar

import (
	"fmt"

	"github.com/paulmach/orb/geo"
)

// Length returns the length of the boundary of the geometry
// using 2d euclidean geometry.
func Length(g geo.Geometry) float64 {
	switch g := g.(type) {
	case geo.Point:
		return 0
	case geo.MultiPoint:
		return 0
	case geo.LineString:
		return lineStringLength(g)
	case geo.MultiLineString:
		sum := 0.0
		for _, ls := range g {
			sum += lineStringLength(ls)
		}

		return sum
	case geo.Ring:
		return lineStringLength(geo.LineString(g))
	case geo.Polygon:
		return polygonLength(g)
	case geo.MultiPolygon:
		sum := 0.0
		for _, p := range g {
			sum += polygonLength(p)
		}

		return sum
	case geo.Collection:
		sum := 0.0
		for _, c := range g {
			sum += Length(c)
		}

		return sum
	case geo.Bound:
		return Length(g.ToRing())
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

func lineStringLength(ls geo.LineString) float64 {
	sum := 0.0
	for i := 1; i < len(ls); i++ {
		sum += Distance(ls[i], ls[i-1])
	}

	return sum
}

func polygonLength(p geo.Polygon) float64 {
	sum := 0.0
	for _, r := range p {
		sum += lineStringLength(geo.LineString(r))
	}

	return sum
}
