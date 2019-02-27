package pathtracer

// Ray describes a path of light with a beginning and a direction.
type Ray struct {
	Origin    Vector
	Direction Vector
}

// Transform applies a transformation matrix to each of its origin and
// its direction, keeping the direction as a unit vector.
func (r Ray) Transform(m Matrix) Ray {
	// Transforming the direction vector essentially rotates it when
	// it's being used as a unit vector.
	mDirection := m.ResetTranslation()

	return Ray{
		r.Origin.Transform(m),
		r.Direction.Transform(mDirection).Normalize(),
	}
}
