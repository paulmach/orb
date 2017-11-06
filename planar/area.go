package planar

import (
	"fmt"
	"math"

	"github.com/paulmach/orb/geo"
)

// Area returns the area of the geometry in the 2d plane.
func Area(g geo.Geometry) float64 {
	// TOOD: make faster non-centroid version.
	_, a := CentroidArea(g)
	return a
}

// CentroidArea returns botht the centroid and the area in the 2d plane.
// Since the area is need for the centroid, return both.
func CentroidArea(g geo.Geometry) (geo.Point, float64) {
	switch g := g.(type) {
	case geo.Point:
		return multiPointCentroid(geo.MultiPoint{g}), 0
	case geo.MultiPoint:
		return multiPointCentroid(g), 0
	case geo.LineString:
		return multiLineStringCentroid(geo.MultiLineString{g}), 0
	case geo.MultiLineString:
		return multiLineStringCentroid(g), 0
	case geo.Ring:
		return ringCentroidArea(g)
	case geo.Polygon:
		return polygonCentroidArea(g)
	case geo.MultiPolygon:
		return multiPolygonCentroidArea(g)
	case geo.Collection:
		panic("TODO")
	case geo.Bound:
		return CentroidArea(g.ToRing())
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

func multiPointCentroid(mp geo.MultiPoint) geo.Point {
	x, y := 0.0, 0.0
	for _, p := range mp {
		x += p[0]
		y += p[1]
	}

	num := float64(len(mp))
	return geo.Point{x / num, y / num}
}

func multiLineStringCentroid(mls geo.MultiLineString) geo.Point {
	point := geo.Point{}
	dist := 0.0

	if len(mls) == 0 {
		return geo.Point{}
	}

	for _, ls := range mls {
		c, d := lineStringCentroidDist(ls)
		if d == math.Inf(1) {
			continue
		}

		point[0] += c[0] * d
		point[1] += c[1] * d

		dist += d
	}

	if dist == math.Inf(1) {
		return geo.Point{}
	}

	point[0] /= dist
	point[1] /= dist

	return point
}

func lineStringCentroidDist(ls geo.LineString) (geo.Point, float64) {
	dist := 0.0
	point := geo.Point{}

	if len(ls) == 0 {
		return geo.Point{}, math.Inf(1)
	}

	// implicitly move everything to near the origin to help with roundoff
	offset := ls[0]
	for i := 0; i < len(ls)-1; i++ {
		p1 := geo.Point{
			ls[i][0] - offset[0],
			ls[i][1] - offset[1],
		}

		p2 := geo.Point{
			ls[i+1][0] - offset[0],
			ls[i+1][1] - offset[1],
		}

		d := Distance(p1, p2)

		point[0] += (p1[0] + p2[0]) / 2.0 * d
		point[1] += (p1[1] + p2[1]) / 2.0 * d
		dist += d
	}

	point[0] /= dist
	point[1] /= dist

	point[0] += ls[0][0]
	point[1] += ls[0][1]
	return point, dist
}

func ringCentroidArea(r geo.Ring) (geo.Point, float64) {
	centroid := geo.Point{}
	area := 0.0

	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		a := (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
		area += a

		centroid[0] += (r[i][0] + r[i+1][0] - 2*offsetX) * a
		centroid[1] += (r[i][1] + r[i+1][1] - 2*offsetY) * a
	}

	// no need to deal with first and last vertex since we "moved"
	// that point the origin (multiply by 0 == 0)

	area /= 2
	centroid[0] /= 6 * area
	centroid[1] /= 6 * area

	centroid[0] += offsetX
	centroid[1] += offsetY

	return centroid, area
}

func polygonCentroidArea(p geo.Polygon) (geo.Point, float64) {
	centroid, area := ringCentroidArea(p[0])
	area = math.Abs(area)

	holeArea := 0.0
	holeCentroid := geo.Point{}
	for i := 1; i < len(p); i++ {
		hc, ha := ringCentroidArea(p[i])

		holeArea += math.Abs(ha)
		holeCentroid[0] += hc[0] * ha
		holeCentroid[1] += hc[1] * ha
	}

	totalArea := area - holeArea

	centroid[0] = (area*centroid[0] - holeArea*holeCentroid[0]) / totalArea
	centroid[1] = (area*centroid[1] - holeArea*holeCentroid[1]) / totalArea

	return centroid, totalArea
}

func multiPolygonCentroidArea(mp geo.MultiPolygon) (geo.Point, float64) {
	point := geo.Point{}
	area := 0.0

	for _, p := range mp {
		c, a := polygonCentroidArea(p)

		point[0] += c[0] * a
		point[1] += c[1] * a

		area += a
	}

	point[0] /= area
	point[1] /= area

	return point, area
}
