orb/maptile/tilecover [![Godoc Reference](https://godoc.org/github.com/paulmach/orb/maptile/tilecover?status.png)](https://godoc.org/github.com/paulmach/orb/maptile/tilecover)
=====================

Package `tilecover` computes the covering set of tiles for an `orb.Geometry`.
It is a a port of the nodejs library [tile-cover](https://github.com/mapbox/tile-cover)
which is a port from Google's S2 library. The same set of tests pass.

### Usage

```
poly := orb.Polygon{}
tiles := tilecover.Geometry(poly, zoom)

for t := range tiles {
	// do something with tile
}

// to merge up to as much as possible to a specific zoom
tiles = tilecover.MergeUp(tiles, 0)
```

#### Similar libraries in other languages:

* [tilecover](https://github.com/mapbox/tile-cover) - Node
