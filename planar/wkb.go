package planar

import "github.com/paulmach/orb/internal/wkb"

// Scan implements the sql.Scanner interface allowing
// point structs to be passed into rows.Scan(...interface{})
// The column must be of type Point and must be fetched in WKB format.
// Will attempt to parse MySQL's SRID+WKB format if the data is of the right size.
// If the column is empty (not null) an empty point (0, 0) will be returned.
func (p *Point) Scan(value interface{}) error {
	data, littleEndian, err := wkb.ValidatePoint(value)
	if err != nil || data == nil {
		return err
	}

	*p, err = unWKBPoint(data, littleEndian)
	return err
}

func readWKBPoint(data []byte, littleEndian bool) Point {
	return Point{
		wkb.ReadFloat64(data[:8], littleEndian),
		wkb.ReadFloat64(data[8:], littleEndian),
	}
}

func unWKBPoint(data []byte, littleEndian bool) (Point, error) {
	return readWKBPoint(data, littleEndian), nil
}

// Scan implements the sql.Scanner interface allowing
// line structs to be passed into rows.Scan(...interface{})
// The column must be of type LineString and contain 2 points,
// or an error will be returned. Data must be fetched in WKB format.
// Will attempt to parse MySQL's SRID+WKB format if the data is of the right size.
// If the column is empty (not null) an empty line [(0, 0), (0, 0)] will be returned.
func (l *Line) Scan(value interface{}) error {

	data, littleEndian, length, err := wkb.ValidateLine(value)
	if err != nil || data == nil {
		return err
	}

	*l, err = unWKBLine(data, littleEndian, length)
	return err
}

func unWKBLine(data []byte, littleEndian bool, length uint32) (Line, error) {
	return Line{
		a: readWKBPoint(data[:16], littleEndian),
		b: readWKBPoint(data[16:], littleEndian),
	}, nil
}

// Scan implements the sql.Scanner interface allowing
// line structs to be passed into rows.Scan(...interface{})
// The column must be of type LineString, Polygon or MultiPoint
// or an error will be returned. Data must be fetched in WKB format.
// Will attempt to parse MySQL's SRID+WKB format if obviously no WKB
// or parsing as WKB fails.
// If the column is empty (not null) an empty point set will be returned.
func (ps *PointSet) Scan(value interface{}) error {

	data, littleEndian, length, err := wkb.ValidatePointSet(value)
	if err != nil || data == nil {
		return err
	}

	*ps, err = unWKBPointSet(data, littleEndian, length)
	return err
}

func unWKBPointSet(data []byte, littleEndian bool, length int) (PointSet, error) {

	points := make([]Point, length, length)
	for i := 0; i < length; i++ {
		points[i] = readWKBPoint(data[i*16:], littleEndian)
	}

	return PointSet(points), nil
}

// Scan implements the sql.Scanner interface allowing
// line structs to be passed into rows.Scan(...interface{})
// The column must be of type LineString, Polygon or MultiPoint
// or an error will be returned. Data must be fetched in WKB format.
// Will attempt to parse MySQL's SRID+WKB format if obviously no WKB
// or parsing as WKB fails.
// If the column is empty (not null) an empty path will be returned.
func (p *Path) Scan(value interface{}) error {
	ps := PointSet{}
	err := ps.Scan(value)
	if err != nil {
		return err
	}

	*p = Path(ps)
	return nil
}
