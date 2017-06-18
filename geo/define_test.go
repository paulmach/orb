package geo

import (
	"math"
	"testing"
)

var epsilon = 1e-6

func TestDefineYesHaversine(t *testing.T) {

	if yesHaversine(make([]bool, 0)) {
		t.Error("should be false, got true")
	}

	if !yesHaversine([]bool{true}) {
		t.Error("should be true, got false")
	}

	if yesHaversine([]bool{false}) {
		t.Error("should be false, got true")
	}
}

func TestDefineDeg2Rad(t *testing.T) {
	if math.Abs(deg2rad(0.0)) > epsilon {
		t.Error("error")
	}

	if math.Abs(deg2rad(180.0)-math.Pi) > epsilon {
		t.Error("error")
	}

	if math.Abs(deg2rad(360.0)-2*math.Pi) > epsilon {
		t.Error("error")
	}
}

func TestDefineRad2Deg(t *testing.T) {
	if math.Abs(rad2deg(0.0)-0.0) > epsilon {
		t.Error("define, rad2deg error")
	}

	if math.Abs(rad2deg(math.Pi)-180.0) > epsilon {
		t.Error("define, rad2deg error")
	}

	if math.Abs(rad2deg(2*math.Pi)-360.0) > epsilon {
		t.Error("define, rad2deg error")
	}
}
