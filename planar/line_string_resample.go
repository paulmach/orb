package planar

// Resample converts the line string into totalPoints-1 evenly spaced segments.
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

// ResampleWithInterval coverts the line string into evenly spaced points of
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

func (ls LineString) resample(dists []float64, totalDistance float64, totalPoints int) LineString {
	points := make([]Point, 1, totalPoints)
	points[0] = ls[0] // start stays the same

	step := 1
	dist := 0.0

	currentDistance := totalDistance / float64(totalPoints-1)
	// declare here and update had nice performance benefits need to retest
	currentSeg := Segment{}
	for i := 0; i < len(ls)-1; i++ {
		currentSeg[0] = ls[i]
		currentSeg[1] = ls[i+1]

		currentSegDistance := dists[i]
		nextDistance := dist + currentSegDistance

		for currentDistance <= nextDistance {
			// need to add a point
			percent := (currentDistance - dist) / currentSegDistance
			points = append(points, Point{
				currentSeg[0][0] + percent*(currentSeg[1][0]-currentSeg[0][0]),
				currentSeg[0][1] + percent*(currentSeg[1][1]-currentSeg[0][1]),
			})

			// move to the next distance we want
			step++
			currentDistance = totalDistance * float64(step) / float64(totalPoints-1)
			if step == totalPoints-1 { // weird round off error on my machine
				currentDistance = totalDistance
			}
		}

		// past the current point in the original segment, so move to the next one
		dist = nextDistance
	}

	// end stays the same, to handle round off errors
	if totalPoints != 1 { // for 1, we want the first point
		points[totalPoints-1] = ls[len(ls)-1]
	}

	return LineString(points)
}

// resampleEdgeCases is used to handle edge case for
// resampling like not enough points and the line string is all the same point.
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
