package planar

import "testing"

func TestSegmentNew(t *testing.T) {
	a := NewPoint(1, 1)
	b := NewPoint(2, 2)

	s := NewSegment(a, b)
	if !s.A().Equal(a) {
		t.Errorf("not equal: %v != %v", s.A(), a)
	}

	if !s.B().Equal(b) {
		t.Errorf("not equal: %v != %v", s.B(), b)
	}
}

func TestSegmentDistanceFrom(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(0, 10))

	cases := []struct {
		point  Point
		result float64
	}{
		{
			point:  NewPoint(1, 5),
			result: 1,
		},
		{
			point:  NewPoint(0, 2),
			result: 0,
		},
		{
			point:  NewPoint(0, -5),
			result: 5,
		},
		{
			point:  NewPoint(0, 13),
			result: 3,
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			if d := s.DistanceFrom(tc.point); d != tc.result {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}

	s = NewSegment(NewPoint(0, 0), NewPoint(0, 0))
	answer := 5.0
	if d := s.DistanceFrom(NewPoint(3, 4)); d != answer {
		t.Errorf("incorrect distance: %v != %v", d, answer)
	}
}

func TestSegmentDistanceFromSquared(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(0, 10))

	cases := []struct {
		point  Point
		result float64
	}{
		{
			point:  NewPoint(1, 5),
			result: 1,
		},
		{
			point:  NewPoint(0, 2),
			result: 0,
		},
		{
			point:  NewPoint(0, -5),
			result: 25,
		},
		{
			point:  NewPoint(0, 13),
			result: 9,
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			if d := s.DistanceFromSquared(tc.point); d != tc.result {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}

	s = NewSegment(NewPoint(0, 0), NewPoint(0, 0))
	answer := 25.0
	if d := s.DistanceFromSquared(NewPoint(3, 4)); d != answer {
		t.Errorf("incorrect distance: %v != %v", d, answer)
	}
}

func TestSegmentProject(t *testing.T) {
	l1 := NewSegment(NewPoint(1, 2), NewPoint(3, 4))

	cases := []struct {
		point  Point
		result float64
	}{
		{
			point:  NewPoint(1, 2),
			result: 0.0,
		},
		{
			point:  NewPoint(3, 4),
			result: 1.0,
		},
		{
			point:  NewPoint(2, 3),
			result: 0.5,
		},
		{
			point:  NewPoint(5, 6),
			result: 2.0,
		},
		{
			point:  NewPoint(-1, 0),
			result: -1.0,
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			proj := l1.Project(tc.point)
			if proj != tc.result {
				t.Errorf("incorrect project: %v != %v", proj, tc.result)
			}
		})
	}

	// point off of segment
	l2 := NewSegment(NewPoint(1, 1), NewPoint(3, 3))
	proj := l2.Project(NewPoint(1, 2))
	expected := 0.25
	if proj != expected {
		t.Errorf("incorrect project: %v != %v", proj, expected)
	}

	// segment of distance 0
	l3 := NewSegment(NewPoint(1, 1), NewPoint(1, 1))
	proj = l3.Project(NewPoint(1, 2))
	expected = 0.0
	if proj != expected {
		t.Errorf("incorrect project: %v != %v", proj, expected)
	}
}

func TestSegmentDistance(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(3, 4))
	if d := s.Distance(); d != 5 {
		t.Errorf("incorrect distance: %v != %v", d, 5)
	}
}

func TestSegmentSquaredDistance(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(3, 4))
	if d := s.DistanceSquared(); d != 25 {
		t.Errorf("incorrect distance: %v != %v", d, 25)
	}
}

func TestSegmentInterpolate(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(5, 10))

	cases := []struct {
		percent float64
		result  Point
	}{
		{
			percent: 0.2,
			result:  NewPoint(1, 2),
		},
		{
			percent: 0.8,
			result:  NewPoint(4, 8),
		},
		{
			percent: -0.2,
			result:  NewPoint(-1, -2),
		},
		{
			percent: 1.2,
			result:  NewPoint(6, 12),
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			if p := s.Interpolate(tc.percent); !p.Equal(tc.result) {
				t.Errorf("interpolate incorrect: %v != %v", p, tc.result)
			}
		})
	}
}

func TestSegmentSide(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(0, 10))

	cases := []struct {
		name   string
		point  Point
		result int
	}{
		{
			name:   "collinear",
			point:  NewPoint(0, -5),
			result: 0,
		},
		{
			name:   "on the segment",
			point:  NewPoint(0, 5),
			result: 0,
		},
		{
			name:   "right",
			point:  NewPoint(1, 5),
			result: -1,
		},
		{
			name:   "left",
			point:  NewPoint(-1, 5),
			result: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if o := s.Side(tc.point); o != tc.result {
				t.Errorf("incorrect side: %d != %d", o, tc.result)
			}
		})
	}
}

