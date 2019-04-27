package pathtracer

// Scene is a collection of geometry.
type Scene interface {
	Intersect(ray Ray) (hit Hit, material Material, ok bool)
	Clone() Scene
}
