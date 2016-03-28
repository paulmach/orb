package planar

import (
	"fmt"
	"math"

	"github.com/paulmach/go.geojson"
)

// A Point is a simple X/Y or Lng/Lat 2d point. [X, Y] or [Lng, Lat]
type Point [2]float64

// NewPoint creates a new point
func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// DistanceFrom returns the Euclidean distance between the points.
func (p Point) DistanceFrom(point Point) float64 {
	d0 := (point[0] - p[0])
	d1 := (point[1] - p[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

// DistanceToSquared returns the squared Euclidean distance between the points.
// This avoids a sqrt computation.
func (p Point) DistanceFromSquared(point Point) float64 {
	d0 := (point[0] - p[0])
	d1 := (point[1] - p[1])
	return d0*d0 + d1*d1
}

// Add a point to the given point.
func (p Point) Add(vector Vector) Point {
	p[0] += vector[0]
	p[1] += vector[1]

	return p
}

// Sub a point from the given point.
func (p Point) Sub(point Point) Vector {
	return Vector{
		p[0] - point[0],
		p[1] - point[1],
	}
}

// Equal checks if the point represents the same point or vector.
func (p Point) Equal(point Point) bool {
	return p == point
}

// X returns the x/horizontal component of the point.
func (p Point) X() float64 {
	return p[0]
}

// Y returns the y/vertical component of the point.
func (p Point) Y() float64 {
	return p[1]
}

// GeoJSON creates a new geojson feature with a point geometry.
func (p Point) GeoJSON() *geojson.Feature {
	return geojson.NewPointFeature([]float64{p[0], p[1]})
}

// WKT returns the point in WKT format, eg. POINT(30.5 10.5)
func (p Point) WKT() string {
	return fmt.Sprintf("POINT(%g %g)", p[0], p[1])
}

// String returns a string representation of the point.
// The format is WKT, e.g. POINT(30.5 10.5)
func (p Point) String() string {
	return p.WKT()
}
