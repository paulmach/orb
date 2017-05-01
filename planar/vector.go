package planar

import "math"

// Vector represents a magnitude and a direction
type Vector [2]float64

// NewVector initializes a new vector.
func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

// Negative returns the negative of the vector.
// eg. Vector{0, 0}.Sub(v)
func (v Vector) Negative() Vector {
	return Vector{0, 0}.Sub(v)
}

// Add adds to two vectors.
func (v Vector) Add(vector Vector) Vector {
	v[0] += vector[0]
	v[1] += vector[1]

	return v
}

// Sub subtracts a point from the given point.
func (v Vector) Sub(vector Vector) Vector {
	v[0] -= vector[0]
	v[1] -= vector[1]

	return v
}

// Normalize treats the point as a vector and
// scales it such that its distance from [0,0] is 1.
func (v Vector) Normalize() Vector {
	dist := math.Sqrt(v[0]*v[0] + v[1]*v[1])

	if dist == 0 {
		v[0] = 0
		v[1] = 0

		return v
	}

	v[0] /= dist
	v[1] /= dist

	return v
}

// Scale each component of the point.
func (v Vector) Scale(factor float64) Vector {
	v[0] *= factor
	v[1] *= factor

	return v
}

// Magnitude return the length/magnitude of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1])
}

// Dot is just x1*x2 + y1*y2
func (v Vector) Dot(vector Vector) float64 {
	return v[0]*vector[0] + v[1]*vector[1]
}

// Equal checks if the point represents the same point or vector.
func (v Vector) Equal(vector Vector) bool {
	return v[0] == vector[0] && v[1] == vector[1]
}

// X returns the x/horizontal component of the vector.
func (v Vector) X() float64 {
	return v[0]
}

// Y returns the y/vertical component of the vector.
func (v Vector) Y() float64 {
	return v[1]
}
