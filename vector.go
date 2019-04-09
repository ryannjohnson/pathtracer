package pathtracer

import "math"

// Axes are often referred to when creating rotation in a single
// dimension.
var (
	AxisX = Vector{1, 0, 0}
	AxisY = Vector{0, 1, 0}
	AxisZ = Vector{0, 0, 1}
)

// NewVector creates a new vector with three dimensions.
func NewVector(x, y, z float64) Vector {
	return Vector{x, y, z}
}

// Vector is a coordinate in 3D space.
type Vector struct {
	X, Y, Z float64
}

// Add results in the combination of two vectors.
func (v Vector) Add(w Vector) Vector {
	return NewVector(v.X+w.X, v.Y+w.Y, v.Z+w.Z)
}

// DotProduct measures the magnitude of similarity between two vectors.
//
// If they're parallel, then a vector representing the product of their
// lengths is returned.
//
// If they're perpendicular, then zero is returned.
//
// https://en.wikipedia.org/wiki/Dot_product
func (v Vector) DotProduct(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

// CrossProduct measures the similarity between two vectors.
//
// If they are parallel, the cross product is zero because there is no
// difference between them.
//
// If they are perpendicular, the cross product will equal a vector with
// a length equal to the area of the triangle created by the two
// vectors.
//
// Eg, if v1 = (1, 0, 0) and v2 = (0, 1, 0), then the cross product will
// be (0, 0, 1).
//
// https://en.wikipedia.org/wiki/Cross_product
func (v Vector) CrossProduct(w Vector) Vector {
	return NewVector(
		v.Y*w.Z-v.Z*w.Y,
		v.Z*w.X-v.X*w.Z,
		v.X*w.Y-v.Y*w.X,
	)
}

// Length describes the vector's distance from the origin.
func (v Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

// Normalize scales a vector into a unit vector, setting its length from
// [0, 0, 0] to 1.
func (v Vector) Normalize() Vector {
	length := v.Length()
	return NewVector(v.X/length, v.Y/length, v.Z/length)
}

// Scale multiplies each dimension of a vector by some value.
func (v Vector) Scale(multiplier float64) Vector {
	return NewVector(v.X*multiplier, v.Y*multiplier, v.Z*multiplier)
}

// Subtract results in the difference of two vectors.
func (v Vector) Subtract(w Vector) Vector {
	return NewVector(v.X-w.X, v.Y-w.Y, v.Z-w.Z)
}

// Transform applies a transformation matrix to a vector, resulting in a
// new vector.
func (v Vector) Transform(m Matrix) Vector {
	return NewVector(
		v.X*m.x00+v.Y*m.x01+v.Z*m.x02+m.x03,
		v.X*m.x10+v.Y*m.x11+v.Z*m.x12+m.x13,
		v.X*m.x20+v.Y*m.x21+v.Z*m.x22+m.x23,
	)
}
