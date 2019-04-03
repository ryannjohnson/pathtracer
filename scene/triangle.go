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
// http://geomalgorithms.com/a06-_intersect-2.html
func IntersectTriangle(ray pathtracer.Ray, triangle Triangle) (distanceAlongRayFromOrigin float64, intersectionPoint, triangleNormal pathtracer.Vector, ok bool) {
	// Edges can be represented by vectors.
	v0v1 := triangle.Vertex0().Sub(triangle.Vertex1())
	v0v2 := triangle.Vertex0().Sub(triangle.Vertex2())

	// Unit vector representing the normal of the triangle.
	triangleNormal = v0v1.CrossProduct(v0v2).Normalize()

	// The ray direction and triangle normal are both unit vectors.
	//
	// If both vectors are parallel, then their dot product will equal
	// 1. If they are perpendicular, their dot product will equal 0. The
	// values in between are essentially the raw cosine of the angle
	// between the two vectors.
	cosineOfRayAndNormal := ray.Direction.DotProduct(triangleNormal)

	// If they're perpendicular and equal zero, then the ray will never
	// intersect the triangle's plane.
	if math.Abs(cosineOfRayAndNormal) < math.SmallestNonzeroFloat64 {
		return
	}

	// To determine how far the ray origin is from the closest point on
	// the triangle's plane, we'll use a known length: the distance from
	// the ray origin to any vertex on the triangle.
	//
	// Since the triangleNormal is a unit vector, its only effect will
	// be scaling the magnitude of the ray origin/triangle vertex line
	// depending on how different their angles are.
	//
	// If they're parallel, this line represents the full distance. If
	// they're at a 60 degree angle away from each other, then the
	// distance will be half of the line's length (cosine of 30
	// degrees).
	v0r0 := triangle.Vertex0().Sub(ray.Origin)
	closestDistanceFromRayOriginToPlane := triangleNormal.DotProduct(v0r0)

	// The closest distance from the ray origin to the triangle plane is
	// the adjacent side of the triangle. To get the hypotenuse, we just
	// divide it by the cosine of the angle.
	//
	// The wider the angle, the longer the hypotenuse is. If the angle
	// is tiny or nonexistant, then our adjacent side is already the
	// hypotenuse, as well.
	distanceAlongRayFromOrigin = closestDistanceFromRayOriginToPlane / cosineOfRayAndNormal

	// We don't want to intersect triangles that are at or behind the
	// ray that's looking for them.
	if distanceAlongRayFromOrigin <= 0 {
		return
	}

	// Point on the triangle's plane, not necessarily inside the
	// triangle yet.
	intersectionPoint = ray.Origin.Add(ray.Direction.Scale(distanceAlongRayFromOrigin))

	var triangleEdge pathtracer.Vector
	var pointEdge pathtracer.Vector
	var edgesCrossProduct pathtracer.Vector

	// Vertex0
	triangleEdge = triangle.Vertex1().Sub(triangle.Vertex0())
	pointEdge = intersectionPoint.Sub(triangle.Vertex0())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if triangleNormal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	// Vertex1
	triangleEdge = triangle.Vertex2().Sub(triangle.Vertex1())
	pointEdge = intersectionPoint.Sub(triangle.Vertex1())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if triangleNormal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	// Vertex2
	triangleEdge = triangle.Vertex0().Sub(triangle.Vertex2())
	pointEdge = intersectionPoint.Sub(triangle.Vertex2())
	edgesCrossProduct = triangleEdge.CrossProduct(pointEdge)
	if triangleNormal.DotProduct(edgesCrossProduct) < 0 {
		return
	}

	ok = true
	return
}
