package planar

import (
	"bytes"
	"fmt"
	"math"
)

// Polygon is a closed area. The first Ring is the outer ring.
// The others are the holes. Each Ring is expected to be closed
// ie. the first point matches the last.
type Polygon []Ring

// NewPolygon creates a new Polygon.
func NewPolygon() Polygon {
	return Polygon{}
}

// DistanceFrom will return the distance from the point to
// the polygon. Returns 0 if the point is within the polygon.
func (p Polygon) DistanceFrom(point Point) float64 {
	dist := p[0].DistanceFrom(point)
	if dist != 0 {
		// outside so just return that distance.
		return dist
	}

	for i := 1; i < len(p); i++ {
		if p[i].Contains(point) {
			return LineString(p[i]).DistanceFrom(point)
		}
	}

	// within the polygon, but not within any of the holes.
	return 0
}

// Centroid computes the area based centroid of the polygon.
// The algorithm removes the contribution of the holes.
func (p Polygon) Centroid() Point {
	point, _ := p.CentroidArea()
	return point
}

// CentroidArea computes the centroid and returns the area.
// If you need both this is faster since we need to area to compute the centroid.
func (p Polygon) CentroidArea() (Point, float64) {
	centroid, area := p[0].CentroidArea()
	area = math.Abs(area)

	holeArea := 0.0
	holeCentroid := Point{}
	for i := 1; i < len(p); i++ {
		ring := p[i]

		hc, ha := ring.CentroidArea()
		holeArea += math.Abs(ha)
		holeCentroid[0] += hc[0] * ha
		holeCentroid[1] += hc[1] * ha
	}

	totalArea := area - holeArea

	centroid[0] = (area*centroid[0] - holeArea*holeCentroid[0]) / totalArea
	centroid[1] = (area*centroid[1] - holeArea*holeCentroid[1]) / totalArea

	return centroid, totalArea
}

// Contains checks if the point is within the polygon.
// Points on the boundary are considered in.
func (p Polygon) Contains(point Point) bool {
	c := p[0].Contains(point)
	if !c {
		return false
	}

	for i := 1; i < len(p); i++ {
		if p[i].Contains(point) {
			return false
		}
	}

	return true
}

// Area computes the positive area of the polygon minus the area
// of the holes.
func (p Polygon) Area() float64 {
	if len(p) == 0 {
		return 0
	}

	area := p[0].Area()

	for i := 1; i < len(p); i++ {
		// minus holes
		area -= p[i].Area()
	}

	return area
}

// Bound returns a bound around the polygon.
func (p Polygon) Bound() Bound {
	return p[0].Bound()
}

// WKT returns the polygon in WKT format, eg. POlYGON((0 0,1 0,1 1,0 0))
// For empty polygons the result will be 'EMPTY'.
func (p Polygon) WKT() string {
	if len(p) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "POLYGON(")
	wktPoints(buff, p[0])

	for i := 1; i < len(p); i++ {
		buff.Write([]byte(","))
		wktPoints(buff, p[i])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns the wkt representation of the polygon.
func (p Polygon) String() string {
	return p.WKT()
}

// Equal compares two polygons. Returns true if lengths are the same
// and all points are Equal.
func (p Polygon) Equal(polygon Polygon) bool {
	if len(p) != len(polygon) {
		return false
	}

	for i, r := range p {
		if !r.Equal(polygon[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the polygon.
// All of the rings are also cloned.
func (p Polygon) Clone() Polygon {
	np := make(Polygon, 0, len(p))
	for _, r := range p {
		np = append(np, r.Clone())
	}

	return np
}
