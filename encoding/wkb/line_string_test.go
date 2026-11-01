package wkb

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/internal/wkbcommon"
)

var (
	testLineString     = orb.LineString{{1, 2}, {3, 4}}
	testLineStringData = []byte{
		1, 2, 0, 0, 0,
		2, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64}
)

func TestLineString(t *testing.T) {
	large := orb.LineString{}
	for i := 0; i < wkbcommon.MaxPointsAlloc+100; i++ {
		large = append(large, orb.Point{float64(i), float64(-i)})
	}

	cases := []struct {
		name     string
		data     []byte
		expected orb.LineString
	}{
		{
			name:     "line string",
			data:     testLineStringData,
			expected: testLineString,
		},
		{
			name:     "large line string",
			data:     MustMarshal(large),
			expected: large,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}

var (
	testMultiLineString = orb.MultiLineString{
		{{10, 10}, {20, 20}, {10, 40}},
		{{40, 40}, {30, 30}, {40, 20}, {30, 10}},
	}
	testMultiLineStringData = []byte{
		0x01, 0x05, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, // Number of Lines 2
		0x01,                   // Encoding Little
		0x02, 0x00, 0x00, 0x00, // Type
		0x03, 0x00, 0x00, 0x00, // Number of points 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y2 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
		0x01,                   // Encoding Little
		0x02, 0x00, 0x00, 0x00, // Type LineString
		0x04, 0x00, 0x00, 0x00, // Number of Points 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X1 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X2 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X3 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y4 10
	}

	testMultiLineStringSingle = orb.MultiLineString{
		{{10, 10}, {20, 20}, {10, 40}},
	}
	testMultiLineStringSingleData = []byte{
		0x01, 0x05, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, // Number of Lines 2
		0x01,                   // Encoding Little
		0x02, 0x00, 0x00, 0x00, // Type
		0x03, 0x00, 0x00, 0x00, // Number of points 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y2 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
	}
)

func TestMultiLineString(t *testing.T) {
	large := orb.MultiLineString{}
	for i := 0; i < wkbcommon.MaxMultiAlloc+100; i++ {
		large = append(large, orb.LineString{})
	}

	cases := []struct {
		name     string
		data     []byte
		expected orb.MultiLineString
	}{
		{
			name:     "multi line string",
			data:     testMultiLineStringData,
			expected: testMultiLineString,
		},
		{
			name:     "single multi line string",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle,
		},
		{
			name:     "large",
			data:     MustMarshal(large),
			expected: large,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}

func TestLineString_pointCountOverflow(t *testing.T) {
	// A crafted little-endian linestring header claims a point count whose
	// byte size (count * 16) overflows a uint32 and wraps to a small value.
	// Before the length check was done in 64-bit space this slipped past the
	// bounds guard and read past the end of the buffer.
	data := []byte{
		0x01,                   // little endian
		0x02, 0x00, 0x00, 0x00, // type: linestring
		0x01, 0x00, 0x00, 0x10, // point count 0x10000001; *16 wraps to 16 in uint32
	}
	data = append(data, make([]byte, 16)...) // only one point's worth of data

	if _, err := Unmarshal(data); err != ErrNotWKB {
		t.Fatalf("expected ErrNotWKB, got %v", err)
	}
}
