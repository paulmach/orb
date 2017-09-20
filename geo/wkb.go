package geo

import "github.com/paulmach/orb/internal/wkb"

// Scan implements the sql.Scanner interface allowing
// Point structs to be passed into rows.Scan(...interface{})
// The column must be of type Point and must be fetched in WKB format.
// Will attempt to parse MySQL's SRID+WKB format if the data is of the right size.
// If the column is empty (not null) an empty point (0, 0) will be returned.
func (p *Point) Scan(value interface{}) error {
	x, y, isNull, err := wkb.ValidatePoint(value)
	if err != nil || isNull {
		return err
	}

	*p = Point{x, y}
	return nil
}

func readWKBPoint(data []byte, littleEndian bool) Point {
	return Point{
		wkb.ReadFloat64(data[:8], littleEndian),
		wkb.ReadFloat64(data[8:], littleEndian),
	}
}

// Scan implements the sql.Scanner interface allowing a
// bound to be read in as the bound of a two point line string.
func (b *Bound) Scan(value interface{}) error {
	var ls LineString
	err := ls.Scan(value)
	if err != nil {
		return err
	}

	*b = ls.Bound()
	return nil
}

// Scan implements the sql.Scanner interface allowing
// LineString structs to be passed into rows.Scan(...interface{})
// The column must be of type LineString or an error will be returned.
// Data must be fetched in WKB format.  Will attempt to parse MySQL's
// SRID+WKB format if obviously not WKB or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (ls *LineString) Scan(value interface{}) error {
	_, err := ls.scan(value)
	return err
}

func (ls *LineString) scan(value interface{}) ([]byte, error) {
	data, littleEndian, length, err := wkb.ValidateLineString(value)
	if err != nil || data == nil {
		return nil, err
	}

	*ls, data, err = scanXYList(data, littleEndian, length)
	return data, err
}

func scanXYList(data []byte, littleEndian bool, length int) (LineString, []byte, error) {
	points := make([]Point, length)
	for i := 0; i < length; i++ {
		points[i] = readWKBPoint(data[16*i:], littleEndian)
	}

	return LineString(points), data[16*length:], nil
}

// Scan implements the sql.Scanner interface allowing
// Polygon structs to be passed into rows.Scan(...interface{})
// The column must be of type Polygon or an error will be returned.
// Data must be fetched in WKB format.  Will attempt to parse MySQL's
// SRID+WKB format if obviously not WKB or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (p *Polygon) Scan(value interface{}) error {
	_, err := p.scan(value)
	return err
}

func (p *Polygon) scan(value interface{}) ([]byte, error) {
	data, littleEndian, rings, err := wkb.ValidatePolygon(value)
	if err != nil || data == nil {
		return nil, err
	}

	poly := make(Polygon, rings)
	for i := 0; i < rings; i++ {
		var ls LineString

		length := wkb.ReadUint32(data, littleEndian)
		ls, data, err = scanXYList(data[4:], littleEndian, int(length))
		if err != nil {
			return nil, err
		}

		poly[i] = Ring(ls)
	}

	*p = poly
	return data, nil
}

// Scan implements the sql.Scanner interface allowing
// MultiPoint structs to be passed into rows.Scan(...interface{})
// The column must be of type MultiPoint or an error will be returned.
// Data must be fetched in WKB format.  Will attempt to parse MySQL's
// SRID+WKB format if obviously not WKB or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (mp *MultiPoint) Scan(value interface{}) error {
	data, littleEndian, length, err := wkb.ValidateMultiPoint(value)
	if err != nil || data == nil {
		return err
	}

	*mp, err = unWKBMultiPoint(data, littleEndian, length)
	return err
}

func unWKBMultiPoint(data []byte, littleEndian bool, length int) (MultiPoint, error) {
	points := make([]Point, length)
	for i := 0; i < length; i++ {
		x, y, err := wkb.ReadPoint(data[21*i:])
		if err != nil {
			return nil, err
		}

		points[i] = Point{x, y}
	}

	return MultiPoint(points), nil
}

// Scan implements the sql.Scanner interface allowing
// MultiLineString to be passed into rows.Scan(...interface{})
// The column must be of type MultiLineString or an error will be returned.
// Data must be fetched in WKB format.  Will attempt to parse MySQL's
// SRID+WKB format if obviously not WKB or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (mls *MultiLineString) Scan(value interface{}) error {
	data, _, length, err := wkb.ValidateMultiLineString(value)
	if err != nil || data == nil {
		return err
	}

	multiline := make(MultiLineString, length)
	for i := 0; i < length; i++ {
		data, err = multiline[i].scan(data)
		if err != nil {
			return err
		}
	}

	*mls = multiline
	return nil
}

// Scan implements the sql.Scanner interface allowing
// MultiPolygon to be passed into rows.Scan(...interface{})
// The column must be of type MultiPolygon or an error will be returned.
// Data must be fetched in WKB format.  Will attempt to parse MySQL's
// SRID+WKB format if obviously not WKB or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (mp *MultiPolygon) Scan(value interface{}) error {
	data, _, length, err := wkb.ValidateMultiPolygon(value)
	if err != nil || data == nil {
		return err
	}

	multipoly := make(MultiPolygon, length)
	for i := 0; i < length; i++ {
		data, err = multipoly[i].scan(data)
		if err != nil {
			return err
		}
	}

	*mp = multipoly
	return nil
}
