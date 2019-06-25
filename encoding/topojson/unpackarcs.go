package topojson

func (t *Topology) unpackArcs() {
	t.arcIndexes = make(map[arcEntry]int)

	// Unpack arcs
	for i, a := range t.arcs {
		t.arcIndexes[arcEntry{a.Start, a.End}] = i
		t.Arcs = append(t.Arcs, t.coordinates[a.Start:a.End+1])
	}
	t.arcs = nil
	t.coordinates = nil
}
