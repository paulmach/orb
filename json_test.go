package orb

import (
	"encoding/json"
	"testing"
)

func TestPointJSON(t *testing.T) {
	p1 := Point{1, 2.1}

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine: %v", err)
	}

	if string(data) != "[1,2.1]" {
		t.Errorf("incorrect json: %v", string(data))
	}

	var p2 Point
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine: %v", err)
	}

	if !p1.Equal(p2) {
		t.Errorf("not equal: %v", p2)
	}
}

func TestLineStringJSON(t *testing.T) {
	ls1 := LineString{{1.5, 2.5}, {3.5, 4.5}, {5.5, 6.5}}

	data, err := json.Marshal(ls1)
	if err != nil {
		t.Fatalf("should marshal just fine: %v", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5],[5.5,6.5]]" {
		t.Errorf("incorrect data: %v", string(data))
	}

	var ls2 LineString
	err = json.Unmarshal(data, &ls2)
	if err != nil {
		t.Fatalf("should unmarshal just fine: %v", err)
	}

	if !ls1.Equal(ls2) {
		t.Errorf("unmarshal not equal: %v", ls2)
	}

	// empty line
	ls1 = LineString{}
	data, err = json.Marshal(ls1)
	if err != nil {
		t.Errorf("should marshal just fine: %v", err)
	}

	if string(data) != "[]" {
		t.Errorf("incorrect json: %v", string(data))
	}
}
