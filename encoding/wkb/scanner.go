package wkb

import (
	"database/sql"
	"database/sql/driver"
	"errors"

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
type GeometryScanner struct {
	g        interface{}
	Geometry orb.Geometry
}

// Scanner will return a GeometryScanner that can scan sql query results.
// The geometryScanner.Geometry attribute will be set to the value.
// If g is non-nil, it MUST be a pointer to an orb.Geometry
// type like a Point or LineString. In that case the value will be written to g.
func Scanner(g interface{}) *GeometryScanner {
	return &GeometryScanner{g: g}
}

// Scan will scan the input []byte data into a geometry.
// This could be into the orb geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	switch g := s.g.(type) {
	case nil:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		s.Geometry = m
		return nil
	case *orb.Point:
		p, err := scanPoint(data)
		if err != nil {
			return err
		}

		*g = p
		return nil
	case *orb.MultiPoint:
		p, err := scanMultiPoint(data)
		if err != nil {
			return err
		}

		*g = p
		return nil
	case *orb.LineString:
		p, err := scanLineString(data)
		if err != nil {
			return err
		}

		*g = p
		return nil
	case *orb.MultiLineString:
		p, err := scanMultiLineString(data)
		if err != nil {
			return err
		}

		*g = p
		return nil
	case *orb.Ring:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if p, ok := m.(orb.Polygon); ok && len(p) == 1 {
			*g = p[0]
			return nil
		}

		return ErrIncorrectGeometry
	case *orb.Polygon:
		m, err := scanPolygon(data)
		if err != nil {
			return err
		}

		*g = m
		return nil
	case *orb.MultiPolygon:
		m, err := scanMultiPolygon(data)
		if err != nil {
			return err
		}

		*g = m
		return nil
	case *orb.Collection:
		m, err := scanCollection(data)
		if err != nil {
			return err
		}

		*g = m
		return nil
	case *orb.Bound:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		*g = m.Bound()
		return nil
	}

	return ErrIncorrectGeometry
}

func scanPoint(data []byte) (orb.Point, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return orb.Point{}, err
	}

	switch p := m.(type) {
	case orb.Point:
		return p, nil
	case orb.MultiPoint:
		if len(p) == 1 {
			return p[0], nil
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
	m, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	switch p := m.(type) {
	case orb.LineString:
		return p, nil
	case orb.MultiLineString:
		if len(p) == 1 {
			return p[0], nil
		}
	}

	return nil, ErrIncorrectGeometry
}

func scanMultiLineString(data []byte) (orb.MultiLineString, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	switch ls := m.(type) {
	case orb.LineString:
		return orb.MultiLineString{ls}, nil
	case orb.MultiLineString:
		return ls, nil
	}

	return nil, ErrIncorrectGeometry
}

func scanPolygon(data []byte) (orb.Polygon, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	switch p := m.(type) {
	case orb.Polygon:
		return p, nil
	case orb.MultiPolygon:
		if len(p) == 1 {
			return p[0], nil
		}
	}

	return nil, ErrIncorrectGeometry
}

func scanMultiPolygon(data []byte) (orb.MultiPolygon, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	switch p := m.(type) {
	case orb.Polygon:
		return orb.MultiPolygon{p}, nil
	case orb.MultiPolygon:
		return p, nil
	}

	return nil, ErrIncorrectGeometry
}

func scanCollection(data []byte) (orb.Collection, error) {
	m, err := Unmarshal(data)
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
	return Marshal(v.v)
}
