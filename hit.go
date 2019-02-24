package pathtracer

// UVCoordinate describes what point on a 2D image should be returned.
// It uses a 1x1 scale so that the underlying texture can be any
// resolution.
type UVCoordinate struct {
	U, V float64
}

// Hit describes the intersection of a Ray with Geometry.
type Hit struct {
	from     Ray
	scene    Scene
	geometry Geometry
	material Material
	point    Vector
	normal   Ray
	uv       UVCoordinate
}
