package g3n

import (
	"io"
	"math"

	objLoader "github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

// NewOBJScene loads from .obj and .mtl files to produce a scene.
func NewOBJScene(objReader, mtlReader io.Reader) (*OBJScene, error) {
	decoder, err := objLoader.DecodeReader(objReader, mtlReader)
	if err != nil {
		return nil, err
	}

	triangles := make([]objTriangle, 0)

	for _, object := range decoder.Objects {
		geom, err := decoder.NewGeometry(&object)
		if err != nil {
			return nil, err
		}

		geom.ReadFaces(func(vA, vB, vC math32.Vector3) bool {
			v0 := pathtracer.NewVector(float64(vA.X), float64(vA.Y), float64(vA.Z))
			v1 := pathtracer.NewVector(float64(vB.X), float64(vB.Y), float64(vB.Z))
			v2 := pathtracer.NewVector(float64(vC.X), float64(vC.Y), float64(vC.Z))
			triangles = append(triangles, objTriangle{v0, v1, v2})
			return false
		})
	}

	return &OBJScene{decoder, triangles}, nil
}

// OBJScene contains obj geometry and materials loaded from the g3n game
// engine library.
type OBJScene struct {
	decoder   *objLoader.Decoder
	triangles []objTriangle
}

type objTriangle [3]pathtracer.Vector

func (t objTriangle) Vertex0() pathtracer.Vector {
	return t[0]
}

func (t objTriangle) Vertex1() pathtracer.Vector {
	return t[1]
}

func (t objTriangle) Vertex2() pathtracer.Vector {
	return t[2]
}

type dummyMaterial struct{}

func (m dummyMaterial) Sample(hit pathtracer.Hit, nextSample pathtracer.Sampler) pathtracer.Color {
	return pathtracer.NewColor(1, 1, 1)
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
//
// TODO: Optimize this method for performance.
// TODO: Figure a way to return a material with the hit.
func (s *OBJScene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, material pathtracer.Material, ok bool) {
	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64

	for _, triangle := range s.triangles {
		distance, point, normal, didIntersect := scene.IntersectTriangle(ray, triangle)
		if !didIntersect {
			continue
		}

		if closestDistance < distance {
			continue
		}

		closestPoint = point
		closestNormal = normal
		closestDistance = distance
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
	material = dummyMaterial{}
	return
}
