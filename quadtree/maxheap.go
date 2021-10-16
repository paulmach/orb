package quadtree

import "github.com/paulmach/orb"

// maxHeap is used for the knearest list. We need a way to maintain
// the furthest point from the query point in the list, hence maxHeap.
// When we find a point closer than the furthest away, we remove
// furthest and add the new point to the heap.
type maxHeap []*heapItem

type heapItem struct {
	point    orb.Pointer
	distance float64
}

func (h *maxHeap) Push(point orb.Pointer, distance float64) {
	// Common usage is Push followed by a Pop if we have > k points.
	// We're reusing the k+1 heapItem object to reduce memory allocations.
	// First we manaully lengthen the slice,
	// then we see if the last item has been allocated already.

	prevLen := len(*h)
	*h = (*h)[:prevLen+1]
	if (*h)[prevLen] == nil {
		(*h)[prevLen] = &heapItem{point: point, distance: distance}
	} else {
		(*h)[prevLen].point = point
		(*h)[prevLen].distance = distance
	}

	i := len(*h) - 1
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := (*h)[up]

		if distance < parent.distance {
			// parent is further so we're done fixing up the heap.
			break
		}

		// swap nodes
		// (*h)[i] = parent
		(*h)[i].point = parent.point
		(*h)[i].distance = parent.distance

		// (*h)[up] = item
		(*h)[up].point = point
		(*h)[up].distance = distance

		i = up
	}
}

// Pop returns the "greatest" item in the list.
// The returned item should not be saved across push/pop operations.
func (h *maxHeap) Pop() *heapItem {
	removed := (*h)[0]
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]

	mh := (*h)
	if len(mh) == 0 {
		return removed
	}

	// move the last item to the top and reset the heap
	mh[0] = lastItem

	i := 0
	current := mh[i]
	for {
		right := (i + 1) << 1
		left := right - 1

		childIndex := i
		child := mh[childIndex]

		// swap with biggest child
		if left < len(mh) && child.distance < mh[left].distance {
			childIndex = left
			child = mh[left]
		}

		if right < len(mh) && child.distance < mh[right].distance {
			childIndex = right
			child = mh[right]
		}

		// non bigger, so quit
		if childIndex == i {
			break
		}

		// swap the nodes
		mh[i] = child
		mh[childIndex] = current

		i = childIndex
	}

	return removed
}
