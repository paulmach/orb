package topojson

func (t *Topology) delta() {
	if t.opts.PostQuantize == 0 {
		return
	}

	for i, arc := range t.Arcs {
		x0 := arc[0][0]
		y0 := arc[0][1]
		for j, point := range arc {
			x1 := point[0]
			y1 := point[1]
			if j == 0 {
				t.Arcs[i][j] = []float64{x1, y1}
			} else {
				t.Arcs[i][j] = []float64{x1 - x0, y1 - y0}
			}
			x0 = x1
			y0 = y1
		}
	}
}
