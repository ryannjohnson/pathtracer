package pathtracer

import "math"

// Axes are often referred to when creating rotation in a single
// dimension.
var (
	AxisX = Vector{1, 0, 0}
	AxisY = Vector{0, 1, 0}
	AxisZ = Vector{0, 0, 1}
)

// Vector is a coordinate in 3D space.
type Vector struct {
	X, Y, Z float64
}

// Add results in the combination of two vectors.
func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

// Length describes the vector's distance from the origin.
func (v Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

// Normalize scales a vector into a unit vector, setting its length from
// [0, 0, 0] to 1.
func (v Vector) Normalize() Vector {
	length := v.Length()
	return Vector{v.X / length, v.Y / length, v.Z / length}
}

// Scale multiplies each dimension of a vector by some value.
func (v Vector) Scale(multiplier float64) Vector {
	return Vector{v.X * multiplier, v.Y * multiplier, v.Z * multiplier}
}

// Transform applies a transformation matrix to a vector, resulting in a
// new vector.
func (v Vector) Transform(m Matrix) Vector {
	return Vector{
		v.X*m.x00 + v.Y*m.x01 + v.Z*m.x02 + m.x03,
		v.X*m.x10 + v.Y*m.x11 + v.Z*m.x12 + m.x13,
		v.X*m.x20 + v.Y*m.x21 + v.Z*m.x22 + m.x23,
	}
}
