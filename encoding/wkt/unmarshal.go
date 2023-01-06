package wkt

import (
	"errors"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
)

var (
	// ErrNotWKT is returned when unmarshalling WKT and the data is not valid.
	ErrNotWKT = errors.New("wkt: invalid data")

	// ErrIncorrectGeometry is returned when unmarshalling WKT data into the wrong type.
	// For example, unmarshaling linestring data into a point.
	ErrIncorrectGeometry = errors.New("wkt: incorrect geometry")

	// ErrUnsupportedGeometry is returned when geometry type is not supported by this lib.
	ErrUnsupportedGeometry = errors.New("wkt: unsupported geometry")
)

// UnmarshalPoint returns the point represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a point.
func UnmarshalPoint(s string) (p orb.Point, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return orb.Point{}, err
	}
	g, ok := geom.(orb.Point)
	if !ok {
		return orb.Point{}, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalMultiPoint returns the multi-point represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a multi-point.
func UnmarshalMultiPoint(s string) (p orb.MultiPoint, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}

	g, ok := geom.(orb.MultiPoint)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalLineString returns the linestring represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a linestring.
func UnmarshalLineString(s string) (p orb.LineString, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}
	g, ok := geom.(orb.LineString)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalMultiLineString returns the multi-linestring represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a multi-linestring.
func UnmarshalMultiLineString(s string) (p orb.MultiLineString, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}
	g, ok := geom.(orb.MultiLineString)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalPolygon returns the polygon represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a polygon.
func UnmarshalPolygon(s string) (p orb.Polygon, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}
	g, ok := geom.(orb.Polygon)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalMultiPolygon returns the multi-polygon represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a multi-polygon.
func UnmarshalMultiPolygon(s string) (p orb.MultiPolygon, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}
	g, ok := geom.(orb.MultiPolygon)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// UnmarshalCollection returns the geometry collection represented by the wkt string.
// Will return ErrIncorrectGeometry if the wkt is not a geometry collection.
func UnmarshalCollection(s string) (p orb.Collection, err error) {
	geom, err := Unmarshal(s)
	if err != nil {
		return orb.Collection{}, err
	}
	g, ok := geom.(orb.Collection)
	if !ok {
		return nil, ErrIncorrectGeometry
	}
	return g, nil
}

// trimSpaceBrackets trim space and brackets
func trimSpaceBrackets(s string) string {
	s = strings.Trim(s, " ")
	if len(s) == 0 {
		return ""
	}
	if s[0] == '(' {
		s = s[1:]
	}
	if s[len(s)-1] == ')' {
		s = s[:len(s)-1]
	}

	return strings.Trim(s, " ")
}

// parsePoint pase point by (x y)
func parsePoint(s string) (p orb.Point, err error) {
	ps := strings.Split(s, " ")
	if len(ps) != 2 {
		return orb.Point{}, ErrNotWKT
	}

	x, err := strconv.ParseFloat(ps[0], 64)
	if err != nil {
		return orb.Point{}, err
	}

	y, err := strconv.ParseFloat(ps[1], 64)
	if err != nil {
		return orb.Point{}, err
	}

	return orb.Point{x, y}, nil
}

// splitGeometryCollection split GEOMETRYCOLLECTION to more geometry
func splitGeometryCollection(s string) (r []string) {
	r = make([]string, 0)
	stack := make([]rune, 0)
	l := len(s)
	for i, v := range s {
		if !strings.Contains(string(stack), "(") {
			stack = append(stack, v)
			continue
		}
		if v >= 'A' && v < 'Z' {
			t := string(stack)
			r = append(r, t[:len(t)-1])
			stack = make([]rune, 0)
			stack = append(stack, v)
			continue
		}
		if i == l-1 {
			r = append(r, string(stack))
			continue
		}
		stack = append(stack, v)
	}
	return
}

// Unmarshal return a geometry by parsing the WKT string.
func Unmarshal(s string) (geom orb.Geometry, err error) {
	s = strings.ToUpper(strings.Trim(s, " "))
	switch {
	case strings.Contains(s, "GEOMETRYCOLLECTION"):
		if s == "GEOMETRYCOLLECTION EMPTY" {
			return orb.Collection{}, nil
		}
		s = strings.Replace(s, "GEOMETRYCOLLECTION", "", -1)
		if len(s) == 0 {
			return nil, ErrNotWKT
		}
		c := orb.Collection{}
		ms := splitGeometryCollection(s)
		if len(ms) == 0 {
			return nil, err
		}
		for _, v := range ms {
			if len(v) == 0 {
				continue
			}
			g, err := Unmarshal(v)
			if err != nil {
				return nil, err
			}
			c = append(c, g)
		}
		geom = c

	case strings.Contains(s, "MULTIPOINT"):
		if s == "MULTIPOINT EMPTY" {
			return orb.MultiPoint{}, nil
		}
		s = strings.Replace(s, "MULTIPOINT", "", -1)
		s = trimSpaceBrackets(s)
		ps := strings.Split(s, ",")
		mp := orb.MultiPoint{}
		for _, p := range ps {
			tp, err := parsePoint(trimSpaceBrackets(p))
			if err != nil {
				return nil, err
			}
			mp = append(mp, tp)
		}
		geom = mp

	case strings.Contains(s, "POINT"):
		s = strings.Replace(s, "POINT", "", -1)
		tp, err := parsePoint(trimSpaceBrackets(s))
		if err != nil {
			return nil, err
		}
		geom = tp

	case strings.Contains(s, "MULTILINESTRING"):
		if s == "MULTILINESTRING EMPTY" {
			return orb.MultiLineString{}, nil
		}
		s = strings.Replace(s, "MULTILINESTRING", "", -1)
		ml := orb.MultiLineString{}
		for _, l := range strings.Split(trimSpaceBrackets(s), "),(") {
			tl := orb.LineString{}
			for _, p := range strings.Split(trimSpaceBrackets(l), ",") {
				tp, err := parsePoint(trimSpaceBrackets(p))
				if err != nil {
					return nil, err
				}
				tl = append(tl, tp)
			}
			ml = append(ml, tl)
		}
		geom = ml

	case strings.Contains(s, "LINESTRING"):
		if s == "LINESTRING EMPTY" {
			return orb.LineString{}, nil
		}
		s = strings.Replace(s, "LINESTRING", "", -1)
		s = trimSpaceBrackets(s)
		ps := strings.Split(s, ",")
		ls := orb.LineString{}
		for _, p := range ps {
			tp, err := parsePoint(trimSpaceBrackets(p))
			if err != nil {
				return nil, err
			}
			ls = append(ls, tp)
		}
		geom = ls

	case strings.Contains(s, "MULTIPOLYGON"):
		if s == "MULTIPOLYGON EMPTY" {
			return orb.MultiPolygon{}, nil
		}
		s = strings.Replace(s, "MULTIPOLYGON", "", -1)
		mpol := orb.MultiPolygon{}
		for _, ps := range strings.Split(trimSpaceBrackets(s), ")),((") {
			pol := orb.Polygon{}
			for _, ls := range strings.Split(trimSpaceBrackets(ps), "),(") {
				ring := orb.Ring{}
				for _, p := range strings.Split(ls, ",") {
					tp, err := parsePoint(trimSpaceBrackets(p))
					if err != nil {
						return nil, err
					}
					ring = append(ring, tp)
				}
				pol = append(pol, ring)
			}
			mpol = append(mpol, pol)
		}
		geom = mpol

	case strings.Contains(s, "POLYGON"):
		if s == "POLYGON EMPTY" {
			return orb.Polygon{}, nil
		}
		s = strings.Replace(s, "POLYGON", "", -1)
		s = trimSpaceBrackets(s)

		rs := strings.Split(s, "),(")
		pol := make(orb.Polygon, 0, len(rs))
		for _, r := range rs {
			ps := strings.Split(trimSpaceBrackets(r), ",")
			ring := orb.Ring{}
			for _, p := range ps {
				tp, err := parsePoint(trimSpaceBrackets(p))
				if err != nil {
					return nil, err
				}
				ring = append(ring, tp)
			}
			pol = append(pol, ring)
		}
		geom = pol
	default:
		return nil, ErrUnsupportedGeometry
	}

	return
}
