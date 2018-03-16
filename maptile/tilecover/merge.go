package tilecover

import "github.com/paulmach/orb/maptile"

// MergeUp will merge up the tiles in a given set up to the
// the give min zoom. The tiles in the input set are expected
// to all be of the same zoom, e.g. outputs of the Geometry function.
func MergeUp(set maptile.Set, min maptile.Zoom) maptile.Set {
	max := maptile.Zoom(1)
	for t, v := range set {
		if v {
			max = t.Z
			break
		}
	}

	if min == max {
		return set
	}

	merged := make(maptile.Set)
	for z := max; z > min; z-- {
		parentSet := make(maptile.Set)
		for t, v := range set {
			if !v {
				continue
			}

			if t.X%2 == 0 && t.Y%2 == 0 {
				sibs := t.Siblings()

				if set[sibs[1]] && set[sibs[2]] && set[sibs[3]] {
					set[sibs[0]] = false
					set[sibs[1]] = false
					set[sibs[2]] = false
					set[sibs[3]] = false

					parent := t.Parent()
					if z-1 == min {
						merged[parent] = true
					} else {
						parentSet[parent] = true
					}
				}
			}
		}

		for t, v := range set {
			if v {
				merged[t] = true
			}
		}

		set = parentSet
	}

	return merged
}
