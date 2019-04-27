package obj

import (
	"io"
	"math"

	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

// Scene contains obj geometry and materials loaded from the g3n game
// engine library.
type Scene struct {
	tree      *scene.TreeNode
	triangles []triangle
	shapes    []scene.Shape
}

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

	// Get the bounding box for the entire canvas so we can include all
	// the geometry in the scene.
	boxMin := pathtracer.NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	boxMax := pathtracer.NewVector(math.MaxFloat64*-1, math.MaxFloat64*-1, math.MaxFloat64*-1)
	for _, currentTriangle := range triangles {
		boxMin.X = math.Min(boxMin.X, math.Min(math.Min(currentTriangle.Vertex0().X, currentTriangle.Vertex1().X), currentTriangle.Vertex2().X))
		boxMin.Y = math.Min(boxMin.Y, math.Min(math.Min(currentTriangle.Vertex0().Y, currentTriangle.Vertex1().Y), currentTriangle.Vertex2().Y))
		boxMin.Z = math.Min(boxMin.Z, math.Min(math.Min(currentTriangle.Vertex0().Z, currentTriangle.Vertex1().Z), currentTriangle.Vertex2().Z))
		boxMax.X = math.Max(boxMax.X, math.Max(math.Max(currentTriangle.Vertex0().X, currentTriangle.Vertex1().X), currentTriangle.Vertex2().X))
		boxMax.Y = math.Max(boxMax.Y, math.Max(math.Max(currentTriangle.Vertex0().Y, currentTriangle.Vertex1().Y), currentTriangle.Vertex2().Y))
		boxMax.Z = math.Max(boxMax.Z, math.Max(math.Max(currentTriangle.Vertex0().Z, currentTriangle.Vertex1().Z), currentTriangle.Vertex2().Z))
	}
	rootBox := scene.NewBox(boxMin, boxMax)

	triangleIndexes := make([]int, len(triangles))
	shapes := make([]scene.Shape, len(triangles))
	for i := range triangles {
		triangleIndexes[i] = i
		shapes[i] = triangles[i]
	}

	rootNode := scene.BuildTreeNode(shapes, triangleIndexes, rootBox)

	return &Scene{rootNode, triangles, shapes}, nil
}

func (s *Scene) Clone() pathtracer.Scene {
	output := &Scene{
		tree:      s.tree.Clone(),
		triangles: make([]triangle, len(s.triangles)),
		shapes:    make([]scene.Shape, len(s.shapes)),
	}
	for i := range s.triangles {
		output.triangles[i] = s.triangles[i]
	}
	for i := range s.shapes {
		output.shapes[i] = s.shapes[i]
	}
	return output
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
func (s *Scene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, hitMaterial pathtracer.Material, ok bool) {
	var closestTriangleIndex int
	hit, closestTriangleIndex, ok = scene.IntersectTreeNode(s.shapes, s.tree, ray)

	if ok {
		hitMaterial = s.triangles[closestTriangleIndex].material
	}
	return
}
