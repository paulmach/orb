package topojson

func (t *Topology) cut() {
	junctions := t.join()

	for _, line := range t.lines {
		mid := line.Start
		end := line.End

		mid += 1
		for mid < end {
			if junctions.Has(t.coordinates[mid]) {
				next := &arc{Start: mid, End: line.End}
				line.End = mid
				line.Next = next
				line = next
			}
			mid += 1
		}
	}

	for _, ring := range t.rings {
		start := ring.Start
		mid := start
		end := ring.End
		fixed := junctions.Has(t.coordinates[start])

		mid += 1
		for mid < end {
			if junctions.Has(t.coordinates[mid]) {
				if fixed {
					next := &arc{Start: mid, End: ring.End}
					ring.End = mid
					ring.Next = next
					ring = next
				} else {
					// For the first junction, we can rotate rather than cut
					t.rotateCoordinates(start, end, end-mid)
					t.coordinates[end] = t.coordinates[start]
					fixed = true
					mid = start // restart; we may have skipped junctions
				}
			}
			mid += 1
		}
	}
}

func (t *Topology) rotateCoordinates(start, end, offset int) {
	t.reverseCoordinates(start, end)
	t.reverseCoordinates(start, start+offset)
	t.reverseCoordinates(start+offset, end)
}

func (t *Topology) reverseCoordinates(start, end int) {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		t.coordinates[i], t.coordinates[j] = t.coordinates[j], t.coordinates[i]
	}
}
