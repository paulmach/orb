package wkt

import (
	"testing"

	"github.com/paulmach/orb"
)

// TestUnmarshalCollection_roundTrip checks the self-oracle: whatever Marshal
// emits for a GEOMETRYCOLLECTION, Unmarshal must read back to an equal value.
// EMPTY members have no parentheses, which the member splitter used to drop.
func TestUnmarshalCollection_roundTrip(t *testing.T) {
	cases := []struct {
		name string
		gc   orb.Collection
	}{
		// each empty-capable geometry type as an EMPTY sole member
		{"linestring empty", orb.Collection{orb.LineString{}}},
		{"polygon empty", orb.Collection{orb.Polygon{}}},
		{"multipoint empty", orb.Collection{orb.MultiPoint{}}},
		{"multilinestring empty", orb.Collection{orb.MultiLineString{}}},
		{"multipolygon empty", orb.Collection{orb.MultiPolygon{}}},
		{"collection empty", orb.Collection{orb.Collection{}}},

		// mixed empty and non-empty members, in both orders
		{"non-empty then empty", orb.Collection{
			orb.Point{1, 2}, orb.LineString{},
		}},
		{"empty then non-empty", orb.Collection{
			orb.LineString{}, orb.Point{1, 2},
		}},
		{"several empties between non-empties", orb.Collection{
			orb.Point{1, 2}, orb.Polygon{}, orb.MultiPoint{},
			orb.LineString{{3, 4}, {5, 6}},
		}},

		// nested collections carrying empty members
		{"nested collection with empty", orb.Collection{
			orb.Collection{orb.LineString{}},
			orb.Point{7, 8},
		}},

		// pure non-empty members (regression: must stay unchanged)
		{"non-empty only", orb.Collection{
			orb.Point{1, 2},
			orb.LineString{{3, 4}, {5, 6}},
			orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
		}},
		{"nested non-empty", orb.Collection{
			orb.Collection{orb.Point{1, 2}, orb.LineString{{3, 4}, {5, 6}}},
			orb.MultiPolygon{{{{1, 2}, {3, 4}}}},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := MarshalString(tc.gc)

			got, err := Unmarshal(s)
			if err != nil {
				t.Fatalf("own output %q rejected by Unmarshal: %v", s, err)
			}

			gc, ok := got.(orb.Collection)
			if !ok {
				t.Fatalf("expected orb.Collection, got %T", got)
			}

			if len(gc) != len(tc.gc) {
				t.Fatalf("member count changed: marshaled %q\n  in  %d members: %v\n  out %d members: %v",
					s, len(tc.gc), tc.gc, len(gc), gc)
			}

			if !gc.Equal(tc.gc) {
				t.Fatalf("round-trip inequality: marshaled %q\n  in  %v\n  out %v", s, tc.gc, gc)
			}

			// a second marshal must be byte-identical: no member added, lost or reshaped.
			if again := MarshalString(gc); again != s {
				t.Fatalf("marshal not idempotent after round-trip:\n  first  %q\n  second %q", s, again)
			}
		})
	}
}

// TestSplitGeometryCollection covers the member splitter directly, including
// EMPTY members (no parentheses) and top-level commas inside nested members.
func TestSplitGeometryCollection(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single empty member",
			input:    "(LINESTRING EMPTY)",
			expected: []string{"LINESTRING EMPTY"},
		},
		{
			name:     "non-empty then empty",
			input:    "(POINT(1 2),LINESTRING EMPTY)",
			expected: []string{"POINT(1 2)", "LINESTRING EMPTY"},
		},
		{
			name:     "two non-empty members",
			input:    "(POINT(1 2),LINESTRING(3 4,5 6))",
			expected: []string{"POINT(1 2)", "LINESTRING(3 4,5 6)"},
		},
		{
			name:     "nested collection member is not split on its inner comma",
			input:    "(GEOMETRYCOLLECTION(POINT(1 2),POINT(3 4)),POINT(5 6))",
			expected: []string{"GEOMETRYCOLLECTION(POINT(1 2),POINT(3 4))", "POINT(5 6)"},
		},
		{
			name:     "multipolygon member with nested commas",
			input:    "(MULTIPOLYGON(((1 2,3 4)),((5 6,7 8))),POINT(9 0))",
			expected: []string{"MULTIPOLYGON(((1 2,3 4)),((5 6,7 8)))", "POINT(9 0)"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := splitGeometryCollection(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.expected) {
				t.Fatalf("member count: got %d %v, want %d %v",
					len(got), got, len(tc.expected), tc.expected)
			}

			for i := range tc.expected {
				if got[i] != tc.expected[i] {
					t.Errorf("member %d: got %q, want %q", i, got[i], tc.expected[i])
				}
			}
		})
	}
}
