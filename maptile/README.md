orb/tile [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/maptile?status.svg)](https://godoc.org/github.com/paulmach/orb/maptile)
========

Package orb/maptile provides types and methods for working with
[web mercator map tiles](https://www.google.com/search?q=web+mercator+map+tiles).
It defines a tile as:

	type Tile struct {
		X, Y uint32
		Z    Zoom
	}

	type Zoom uint32

Functions are provided to create tiles from lon/lat points as well as
[quadkeys](https://msdn.microsoft.com/en-us/library/bb259689.aspx).
The tile defines helper methods such as `Parent()`, `Children()`, `Siblings()`, etc.

### tilecover sub-package

Still a work in progress but the goal is to provide geo.Geometry -> covering tiles.

#### Similar libraries in other languages:

* [mercantile](https://github.com/mapbox/mercantile) - Python
* [sphericalmercator](https://github.com/mapbox/sphericalmercator) - Node
* [tilebelt](https://github.com/mapbox/tilebelt) - Node
