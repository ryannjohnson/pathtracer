package pathtracer

// Geometry represents anything that can intercept rays in 3D space.
type Geometry interface {
	Intersect(Ray) *Hit
}
