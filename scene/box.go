package scene

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

// Box is a 3D construct aligned to the 3 axes in the scene. It cannot
// be rotated away from this orientation.
type Box struct {
	max pathtracer.Vector
	min pathtracer.Vector
}

// NewBox creates a box according to any two vectors as its furthest
// corners.
func NewBox(a, b pathtracer.Vector) Box {
	return Box{
		max: pathtracer.NewVector(
			math.Max(a.X, b.X),
			math.Max(a.Y, b.Y),
			math.Max(a.Z, b.Z),
		),
		min: pathtracer.NewVector(
			math.Min(a.X, b.X),
			math.Min(a.Y, b.Y),
			math.Min(a.Z, b.Z),
		),
	}
}

// IntersectsRay determines whether a ray intersects a box at all. This
// includes if the ray's origin is inside the box.
//
// https://gamedev.stackexchange.com/a/18459
func (b Box) IntersectsRay(ray pathtracer.Ray) (tmin, tmax float64, ok bool) {
	// Each of these six variables below represent how far along the ray
	// an intersection occurs with each plane of the box.
	//
	// For example, tx0 represents how far along the ray you'd have to
	// travel until you ran into the X value of the box's "min" corner.
	//
	// If the ray starts from within the box, then t_0 will be negative
	// and t_1 will be positive. If the ray never intersects, then all
	// values will be negative.
	tx0 := (b.min.X - ray.Origin.X) / ray.Direction.X
	tx1 := (b.max.X - ray.Origin.X) / ray.Direction.X
	ty0 := (b.min.Y - ray.Origin.Y) / ray.Direction.Y
	ty1 := (b.max.Y - ray.Origin.Y) / ray.Direction.Y
	tz0 := (b.min.Z - ray.Origin.Z) / ray.Direction.Z
	tz1 := (b.max.Z - ray.Origin.Z) / ray.Direction.Z

	// The direction of the ray may have flipped the sign of any of the
	// axes. The inner `math.Min` and `math.Max` calls correct this,
	// preserving the magnitude of t_0 and t_1, respectively.
	//
	// For tmax, we find the closest intersection distance of any of the
	// boxes farthest 3 planes.
	//
	// For tmin, we find the farthest intersection distance of any of
	// the closest 3 planes.
	tmin = math.Max(math.Max(math.Min(tx0, tx1), math.Min(ty0, ty1)), math.Min(tz0, tz1))
	tmax = math.Min(math.Min(math.Max(tx0, tx1), math.Max(ty0, ty1)), math.Max(tz0, tz1))

	// If the farthest intersection is negative, then the box is behind
	// the ray.
	if tmax < 0 {
		return
	}

	ok = tmin <= tmax
	return
}
