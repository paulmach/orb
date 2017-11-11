package orb

// A Point is a Lon/Lat 2d point.
type Point [2]float64

// NewPoint creates a new point.
func NewPoint(lon, lat float64) Point {
	return Point{lon, lat}
}

// GeoJSONType returns the GeoJSON type for the object.
func (p Point) GeoJSONType() string {
	return "Point"
}

// Dimensions returns 0 because a point is a 0d object.
func (p Point) Dimensions() int {
	return 0
}

// Bound returns a single point bound of the point.
func (p Point) Bound() Bound {
	return Bound{p, p}
}

// Equal checks if the point represents the same point or vector.
func (p Point) Equal(point Point) bool {
	return p[0] == point[0] && p[1] == point[1]
}
