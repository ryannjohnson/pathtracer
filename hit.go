package pathtracer

// UVCoordinate describes what point on a 2D image should be returned.
// It uses a 1x1 scale so that the underlying texture can be any
// resolution.
type UVCoordinate interface {
	U() float64
	V() float64
}

// Hit describes the intersection of a Ray with Geometry.
type Hit struct {
	From     Ray
	Position Vector
	Normal   Ray
	UV       UVCoordinate
	Material Material
}
