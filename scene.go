package pathtracer

// Scene is a collection of geometry that can be rendered into color.
type Scene interface {
	Intersect(Ray) *Hit
}
