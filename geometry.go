package orb

// Geometry is an interface that represents the shared attributes
// of a geometry.
type Geometry interface {
	GeoJSONType() string
	Dimensions() int // e.g. 0d, 1d, 2d
	Bound() Bound

	// requiring because sub package type switch over all possible types.
	private()
}

// compile time checks
var (
	_ Geometry = Point{}
	_ Geometry = MultiPoint{}
	_ Geometry = LineString{}
	_ Geometry = MultiLineString{}
	_ Geometry = Ring{}
	_ Geometry = Polygon{}
	_ Geometry = MultiPolygon{}
	_ Geometry = Bound{}

	_ Geometry = Collection{}
)

func (p Point) private()             {}
func (mp MultiPoint) private()       {}
func (ls LineString) private()       {}
func (mls MultiLineString) private() {}
func (r Ring) private()              {}
func (p Polygon) private()           {}
func (mp MultiPolygon) private()     {}
func (b Bound) private()             {}
func (c Collection) private()        {}

// A Collection is a collection of geometries that is also a Geometry.
type Collection []Geometry

// GeoJSONType returns the geometry collection type.
func (c Collection) GeoJSONType() string {
	return "GeometryCollection"
}

// Dimensions returns the max of the dimensions of the collection.
func (c Collection) Dimensions() int {
	max := -1
	for _, g := range c {
		if d := g.Dimensions(); d > max {
			max = d
		}
	}

	return max
}

// Bound returns the bounding box of all the Geometries combined.
func (c Collection) Bound() Bound {
	r := c[0].Bound()
	for i := 1; i < len(c); i++ {
		r = r.Union(c[i].Bound())
	}

	return r
}

// Clone returns a deep copy of the collection.
func (c Collection) Clone() Collection {
	if c == nil {
		return nil
	}

	nc := make(Collection, len(c))
	for i, g := range c {
		nc[i] = Clone(g)
	}

	return nc
}
