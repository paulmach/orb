package planar

import "testing"

func TestVectorAdd(t *testing.T) {
	v := NewPoint(1, 2)

	answer := NewPoint(4, 6)
	v = v.Add(NewVector(3, 4))
	if !v.Equal(answer) {
		t.Errorf("vector, add expect %v == %v", v, answer)
	}
}

func TestVectorSub(t *testing.T) {
	v := NewVector(3, 4)

	answer := NewVector(2, 1)
	v = v.Sub(NewVector(1, 3))
	if !v.Equal(answer) {
		t.Errorf("vector, subtract expect %v == %v", v, answer)
	}
}

func TestVectorNormalize(t *testing.T) {
	v := NewVector(5, 0)
	answer := NewVector(1, 0)

	if v = v.Normalize(); !v.Equal(answer) {
		t.Errorf("vector, normalize expect %v == %v", v, answer)
	}

	v = NewVector(0, 5)
	answer = NewVector(0, 1)
	if v = v.Normalize(); !v.Equal(answer) {
		t.Errorf("vector, normalize expect %v == %v", v, answer)
	}

	v = NewVector(0, 0)
	answer = NewVector(0, 0)
	if v = v.Normalize(); !v.Equal(answer) {
		t.Errorf("vector, normalize expect %v == %v", v, answer)
	}
}

func TestVectorScale(t *testing.T) {
	v := NewVector(5, 0)
	answer := NewVector(10, 0)
	if v = v.Scale(2.0); !v.Equal(answer) {
		t.Errorf("vector, scale expect %v == %v", v, answer)
	}

	v = NewVector(0, 5)
	answer = NewVector(0, 15)
	if v = v.Scale(3.0); !v.Equal(answer) {
		t.Errorf("vector, scale expect %v == %v", v, answer)
	}

	v = NewVector(2, 3)
	answer = NewVector(10, 15)
	if v = v.Scale(5.0); !v.Equal(answer) {
		t.Errorf("vector, scale expect %v == %v", v, answer)
	}

	v = NewVector(2, 3)
	answer = NewVector(-10, -15)
	if v = v.Scale(-5.0); !v.Equal(answer) {
		t.Errorf("vector, scale expect %v == %v", v, answer)
	}
}

func TestVectorDot(t *testing.T) {
	v1 := NewVector(0, 0)
	v2 := NewVector(1, 2)
	answer := 0.0
	if d := v1.Dot(v2); d != answer {
		t.Errorf("vector, dot expteced %v == %v", d, answer)
	}

	v1 = NewVector(4, 5)
	answer = 14.0
	if d := v1.Dot(v2); d != answer {
		t.Errorf("vector, dot expteced %v == %v", d, answer)
	}

	// reverse version
	if d := v2.Dot(v1); d != answer {
		t.Errorf("vector, dot expteced %v == %v", d, answer)
	}
}

func TestVectorEqual(t *testing.T) {
	v1 := NewVector(1, 0)
	v2 := NewVector(1, 0)

	v3 := NewVector(2, 3)
	v4 := NewVector(2, 4)

	if !v1.Equal(v2) {
		t.Errorf("vector, equals expect %v == %v", v1, v2)
	}

	if v2.Equal(v3) {
		t.Errorf("vector, equals expect %v != %v", v2, v3)
	}

	if v3.Equal(v4) {
		t.Errorf("vector, equals expect %v != %v", v3, v4)
	}
}
