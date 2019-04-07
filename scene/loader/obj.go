package g3n

import (
	"io"
	"math"

	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/material"
	"github.com/ryannjohnson/pathtracer/scene"
)

// NewOBJScene loads from .obj and .mtl files to produce a scene.
func NewOBJScene(objReader, mtlReader io.Reader) (*OBJScene, error) {
	decoder, err := obj.DecodeReader(objReader, mtlReader)
	if err != nil {
		return nil, err
	}

	triangles := make([]objTriangle, 0)

	for _, object := range decoder.Objects {
		for _, face := range object.Faces {
			readTriangles(decoder, face, func(faceTriangle objTriangle) {
				triangles = append(triangles, faceTriangle)
			})
		}
	}

	return &OBJScene{decoder, triangles}, nil
}

func readTriangles(decoder *obj.Decoder, face obj.Face, callback func(objTriangle)) {
	numTriangles := len(face.Vertices) - 2
	for i := 0; i < numTriangles; i++ {
		vertexes := [3]pathtracer.Vector{
			pathtracer.NewVector(
				float64(decoder.Vertices[face.Vertices[0]*3]),
				float64(decoder.Vertices[face.Vertices[0]*3+1]),
				float64(decoder.Vertices[face.Vertices[0]*3+2]),
			),
			pathtracer.NewVector(
				float64(decoder.Vertices[face.Vertices[i+1]*3]),
				float64(decoder.Vertices[face.Vertices[i+1]*3+1]),
				float64(decoder.Vertices[face.Vertices[i+1]*3+2]),
			),
			pathtracer.NewVector(
				float64(decoder.Vertices[face.Vertices[i+2]*3]),
				float64(decoder.Vertices[face.Vertices[i+2]*3+1]),
				float64(decoder.Vertices[face.Vertices[i+2]*3+2]),
			),
		}

		normals := [3]pathtracer.Vector{
			pathtracer.NewVector(
				float64(decoder.Normals[face.Normals[0]*3]),
				float64(decoder.Normals[face.Normals[0]*3+1]),
				float64(decoder.Normals[face.Normals[0]*3+2]),
			),
			pathtracer.NewVector(
				float64(decoder.Normals[face.Normals[i+1]*3]),
				float64(decoder.Normals[face.Normals[i+1]*3+1]),
				float64(decoder.Normals[face.Normals[i+1]*3+2]),
			),
			pathtracer.NewVector(
				float64(decoder.Normals[face.Normals[i+2]*3]),
				float64(decoder.Normals[face.Normals[i+2]*3+1]),
				float64(decoder.Normals[face.Normals[i+2]*3+2]),
			),
		}

		uvs := [3]objUVCoordinate{}
		if len(decoder.Uvs) != 0 && face.Uvs[0] < len(decoder.Uvs) {
			uvs[0] = objUVCoordinate{
				float64(decoder.Uvs[face.Uvs[0]*2]),
				float64(decoder.Uvs[face.Uvs[0]*2+1]),
			}
			uvs[1] = objUVCoordinate{
				float64(decoder.Uvs[face.Uvs[i+1]*2]),
				float64(decoder.Uvs[face.Uvs[i+1]*2+1]),
			}
			uvs[2] = objUVCoordinate{
				float64(decoder.Uvs[face.Uvs[i+2]*2]),
				float64(decoder.Uvs[face.Uvs[i+2]*2+1]),
			}
		}

		triangle := objTriangle{
			vertexes: vertexes,
			normals:  normals,
			uvs:      uvs,
			material: objMaterial{source: decoder.Materials[face.Material]},
			smooth:   face.Smooth,
		}

		callback(triangle)
	}
}

// OBJScene contains obj geometry and materials loaded from the g3n game
// engine library.
type OBJScene struct {
	decoder   *obj.Decoder
	triangles []objTriangle
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
//
// TODO: Optimize this method for performance.
func (s *OBJScene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, material pathtracer.Material, ok bool) {
	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64
	var closestTriangle objTriangle

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
		closestTriangle = triangle
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

type objUVCoordinate [2]float64

func (o objUVCoordinate) U() float64 { return o[0] }
func (o objUVCoordinate) V() float64 { return o[1] }

type objTriangle struct {
	vertexes [3]pathtracer.Vector
	uvs      [3]objUVCoordinate
	normals  [3]pathtracer.Vector
	material objMaterial
	smooth   bool
}

func (t objTriangle) Vertex0() pathtracer.Vector { return t.vertexes[0] }
func (t objTriangle) Vertex1() pathtracer.Vector { return t.vertexes[1] }
func (t objTriangle) Vertex2() pathtracer.Vector { return t.vertexes[2] }

type objMaterial struct {
	source *obj.Material
}

func (m objMaterial) Sample(hit pathtracer.Hit, nextSample pathtracer.Sampler) pathtracer.Color {
	color := pathtracer.NewColor(0, 0, 0)

	if m.source.Diffuse.R >= pathtracer.EPS || m.source.Diffuse.G >= pathtracer.EPS || m.source.Diffuse.B >= pathtracer.EPS {
		ray := pathtracer.Ray{
			Origin:    hit.Position,
			Direction: material.DiffuseBounce(hit.Normal),
		}

		colorFromScene := nextSample(ray)
		colorToCamera := colorFromScene.Multiply(math32ToColor(m.source.Diffuse))
		color = color.Add(colorToCamera)
	}

	if m.source.Specular.R >= pathtracer.EPS || m.source.Specular.G >= pathtracer.EPS || m.source.Specular.B >= pathtracer.EPS {
		// TODO: Account for specular glossiness.
		ray := pathtracer.Ray{
			Origin:    hit.Position,
			Direction: material.SpecularBounce(hit.Normal, hit.From.Direction),
		}

		colorFromScene := nextSample(ray)
		colorToCamera := colorFromScene.Multiply(math32ToColor(m.source.Diffuse))
		color = color.Add(colorToCamera)
	}

	color = color.Add(math32ToColor(m.source.Emissive))

	return color
}

func math32ToColor(m32 math32.Color) pathtracer.Color {
	return pathtracer.NewColor(
		float64(m32.R),
		float64(m32.G),
		float64(m32.B),
	)
}
