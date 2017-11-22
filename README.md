orb [![Build Status](https://travis-ci.org/paulmach/orb.png?branch=master)](https://travis-ci.org/paulmach/orb) [![Godoc Reference](https://godoc.org/github.com/paulmach/orb?status.png)](https://godoc.org/github.com/paulmach/orb)
======

Orb is a package for working with geo and planar/projected geometric data is golang.
It supports the following types:

	type Point [2]float64
	type MultiPoint []Point

	type LineString []Point
	type MultiLineString []LineString

	type Ring LineString
	type Polygon []Ring
	type MultiPolygon []Polygon

	type Collection []Geometry

	type Bound struct { Min, Max Point }

All of these types match the `orb.Geometry` interface which is defined as:

	type Geometry interface {
		GeoJSONType() string
		Dimensions() int // e.g. 0d, 1d, 2d
		Bound() Bound
	}

Only a few methods are defined directly on these type, for example `Clone`, `Equal`, `GeoJSONType`.
Other operation that depend on geo vs. planar contexts are defined in the respective sub package.
For example:

* Computing the geo distance between two point:

		p1 := orb.Point{-72.796408, -45.407131}
		p2 := orb.Point{-72.688541, -45.384987}

		geo.Distance(p1, p2)

* Compute the planar area and centroid of a polygon:

		poly := orb.Polygon{...}
		centroid, area := planar.CentroidArea(poly)

### Other Subpackages

* [`geojson`](geojson) - working with geojson and the types in this package.
* [`clip`](clip) - clipping geometry to a bounding box
* [`maptile`](maptile) - working with mercator map tiles
* [`resample`](resample) - resample points in a line string geometry.
* [`project`](project) - project geometries between geo and planar contexts
