package pathtracer

// Hit describes the intersection of a Ray with Geometry.
type Hit struct {
	from     Ray
	scene    Scene
	geometry Geometry
	material Material
	point    Point
	normal   Ray
	uv       Point
}
