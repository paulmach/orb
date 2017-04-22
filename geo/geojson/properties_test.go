package geojson

import (
	"testing"
)

func propertiesTestFeature() *Feature {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	    "properties": {"bool":true,"falsebool":false,"int": 1,"float64": 1.2,"string":"text"}
	  }`

	f, _ := UnmarshalFeature([]byte(rawJSON))
	return f
}

func TestPropertiesMustBool(t *testing.T) {
	f := propertiesTestFeature()

	b := f.Properties.MustBool("random", true)
	if b != true {
		t.Errorf("should return default if property doesn't exist")
	}

	b = f.Properties.MustBool("falsebool", true)
	if b != false {
		t.Errorf("should return proper property, with default")
	}

	b = f.Properties.MustBool("falsebool")
	if b != false {
		t.Errorf("should return proper property, without default")
	}
}

func TestPropertiesMustInt(t *testing.T) {
	f := propertiesTestFeature()

	i := f.Properties.MustInt("random", 10)
	if i != 10 {
		t.Errorf("should return default if property doesn't exist")
	}

	i = f.Properties.MustInt("int", 10)
	if i != 1 {
		t.Errorf("should return proper property, with default")
	}

	i = f.Properties.MustInt("int")
	if i != 1 {
		t.Errorf("should return proper property, without default")
	}

	f.Properties["true_int"] = 5
	i = f.Properties.MustInt("true_int")
	if i != 5 {
		// json decode makes all things float64,
		// but manually setting will be a true int
		t.Errorf("should work for true integer types")
	}

	i = f.Properties.MustInt("float64")
	if i != 1 {
		t.Errorf("should convert float64 to int")
	}
}

func TestPropertiesMustFloat64(t *testing.T) {
	f := propertiesTestFeature()

	i := f.Properties.MustFloat64("random", 10)
	if i != 10 {
		t.Errorf("should return default if property doesn't exist")
	}

	i = f.Properties.MustFloat64("float64", 10.0)
	if i != 1.2 {
		t.Errorf("should return proper property, with default")
	}

	i = f.Properties.MustFloat64("float64")
	if i != 1.2 {
		t.Errorf("should return proper property, without default")
	}
}

func TestPropertiesMustString(t *testing.T) {
	f := propertiesTestFeature()

	s := f.Properties.MustString("random", "something")
	if s != "something" {
		t.Errorf("should return default if property doesn't exist")
	}

	s = f.Properties.MustString("string", "something")
	if s != "text" {
		t.Errorf("should return proper property, with default")
	}

	s = f.Properties.MustString("string")
	if s != "text" {
		t.Errorf("should return proper property, without default")
	}
}
