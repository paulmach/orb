package orb

import (
	"testing"
)

func TestReadUint32(t *testing.T) {
	if v := readUint32([]byte{1, 0, 0, 0}, true); v != 1 {
		t.Errorf("parsed to wrong value, got %v", v)
	}

	if v := readUint32([]byte{1, 0, 0, 1}, true); v != 16777217 {
		t.Errorf("parsed to wrong value, got %v", v)
	}

	if v := readUint32([]byte{1, 0, 0, 0}, false); v != 16777216 {
		t.Errorf("parsed to wrong value, got %v", v)
	}
}

func TestReadFloat64(t *testing.T) {
	if v := readFloat64([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, true); v != 0 {
		t.Errorf("parsed to wrong value, got %v", v)
	}

	if v := readFloat64([]byte{192, 94, 157, 24, 227, 60, 152, 15}, false); v != -122.4546440212 {
		t.Errorf("parsed to wrong value, got %v", v)
	}

	if v := readFloat64([]byte{15, 152, 60, 227, 24, 157, 94, 192}, true); v != -122.4546440212 {
		t.Errorf("parsed to wrong value, got %v", v)
	}
}
