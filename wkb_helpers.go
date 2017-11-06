package orb

import (
	"math"
)

// source of information
// https://www.ibm.com/support/knowledgecenter/en/SS6NHC/com.ibm.db2.luw.spatial.topics.doc/doc/rsbp4121.html

const (
	wkbPoint           = 1
	wkbLineString      = 2
	wkbPolygon         = 3
	wkbMultiPoint      = 4
	wkbMultiLineString = 5
	wkbMultiPolygon    = 6
)

// validatePoint checks the wkb input and returns x, y, isnull, err.
func validatePoint(value interface{}) (float64, float64, bool, error) {
	data, ok := value.([]byte)
	if !ok {
		return 0, 0, false, ErrUnsupportedDataType
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
		return 0, 0, false, ErrIncorrectGeometry
	}

	x, y, err := readPoint(data)
	return x, y, false, err
}

// readPoint reads the beginning of the data to find the wkb point.
func readPoint(data []byte) (float64, float64, error) {
	if len(data) < 21 {
		return 0, 0, ErrIncorrectGeometry
	}

	littleEndian, typeCode, err := readHeader(data)
	if err != nil {
		return 0, 0, err
	}

	if typeCode != wkbPoint {
		return 0, 0, ErrIncorrectGeometry
	}

	return readFloat64(data[5:13], littleEndian),
		readFloat64(data[13:], littleEndian),
		nil
}

// validateLineString checks the wkb for a linestring geometry.
func validateLineString(value interface{}) ([]byte, bool, int, error) {
	data, littleEndian, err := validateSet(value, wkbLineString)
	if err != nil || data == nil {
		return nil, false, 0, err
	}

	length := readUint32(data[5:9], littleEndian)
	return data[9:], littleEndian, int(length), nil
}

// validatePolygon checks the wkb for a polygon geometry.
func validatePolygon(value interface{}) ([]byte, bool, int, error) {
	data, littleEndian, err := validateSet(value, wkbPolygon)
	if err != nil || data == nil {
		return nil, false, 0, err
	}

	length := readUint32(data[5:9], littleEndian)
	return data[9:], littleEndian, int(length), nil
}

// validateMultiPoint checks wkb for a multipoint geometry.
func validateMultiPoint(value interface{}) ([]byte, bool, int, error) {
	data, littleEndian, err := validateSet(value, wkbMultiPoint)
	if err != nil || data == nil {
		return nil, false, 0, err
	}

	length := int(readUint32(data[5:9], littleEndian))
	if len(data) != 9+21*length {
		return nil, false, 0, ErrNotWKB
	}

	return data[9:], littleEndian, length, nil
}

// validateMultiLineString checks the wkb for a multilinestring geometry.
func validateMultiLineString(value interface{}) ([]byte, bool, int, error) {
	data, littleEndian, err := validateSet(value, wkbMultiLineString)
	if err != nil || data == nil {
		return nil, false, 0, err
	}

	length := readUint32(data[5:9], littleEndian)
	return data[9:], littleEndian, int(length), nil
}

// validateMultiPolygon checks the wkb for a multipolygon geometry.
func validateMultiPolygon(value interface{}) ([]byte, bool, int, error) {
	data, littleEndian, err := validateSet(value, wkbMultiPolygon)
	if err != nil || data == nil {
		return nil, false, 0, err
	}

	length := readUint32(data[5:9], littleEndian)
	return data[9:], littleEndian, int(length), nil
}

func validateSet(value interface{}, t uint32) ([]byte, bool, error) {
	data, ok := value.([]byte)
	if !ok {
		return nil, false, ErrUnsupportedDataType
	}

	if len(data) == 0 {
		return nil, false, nil
	}

	if len(data) < 6 {
		return nil, false, ErrNotWKB
	}

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
		littleEndian, typeCode, err = readHeader(data)
		if err != nil {
			return nil, false, err
		}

		if typeCode == t {
			break
		}

		data = data[4:]
	}

	if typeCode != t {
		return nil, false, ErrIncorrectGeometry
	}

	return data, littleEndian, nil
}

// readHeader reads the beginning of the data and returns the header.
func readHeader(data []byte) (bool, uint32, error) {
	if len(data) < 6 {
		return false, 0, ErrNotWKB
	}

	if data[0] == 0 {
		return false, readUint32(data[1:5], false), nil
	}

	if data[0] == 1 {
		return true, readUint32(data[1:5], true), nil
	}

	return false, 0, ErrNotWKB
}

// readUint32 reads the data and returns a uint32.
func readUint32(data []byte, littleEndian bool) uint32 {
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

// readFloat64 reads the data and returns a float64.
func readFloat64(data []byte, littleEndian bool) float64 {
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
