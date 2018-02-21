package quadtree

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

func TestNew(t *testing.T) {
	bound := orb.Bound{Min: orb.Point{0, 2}, Max: orb.Point{1, 3}}
	qt := New(bound)

	if !qt.Bound().Equal(bound) {
		t.Errorf("should use provided bound, got %v", qt.Bound())
	}
}

func TestQuadtreeAdd(t *testing.T) {
	p := orb.Point{}
	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	for i := 0; i < 10; i++ {
		// should be able to insert the same point over and over.
		qt.Add(p)
	}
}

func TestQuadtreeRemove(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	mp := orb.MultiPoint{}
	for i := 0; i < 1000; i++ {
		mp = append(mp, orb.Point{r.Float64(), r.Float64()})
		qt.Add(mp[i])
	}

	for i := 0; i < 1000; i += 3 {
		qt.Remove(mp[i], nil)
		mp[i] = orb.Point{-10000, -10000}
	}

	// make sure finding still works for 1000 random points
	for i := 0; i < 1000; i++ {
		p := orb.Point{r.Float64(), r.Float64()}

		f := qt.Find(p)
		_, j := planar.DistanceFromWithIndex(mp, p)

		if e := mp[j]; !e.Equal(f.Point()) {
			t.Errorf("index: %d, unexpected point %v != %v", i, e, f.Point())
		}
	}
}

func TestQuadtreeFind(t *testing.T) {
	points := orb.MultiPoint{}
	dim := 17

	for i := 0; i < dim*dim; i++ {
		points = append(points, orb.Point{float64(i % dim), float64(i / dim)})
	}

	qt := New(points.Bound())
	for _, p := range points {
		qt.Add(p)
	}

	cases := []struct {
		point    orb.Point
		expected orb.Point
	}{
		{point: orb.Point{0.1, 0.1}, expected: orb.Point{0, 0}},
		{point: orb.Point{3.1, 2.9}, expected: orb.Point{3, 3}},
		{point: orb.Point{7.1, 7.1}, expected: orb.Point{7, 7}},
		{point: orb.Point{0.1, 15.9}, expected: orb.Point{0, 16}},
		{point: orb.Point{15.9, 15.9}, expected: orb.Point{16, 16}},
	}

	for i, tc := range cases {
		if v := qt.Find(tc.point); !v.Point().Equal(tc.expected) {
			t.Errorf("incorrect point on %d, got %v", i, v)
		}
	}
}

func TestQuadtreeFind_Random(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	mp := orb.MultiPoint{}
	for i := 0; i < 1000; i++ {
		mp = append(mp, orb.Point{r.Float64(), r.Float64()})
		qt.Add(mp[i])
	}

	for i := 0; i < 1000; i++ {
		p := orb.Point{r.Float64(), r.Float64()}

		f := qt.Find(p)
		_, j := planar.DistanceFromWithIndex(mp, p)

		if e := mp[j]; !e.Equal(f.Point()) {
			t.Errorf("index: %d, unexpected point %v != %v", i, e, f.Point())
		}
	}
}

func TestQuadtreeMatching(t *testing.T) {
	type dataPointer struct {
		orb.Pointer
		visible bool
	}

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	qt.Add(dataPointer{orb.Point{0, 0}, false})
	qt.Add(dataPointer{orb.Point{1, 1}, true})

	cases := []struct {
		name     string
		filter   FilterFunc
		point    orb.Point
		expected orb.Point
	}{
		{
			name:     "no filtred",
			point:    orb.Point{0.1, 0.1},
			expected: orb.Point{0, 0},
		},
		{
			name:     "with filter",
			filter:   func(p orb.Pointer) bool { return p.(dataPointer).visible },
			point:    orb.Point{0.1, 0.1},
			expected: orb.Point{1, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := qt.Matching(tc.point, tc.filter)
			if !v.Point().Equal(tc.expected) {
				t.Errorf("incorrect point %v != %v", v, tc.expected)
			}
		})
	}
}

func TestQuadtreeKNearest(t *testing.T) {
	type dataPointer struct {
		orb.Pointer
		visible bool
	}

	q := New(orb.Bound{Max: orb.Point{5, 5}})
	q.Add(dataPointer{orb.Point{0, 0}, false})
	q.Add(dataPointer{orb.Point{1, 1}, true})
	q.Add(dataPointer{orb.Point{2, 2}, false})
	q.Add(dataPointer{orb.Point{3, 3}, true})
	q.Add(dataPointer{orb.Point{4, 4}, false})
	q.Add(dataPointer{orb.Point{5, 5}, true})

	filters := map[bool]FilterFunc{
		false: nil,
		true:  func(p orb.Pointer) bool { return p.(dataPointer).visible },
	}

	cases := []struct {
		name     string
		filtered bool
		point    orb.Point
		expected []orb.Point
	}{
		{
			name:     "unfiltered",
			filtered: false,
			point:    orb.Point{0.1, 0.1},
			expected: []orb.Point{{0, 0}, {1, 1}},
		},
		{
			name:     "filtered",
			filtered: true,
			point:    orb.Point{0.1, 0.1},
			expected: []orb.Point{{1, 1}, {3, 3}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.filtered {
				v := q.KNearest(nil, tc.point, 2)
				if len(v) != len(tc.expected) {
					t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
				}
			}

			v := q.KNearestMatching(nil, tc.point, 2, filters[tc.filtered])
			if len(v) != len(tc.expected) {
				t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
			}

			result := make([]orb.Point, 0)
			for _, p := range v {
				result = append(result, p.Point())
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i][0] < result[j][0]
			})

			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i][0] < tc.expected[j][0]
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Log(result)
				t.Log(tc.expected)
				t.Errorf("incorrect results")
			}
		})
	}
}

