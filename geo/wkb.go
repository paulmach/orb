package geo

import "github.com/paulmach/orb/internal/wkb"

// Scan implements the sql.Scanner interface allowing
// point structs to be passed into rows.Scan(...interface{})
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
		x, y, err := wkb.ReadPoint(data[21*i:])
		if err != nil {
			return nil, err
		}

		points[i] = Point{x, y}
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
	data, littleEndian, length, err := wkb.ValidatePath(value)
	if err != nil || data == nil {
		return err
	}

	*p, err = unWKBPath(data, littleEndian, length)
	return err
}

func unWKBPath(data []byte, littleEndian bool, length int) (Path, error) {
	points := make([]Point, length, length)
	for i := 0; i < length; i++ {
		points[i] = readWKBPoint(data[16*i:], littleEndian)
	}

	return Path(points), nil
}
