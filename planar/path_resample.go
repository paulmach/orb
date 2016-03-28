package planar

// Resample converts the path into totalPoints-1 evenly spaced segments.
// Assumes euclidean geometry.
func (p Path) Resample(totalPoints int) Path {
	if totalPoints <= 0 {
		return Path(make([]Point, 0))
	}

	p, ret := p.resampleEdgeCases(totalPoints)
	if ret {
		return p
	}

	// precomputes the total distance and intermediate distances
	total, dists := precomputeDistances([]Point(p))
	return p.resample(dists, total, totalPoints)
}

// ResampleWithInterval coverts the path into evenly spaced points of
// about the given distance. The total distance is computed using euclidean
// geometry and then divided by the given distance to get the number of segments.
func (p Path) ResampleWithInterval(dist float64) Path {
	if dist <= 0 {
		return Path(make([]Point, 0))
	}

	// precomputes the total distance and intermediate distances
	total, dists := precomputeDistances([]Point(p))

	totalPoints := int(total/dist) + 1
	p, ret := p.resampleEdgeCases(totalPoints)
	if ret {
		return p
	}

	return p.resample(dists, total, totalPoints)
}

func (p Path) resample(distances []float64, totalDistance float64, totalPoints int) Path {
	points := make([]Point, 1, totalPoints)
	points[0] = p[0] // start stays the same

	step := 1
	distance := 0.0

	currentDistance := totalDistance / float64(totalPoints-1)
	currentLine := Line{} // declare here and update has nice performance benefits
	for i := 0; i < len(p)-1; i++ {
		currentLine.a = p[i]
		currentLine.b = p[i+1]

		currentLineDistance := distances[i]
		nextDistance := distance + currentLineDistance

		for currentDistance <= nextDistance {
			// need to add a point
			percent := (currentDistance - distance) / currentLineDistance
			points = append(points, Point{
				currentLine.a[0] + percent*(currentLine.b[0]-currentLine.a[0]),
				currentLine.a[1] + percent*(currentLine.b[1]-currentLine.a[1]),
			})

			// move to the next distance we want
			step++
			currentDistance = totalDistance * float64(step) / float64(totalPoints-1)
			if step == totalPoints-1 { // weird round off error on my machine
				currentDistance = totalDistance
			}
		}

		// past the current point in the original line, so move to the next one
		distance = nextDistance
	}

	// end stays the same, to handle round off errors
	if totalPoints != 1 { // for 1, we want the first point
		points[totalPoints-1] = p[len(p)-1]
	}

	return Path(points)
}

// resampleEdgeCases is used to handle edge case for
// resampling like not enough points and the path is all the same point.
// will return nil if there are no edge cases. If return true if
// one of these edge cases was found and handled.
func (p Path) resampleEdgeCases(totalPoints int) (Path, bool) {
	// degenerate case
	if len(p) <= 1 {
		return p, true
	}

	// if all the points are the same, treat as special case.
	equal := true
	for _, point := range p {
		if !p[0].Equal(point) {
			equal = false
			break
		}
	}

	if equal {
		if totalPoints > len(p) {
			// extend to be requested length
			for len(p) != totalPoints {
				p = append(p, p[0])
			}

			return p, true
		}

		// contract to be requested length
		p = p[:totalPoints]
		return p, true
	}

	return p, false
}

// precomputeDistances precomputes the total distance and intermediate distances.
func precomputeDistances(p []Point) (float64, []float64) {
	total := 0.0
	dists := make([]float64, len(p)-1)
	for i := 0; i < len(p)-1; i++ {
		dists[i] = p[i].DistanceFrom(p[i+1])
		total += dists[i]
	}

	return total, dists
}
