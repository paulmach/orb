package wkb

import (
	"math"

	"github.com/paulmach/orb"
)

func ValidatePoint(value interface{}) ([]byte, bool, error) {
	var err error

	data, ok := value.([]byte)
	if !ok {
		return nil, false, orb.ErrUnsupportedDataType
	}

	switch len(data) {
	case 0:
		// empty data, return empty go struct which in this case
		// would be [0,0]
		return nil, false, nil
	case 21:
		// the length of a point type in WKB
	case 25:
		// Most likely MySQL's SRID+WKB format.
		// However, could be a line string or multipoint with only one point.
		// But those would be invalid for parsing a point.
		data = data[4:]
	default:
		return nil, false, orb.ErrIncorrectGeometry
	}

	littleEndian, typeCode, err := ReadPrefix(data)
	if err != nil {
		return nil, false, err
	}

	if typeCode != 1 {
		return nil, false, orb.ErrIncorrectGeometry
	}

	return data[5:], littleEndian, nil

}

func ValidateLine(value interface{}) ([]byte, bool, uint32, error) {
	data, ok := value.([]byte)
	if !ok {
		return nil, false, 0, orb.ErrUnsupportedDataType
	}

	switch len(data) {
	case 0:
		return nil, false, 0, nil
	case 41:
		// the length of a 2 point linestring type in WKB
	case 45:
		// Most likely MySQL's SRID+WKB format.
		// However, could be some encoding of another type.
		// But those would be invalid for parsing a line.
		data = data[4:]
	default:
		return nil, false, 0, orb.ErrIncorrectGeometry
	}

	littleEndian, typeCode, err := ReadPrefix(data)
	if err != nil {
		return nil, false, 0, err
	}

	if typeCode != 2 {
		return nil, false, 0, orb.ErrIncorrectGeometry
	}

	length := ReadUint32(data[5:9], littleEndian)
	if length != 2 {
		return nil, false, 0, orb.ErrIncorrectGeometry
	}

	return data[9:], littleEndian, length, nil
}

func ValidatePointSet(value interface{}) ([]byte, bool, int, error) {
	data, ok := value.([]byte)
	if !ok {
		return nil, false, 0, orb.ErrUnsupportedDataType
	}

	if len(data) == 0 {
		return nil, false, 0, nil
	}

	if len(data) < 6 {
		return nil, false, 0, orb.ErrNotWKB
	}

	// first byte of real WKB data indicates endian and should 1 or 0.
	if data[0] != 0 && data[0] != 1 {
		data = data[4:] // possibly mysql srid+wkb
	}

	littleEndian, typeCode, err := ReadPrefix(data)
	if err != nil {
		return nil, false, 0, err
	}

	// must be LineString, Polygon or MultiPoint
	if typeCode != 2 && typeCode != 3 && typeCode != 4 {
		return nil, false, 0, orb.ErrIncorrectGeometry
	}

	if typeCode == 3 {
		// For polygons there is a ring count.
		// We only allow one ring here.
		rings := int(ReadUint32(data[5:9], littleEndian))
		if rings != 1 {
			return nil, false, 0, orb.ErrIncorrectGeometry
		}

		data = data[9:]
	} else {
		data = data[5:]
	}

	length := int(ReadUint32(data[:4], littleEndian))
	if len(data) != 4+16*length {
		return nil, false, 0, orb.ErrNotWKB
	}

	return data[4:], littleEndian, length, nil
}

func ReadPrefix(data []byte) (bool, uint32, error) {
	if len(data) < 6 {
		return false, 0, orb.ErrNotWKB
	}

	if data[0] == 0 {
		return false, ReadUint32(data[1:5], false), nil
	}

	if data[0] == 1 {
		return true, ReadUint32(data[1:5], true), nil
	}

	return false, 0, orb.ErrNotWKB
}

func ReadUint32(data []byte, littleEndian bool) uint32 {
	var v uint32

	if littleEndian {
		for i := 3; i >= 0; i-- {
			v <<= 8
			v |= uint32(data[i])
		}
	} else {
		for i := 0; i < 4; i++ {
			v <<= 8
			v |= uint32(data[i])
		}
	}

	return v
}

func ReadFloat64(data []byte, littleEndian bool) float64 {
	var v uint64

	if littleEndian {
		for i := 7; i >= 0; i-- {
			v <<= 8
			v |= uint64(data[i])
		}
	} else {
		for i := 0; i < 8; i++ {
			v <<= 8
			v |= uint64(data[i])
		}
	}

	return math.Float64frombits(v)
}
