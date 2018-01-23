package wkb

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"testing"

	"github.com/paulmach/orb"
)

func TestMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		Marshal(g, binary.BigEndian)
	}
}

func BenchmarkEncode_Point(b *testing.B) {
	g := orb.Point{1, 2}
	e := NewEncoder(ioutil.Discard)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(g)
	}
}

func BenchmarkEncode_LineString(b *testing.B) {
	g := orb.LineString{
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
	}
	e := NewEncoder(ioutil.Discard)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(g)
	}
}

func compare(t testing.TB, e orb.Geometry, b []byte) {
	t.Helper()

	g, err := NewDecoder(bytes.NewReader(b)).Decode()
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if !orb.Equal(g, e) {
		t.Errorf("incorrect geometry: %v != %v", g, e)
	}

	data, err := Marshal(g)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("incorrent encoding")
	}
}
