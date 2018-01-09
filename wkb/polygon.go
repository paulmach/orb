package wkb

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/paulmach/orb"
)

func readPolygon(r io.Reader, bom binary.ByteOrder) (orb.Polygon, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	result := make(orb.Polygon, 0, num)
	for i := 0; i < int(num); i++ {
		ls, err := readLineString(r, bom)
		if err != nil {
			return nil, err
		}

		result = append(result, orb.Ring(ls))
	}

	return result, nil
}

func readMultiPolygon(r io.Reader, bom binary.ByteOrder) (orb.MultiPolygon, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	result := make(orb.MultiPolygon, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != polygonType {
			return nil, errors.New("expect multipolygons to contains polygons, did not find a polygon")
		}

		p, err := readPolygon(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}
