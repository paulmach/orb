package wkbcommon

import (
	"bytes"
	"errors"
	"io"

	"github.com/paulmach/orb"
)

var (
	// ErrUnsupportedDataType is returned by Scan methods when asked to scan
	// non []byte data from the database. This should never happen
	// if the driver is acting appropriately.
	ErrUnsupportedDataType = errors.New("wkb: scan value must be []byte")

	// ErrNotWKB is returned when unmarshalling WKB and the data is not valid.
	ErrNotWKB = errors.New("wkb: invalid data")

	// ErrIncorrectGeometry is returned when unmarshalling WKB data into the wrong type.
	// For example, unmarshaling linestring data into a point.
	ErrIncorrectGeometry = errors.New("wkb: incorrect geometry")

	// ErrUnsupportedGeometry is returned when geometry type is not supported by this lib.
	ErrUnsupportedGeometry = errors.New("wkb: unsupported geometry")
)

// ScanPoint takes binary wkb and decodes it into a point.
func ScanPoint(data []byte) (orb.Point, error) {
	order, typ, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return orb.Point{}, err
	}

	// Checking for MySQL's SRID+WKB format where the SRID is 0.
	// Common SRIDs would be handled in the unmarshalByteOrderType above.
	if len(data) == 25 &&
		data[0] == 0 && data[1] == 0 && data[2] == 0 && data[3] == 0 {

		data = data[4:]
		order, typ, data, err = unmarshalByteOrderType(data)
		if err != nil {
			return orb.Point{}, err
		}
	}

	switch typ {
	case pointType:
		return unmarshalPoint(order, data[5:])
	case multiPointType:
		mp, err := unmarshalMultiPoint(order, data[5:])
		if err != nil {
			return orb.Point{}, err
		}
		if len(mp) == 1 {
			return mp[0], nil
		}
	}

	return orb.Point{}, ErrIncorrectGeometry
}

// ScanMultiPoint takes binary wkb and decodes it into a multi-point.
func ScanMultiPoint(data []byte) (orb.MultiPoint, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	switch p := m.(type) {
	case orb.Point:
		return orb.MultiPoint{p}, nil
	case orb.MultiPoint:
		return p, nil
	}

	return nil, ErrIncorrectGeometry
}

// ScanLineString takes binary wkb and decodes it into a line string.
func ScanLineString(data []byte) (orb.LineString, error) {
	order, typ, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, err
	}

	switch typ {
	case lineStringType:
		return unmarshalLineString(order, data[5:])
	case multiLineStringType:
		mls, err := unmarshalMultiLineString(order, data[5:])
		if err != nil {
			return nil, err
		}
		if len(mls) == 1 {
			return mls[0], nil
		}
	}

	return nil, ErrIncorrectGeometry
}

// ScanMultiLineString takes binary wkb and decodes it into a multi-line string.
func ScanMultiLineString(data []byte) (orb.MultiLineString, error) {
	order, typ, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, err
	}

	switch typ {
	case lineStringType:
		ls, err := unmarshalLineString(order, data[5:])
		if err != nil {
			return nil, err
		}

		return orb.MultiLineString{ls}, nil
	case multiLineStringType:
		return unmarshalMultiLineString(order, data[5:])
	}

	return nil, ErrIncorrectGeometry
}

// ScanPolygon takes binary wkb and decodes it into a polygon.
func ScanPolygon(data []byte) (orb.Polygon, error) {
	order, typ, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, err
	}

	switch typ {
	case polygonType:
		return unmarshalPolygon(order, data[5:])
	case multiPolygonType:
		mp, err := unmarshalMultiPolygon(order, data[5:])
		if err != nil {
			return nil, err
		}
		if len(mp) == 1 {
			return mp[0], nil
		}
	}

	return nil, ErrIncorrectGeometry
}

// ScanMultiPolygon takes binary wkb and decodes it into a multi-polygon.
func ScanMultiPolygon(data []byte) (orb.MultiPolygon, error) {
	order, typ, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, err
	}

	switch typ {
	case polygonType:
		p, err := unmarshalPolygon(order, data[5:])
		if err != nil {
			return nil, err
		}
		return orb.MultiPolygon{p}, nil
	case multiPolygonType:
		return unmarshalMultiPolygon(order, data[5:])
	}

	return nil, ErrIncorrectGeometry
}

// ScanCollection takes binary wkb and decodes it into a collection.
func ScanCollection(data []byte) (orb.Collection, error) {
	m, err := NewDecoder(bytes.NewReader(data)).Decode()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, ErrNotWKB
	}

	if err != nil {
		return nil, err
	}

	switch p := m.(type) {
	case orb.Collection:
		return p, nil
	}

	return nil, ErrIncorrectGeometry
}
