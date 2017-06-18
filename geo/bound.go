package geo

import (
	"math"

	"github.com/paulmach/orb"
)

// A Bound represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Bound [2]Point

// NewBound creates a new bound given the parameters.
func NewBound(west, east, south, north float64) Bound {
	return Bound{
		Point{west, south},
		Point{east, north},
	}
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{corner, corner}.Extend(oppositeCorner)
}

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point in meters.
func NewBoundAroundPoint(center Point, distance float64) Bound {
	radDist := distance / orb.EarthRadius
	radLat := deg2rad(center.Lat())
	radLon := deg2rad(center.Lon())
	minLat := radLat - radDist
	maxLat := radLat + radDist

	var minLon, maxLon float64
	if minLat > minLatitude && maxLat < maxLatitude {
		deltaLon := math.Asin(math.Sin(radDist) / math.Cos(radLat))
		minLon = radLon - deltaLon
		if minLon < minLongitude {
			minLon += 2 * math.Pi
		}
		maxLon = radLon + deltaLon
		if maxLon > maxLongitude {
			maxLon -= 2 * math.Pi
		}
	} else {
		minLat = math.Max(minLat, minLatitude)
		maxLat = math.Min(maxLat, maxLatitude)
		minLon = minLongitude
		maxLon = maxLongitude
	}

	return Bound{
		Point{rad2deg(minLon), rad2deg(minLat)},
		Point{rad2deg(maxLon), rad2deg(maxLat)},
	}
}

// GeoJSONType returns the GeoJSON type for the object.
func (b Bound) GeoJSONType() string {
	return "Polygon"
}

// ToPolygon converts the bound into a Polygon object.
func (b Bound) ToPolygon() Polygon {
	return Polygon{b.ToRing()}
}

// ToRing converts the bound into a loop defined
// by the boundary of the box.
func (b Bound) ToRing() Ring {
	r := make(Ring, 5)
	r[0] = b[0]
	r[1] = Point{b[0][0], b[1][1]}
	r[2] = b[1]
	r[3] = Point{b[1][0], b[0][1]}
	r[4] = b[0]

	return r
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b[0][0] = math.Min(b[0].Lon(), point.Lon())
	b[1][0] = math.Max(b[1].Lon(), point.Lon())

	b[0][1] = math.Min(b[0].Lat(), point.Lat())
	b[1][1] = math.Max(b[1].Lat(), point.Lat())

	return b
}

// Union extends this bound to contain the union of this and the given bound.
func (b Bound) Union(other Bound) Bound {
	b = b.Extend(other[0])
	b = b.Extend(other[1])
	b = b.Extend(Point{other[0][0], other[1][1]})
	b = b.Extend(Point{other[1][0], other[0][1]})

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	if point.Lat() < b[0].Lat() || b[1].Lat() < point.Lat() {
		return false
	}

	if point.Lon() < b[0].Lon() || b[1].Lon() < point.Lon() {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	if (b[1][0] < bound[0][0]) ||
		(b[0][0] > bound[1][0]) ||
		(b[1][1] < bound[0][1]) ||
		(b[0][1] > bound[1][1]) {
		return false
	}

	return true
}

// Pad expands the bound in all directions by the given amount of meters.
func (b Bound) Pad(meters float64) Bound {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(b[1].Lat()))
	dx = math.Max(dx, dy/math.Cos(deg2rad(b[0].Lat())))

	b[0][0] -= dx
	b[0][1] -= dy

	b[1][0] += dx
	b[1][1] += dy

	return b
}

// Center returns the center of the bounds by "averaging" the x and y coords.
func (b Bound) Center() Point {
	return Point{
		(b[0][0] + b[1][0]) / 2.0,
		(b[0][1] + b[1][1]) / 2.0,
	}
}

// Height returns the approximate height in meters.
func (b Bound) Height() float64 {
	return 111131.75 * (b[1][1] - b[0][1])
}

// Width returns the approximate width in meters
// of the center of the bound.
func (b Bound) Width(haversine ...bool) float64 {
	c := (b[0][1] + b[1][1]) / 2.0

	s1 := Point{b[0][0], c}
	s2 := Point{b[1][0], c}

	return s1.DistanceFrom(s2, yesHaversine(haversine))
}

// North returns the top of the bound.
func (b Bound) North() float64 {
	return b[1][1]
}

// South returns the bottom of the bound.
func (b Bound) South() float64 {
	return b[0][1]
}

// East returns the right of the bound.
func (b Bound) East() float64 {
	return b[1][0]
}

// West returns the left of the bound.
func (b Bound) West() float64 {
	return b[0][0]
}

// SouthWest returns the lower left point of the bound.
func (b Bound) SouthWest() Point {
	return NewPoint(b[0][0], b[0][1])
}

// NorthEast return the upper right point of the bound.
func (b Bound) NorthEast() Point {
	return NewPoint(b[1][0], b[1][1])
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b[0].Lon() > b[1].Lon() || b[0].Lat() > b[1].Lat()
}

// IsZero return true if the bound just includes just null island.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}

// String returns the string respentation of the bound in WKT format.
func (b Bound) String() string {
	return b.WKT()
}

// WKT returns the string respentation of the bound in WKT format.
// POLYGON(west, south, west, north, east, north, east, south, west, south)
func (b Bound) WKT() string {
	return b.ToPolygon().WKT()
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b[0] == c[0] && b[1] == c[1]
}
