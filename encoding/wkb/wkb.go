// Package wkb is for decoding ESRI's Well Known Binary (WKB) format
// sepcification at https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Well-known_binary
package wkb

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/internal/wkbcommon"
)

// byteOrder represents little or big endian encoding.
// We don't use binary.ByteOrder because that is an interface
// that leaks to the heap all over the place.
type byteOrder int

const bigEndian byteOrder = 0
const littleEndian byteOrder = 1

const (
	pointType              uint32 = 1
	lineStringType         uint32 = 2
	polygonType            uint32 = 3
	multiPointType         uint32 = 4
	multiLineStringType    uint32 = 5
	multiPolygonType       uint32 = 6
	geometryCollectionType uint32 = 7
)

const (
	// limits so that bad data can't come in and preallocate tons of memory.
	// Well formed data with less elements will allocate the correct amount just fine.
	maxPointsAlloc = 10000
	maxMultiAlloc  = 100
)

// DefaultByteOrder is the order used for marshalling or encoding
// is none is specified.
var DefaultByteOrder binary.ByteOrder = binary.LittleEndian

// An Encoder will encode a geometry as WKB to the writer given at
// creation time.
type Encoder struct {
	e *wkbcommon.Encoder
}

// MustMarshal will encode the geometry and panic on error.
// Currently there is no reason to error during geometry marshalling.
func MustMarshal(geom orb.Geometry, byteOrder ...binary.ByteOrder) []byte {
	d, err := Marshal(geom, byteOrder...)
	if err != nil {
		panic(err)
	}

	return d
}

// Marshal encodes the geometry with the given byte order.
func Marshal(geom orb.Geometry, byteOrder ...binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, wkbcommon.GeomLength(geom)))

	e := NewEncoder(buf)
	if len(byteOrder) > 0 {
		e.SetByteOrder(byteOrder[0])
	}

	err := e.Encode(geom)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, nil
	}

	return buf.Bytes(), nil
}

// NewEncoder creates a new Encoder for the given writer.
func NewEncoder(w io.Writer) *Encoder {
	e := wkbcommon.NewEncoder(w)
	e.SetByteOrder(DefaultByteOrder)
	return &Encoder{e: e}
}

// SetByteOrder will override the default byte order set when
// the encoder was created.
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) {
	e.e.SetByteOrder(bo)
}

// Encode will write the geometry encoded as WKB to the given writer.
func (e *Encoder) Encode(geom orb.Geometry) error {
	return e.e.Encode(geom)
}

// Decoder can decoder WKB geometry off of the stream.
type Decoder struct {
	d *wkbcommon.Decoder
}

// Unmarshal will decode the type into a Geometry.
func Unmarshal(data []byte) (orb.Geometry, error) {
	g, err := wkbcommon.Unmarshal(data)
	if err != nil {
		return nil, mapCommonError(err)
	}

	return g, nil
}

// NewDecoder will create a new WKB decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		d: wkbcommon.NewDecoder(r),
	}
}

// Decode will decode the next geometry off of the stream.
func (d *Decoder) Decode() (orb.Geometry, error) {
	g, err := d.d.Decode()
	if err != nil {
		return nil, mapCommonError(err)
	}

	return g, nil
}
