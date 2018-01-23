encoding/mvt [![Godoc Reference](https://godoc.org/github.com/paulmach/orb?status.png)](https://godoc.org/github.com/paulmach/orb/encoding/mvt)
============

Package mvt provides functions for encoding and decoding
[Mapbox Vector Tiles](https://www.mapbox.com/vector-tiles/specification/).
The interface is defined as:

	type Layer struct {
		Name     string
		Version  uint32
		Extent   uint32
		Features []*geojson.Feature
	}

	func MarshalGzipped(layers Layers) ([]byte, error)
	func Marshal(layers Layers) ([]byte, error)

	func UnmarshalGzipped(data []byte) (Layers, error)
	func Unmarshal(data []byte) (Layers, error)

These function decode the geometry and leave it in the "tile coordinates".
To project it to and from WGS84 (standard lon/lat) use:

	func (l Layer) ProjectToTile(tile maptile.Tile)
	func (l Layer) ProjectToWGS84(tile maptile.Tile)
