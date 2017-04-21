package wkb

import (
	"math"

	"github.com/paulmach/orb"
)

const (
	wkbPoint      = 1
	wkbLineString = 2
	wkbMultiPoint = 4
)

// ValidatePoint checks the wkb input and returns x, y, isnull, err.
func ValidatePoint(value interface{}) (float64, float64, bool, error) {
	data, ok := value.([]byte)
	if !ok {
		return 0, 0, false, orb.ErrUnsupportedDataType
	}

	switch len(data) {
	case 0:
		// empty data, return empty go struct which in this case
		// would be [0,0]
		return 0, 0, true, nil
	case 21:
		// the length of a point type in WKB
	case 25:
		// Most likely MySQL's SRID+WKB format.
		// However, could be a line string or multipoint with only one point.
		// But those would be invalid for parsing a point.
		data = data[4:]
	default:
		return 0, 0, false, orb.ErrIncorrectGeometry
	}

	x, y, err := ReadPoint(data)
	return x, y, false, err
}

// ReadPoint reads the beginning of the data to find the wkb point.
func ReadPoint(data []byte) (float64, float64, error) {
	if len(data) < 21 {
		return 0, 0, orb.ErrIncorrectGeometry
	}

	littleEndian, typeCode, err := ReadHeader(data)
	if err != nil {
		return 0, 0, err
	}

	if typeCode != wkbPoint {
		return 0, 0, orb.ErrIncorrectGeometry
	}

	return ReadFloat64(data[5:13], littleEndian),
		ReadFloat64(data[13:], littleEndian),
		nil
}

// ValidateLine checks the wkb input for a two point linestring.
func ValidateLine(value interface{}) ([]byte, bool, error) {
	data, littleEndian, length, err := ValidatePath(value)
	if err != nil {
		return data, littleEndian, err
	}

	if length != 2 {
		return nil, false, orb.ErrIncorrectGeometry
	}

	return data, littleEndian, err
}

// ValidateMultiPoint checks wkb for a multipoint geometry.
func ValidateMultiPoint(value interface{}) ([]byte, bool, int, error) {
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

	data, littleEndian, err := validateSet(data, wkbMultiPoint)
	if err != nil {
		return nil, false, 0, err
	}

	length := int(ReadUint32(data[5:9], littleEndian))
	if len(data) != 9+21*length {
		return nil, false, 0, orb.ErrNotWKB
	}

	return data[9:], littleEndian, length, nil
}

// ValidatePath checks the wkb for a linestring geometry.
func ValidatePath(value interface{}) ([]byte, bool, int, error) {
	data, ok := value.([]byte)
	if !ok {
		return nil, false, 0, orb.ErrUnsupportedDataType
	}

	if len(data) == 0 {
		return nil, false, 0, nil
	}

	data, littleEndian, err := validateSet(data, wkbLineString)
	if err != nil {
		return nil, false, 0, err
	}

	length := ReadUint32(data[5:9], littleEndian)
	return data[9:], littleEndian, int(length), nil
}

func validateSet(data []byte, t uint32) ([]byte, bool, error) {
	var (
		littleEndian bool
		typeCode     uint32
		err          error
	)

	if data[0] != 0 && data[0] != 1 {
		data = data[4:]
	}

	// To try and detect if this is direct from mysql
	// we try to see if cropping the first 4 values helps
	// make a valid WKB.
	for i := 0; i < 2; i++ {
		littleEndian, typeCode, err = ReadHeader(data)
		if err != nil {
			return nil, false, err
		}

		if typeCode == t {
			break
		}

		data = data[4:]
	}

	if typeCode != t {
		return nil, false, orb.ErrIncorrectGeometry
	}

	return data, littleEndian, nil
}

// ReadHeader reads the beginning of the data and returns the header.
func ReadHeader(data []byte) (bool, uint32, error) {
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

// ReadUint32 reads the data and returns a uint32.
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

// ReadFloat64 reads the data and returns a float64.
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
