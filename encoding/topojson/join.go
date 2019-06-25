package topojson

type junctionMap map[point]bool

func (j junctionMap) Has(c []float64) bool {
	_, ok := j[point{c[0], c[1]}]
	return ok
}

func (t *Topology) join() junctionMap {
	n := len(t.coordinates)
	visited := make([]int, n)
	left := make([]int, n)
	right := make([]int, n)
	junctions := make([]bool, n)
	indexes := make([]int, n)

	junctionCount := 0

	indexByPoint := make(map[point]int)
	for i, c := range t.coordinates {
		p := point{c[0], c[1]}
		if v, ok := indexByPoint[p]; !ok {
			indexByPoint[p] = i
			indexes[i] = i
		} else {
			indexes[i] = v
		}
	}
	indexByPoint = nil

	sequence := func(i, previousIndex, currentIndex, nextIndex int) {
		if visited[currentIndex] == i {
			return // Ignore self-intersection
		}
		visited[currentIndex] = i
		leftIndex := left[currentIndex]
		if leftIndex >= 0 {
			rightIndex := right[currentIndex]
			if (leftIndex != previousIndex || rightIndex != nextIndex) &&
				(leftIndex != nextIndex || rightIndex != previousIndex) {
				junctionCount += 1
				junctions[currentIndex] = true
			}
		} else {
			left[currentIndex] = previousIndex
			right[currentIndex] = nextIndex
		}
	}

	for i := 0; i < n; i++ {
		visited[i] = -1
		left[i] = -1
		right[i] = -1
	}

	for i, l := range t.lines {
		start := l.Start
		end := l.End

		previous := 0
		current := indexes[start]
		start += 1
		next := indexes[start]

		junctionCount += 1
		junctions[current] = true

		start += 1
		for start <= end {
			previous = current
			current = next
			next = indexes[start]
			sequence(i, previous, current, next)
			start += 1
		}

		junctionCount += 1
		junctions[next] = true
	}

	for i := 0; i < n; i++ {
		visited[i] = -1
	}

	for i, r := range t.rings {
		start := r.Start + 1
		end := r.End

		previous := indexes[(end-1)%n]
		current := indexes[start-1]
		next := indexes[start%n]

		sequence(i, previous, current, next)

		start += 1
		for start <= end {
			previous = current
			current = next
			next = indexes[start]
			sequence(i, previous, current, next)
			start += 1
		}
	}

	result := junctionMap{}
	for k, v := range junctions {
		if v {
			c := t.coordinates[k]
			result[point{c[0], c[1]}] = true
		}
	}
	return result
}
