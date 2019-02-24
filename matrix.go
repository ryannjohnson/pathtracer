package pathtracer

import "math"

// NewMatrix creates a new identity matrix, which would not produce any
// changes when applied to a vector.
func NewMatrix() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Matrix represents transformations that can be applied to Vectors.
//
// https://github.com/fogleman/pt/blob/master/pt/matrix.go
type Matrix struct {
	x00, x01, x02, x03 float64
	x10, x11, x12, x13 float64
	x20, x21, x22, x23 float64
	x30, x31, x32, x33 float64
}

// Multiply produces a matrix that represents the transformations of
// both matrices. Order matters.
//
// https://github.com/fogleman/pt/blob/master/pt/matrix.go
// https://www.mathsisfun.com/algebra/matrix-multiplying.html
func (m Matrix) Multiply(b Matrix) Matrix {
	output := Matrix{}
	output.x00 = m.x00*b.x00 + m.x01*b.x10 + m.x02*b.x20 + m.x03*b.x30
	output.x10 = m.x10*b.x00 + m.x11*b.x10 + m.x12*b.x20 + m.x13*b.x30
	output.x20 = m.x20*b.x00 + m.x21*b.x10 + m.x22*b.x20 + m.x23*b.x30
	output.x30 = m.x30*b.x00 + m.x31*b.x10 + m.x32*b.x20 + m.x33*b.x30
	output.x01 = m.x00*b.x01 + m.x01*b.x11 + m.x02*b.x21 + m.x03*b.x31
	output.x11 = m.x10*b.x01 + m.x11*b.x11 + m.x12*b.x21 + m.x13*b.x31
	output.x21 = m.x20*b.x01 + m.x21*b.x11 + m.x22*b.x21 + m.x23*b.x31
	output.x31 = m.x30*b.x01 + m.x31*b.x11 + m.x32*b.x21 + m.x33*b.x31
	output.x02 = m.x00*b.x02 + m.x01*b.x12 + m.x02*b.x22 + m.x03*b.x32
	output.x12 = m.x10*b.x02 + m.x11*b.x12 + m.x12*b.x22 + m.x13*b.x32
	output.x22 = m.x20*b.x02 + m.x21*b.x12 + m.x22*b.x22 + m.x23*b.x32
	output.x32 = m.x30*b.x02 + m.x31*b.x12 + m.x32*b.x22 + m.x33*b.x32
	output.x03 = m.x00*b.x03 + m.x01*b.x13 + m.x02*b.x23 + m.x03*b.x33
	output.x13 = m.x10*b.x03 + m.x11*b.x13 + m.x12*b.x23 + m.x13*b.x33
	output.x23 = m.x20*b.x03 + m.x21*b.x13 + m.x22*b.x23 + m.x23*b.x33
	output.x33 = m.x30*b.x03 + m.x31*b.x13 + m.x32*b.x23 + m.x33*b.x33
	return output
}

// Rotate creates a transformation matrix that rotates vectors around an
// arbitrary axis. For example, to rotate around the X axis, the input
// vector would be [1, 0, 0].
//
// https://github.com/fogleman/pt/blob/master/pt/matrix.go
// http://ksuweb.kennesaw.edu/~plaval/math4490/rotgen.pdf
func (m Matrix) Rotate(v Vector, radians float64) Matrix {
	v = v.Normalize()
	s := math.Sin(radians)
	c := math.Cos(radians)
	t := 1 - c
	r := Matrix{
		t*v.X*v.X + c, t*v.X*v.Y + v.Z*s, t*v.Z*v.X - v.Y*s, 0,
		t*v.X*v.Y - v.Z*s, t*v.Y*v.Y + c, t*v.Y*v.Z + v.X*s, 0,
		t*v.Z*v.X + v.Y*s, t*v.Y*v.Z - v.X*s, t*v.Z*v.Z + c, 0,
		0, 0, 0, 1}
	return m.Multiply(r)
}
