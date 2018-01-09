package wkb

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/paulmach/orb"
)

func readLineString(r io.Reader, bom binary.ByteOrder) (orb.LineString, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	result := make(orb.LineString, 0, num)
	for i := 0; i < int(num); i++ {
		p, err := readPoint(r, bom)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func readMultiLineString(r io.Reader, bom binary.ByteOrder) (orb.MultiLineString, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	result := make(orb.MultiLineString, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != lineStringType {
			return nil, errors.New("expect multilines to contains lines, did not find a line")
		}

		ls, err := readLineString(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, ls)
	}

	return result, nil
}