func TestQuadtreeKNearest_DistanceLimit(t *testing.T) {
	type dataPointer struct {
		orb.Pointer
		visible bool
	}

	q := New(orb.Bound{Max: orb.Point{5, 5}})
	q.Add(dataPointer{orb.Point{0, 0}, false})
	q.Add(dataPointer{orb.Point{1, 1}, true})
	q.Add(dataPointer{orb.Point{2, 2}, false})
	q.Add(dataPointer{orb.Point{3, 3}, true})
	q.Add(dataPointer{orb.Point{4, 4}, false})
	q.Add(dataPointer{orb.Point{5, 5}, true})

	filters := map[bool]FilterFunc{
		false: nil,
		true:  func(p orb.Pointer) bool { return p.(dataPointer).visible },
	}

	cases := []struct {
		name     string
		filtered bool
		distance float64
		point    orb.Point
		expected []orb.Point
	}{
		{
			name:     "filtered",
			filtered: true,
			distance: 5,
			point:    orb.Point{0.1, 0.1},
			expected: []orb.Point{{1, 1}, {3, 3}},
		},
		{
			name:     "unfiltered",
			filtered: false,
			distance: 1,
			point:    orb.Point{0.1, 0.1},
			expected: []orb.Point{{0, 0}},
		},
	}

	var v []orb.Pointer
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v = q.KNearestMatching(v, tc.point, 5, filters[tc.filtered], tc.distance)
			if len(v) != len(tc.expected) {
				t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
			}

			result := make([]orb.Point, 0)
			for _, p := range v {
				result = append(result, p.Point())
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i][0] < result[j][0]
			})

			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i][0] < tc.expected[j][0]
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Log(result)
				t.Log(tc.expected)
				t.Errorf("incorrect results")
			}
		})
	}
}

func TestQuadtreeInBoundMatching(t *testing.T) {
	type dataPointer struct {
		orb.Pointer
		visible bool
	}

	q := New(orb.Bound{Max: orb.Point{5, 5}})
	q.Add(dataPointer{orb.Point{0, 0}, false})
	q.Add(dataPointer{orb.Point{1, 1}, true})
	q.Add(dataPointer{orb.Point{2, 2}, false})
	q.Add(dataPointer{orb.Point{3, 3}, true})
	q.Add(dataPointer{orb.Point{4, 4}, false})
	q.Add(dataPointer{orb.Point{5, 5}, true})

	filters := map[bool]FilterFunc{
		false: nil,
		true:  func(p orb.Pointer) bool { return p.(dataPointer).visible },
	}

	cases := []struct {
		name     string
		filtered bool
		expected []orb.Point
	}{
		{
			name:     "unfiltered",
			filtered: false,
			expected: []orb.Point{{0, 0}, {1, 1}, {2, 2}},
		},
		{
			name:     "filtered",
			filtered: true,
			expected: []orb.Point{{1, 1}},
		},
	}

	bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{2, 2}}

	var v []orb.Pointer
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v = q.InBoundMatching(v, bound, filters[tc.filtered])
			if len(v) != len(tc.expected) {
				t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
			}

			result := make([]orb.Point, 0)
			for _, p := range v {
				result = append(result, p.Point())
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i][0] < result[j][0]
			})

			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i][0] < tc.expected[j][0]
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Log(result)
				t.Log(tc.expected)
				t.Errorf("incorrect results")
			}
		})
	}

}

func TestQuadtreeInBound_Random(t *testing.T) {
	r := rand.New(rand.NewSource(43))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	mp := orb.MultiPoint{}
	for i := 0; i < 1000; i++ {
		mp = append(mp, orb.Point{r.Float64(), r.Float64()})
		qt.Add(mp[i])
	}

	for i := 0; i < 1000; i++ {
		p := orb.Point{r.Float64(), r.Float64()}

		b := orb.Bound{Min: p, Max: p}
		b = b.Pad(0.1)
		ps := qt.InBound(nil, b)

		// find the right answer brute force
		var list []orb.Pointer
		for _, p := range mp {
			if b.Contains(p) {
				list = append(list, p)
			}
		}

		if len(list) != len(ps) {
			t.Errorf("index: %d, lengths not equal %v != %v", i, len(list), len(ps))
		}
	}
}
