package wkb

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/paulmach/orb"
)

var (
	_ sql.Scanner  = &GeometryScanner{}
	_ driver.Value = value{}
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

// GeometryScanner is a thing that can scan in sql query results.
// It can be used as a scan destination:
//
//	var s wkb.GeometryScanner
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use s.Geometry
//	} else {
//	  // NULL value
//	}
type GeometryScanner struct {
	g        interface{}
	Geometry orb.Geometry
	Valid    bool // Valid is true if the geometry is not NULL
}

// Scanner will return a GeometryScanner that can scan sql query results.
// The geometryScanner.Geometry attribute will be set to the value.
// If g is non-nil, it MUST be a pointer to an orb.Geometry
// type like a Point or LineString. In that case the value will be written to
// g and the Geometry attribute.
//
//	var p orb.Point
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(wkb.Scanner(&p))
//	...
//	// use p
//
// If the value may be null check Valid first:
//
//	var point orb.Point
//	s := wkb.Scanner(&point)
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use p
//	} else {
//	  // NULL value
//	}
//
// Scanning directly from MySQL columns is supported. By default MySQL returns geometry
// data as WKB but prefixed with a 4 byte SRID. To support this, if the data is not
// valid WKB, the code will strip the first 4 bytes and try again.
// This works for most use cases.
func Scanner(g interface{}) *GeometryScanner {
	return &GeometryScanner{g: g}
}

// MySQLGeometryScanner should be used when scanning from a mysql/mariadb database.
// It will take the first 4 bytes and little endian decode them into the SRID.
type MySQLGeometryScanner struct {
	GeometryScanner
	SRID int
}

// MySQLScanner will handle the unique format of a mysql/mariadb geometry data.
// The wkb data is prefixed with a 4 byte SRID.
func MySQLScanner(g interface{}) *MySQLGeometryScanner {
	return &MySQLGeometryScanner{GeometryScanner: GeometryScanner{g: g}}
}

// Scan will scan the input []byte data into a geometry.
// Since it's a MySQL Scanner it will read the first 4 bytes as the SRID
// before decoding the rest of the data as wkb.
func (s *MySQLGeometryScanner) Scan(d interface{}) error {
	if d == nil {
		return nil
	}

	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	if data == nil {
		return nil
	}

	if len(data) < 4 {
		return ErrNotWKB
	}

	s.SRID = int(binary.LittleEndian.Uint32(data))
	return s.GeometryScanner.Scan(data[4:])
}

// Scan will scan the input []byte data into a geometry.
// This could be into the orb geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	s.Geometry = nil
	s.Valid = false

	if d == nil {
		return nil
	}

	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	if data == nil {
		return nil
	}

	// go-pg will return ST_AsBinary(*) data as `\xhexencoded` which
	// needs to be converted to true binary for further decoding.
	// Code detects the \x prefix and then converts the rest from Hex to binary.
	if len(data) > 2 && data[0] == byte('\\') && data[1] == byte('x') {
		n, err := hex.Decode(data, data[2:])
		if err != nil {
			return fmt.Errorf("thought the data was hex, but it is not: %v", err)
		}
		data = data[:n]
	}

	switch g := s.g.(type) {
	case nil:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Point:
		p, err := scanPoint(data)
		if err != nil {
			return err
		}

		*g = p
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.MultiPoint:
		p, err := scanMultiPoint(data)
		if err != nil {
			return err
		}

		*g = p
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.LineString:
		p, err := scanLineString(data)
		if err != nil {
			return err
		}

		*g = p
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.MultiLineString:
		p, err := scanMultiLineString(data)
		if err != nil {
			return err
		}

		*g = p
		s.Geometry = p
		s.Valid = true
		return nil
	case *orb.Ring:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if p, ok := m.(orb.Polygon); ok && len(p) == 1 {
			*g = p[0]
			s.Geometry = p[0]
			s.Valid = true
			return nil
		}

		return ErrIncorrectGeometry
	case *orb.Polygon:
		m, err := scanPolygon(data)
		if err != nil {
			return err
		}

		*g = m
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.MultiPolygon:
		m, err := scanMultiPolygon(data)
		if err != nil {
			return err
		}

		*g = m
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Collection:
		m, err := scanCollection(data)
		if err != nil {
			return err
		}

		*g = m
		s.Geometry = m
		s.Valid = true
		return nil
	case *orb.Bound:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		b := m.Bound()
		*g = b
		s.Geometry = b
		s.Valid = true
		return nil
	}

	return ErrIncorrectGeometry
}

func scanPoint(data []byte) (orb.Point, error) {
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

func scanMultiPoint(data []byte) (orb.MultiPoint, error) {
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

func scanLineString(data []byte) (orb.LineString, error) {
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

func scanMultiLineString(data []byte) (orb.MultiLineString, error) {
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

func scanPolygon(data []byte) (orb.Polygon, error) {
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

func scanMultiPolygon(data []byte) (orb.MultiPolygon, error) {
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

func scanCollection(data []byte) (orb.Collection, error) {
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

type value struct {
	v orb.Geometry
}

// Value will create a driver.Valuer that will WKB the geometry
// into the database query.
func Value(g orb.Geometry) driver.Valuer {
	return value{v: g}
}

func (v value) Value() (driver.Value, error) {
	val, err := Marshal(v.v)
	if val == nil {
		return nil, err
	}
	return val, err
}

type mysqlValue struct {
	v    orb.Geometry
	srid uint32
}

// MySQLValue should be used for a value to be inserted into a mysql/mariadb database.
// It will prefix the geometry with 4 bytes representing the SRID. Zero if not provided.
func MySQLValue(g orb.Geometry, srid ...int) driver.Valuer {
	v := mysqlValue{v: g}
	if len(srid) != 0 {
		v.srid = uint32(srid[0])
	}

	return v
}

func (v mysqlValue) Value() (driver.Value, error) {
	data := make([]byte, 4, geomLength(v.v)+4)
	DefaultByteOrder.PutUint32(data, v.srid)

	buf := bytes.NewBuffer(data)

	e := NewEncoder(buf)
	err := e.Encode(v.v)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, nil
	}

	return buf.Bytes(), nil
}
