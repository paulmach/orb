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
	ErrUnsupportedDataType = errors.New("wkbcommon: scan value must be []byte")

	// ErrNotWKB is returned when unmarshalling WKB and the data is not valid.
	ErrNotWKB = errors.New("wkbcommon: invalid data")

	// ErrIncorrectGeometry is returned when unmarshalling WKB data into the wrong type.
	// For example, unmarshaling linestring data into a point.
	ErrIncorrectGeometry = errors.New("wkbcommon: incorrect geometry")

	// ErrUnsupportedGeometry is returned when geometry type is not supported by this lib.
	ErrUnsupportedGeometry = errors.New("wkbcommon: unsupported geometry")
)

// ScanPoint takes binary wkb and decodes it into a point.
func ScanPoint(data []byte) (orb.Point, int, error) {
	order, typ, srid, geomData, err := unmarshalByteOrderType(data)
	if err != nil {
		return orb.Point{}, 0, err
	}

	// Checking for MySQL's SRID+WKB format where the SRID is 0.
	// Common SRIDs would be handled in the unmarshalByteOrderType above.
	if len(data) == 25 &&
		data[0] == 0 && data[1] == 0 && data[2] == 0 && data[3] == 0 {

		data = data[4:]
		order, typ, srid, geomData, err = unmarshalByteOrderType(data)
		if err != nil {
			return orb.Point{}, 0, err
		}
	}

	switch typ {
	case pointType:
		p, err := unmarshalPoint(order, geomData)
		if err != nil {
			return orb.Point{}, 0, err
		}

		return p, srid, nil
	case multiPointType:
		mp, err := unmarshalMultiPoint(order, geomData)
		if err != nil {
			return orb.Point{}, 0, err
		}
		if len(mp) == 1 {
			return mp[0], 0, nil
		}
	}

	return orb.Point{}, 0, ErrIncorrectGeometry
}

// ScanMultiPoint takes binary wkb and decodes it into a multi-point.
func ScanMultiPoint(data []byte) (orb.MultiPoint, int, error) {
	m, srid, err := Unmarshal(data)
	if err != nil {
		return nil, 0, err
	}

	switch p := m.(type) {
	case orb.Point:
		return orb.MultiPoint{p}, srid, nil
	case orb.MultiPoint:
		return p, srid, nil
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanLineString takes binary wkb and decodes it into a line string.
func ScanLineString(data []byte) (orb.LineString, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case lineStringType:
		ls, err := unmarshalLineString(order, data)
		if err != nil {
			return nil, 0, err
		}

		return ls, srid, nil
	case multiLineStringType:
		mls, err := unmarshalMultiLineString(order, data)
		if err != nil {
			return nil, 0, err
		}
		if len(mls) == 1 {
			return mls[0], srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanMultiLineString takes binary wkb and decodes it into a multi-line string.
func ScanMultiLineString(data []byte) (orb.MultiLineString, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case lineStringType:
		ls, err := unmarshalLineString(order, data)
		if err != nil {
			return nil, 0, err
		}

		return orb.MultiLineString{ls}, srid, nil
	case multiLineStringType:
		ls, err := unmarshalMultiLineString(order, data)
		if err != nil {
			return nil, 0, err
		}

		return ls, srid, nil
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanPolygon takes binary wkb and decodes it into a polygon.
func ScanPolygon(data []byte) (orb.Polygon, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case polygonType:
		p, err := unmarshalPolygon(order, data)
		if err != nil {
			return nil, 0, err
		}

		return p, srid, nil
	case multiPolygonType:
		mp, err := unmarshalMultiPolygon(order, data)
		if err != nil {
			return nil, 0, err
		}
		if len(mp) == 1 {
			return mp[0], srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanMultiPolygon takes binary wkb and decodes it into a multi-polygon.
func ScanMultiPolygon(data []byte) (orb.MultiPolygon, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case polygonType:
		p, err := unmarshalPolygon(order, data)
		if err != nil {
			return nil, 0, err
		}
		return orb.MultiPolygon{p}, srid, nil
	case multiPolygonType:
		mp, err := unmarshalMultiPolygon(order, data)
		if err != nil {
			return nil, 0, err
		}

		return mp, srid, nil
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanCollection takes binary wkb and decodes it into a collection.
func ScanCollection(data []byte) (orb.Collection, int, error) {
	m, srid, err := NewDecoder(bytes.NewReader(data)).Decode()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, 0, ErrNotWKB
	}

	if err != nil {
		return nil, 0, err
	}

	switch p := m.(type) {
	case orb.Collection:
		return p, srid, nil
	}

	return nil, 0, ErrIncorrectGeometry
}
