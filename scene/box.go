package scene

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

// Box is a 3D construct aligned to the 3 axes in the scene. It cannot
// be rotated away from this orientation.
type Box struct {
	min pathtracer.Vector
	max pathtracer.Vector
}

// NewBox creates a box according to any two vectors as its furthest
// corners.
func NewBox(min, max pathtracer.Vector) Box {
	return Box{min, max}
}

// Min returns the corner of the box corresponding to the negative
// direction of each axis.
func (b Box) Min() pathtracer.Vector { return b.min }

// Max returns the corner of the box corresponding to the positive
// direction of each axis.
func (b Box) Max() pathtracer.Vector { return b.max }

// Vertexes returns each of the box's eight vertexes.
func (b Box) Vertexes() []pathtracer.Vector {
	return []pathtracer.Vector{
		pathtracer.NewVector(b.min.X, b.min.Y, b.min.Z),
		pathtracer.NewVector(b.max.X, b.min.Y, b.min.Z),
		pathtracer.NewVector(b.min.X, b.max.Y, b.min.Z),
		pathtracer.NewVector(b.min.X, b.min.Y, b.max.Z),
		pathtracer.NewVector(b.min.X, b.max.Y, b.max.Z),
		pathtracer.NewVector(b.max.X, b.min.Y, b.max.Z),
		pathtracer.NewVector(b.max.X, b.max.Y, b.min.Z),
		pathtracer.NewVector(b.max.X, b.max.Y, b.max.Z),
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

// IntersectsTriangle returns true if there is no separating axis
// between the box and the triangle.
//
// https://stackoverflow.com/questions/17458562/efficient-aabb-triangle-intersection-in-c-sharp
// http://fileadmin.cs.lth.se/cs/Personal/Tomas_Akenine-Moller/pubs/tribox.pdf
func (b Box) IntersectsTriangle(triangle Triangle) bool {
	// Try to find separating planes along the X, Y, and Z axes.
	boxEdges := []pathtracer.Vector{pathtracer.AxisX, pathtracer.AxisY, pathtracer.AxisZ}
	triangleVertexes := []pathtracer.Vector{triangle.Vertex0(), triangle.Vertex1(), triangle.Vertex2()}
	triangleMin, triangleMax := projectDistanceAlongAxis(triangleVertexes, pathtracer.AxisX)
	if triangleMin > b.max.X || triangleMax < b.min.X {
		return false
	}
	triangleMin, triangleMax = projectDistanceAlongAxis(triangleVertexes, pathtracer.AxisY)
	if triangleMin > b.max.Y || triangleMax < b.min.Y {
		return false
	}
	triangleMin, triangleMax = projectDistanceAlongAxis(triangleVertexes, pathtracer.AxisZ)
	if triangleMin > b.max.Z || triangleMax < b.min.Z {
		return false
	}

	// Try to see if the triangle is "parallel" to the box in the way
	// that its separating plane is parallel to the triangle.
	v0v1 := triangle.Vertex0().Subtract(triangle.Vertex1())
	v0v2 := triangle.Vertex0().Subtract(triangle.Vertex2())
	triangleNormal := v0v1.CrossProduct(v0v2)
	triangleDistanceFromOrigin := triangleNormal.DotProduct(triangle.Vertex0())
	boxVertexes := b.Vertexes()

	boxMin, boxMax := projectDistanceAlongAxis(boxVertexes, triangleNormal)
	if boxMin > triangleDistanceFromOrigin || boxMax < triangleDistanceFromOrigin {
		return false
	}

	// When the triangle is parallel-ish to the box _and_ strattles the
	// boxes axes, then we find gaps by projecting the shapes on planes
	// parallel to both shapes.
	//
	// The cross product finds the axis perpendicular to the input
	// vectors.
	v1v2 := triangle.Vertex1().Subtract(triangle.Vertex2())
	triangleEdges := []pathtracer.Vector{v0v1, v0v2, v1v2}
	for _, triangleEdge := range triangleEdges {
		for _, boxEdge := range boxEdges {
			axis := triangleEdge.CrossProduct(boxEdge)
			triangleMin, triangleMax := projectDistanceAlongAxis(triangleVertexes, axis)
			boxMin, boxMax = projectDistanceAlongAxis(boxVertexes, axis)

			if boxMin > triangleMax || boxMax < triangleMin {
				return false
			}
		}
	}

	return true
}

// projectDistanceAlongAxis determines the closest and farthest
// distances along an arbitrary axis that any of the included vectors
// reach. Distances are all relative to the origin.
func projectDistanceAlongAxis(vertexes []pathtracer.Vector, axis pathtracer.Vector) (tmin, tmax float64) {
	tmin = math.MaxFloat64
	tmax = math.MaxFloat64 * -1

	for _, vertex := range vertexes {
		distance := vertex.DotProduct(axis)

		if tmin > distance {
			tmin = distance
		}
		if tmax < distance {
			tmax = distance
		}
	}

	return
}
