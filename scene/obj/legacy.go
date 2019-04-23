package obj

import (
	"io"
	"math"

	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

// NewLegacyScene loads from .obj and .mtl files to produce a scene.
func NewLegacyScene(objReader, mtlReader io.Reader) (*LegacyScene, error) {
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

	return &LegacyScene{decoder, triangles}, nil
}

// LegacyScene contains obj geometry and materials loaded from the g3n
// game engine library.
type LegacyScene struct {
	decoder   *obj.Decoder
	triangles []triangle
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
func (s *LegacyScene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, material pathtracer.Material, ok bool) {
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
