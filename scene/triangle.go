package scene

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

// Triangle is a face with three points and three edges.
type Triangle interface {
	Vertex0() pathtracer.Vector
	Vertex1() pathtracer.Vector
	Vertex2() pathtracer.Vector
}

// IntersectTriangle determins if a ray passes through a triangle and at
// what distance from the origin if so.
//
// https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-rendering-a-triangle/ray-triangle-intersection-geometric-solution
func IntersectTriangle(ray pathtracer.Ray, triangle Triangle) (planeDistanceFromRayOrigin float64, intersectionPoint, intersectionNormal pathtracer.Vector, ok bool) {
	// Edges can be represented by vectors.
	v0v1 := triangle.Vertex0().Sub(triangle.Vertex1())
	v0v2 := triangle.Vertex0().Sub(triangle.Vertex2())
	normal := v0v1.CrossProduct(v0v2).Normalize()

	// Parallel triangles and rays do not intersect. We measure this by
	// seeing if the triangle's normal is perpendicular to the ray's
	// direction.
	rayDirectionDotNormal := ray.Direction.DotProduct(normal)
	if math.Abs(rayDirectionDotNormal) < math.SmallestNonzeroFloat64 {
		return
	}

	// http://geomalgorithms.com/a06-_intersect-2.html
	//
	// TODO: Rename the planeDistanceFromRayOrigin and try to collapse
	// these calculations into a simpler form.
	v0p0 := triangle.Vertex0().Sub(ray.Origin)
	planeDistanceFromRayOrigin = normal.DotProduct(v0p0) / normal.DotProduct(ray.Direction)
	if planeDistanceFromRayOrigin <= 0 {
		return
	}

	intersectionPoint = ray.Origin.Add(ray.Direction.Scale(planeDistanceFromRayOrigin))
	intersectionNormal = normal.Normalize()

	var triangleEdge pathtracer.Vector
	var pointEdge pathtracer.Vector
	var edgesCrossProduct pathtracer.Vector

	// Vertex0
	triangleEdge = triangle.Vertex1().Sub(triangle.Vertex0())
	pointEdge = intersectionPoint.Sub(triangle.Vertex0())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if normal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	// Vertex1
	triangleEdge = triangle.Vertex2().Sub(triangle.Vertex1())
	pointEdge = intersectionPoint.Sub(triangle.Vertex1())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if normal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	// Vertex2
	triangleEdge = triangle.Vertex0().Sub(triangle.Vertex2())
	pointEdge = intersectionPoint.Sub(triangle.Vertex2())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if normal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	ok = true
	return
}
