package orb

// EarthRadius is the radius of the earth in meters. It is used in geo distance calculations.
// To keep things consistent, this value matches WGS84 Web Mercator (EPSG:3857).
const EarthRadius = 6378137.0 // meters

// Pointer is an interface for things that can express themselves as
// generic x, y coordinates. The function is not named Point so that
// {geo,planar}.Point types can be embedded.
type Pointer interface {
	XY() (x, y float64)
}

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
