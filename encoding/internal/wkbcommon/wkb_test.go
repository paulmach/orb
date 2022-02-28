package wkbcommon

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/paulmach/orb"
)

func TestMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		Marshal(g, binary.BigEndian)
	}
}

func TestMustMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		MustMarshal(g, binary.BigEndian)
	}
}

func compare(t testing.TB, e orb.Geometry, b []byte) {
	t.Helper()

	// Decoder
	g, err := NewDecoder(bytes.NewReader(b)).Decode()
	if err != nil {
		t.Fatalf("decoder: read error: %v", err)
	}

	if !orb.Equal(g, e) {
		t.Errorf("decoder: incorrect geometry: %v != %v", g, e)
	}

	// Umarshal
	g, err = Unmarshal(b)
	if err != nil {
		t.Fatalf("unmarshal: read error: %v", err)
	}

	if !orb.Equal(g, e) {
		t.Errorf("unmarshal: incorrect geometry: %v != %v", g, e)
	}

	var data []byte
	if b[0] == 0 {
		data, err = Marshal(g, binary.BigEndian)
	} else {
		data, err = Marshal(g, binary.LittleEndian)
	}
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("marshal: incorrent encoding")
	}

	// preallocation
	if len(data) != GeomLength(e) {
		t.Errorf("prealloc length: %v != %v", len(data), GeomLength(e))
	}
}
