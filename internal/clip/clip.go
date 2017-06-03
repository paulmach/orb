package clip

// Code based on https://github.com/mapbox/lineclip

// LineString is something that behaves like a line and
// can be clipped.
type LineString interface {
	Len() int
	Get(int) (x, y float64)
	Append(x, y float64)
	Clear()
}

// MultiLineString is the interface for the output of a clipped line.
type MultiLineString interface {
	Append(i int, x, y float64)
}

// Bound is shared type to represent a bound.
type Bound struct {
	Left, Right, Bottom, Top float64
}

// Line will clip a line into a set of lines
// along the bounding box boundary.
func Line(box Bound, in LineString, out MultiLineString) {
	line := 0

	codeA := box.bitCode(in.Get(0))
	loopTo := in.Len()
	for i := 1; i < loopTo; i++ {
		ax, ay := in.Get(i - 1)
		bx, by := in.Get(i)

		codeB := box.bitCode(bx, by)
		endCode := codeB

		// loops through all the intersection of the line and box.
		// eg. across a corner could have two intersections.
		for {
			if codeA|codeB == 0 {
				// both points are in the box, accept
				out.Append(line, ax, ay)
				if codeB != endCode { // segment went outside
					out.Append(line, bx, by)
					if i < loopTo-1 {
						line++
					}
				} else if i == loopTo-1 {
					out.Append(line, bx, by)
				}
				break
			} else if codeA&codeB != 0 {
				// both on one side of the box.
				// segment not part of the final result.
				break
			} else if codeA != 0 {
				// A is outside, B is inside, clip edge
				ax, ay = box.intersect(codeA, ax, ay, bx, by)
				codeA = box.bitCode(ax, ay)
			} else {
				// B is outside, A is inside, clip edge
				bx, by = box.intersect(codeB, ax, ay, bx, by)
				codeB = box.bitCode(bx, by)
			}
		}

		codeA = endCode // new start is the old end
	}
}

// Ring will clip the Ring into a smaller ring around the bounding box boundary.
func Ring(box Bound, in LineString, out LineString) {
	if in.Len() == 0 {
		return
	}

	cache := out
	out = &lineString{}

	for edge := 1; edge <= 8; edge *= 2 {
		out.Clear()

		loopTo := in.Len()
		prevX, prevY := in.Get(loopTo - 1)
		prevInside := box.bitCode(prevX, prevY)&edge == 0

		for i := 0; i < loopTo; i++ {
			px, py := in.Get(i)
			inside := box.bitCode(px, py)&edge == 0

			// if segment goes through the clip window, add an intersection
			if inside != prevInside {
				ix, iy := box.intersect(edge, prevX, prevY, px, py)
				out.Append(ix, iy)
			}
			if inside {
				out.Append(px, py)
			}

			prevX = px
			prevY = py
			prevInside = inside
		}

		if out.Len() == 0 {
			return
		}

		in = out
		cache, out = out, cache
	}
}

// bit code reflects the point position relative to the bbox:

//         left  mid  right
//    top  1001  1000  1010
//    mid  0001  0000  0010
// bottom  0101  0100  0110

func (b Bound) bitCode(x, y float64) int {
	code := 0
	if x < b.Left {
		code |= 1
	} else if x > b.Right {
		code |= 2
	}

	if y < b.Bottom {
		code |= 4
	} else if y > b.Top {
		code |= 8
	}

	return code
}

// intersect a segment against one of the 4 lines that make up the bbox

func (b Bound) intersect(edge int, ax, ay, bx, by float64) (x, y float64) {
	if edge&8 != 0 {
		// top
		return ax + (bx-ax)*(b.Top-ay)/(by-ay), b.Top
	} else if edge&4 != 0 {
		// bottom
		return ax + (bx-ax)*(b.Bottom-ay)/(by-ay), b.Bottom
	} else if edge&2 != 0 {
		// right
		return b.Right, ay + (by-ay)*(b.Right-ax)/(bx-ax)
	} else if edge&1 != 0 {
		// left
		return b.Left, ay + (by-ay)*(b.Left-ax)/(bx-ax)
	}

	panic("no edge??")
}

type lineString [][2]float64

func (ls *lineString) Len() int {
	return len(*ls)
}

func (ls *lineString) Get(i int) (x, y float64) {
	return (*ls)[i][0], (*ls)[i][1]
}

func (ls *lineString) Append(x, y float64) {
	*ls = append(*ls, [2]float64{x, y})
}

func (ls *lineString) Clear() {
	*ls = (*ls)[:0]
}