func TestSegmentIntersection(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(1, 1))

	cases := []struct {
		name    string
		segment Segment
		result  Point
	}{
		{
			name:    "end point match",
			segment: NewSegment(NewPoint(1, 1), NewPoint(2, 3)),
			result:  NewPoint(1, 1),
		},
		{
			name:    "start point match",
			segment: NewSegment(NewPoint(1, 10), NewPoint(0, 0)),
			result:  NewPoint(0, 0),
		},
		{
			name:    "through the middle",
			segment: NewSegment(NewPoint(0, 1), NewPoint(1, 0)),
			result:  NewPoint(0.5, 0.5),
		},
		{
			name:    "through the middle, longer",
			segment: NewSegment(NewPoint(0, 1), NewPoint(2, -1)),
			result:  NewPoint(0.5, 0.5),
		},
		{
			name:    "through the middle, longer",
			segment: NewSegment(NewPoint(0.5, 0.5), NewPoint(2, -1)),
			result:  NewPoint(0.5, 0.5),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if p, _ := s.Intersection(tc.segment); !p.Equal(tc.result) {
				t.Errorf("intersection expected: %v != %v", p, tc.result)
			}
		})
	}
}

func TestSegmentIntersectionErrors(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(1, 1))

	cases := []struct {
		name    string
		segment Segment
	}{
		{
			name:    "parallel",
			segment: NewSegment(NewPoint(1, 0), NewPoint(2, 1)),
		},
		{
			name:    "not intersecting",
			segment: NewSegment(NewPoint(1, 0), NewPoint(3, 1)),
		},
		{
			// TODO: this case is a bit weird.
			name:    "share just endpoint, but collinear",
			segment: NewSegment(NewPoint(1, 1), NewPoint(2, 2)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if p, err := s.Intersection(tc.segment); err == nil {
				t.Errorf("no intersection expected: %v", p)
			}
		})
	}
}

func TestSegmentIntersects(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(1, 1))

	cases := []struct {
		name    string
		segment Segment
		result  bool
	}{
		{
			name:    "parallel",
			segment: NewSegment(NewPoint(1, 0), NewPoint(2, 1)),
			result:  false,
		},
		{
			name:    "cross in the middle",
			segment: NewSegment(NewPoint(1, 0), NewPoint(0, 1)),
			result:  true,
		},
		{
			name:    "cross in the middle, longer",
			segment: NewSegment(NewPoint(1, 1), NewPoint(2, 1)),
			result:  true,
		},
		{
			name:    "share endpoint, parallel",
			segment: NewSegment(NewPoint(1, 1), NewPoint(2, 2)),
			result:  true,
		},
		{
			name:    "on the segment",
			segment: NewSegment(NewPoint(0.5, 0.5), NewPoint(2, 2)),
			result:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if r := s.Intersects(tc.segment); r != tc.result {
				t.Errorf("incorrect intersect: %v != %v", r, tc.result)
			}
		})
	}

	// segment with endpoint on segment should intersect each other both ways
	s2 := NewSegment(NewPoint(0.5, 0.5), NewPoint(2, 2))

	if !s.Intersects(s2) {
		t.Errorf("should intersect")
	}

	if !s2.Intersects(s) {
		t.Errorf("should intersect")
	}
}

func TestSegmentMidpoint(t *testing.T) {
	s := NewSegment(NewPoint(0, 0), NewPoint(10, 20))

	expected := NewPoint(5, 10)
	if p := s.Midpoint(); !p.Equal(expected) {
		t.Errorf("midpoint incorrect: %v != %v", p, expected)
	}
}

func TestSegmentBound(t *testing.T) {
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	s := NewSegment(a, b)

	expected := NewRect(1, 3, 2, 4)
	if b := s.Bound(); !b.Equal(expected) {
		t.Errorf("bounds incorrect: %v != %v", b, expected)
	}

	if b := s.Reverse().Bound(); !b.Equal(expected) {
		t.Errorf("bounds incorrect: %v != %v", b, expected)
	}
}

func TestSegmentReverse(t *testing.T) {
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	s := NewSegment(a, b).Reverse()

	if !s.A().Equal(b) || !s.B().Equal(a) {
		t.Errorf("reverse did not work: %v", s)
	}
}

func TestSegmentEqual(t *testing.T) {
	s1 := NewSegment(NewPoint(1, 2), NewPoint(3, 4))
	s2 := NewSegment(NewPoint(1, 2), NewPoint(3, 4))

	if !s1.Equal(s2) || !s2.Equal(s1) {
		t.Errorf("segments not equal: %v != %v", s1, s2)
	}
}

func TestSegmentWKT(t *testing.T) {
	l := NewSegment(NewPoint(1, 2), NewPoint(3, 4))

	answer := "LINESTRING(1 2,3 4)"
	if s := l.WKT(); s != answer {
		t.Errorf("incorrect wkt: %v != %v", s, answer)
	}
}
