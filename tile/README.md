orb/tile [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/tile?status.png)](https://godoc.org/github.com/paulmach/orb/tile)
========

orb/tile is a library for working with [web mercator map tiles](https://www.google.com/search?q=web+mercator+map+tiles).
It defines a tile as:
```
	type Tile struct {
		X, Y, Z uint64
	}
```

Functions are provided to create tiles from lon/lat points as well as
[quadkeys](https://msdn.microsoft.com/en-us/library/bb259689.aspx).
The tiles define helper methods such as `Parent()`, `Children()`, `Siblings()`, etc.

### Bound

The package also defines a `tile.Bound` to represent a rectangle of tiles at a given zoom.
In some cases this can be more useful than using a lon/lat bound, i.e. which tile does a
corner point represent? Some methods provided:

* `Contains(t Tile)`
* `Covering(zoom uint64) Tiles` - returns the set of tiles at the give zoom that cover the bound.

#### Similar libraries in other languages:

* [mercantile](https://github.com/mapbox/mercantile) - Python
* [sphericalmercator](https://github.com/mapbox/sphericalmercator) - Node
* [tilebelt](https://github.com/mapbox/tilebelt) - Node
