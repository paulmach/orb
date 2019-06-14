package geojson

import (
	"testing"
	"time"
)

func TestWhen1(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Errorf("Error instantiating date for test. err=%s", err)
		return
	}
	when1 := When{"Instant", &t1}
	if !when1.Valid() {
		t.Errorf("When is invalid. when=%v", when1)
	}
}

func TestWhen2(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2002-01-01T22:08:41+00:00")
	if err != nil {
		t.Errorf("Error instantiating date for test. err=%s", err)
		return
	}
	when1 := When{"StrangeType", &t1}
	if when1.Valid() {
		t.Errorf("When should be marked invalid. when=%v", when1)
	}
}
