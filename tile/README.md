orb/tile [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/tile?status.png)](https://godoc.org/github.com/paulmach/orb/tile)
========

Package orb/tile provides types and methods for working with
[web mercator map tiles](https://www.google.com/search?q=web+mercator+map+tiles).
It defines a tile as:

	type Tile struct {
		X, Y, Z uint32
	}

Functions are provided to create tiles from lon/lat points as well as
[quadkeys](https://msdn.microsoft.com/en-us/library/bb259689.aspx).
The tile define helper methods such as `Parent()`, `Children()`, `Siblings()`, etc.

### tilecover sub-package

Still a work in progress but the goal is to provide geo.Geometry -> covering tile functions.

#### Similar libraries in other languages:

* [mercantile](https://github.com/mapbox/mercantile) - Python
* [sphericalmercator](https://github.com/mapbox/sphericalmercator) - Node
* [tilebelt](https://github.com/mapbox/tilebelt) - Node
