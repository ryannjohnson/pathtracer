package obj

import (
	"io"
	"math"

	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

// NewScene loads from .obj and .mtl files to produce a scene.
func NewScene(objReader, mtlReader io.Reader) (*Scene, error) {
	decoder, err := obj.DecodeReader(objReader, mtlReader)
	if err != nil {
		return nil, err
	}

	triangles := make([]triangle, 0)

	for _, object := range decoder.Objects {
		for _, face := range object.Faces {
			readTriangles(decoder, face, func(faceTriangle triangle) {
				triangles = append(triangles, faceTriangle)
			})
		}
	}

	return &Scene{decoder, triangles}, nil
}

// Scene contains obj geometry and materials loaded from the g3n game
// engine library.
type Scene struct {
	decoder   *obj.Decoder
	triangles []triangle
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
//
// TODO: Optimize this method for performance.
func (s *Scene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, material pathtracer.Material, ok bool) {
	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64
	var closestTriangle triangle

	for _, faceTriangle := range s.triangles {
		distance, point, normal, didIntersect := scene.IntersectTriangle(ray, faceTriangle)
		if !didIntersect {
			continue
		}

		if closestDistance < distance {
			continue
		}

		closestPoint = point
		closestNormal = normal
		closestDistance = distance
		closestTriangle = faceTriangle
	}

	if closestDistance == math.MaxFloat64 {
		return
	}

	ok = true
	hit = pathtracer.Hit{
		From:     ray,
		Position: closestPoint,
		Normal:   closestNormal,
	}
	material = closestTriangle.material
	return
}
