# encoding/wkt [![Godoc Reference](https://pkg.go.dev/badge/github.com/paulmach/orb)](https://pkg.go.dev/github.com/paulmach/orb/encoding/wkt)

This package provides encoding and decoding of [WKT](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry)
data. The interface is defined as:

```go
func MarshalString(g orb.Geometry) string

func UnmarshalCollection(s string) (p orb.Collection, err error)
func UnmarshalLineString(s string) (p orb.LineString, err error)
func UnmarshalMultiLineString(s string) (p orb.MultiLineString, err error)
func UnmarshalMultiPoint(s string) (p orb.MultiPoint, err error)
func UnmarshalMultiPolygon(s string) (p orb.MultiPolygon, err error)
func UnmarshalPoint(s string) (p orb.Point, err error)
func UnmarshalPolygon(s string) (p orb.Polygon, err error)
```
