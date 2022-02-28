# encoding/ewkb [![Godoc Reference](https://pkg.go.dev/badge/github.com/paulmach/orb)](https://pkg.go.dev/github.com/paulmach/orb/encoding/ewkb)

This package provides encoding and decoding of extended [WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Format_variations)
data. This format includes the SRID in the data. If you don't need the SRID you can use the
[wkb](../wkb) package for a simpler interface.

The interface is defined as:

```go
func Marshal(geom orb.Geometry, srid int, byteOrder ...binary.ByteOrder) ([]byte, error)
func MustMarshal(geom orb.Geometry, srid int, byteOrder ...binary.ByteOrder) []byte

func NewEncoder(w io.Writer) *Encoder
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) *Encoder
func (e *Encoder) SetSRID(srid int) *Encoder
func (e *Encoder) Encode(geom orb.Geometry) error

func Unmarshal(b []byte) (orb.Geometry, int, error)

func NewDecoder(r io.Reader) *Decoder
func (d *Decoder) Decode() (orb.Geometry, int, error)
```
