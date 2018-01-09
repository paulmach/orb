package wkb

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/paulmach/orb"
)

func TestLineString(t *testing.T) {
	cases := []struct {
		name     string
		bytes    []byte
		bom      binary.ByteOrder
		expected orb.LineString
	}{
		{
			name:     "line string",
			bytes:    []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64},
			bom:      binary.LittleEndian,
			expected: orb.LineString{{1, 2}, {3, 4}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			ls, err := readLineString(bytes.NewReader(tc.bytes), tc.bom)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if len(ls) != len(tc.expected) {
				t.Fatalf("incorrect length: %d != %d", len(ls), len(tc.expected))
			}

			if !ls.Equal(tc.expected) {
				t.Errorf("incorrect line: %v != %v", ls, tc.expected)
			}
		})
	}
}

func TestMultiLineString(t *testing.T) {
	cases := []struct {
		name     string
		bytes    []byte
		bom      binary.ByteOrder
		expected orb.MultiLineString
	}{
		{
			name: "multi line string",
			bytes: []byte{
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
			},
			bom: binary.LittleEndian,
			expected: orb.MultiLineString{
				{{10, 10}, {20, 20}, {10, 40}},
				{{40, 40}, {30, 30}, {40, 20}, {30, 10}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mls, err := readMultiLineString(bytes.NewReader(tc.bytes), tc.bom)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if len(mls) != len(tc.expected) {
				t.Fatalf("incorrect length: %d != %d", len(mls), len(tc.expected))
			}

			for i := range mls {
				if !mls[i].Equal(tc.expected[i]) {
					t.Errorf("expected[%v]: %v != %v", i, mls[i], tc.expected[i])
				}
			}
		})
	}
}
