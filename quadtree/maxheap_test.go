package quadtree

import (
	"math/rand"
	"testing"
)

func TestMaxHeap(t *testing.T) {
	r := rand.New(rand.NewSource(22))

	for i := 1; i < 100; i++ {
		h := make(maxHeap, 0, i)
		for j := 0; j < i; j++ {
			h.Push(&heapItem{distance: r.Float64()})
		}

		current := h.Pop().distance
		for len(h) > 0 {
			next := h.Pop().distance
			if next > current {
				t.Errorf("incorrect")
			}

			current = next
		}
	}
}
