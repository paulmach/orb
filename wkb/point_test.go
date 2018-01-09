package wkb

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/paulmach/orb"
)

func TestPoint(t *testing.T) {
	cases := []struct {
		name     string
		bytes    []byte
		bom      binary.ByteOrder
		expected orb.Point
	}{
		{
			name: "point",
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x46, 0x81, 0xF6, 0x23, 0x2E, 0x4A, 0x5D, 0xC0,
				0x03, 0x46, 0x1B, 0x3C, 0xAF, 0x5B, 0x40, 0x40,
			},
			bom:      binary.LittleEndian,
			expected: orb.Point{-117.15906619141342, 32.71628524142945},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := readPoint(bytes.NewReader(tc.bytes), tc.bom)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("incorrect point: %v != %v", p, tc.expected)
			}
		})
	}
}

func TestMultiPoint(t *testing.T) {
	cases := []struct {
		name     string
		bytes    []byte
		bom      binary.ByteOrder
		expected orb.MultiPoint
	}{
		{
			name: "multi point",
			bytes: []byte{
				0x04, 0x00, 0x00, 0x00, // Number of Points (4)
				0x01,                   // Byte Order Little
				0x01, 0x00, 0x00, 0x00, // Type Point (1)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X1 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
				0x01,                   // Byte Order Little
				0x01, 0x00, 0x00, 0x00, // Type Point (1)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 30
				0x01,                   // Byte Order Little
				0x01, 0x00, 0x00, 0x00, // Type Point (1)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
				0x01,                   // Byte Order Little
				0x01, 0x00, 0x00, 0x00, // Type Point (1)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y4 10
			},
			bom:      binary.LittleEndian,
			expected: orb.MultiPoint{{10, 40}, {40, 30}, {20, 20}, {30, 10}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mp, err := readMultiPoint(bytes.NewReader(tc.bytes), tc.bom)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if len(mp) != len(tc.expected) {
				t.Fatalf("incorrect length: %d != %d", len(mp), len(tc.expected))
			}

			if !mp.Equal(tc.expected) {
				t.Errorf("incorrect values: %v != %v", mp, tc.expected)
			}
		})
	}
}
