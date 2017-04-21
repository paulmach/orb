package planar

// Resample converts the path into totalPoints-1 evenly spaced segments.
func (ls LineString) Resample(totalPoints int) LineString {
	if totalPoints <= 0 {
		return LineString{}
	}

	ls, ret := ls.resampleEdgeCases(totalPoints)
	if ret {
		return ls
	}

	// precomputes the total distance and intermediate distances
	total, dists := precomputeDistances(ls)
	return ls.resample(dists, total, totalPoints)
}

// ResampleWithInterval coverts the path into evenly spaced points of
// about the given distance. The total distance is computed using euclidean
// geometry and then divided by the given distance to get the number of segments.
func (ls LineString) ResampleWithInterval(dist float64) LineString {
	if dist <= 0 {
		return LineString{}
	}

	// precomputes the total distance and intermediate distances
	total, dists := precomputeDistances(ls)

	totalPoints := int(total/dist) + 1
	ls, ret := ls.resampleEdgeCases(totalPoints)
	if ret {
		return ls
	}

	return ls.resample(dists, total, totalPoints)
}

func (ls LineString) resample(distances []float64, totalDistance float64, totalPoints int) LineString {
	points := make([]Point, 1, totalPoints)
	points[0] = ls[0] // start stays the same

	step := 1
	distance := 0.0

	currentDistance := totalDistance / float64(totalPoints-1)
	currentLine := Line{} // declare here and update has nice performance benefits
	for i := 0; i < len(ls)-1; i++ {
		currentLine.a = ls[i]
		currentLine.b = ls[i+1]

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
		points[totalPoints-1] = ls[len(ls)-1]
	}

	return LineString(points)
}

// resampleEdgeCases is used to handle edge case for
// resampling like not enough points and the path is all the same point.
// will return nil if there are no edge cases. If return true if
// one of these edge cases was found and handled.
func (ls LineString) resampleEdgeCases(totalPoints int) (LineString, bool) {
	// degenerate case
	if len(ls) <= 1 {
		return ls, true
	}

	// if all the points are the same, treat as special case.
	equal := true
	for _, point := range ls {
		if !ls[0].Equal(point) {
			equal = false
			break
		}
	}

	if equal {
		if totalPoints > len(ls) {
			// extend to be requested length
			for len(ls) != totalPoints {
				ls = append(ls, ls[0])
			}

			return ls, true
		}

		// contract to be requested length
		ls = ls[:totalPoints]
		return ls, true
	}

	return ls, false
}

// precomputeDistances precomputes the total distance and intermediate distances.
func precomputeDistances(ls LineString) (float64, []float64) {
	total := 0.0
	dists := make([]float64, len(ls)-1)
	for i := 0; i < len(ls)-1; i++ {
		dists[i] = ls[i].DistanceFrom(ls[i+1])
		total += dists[i]
	}

	return total, dists
}
