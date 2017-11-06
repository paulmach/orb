package geo

import (
	"testing"
)

func TestLineStringReverse(t *testing.T) {
	t.Run("1 point line", func(t *testing.T) {
		ls := append(NewLineString(), NewPoint(1, 2))
		rs := ls.Clone()
		rs.Reverse()

		if !rs.Equal(ls) {
			t.Errorf("1 point lines should be equal if reversed")
		}
	})

	cases := []struct {
		name   string
		input  LineString
		output LineString
	}{
		{
			name:   "2 point line",
			input:  append(NewLineString(), NewPoint(1, 2), NewPoint(3, 4)),
			output: append(NewLineString(), NewPoint(3, 4), NewPoint(1, 2)),
		},
		{
			name:   "3 point line",
			input:  append(NewLineString(), NewPoint(1, 2), NewPoint(3, 4), NewPoint(5, 6)),
			output: append(NewLineString(), NewPoint(5, 6), NewPoint(3, 4), NewPoint(1, 2)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reversed := tc.input
			reversed.Reverse()

			if !reversed.Equal(tc.output) {
				t.Errorf("line should be reversed: %v", reversed)
			}

			if !tc.input.Equal(reversed) {
				t.Errorf("should reverse inplace")
			}
		})
	}
}

func TestLineStringWKT(t *testing.T) {
	ls := NewLineString()

	answer := "EMPTY"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	ls = append(ls, NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}

	ls = append(ls, NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := ls.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}
}
