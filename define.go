package orb

// EarthRadius is the radius of the earth in meters. It is used in geo distance calculations.
// To keep things consistent, this value matches WGS84 Web Mercator (EPSG:3857).
const EarthRadius = 6378137.0 // meters

// Orientation defines the order of the points in a polygon
// or closed ring.
type Orientation int8

// Constants to define orientation.
// They follow the right hand rule for orientation.
const (
	// CCW stands for Counter Clock Wise
	CCW Orientation = 1

	// CW stands for Clock Wise
	CW Orientation = -1
)

// A DistanceFunc is a function that computes the distance between two points.
type DistanceFunc func(Point, Point) float64

// A Projection a function that moves a point from one space to another.
type Projection func(Point) Point
