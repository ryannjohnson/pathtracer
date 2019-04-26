package obj

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

type tree struct {
	root      *scene.TreeNode
	triangles []triangle
	shapes    []scene.Shape
}

func (t *tree) intersect(ray pathtracer.Ray) (hit pathtracer.Hit, hitMaterial material, ok bool) {
	var closestTriangleIndex int
	hit, closestTriangleIndex, ok = scene.IntersectTreeNode(t.shapes, t.root, ray)

	if ok {
		hitMaterial = t.triangles[closestTriangleIndex].material
	}
	return
}

func newTree(triangles []triangle) *tree {
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

	return &tree{rootNode, triangles, shapes}
}
