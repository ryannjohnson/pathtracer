package obj

import (
	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
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
