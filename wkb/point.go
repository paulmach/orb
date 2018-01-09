package wkb

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/paulmach/orb"
)

func readPoint(r io.Reader, bom binary.ByteOrder) (orb.Point, error) {
	var p orb.Point

	if err := binary.Read(r, bom, &p[0]); err != nil {
		return orb.Point{}, err
	}

	if err := binary.Read(r, bom, &p[1]); err != nil {
		return orb.Point{}, err
	}

	return p, nil
}

func readMultiPoint(r io.Reader, bom binary.ByteOrder) (orb.MultiPoint, error) {
	var num uint32 // Number of points.
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	result := make(orb.MultiPoint, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != pointType {
			return nil, errors.New("expect multipoint to contains points, did not find a point")
		}

		p, err := readPoint(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}
