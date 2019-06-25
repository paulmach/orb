package topojson

func (t *Topology) dedup() {
	arcsByEnd := make(map[point][]*arc)

	t.arcs = make([]*arc, 0)

	dedupLine := func(arc *arc) {
		// Does this arc match an existing arc in order?
		startPoint := newPoint(t.coordinates[arc.Start])
		startArcs, startOk := arcsByEnd[startPoint]
		if startOk {
			for _, startArc := range startArcs {
				if t.lineEqual(arc, startArc) {
					arc.Start = startArc.Start
					arc.End = startArc.End
					return
				}
			}
		}

		// Does this arc match an existing arc in reverse order?
		endPoint := newPoint(t.coordinates[arc.End])
		endArcs, endOk := arcsByEnd[endPoint]
		if endOk {
			for _, endArc := range endArcs {
				if t.lineEqualReverse(arc, endArc) {
					arc.Start = endArc.End
					arc.End = endArc.Start
					return
				}
			}
		}

		arcsByEnd[startPoint] = append(startArcs, arc)
		arcsByEnd[endPoint] = append(endArcs, arc)
		t.arcs = append(t.arcs, arc)
	}

	dedupRing := func(arc *arc) {
		// Does this arc match an existing line in order, or reverse order?
		// Rings are closed, so their start point and end point is the same.
		endPoint := newPoint(t.coordinates[arc.Start])
		endArcs, endOk := arcsByEnd[endPoint]
		if endOk {
			for _, endArc := range endArcs {
				if t.ringEqual(arc, endArc) {
					arc.Start = endArc.Start
					arc.End = endArc.End
					return
				}

				if t.ringEqualReverse(arc, endArc) {
					arc.Start = endArc.End
					arc.End = endArc.Start
					return
				}
			}
		}

		// Otherwise, does this arc match an existing ring in order, or reverse order?
		endPoint = newPoint(t.coordinates[arc.Start+t.findMinimumOffset(arc)])
		endArcs, endOk = arcsByEnd[endPoint]
		if endOk {
			for _, endArc := range endArcs {
				if t.ringEqual(arc, endArc) {
					arc.Start = endArc.Start
					arc.End = endArc.End
					return
				}

				if t.ringEqualReverse(arc, endArc) {
					arc.Start = endArc.End
					arc.End = endArc.Start
					return
				}
			}
		}

		arcsByEnd[endPoint] = append(endArcs, arc)
		t.arcs = append(t.arcs, arc)
	}

	for _, line := range t.lines {
		for line != nil {
			dedupLine(line)
			line = line.Next
		}
	}

	for _, ring := range t.rings {
		if ring.Next != nil {
			// arc is no longer closed
			for ring != nil {
				dedupLine(ring)
				ring = ring.Next
			}
		} else {
			dedupRing(ring)
		}
	}

	t.lines = nil
	t.rings = nil
}

func (t *Topology) lineEqual(a, b *arc) bool {
	ia := a.Start
	ib := b.Start
	ja := a.End
	jb := b.End
	if ia-ja != ib-jb {
		return false
	}

	for ia <= ja {
		if !pointEquals(t.coordinates[ia], t.coordinates[ib]) {
			return false
		}
		ia += 1
		ib += 1
	}

	return true
}

func (t *Topology) lineEqualReverse(a, b *arc) bool {
	ia := a.Start
	ib := b.Start
	ja := a.End
	jb := b.End
	if ia-ja != ib-jb {
		return false
	}

	for ia <= ja {
		if !pointEquals(t.coordinates[ia], t.coordinates[jb]) {
			return false
		}
		ia += 1
		jb -= 1
	}

	return true
}

func (t *Topology) ringEqual(a, b *arc) bool {
	ia := a.Start
	ib := b.Start
	ja := a.End
	jb := b.End
	n := ja - ia
	if n != jb-ib {
		return false
	}

	ka := t.findMinimumOffset(a)
	kb := t.findMinimumOffset(b)

	for i := 0; i < n; i++ {
		pa := t.coordinates[ia+(i+ka)%n]
		pb := t.coordinates[ib+(i+kb)%n]
		if !pointEquals(pa, pb) {
			return false
		}
	}

	return true
}

func (t *Topology) ringEqualReverse(a, b *arc) bool {
	ia := a.Start
	ib := b.Start
	ja := a.End
	jb := b.End
	n := ja - ia
	if n != jb-ib {
		return false
	}

	ka := t.findMinimumOffset(a)
	kb := n - t.findMinimumOffset(b)

	for i := 0; i < n; i++ {
		pa := t.coordinates[ia+(i+ka)%n]
		pb := t.coordinates[jb-(i+kb)%n]
		if !pointEquals(pa, pb) {
			return false
		}
	}

	return true
}

// Rings are rotated to a consistent, but arbitrary, start point.
// This is necessary to detect when a ring and a rotated copy are dupes.
func (t *Topology) findMinimumOffset(arc *arc) int {
	start := arc.Start
	end := arc.End
	mid := start
	minimum := mid
	minimumPoint := t.coordinates[mid]

	mid += 1
	for mid < end {
		point := t.coordinates[mid]
		if point[0] < minimumPoint[0] || point[0] == minimumPoint[0] && point[1] < minimumPoint[1] {
			minimum = mid
			minimumPoint = point
		}
		mid += 1
	}

	return minimum - start
}
