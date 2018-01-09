// Package wkb is for decoding ESRI's Well Known Binary (WKB) format for OGC geometry (WKBGeometry)
// sepcification at http://edndoc.esri.com/arcsde/9.1/general_topics/wkb_representation.htm
// There are a few types supported by the specification. Each general type is in it's own file.
// So, to find the implementation of Point (and MultiPoint) it will be located in the point.go
// file. Each of the basic type here adhere to the tegola.Geometry interface. So, a wkb point
// is, also, a tegola.Point
package wkb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/paulmach/orb"
)

// geometry types
// http://edndoc.esri.com/arcsde/9.1/general_topics/wkb_representation.htm
const (
	pointType              uint32 = 1
	lineStringType         uint32 = 2
	polygonType            uint32 = 3
	multiPointType         uint32 = 4
	multiLineStringType    uint32 = 5
	multiPolygonType       uint32 = 6
	geometryCollectionType uint32 = 7
)

func encode(bom binary.ByteOrder, geometry orb.Geometry) []interface{} {
	var data []interface{}
	if bom == binary.LittleEndian {
		data = append(data, byte(1))
	} else {
		data = append(data, byte(0))
	}

	switch g := geometry.(type) {
	case orb.Point:
		return append(data, pointType, g.X(), g.Y())
	case orb.MultiPoint:
		data = append(data, multiPointType)
		if len(g) == 0 {
			return data
		}
		for _, p := range g {
			data = append(data, encode(bom, p)...)
		}
		return data
	case orb.LineString:
		data = append(data, lineStringType, uint32(len(g)))
		for i := range g {
			data = append(data, g[i]) // The points.
		}
		return data

	case orb.MultiLineString:
		data = append(data, multiLineStringType)
		data = append(data, len(g)) // Number of lines in the Multi line string
		for _, l := range g {
			ld := encode(bom, l)
			if ld == nil {
				return nil
			}
			data = append(data, ld...)
		}
		return data
	case orb.Ring:
		return encode(bom, orb.Polygon{g})
	case orb.Polygon:
		data = append(data, polygonType)
		data = append(data, uint32(len(g))) // Number of rings in the polygon
		for _, r := range g {
			data = append(data, uint32(len(r))) // Number of points in the ring
			for _, p := range r {
				data = append(data, p) // The points in the ring
			}
		}
		return data
	case orb.MultiPolygon:
		data = append(data, multiPolygonType)
		data = append(data, uint32(len(g))) // Number of Polygons in the Multi.
		for _, p := range g {
			pd := encode(bom, p)
			if pd == nil {
				return nil
			}
			data = append(data, pd...)
		}
		return data
	case orb.Collection:
		data = append(data, geometryCollectionType)
		data = append(data, uint32(len(g))) // Number of Geometries
		for _, geom := range g {
			gd := encode(bom, geom)
			if gd == nil {
				return nil
			}
			data = append(data, gd...)
		}
		return data
	case orb.Bound:
		return encode(bom, g.ToPolygon())
	}

	panic("unsupported type")
}

// Marshal encodes the geometry with the given byte order.
func Marshal(geom orb.Geometry, byteOrder binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := Write(buf, byteOrder, geom)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Write will encode the given Geometry as a binary representation with the given
// byte order, and write it to the provided io.Writer.
func Write(w io.Writer, byteOrder binary.ByteOrder, geom orb.Geometry) error {
	if geom == nil {
		return nil
	}

	data := encode(byteOrder, geom)
	if data == nil {
		return errors.New("unabled to encode")
	}

	return binary.Write(w, byteOrder, data)
}

// Unmarshal will decode the type into a Geometry
func Unmarshal(b []byte) (orb.Geometry, error) {
	return Read(bytes.NewReader(b))
}

// Read is the main function that given a io.Reader will attempt to decode the
// Geometry from the byte stream.
func Read(r io.Reader) (orb.Geometry, error) {
	byteOrder, typ, err := readByteOrderType(r)
	if err != nil {
		return nil, err
	}

	switch typ {
	case pointType:
		return readPoint(r, byteOrder)
	case multiPointType:
		return readMultiPoint(r, byteOrder)
	case lineStringType:
		return readLineString(r, byteOrder)
	case multiLineStringType:
		return readMultiLineString(r, byteOrder)
	case polygonType:
		return readPolygon(r, byteOrder)
	case multiPolygonType:
		return readMultiPolygon(r, byteOrder)
	case geometryCollectionType:
		return readCollection(r, byteOrder)
	}

	return nil, fmt.Errorf("unknown geometry: %v", typ)
}

func readByteOrderType(r io.Reader) (binary.ByteOrder, uint32, error) {
	var bom = make([]byte, 1)

	// the bom is the first byte
	if _, err := r.Read(bom); err != nil {
		return nil, 0, err
	}

	var byteOrder binary.ByteOrder
	if bom[0] == 0 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}

	// the type which is 4 bytes
	var typ uint32

	err := binary.Read(r, byteOrder, &typ)
	if err != nil {
		return nil, 0, err
	}

	return byteOrder, typ, err
}
