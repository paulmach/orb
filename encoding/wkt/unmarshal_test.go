package wkt

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestTrimSpaceBrackets(t *testing.T) {
	cases := []struct {
		s        string
		expected string
	}{
		{
			s:        "(1 2)",
			expected: "1 2",
		},
		{
			s:        "((1 2),(0.5 1.5))",
			expected: "(1 2),(0.5 1.5)",
		},
		{
			s:        "(1 2,0.5 1.5)",
			expected: "1 2,0.5 1.5",
		},
		{
			s:        "((1 2,3 4),(5 6,7 8))",
			expected: "(1 2,3 4),(5 6,7 8)",
		},
		{
			s:        "(((1 2,3 4)),((5 6,7 8),(1 2,5 4)))",
			expected: "((1 2,3 4)),((5 6,7 8),(1 2,5 4))",
		},
	}

	for _, tc := range cases {
		if trimSpaceBrackets(tc.s) != tc.expected {
			t.Log(trimSpaceBrackets(tc.s))
			t.Log(tc.expected)
			t.Errorf("trim space and brackets error")
		}
	}
}

func TestUnmarshalPoint(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.Point
	}{
		// int
		{
			s:        "POINT(1 2)",
			expected: orb.Point{1, 2},
		},
		// float64
		{
			s:        "POINT(1.34 2.35)",
			expected: orb.Point{1.34, 2.35},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalPoint(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalMultiPoint(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.MultiPoint
	}{
		// empty
		{
			s:        "MULTIPOINT EMPTY",
			expected: orb.MultiPoint{},
		},
		// int
		{
			s:        "MULTIPOINT((1 2),(0.5 1.5))",
			expected: orb.MultiPoint{{1, 2}, {0.5, 1.5}},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalMultiPoint(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalLineString(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.LineString
	}{
		{
			s:        "LINESTRING EMPTY",
			expected: orb.LineString{},
		},

		{
			s:        "LINESTRING(1 2,0.5 1.5)",
			expected: orb.LineString{{1, 2}, {0.5, 1.5}},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalLineString(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalMultiLineString(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.MultiLineString
	}{
		{
			s:        "MULTILINESTRING EMPTY",
			expected: orb.MultiLineString{},
		},

		{
			s:        "MULTILINESTRING((1 2,3 4),(5 6,7 8))",
			expected: orb.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalMultiLineString(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalPolygon(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.Polygon
	}{
		// empty
		{
			s:        "POLYGON EMPTY",
			expected: orb.Polygon{},
		},
		// ring
		// origin: orb.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}}
		// convert: orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}
		{
			s:        "POLYGON((0 0,1 0,1 1,0 0))",
			expected: orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
		},
		// bound
		// origin: orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 2}},
		// convert: orb.Polygon{{{0, 0}, {1, 0}, {1, 2}, {0, 2}, {0, 0}}}
		{
			s:        "POLYGON((0 0,1 0,1 2,0 2,0 0))",
			expected: orb.Polygon{{{0, 0}, {1, 0}, {1, 2}, {0, 2}, {0, 0}}},
		},
		// polygon
		{
			s:        "POLYGON((1 2,3 4),(5 6,7 8))",
			expected: orb.Polygon{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalPolygon(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalMutilPolygon(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.MultiPolygon
	}{
		// empty
		{
			s:        "MULTIPOLYGON EMPTY",
			expected: orb.MultiPolygon{},
		},
		// multi polygon
		{
			s:        "MULTIPOLYGON(((1 2,3 4)),((5 6,7 8),(1 2,5 4)))",
			expected: orb.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}, {{1, 2}, {5, 4}}}},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalMultiPolygon(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}

func TestUnmarshalCollection(t *testing.T) {
	cases := []struct {
		s        string
		expected orb.Collection
	}{
		// empty
		{
			s:        "GEOMETRYCOLLECTION EMPTY",
			expected: orb.Collection{},
		},
		// multi polygon
		{
			s:        "GEOMETRYCOLLECTION(POINT(1 2),LINESTRING(3 4,5 6))",
			expected: orb.Collection{orb.Point{1, 2}, orb.LineString{{3, 4}, {5, 6}}},
		},
		{
			s: "GEOMETRYCOLLECTION(POINT(1 2),LINESTRING(3 4,5 6),MULTILINESTRING((1 2,3 4),(5 6,7 8)),POLYGON((0 0,1 0,1 1,0 0)),POLYGON((1 2,3 4),(5 6,7 8)),MULTIPOLYGON(((1 2,3 4)),((5 6,7 8),(1 2,5 4)))",
			expected: orb.Collection{
				orb.Point{1, 2},
				orb.LineString{{3, 4}, {5, 6}},
				orb.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
				orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
				orb.Polygon{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
				orb.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}, {{1, 2}, {5, 4}}}},
			},
		},
	}

	for _, tc := range cases {
		geom, err := UnmarshalCollection(tc.s)
		if err != nil {
			// t.Fatal(err)
		}
		if !geom.Equal(tc.expected) {
			t.Log(geom)
			t.Log(tc.expected)
			t.Errorf("incorrect wkt unmarshalling")
		}
	}
}
