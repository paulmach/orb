package wkb

import (
	"database/sql"
	"database/sql/driver"
	"errors"

	"github.com/paulmach/orb"
)

var (
	_ sql.Scanner  = scanner{}
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
)

type scanner struct {
	g interface{}
}

// Scanner accepts a pointer to an orb geoemtry like a Point or LineString
// and will scan it from a sql.Query.
func Scanner(g interface{}) sql.Scanner {
	return scanner{g: g}
}

func (s scanner) Scan(d interface{}) error {
	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	switch g := s.g.(type) {
	case *orb.Point:
		if len(data) == 21 {
			// the length of a point type in WKB
		} else if len(data) == 25 {
			// Most likely MySQL's SRID+WKB format.
			// However, could be a line string or multipoint with only one point.
			// But those would be invalid for parsing a point.
			data = data[4:]
		} else if len(data) == 0 {
			// empty data, return empty go struct which in this case
			// would be [0,0]
			*g = orb.Point{}
			return nil
		}

		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if p, ok := m.(orb.Point); ok {
			*g = p
			return nil
		}

		return ErrIncorrectGeometry

	case *orb.MultiPoint:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		switch p := m.(type) {
		case orb.Point:
			*g = orb.MultiPoint{p}
		case orb.MultiPoint:
			*g = p
		default:
			return ErrIncorrectGeometry
		}

		return nil
	case *orb.LineString:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if ls, ok := m.(orb.LineString); ok {
			*g = ls
			return nil
		}

		return ErrIncorrectGeometry

	case *orb.MultiLineString:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		switch ls := m.(type) {
		case orb.LineString:
			*g = orb.MultiLineString{ls}
		case orb.MultiLineString:
			*g = ls
		default:
			return ErrIncorrectGeometry
		}

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
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if p, ok := m.(orb.Polygon); ok {
			*g = p
			return nil
		}

		return ErrIncorrectGeometry
	case *orb.MultiPolygon:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		switch p := m.(type) {
		case orb.Polygon:
			*g = orb.MultiPolygon{p}
		case orb.MultiPolygon:
			*g = p
		default:
			return ErrIncorrectGeometry
		}

		return nil
	case *orb.Collection:
		m, err := Unmarshal(data)
		if err != nil {
			return err
		}

		if c, ok := m.(orb.Collection); ok {
			*g = c
			return nil
		}

		return ErrIncorrectGeometry
	}

	return ErrIncorrectGeometry
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
