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

func TestMustMarshal(t *testing.T) {
	for _, g := range orb.AllGeometries {
		MustMarshal(g, binary.BigEndian)
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
		t.Errorf("incorrent encoding")
	}

	// preallocation
	if len(data) != geomLength(e) {
		t.Errorf("prealloc length: %v != %v", len(data), geomLength(e))
	}

	// Scanner
	var sg orb.Geometry

	switch e.(type) {
	case orb.Point:
		p := orb.Point{}
		err = Scanner(&p).Scan(b)
		sg = p
	case orb.MultiPoint:
		mp := orb.MultiPoint{}
		err = Scanner(&mp).Scan(b)
		sg = mp
	case orb.LineString:
		ls := orb.LineString{}
		err = Scanner(&ls).Scan(b)
		sg = ls
	case orb.MultiLineString:
		mls := orb.MultiLineString{}
		err = Scanner(&mls).Scan(b)
		sg = mls
	case orb.Polygon:
		p := orb.Polygon{}
		err = Scanner(&p).Scan(b)
		sg = p
	case orb.MultiPolygon:
		mp := orb.MultiPolygon{}
		err = Scanner(&mp).Scan(b)
		sg = mp
	case orb.Collection:
		c := orb.Collection{}
		err = Scanner(&c).Scan(b)
		sg = c
	default:
		t.Fatalf("unknown type: %T", e)
	}

	if err != nil {
		t.Errorf("scan error: %v", err)
	}

	if sg.GeoJSONType() != e.GeoJSONType() {
		t.Errorf("scanning to wrong type: %v != %v", sg.GeoJSONType(), e.GeoJSONType())
	}

	if !orb.Equal(sg, e) {
		t.Errorf("incorrect geometry: %v != %v", sg, e)
	}
}
