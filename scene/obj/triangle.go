package obj

import (
	"math"

	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

type objUVCoordinate [2]float64

func (o objUVCoordinate) U() float64 { return o[0] }
func (o objUVCoordinate) V() float64 { return o[1] }

type triangle struct {
	vertexes [3]pathtracer.Vector
	uvs      [3]objUVCoordinate
	normals  [3]pathtracer.Vector
	material material
	smooth   bool
}

func (t triangle) Vertex0() pathtracer.Vector { return t.vertexes[0] }
func (t triangle) Vertex1() pathtracer.Vector { return t.vertexes[1] }
func (t triangle) Vertex2() pathtracer.Vector { return t.vertexes[2] }

func (t triangle) Intersect(ray pathtracer.Ray) (distanceAlongRayFromOrigin float64, intersectionPoint, triangleNormal pathtracer.Vector, ok bool) {
	return scene.IntersectTriangle(ray, t)
}

func (t triangle) IntersectsBox(box scene.Box) bool {
	return box.IntersectsTriangle(t)
}

func (t triangle) Length() float64 {
	v0v1 := t.Vertex0().Subtract(t.Vertex1()).Length()
	v1v2 := t.Vertex1().Subtract(t.Vertex2()).Length()
	v2v0 := t.Vertex2().Subtract(t.Vertex0()).Length()
	return math.Max(math.Max(v0v1, v1v2), v2v0)
}

func readTriangles(decoder *obj.Decoder, face obj.Face, callback func(triangle)) {
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

		faceTriangle := triangle{
			vertexes: vertexes,
			normals:  normals,
			uvs:      uvs,
			material: material{source: decoder.Materials[face.Material]},
			smooth:   face.Smooth,
		}

		callback(faceTriangle)
	}
}
