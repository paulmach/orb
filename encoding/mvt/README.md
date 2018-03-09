encoding/mvt [![Godoc Reference](https://godoc.org/github.com/paulmach/orb?status.svg)](https://godoc.org/github.com/paulmach/orb/encoding/mvt)
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

### Encoding example

```go
// Start with a set of feature collections defining each layer in lon/lat (WGS84).
collections := map[string]*geojson.FeatureCollection{}

// Convert to a layers object and project to tile coordinates.
layers := mvt.NewLayers(collections)
layers.ProjectToTile(maptile.New(x, y, z))

// In order to be used as source for MapboxGL geometries need to be clipped
// to max allowed extent. (uncomment next line)
// layers.Clip(mvt.MapboxGLDefaultExtentBound)

// Simplify the geometry now that it's in the tile coordinate space.
layers.Simplify(simplify.DouglasPeucker(1.0))

// Depending on use-case remove empty geometry, those two small to be
// represented in this tile space.
// In this case lines shorter than 1, and areas smaller than 1.
layers.RemoveEmpty(1.0, 1.0)

// encoding using the Mapbox Vector Tile protobuf encoding.
data, err := layers.Marshal() // this data is NOT gzipped.

// Sometimes MVT data is stored and transfered gzip compressed. In that case:
data, err := layers.MarshalGzipped()
```
