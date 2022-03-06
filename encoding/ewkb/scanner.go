package ewkb

import (
	"database/sql"
	"database/sql/driver"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/internal/wkbcommon"
)

var (
	_ sql.Scanner  = &GeometryScanner{}
	_ driver.Value = value{}
)

// GeometryScanner is a thing that can scan in sql query results.
// It can be used as a scan destination:
//
//	var s wkb.GeometryScanner
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use s.Geometry
//	  // use s.SRID
//	} else {
//	  // NULL value
//	}
type GeometryScanner struct {
	g        interface{}
	SRID     int
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
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(s)
//	...
//	if s.Valid {
//	  // use p
//	} else {
//	  // NULL value
//	}
func Scanner(g interface{}) *GeometryScanner {
	return &GeometryScanner{g: g}
}

// Scan will scan the input []byte data into a geometry.
// This could be into the orb geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	s.Geometry = nil
	s.Valid = false

	g, srid, valid, err := wkbcommon.Scan(s.g, d)
	if err != nil {
		return mapCommonError(err)
	}

	s.Geometry = g
	s.SRID = srid
	s.Valid = valid

	return nil
}

type value struct {
	srid int
	v    orb.Geometry
}

// Value will create a driver.Valuer that will EWKB the geometry into the database query.
//
//	db.Exec("INSERT INTO table (point_column) VALUES (?)", ewkb.Value(p, 4326))
func Value(g orb.Geometry, srid int) driver.Valuer {
	return value{srid: srid, v: g}

}

func (v value) Value() (driver.Value, error) {
	val, err := Marshal(v.v, v.srid)
	if val == nil {
		return nil, err
	}
	return val, err
}
