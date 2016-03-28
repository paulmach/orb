package geo

import (
	"encoding/json"
	"testing"
)

func TestPointJSON(t *testing.T) {
	p1 := NewPoint(1, 2.1)

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[1,2.1]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var p2 Point
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equal(p2) {
		t.Errorf("unmarshal incorrect, got %v", p2)
	}
}

func TestPathJSON(t *testing.T) {
	p1 := NewPath()
	p1 = append(p1,
		NewPoint(1.5, 2.5),
		NewPoint(3.5, 4.5),
		NewPoint(5.5, 6.5),
	)

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5],[5.5,6.5]]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var p2 Path
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equal(p2) {
		t.Errorf("unmarshal incorrect, got %v", p2)
	}

	// empty path
	p1 = NewPath()
	data, err = json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}
}
