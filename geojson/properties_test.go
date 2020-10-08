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

func TestFeaturePropertyBool(t *testing.T) {
	f := propertiesTestFeature()

	_, err := f.PropertyBool("random")
	if err == nil {
		t.Errorf("should return error if invalid key")
	}

	b, err := f.PropertyBool("bool")
	if err != nil {
		t.Errorf("should not return error if valid key")
	}

	if b != true {
		t.Errorf("should return proper property")
	}
}

func TestFeaturePropertyInt(t *testing.T) {
	f := propertiesTestFeature()

	_, err := f.PropertyInt("random")
	if err == nil {
		t.Errorf("should return error if invalid key")
	}

	i, err := f.PropertyInt("int")
	if err != nil {
		t.Errorf("should not return error if valid key")
	}

	if i != 1 {
		t.Errorf("should return proper property")
	}
}

func TestFeaturePropertyFloat64(t *testing.T) {
	f := propertiesTestFeature()

	_, err := f.PropertyFloat64("random")
	if err == nil {
		t.Errorf("should return error if invalid key")
	}

	i, err := f.PropertyFloat64("float64")
	if err != nil {
		t.Errorf("should not return error if valid key")
	}

	if i != 1.2 {
		t.Errorf("should return proper property")
	}
}

func TestFeaturePropertyString(t *testing.T) {
	f := propertiesTestFeature()

	_, err := f.PropertyString("random")
	if err == nil {
		t.Errorf("should return error if invalid key")
	}

	s, err := f.PropertyString("string")
	if err != nil {
		t.Errorf("should not return error if valid key")
	}

	if s != "text" {
		t.Errorf("should return proper property")
	}
}

func TestFeaturePropertyMustBool(t *testing.T) {
	f := propertiesTestFeature()

	b := f.PropertyMustBool("random", true)
	if b != true {
		t.Errorf("should return default if property doesn't exist")
	}

	b = f.PropertyMustBool("falsebool", true)
	if b != false {
		t.Errorf("should return proper property, with default")
	}

	b = f.PropertyMustBool("falsebool")
	if b != false {
		t.Errorf("should return proper property, without default")
	}
}

func TestFeaturePropertyMustInt(t *testing.T) {
	f := propertiesTestFeature()

	i := f.PropertyMustInt("random", 10)
	if i != 10 {
		t.Errorf("should return default if property doesn't exist")
	}

	i = f.PropertyMustInt("int", 10)
	if i != 1 {
		t.Errorf("should return proper property, with default")
	}

	i = f.PropertyMustInt("int")
	if i != 1 {
		t.Errorf("should return proper property, without default")
	}

	f.SetProperty("true_int", 5)
	i = f.PropertyMustInt("true_int")
	if i != 5 {
		// json decode makes all things float64,
		// but manually setting will be a true int
		t.Errorf("should work for true integer types")
	}

	i = f.PropertyMustInt("float64")
	if i != 1 {
		t.Errorf("should convert float64 to int")
	}
}

func TestFeaturePropertyMustFloat64(t *testing.T) {
	f := propertiesTestFeature()

	i := f.PropertyMustFloat64("random", 10)
	if i != 10 {
		t.Errorf("should return default if property doesn't exist")
	}

	i = f.PropertyMustFloat64("float64", 10.0)
	if i != 1.2 {
		t.Errorf("should return proper property, with default")
	}

	i = f.PropertyMustFloat64("float64")
	if i != 1.2 {
		t.Errorf("should return proper property, without default")
	}
}

func TestFeaturePropertyMustString(t *testing.T) {
	f := propertiesTestFeature()

	s := f.PropertyMustString("random", "something")
	if s != "something" {
		t.Errorf("should return default if property doesn't exist")
	}

	s = f.PropertyMustString("string", "something")
	if s != "text" {
		t.Errorf("should return proper property, with default")
	}

	s = f.PropertyMustString("string")
	if s != "text" {
		t.Errorf("should return proper property, without default")
	}
}
