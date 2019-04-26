package obj

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

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
	for i := range triangles {
		triangleIndexes[i] = i
	}

	rootNode := buildTreeNode(triangles, triangleIndexes, rootBox)

	return &tree{rootNode, triangles}
}

func buildTreeNode(triangles []triangle, triangleIndexes []int, box scene.Box) *treeNode {
	triangleIndexesInBox := make([]int, 0)
	for _, index := range triangleIndexes {
		if box.IntersectsTriangle(triangles[index]) {
			triangleIndexesInBox = append(triangleIndexesInBox, index)
		}
	}

	hasNoTriangles := len(triangleIndexesInBox) == 0
	if hasNoTriangles {
		return nil
	}

	if len(triangleIndexesInBox) <= 3 {
		// TODO: Test using 1 and 2 as the thresholds and measure
		// performance between identical renders.
		node := &treeNode{
			box:             box,
			triangleIndexes: triangleIndexesInBox,
		}
		return node
	}

	// No need to subdivide the containing box anymore if the smallest
	// triangle is bigger than the largest side of the box.
	var minTriangleArea = math.MaxFloat64
	for _, triangleIndex := range triangleIndexesInBox {
		currentTriangle := triangles[triangleIndex]

		v0v1 := currentTriangle.Vertex0().Subtract(currentTriangle.Vertex1())
		v0v2 := currentTriangle.Vertex0().Subtract(currentTriangle.Vertex2())
		currentTriangleArea := v0v1.CrossProduct(v0v2).Length() / 2

		if minTriangleArea > currentTriangleArea {
			minTriangleArea = currentTriangleArea
		}
	}

	maxBoxArea := boxLargestSideArea(box)
	if maxBoxArea < minTriangleArea {
		// Boxes are smaller than triangles at this point. Subdivision
		// has diminished its returns by now.
		node := &treeNode{
			box:             box,
			triangleIndexes: triangleIndexesInBox,
		}
		return node
	}

	boxA, boxB := splitBoxByLongestAxis(box)

	nodeA := buildTreeNode(triangles, triangleIndexesInBox, boxA)
	nodeB := buildTreeNode(triangles, triangleIndexesInBox, boxB)

	if nodeA != nil && nodeB != nil {
		return &treeNode{
			box:   box,
			left:  nodeA,
			right: nodeB,
		}
	}

	if nodeA != nil {
		// Shorten the search hierarchy by returning the node directly.
		return nodeA
	}

	if nodeB != nil {
		// Shorten the search hierarchy by returning the node directly.
		return nodeB
	}

	// This should never happen if all the triangles are accounted for
	// in the boxes.
	panic("triangles were dropped in the tree building process")
}

type tree struct {
	root      *treeNode
	triangles []triangle
}

func (t *tree) intersect(ray pathtracer.Ray) (pathtracer.Hit, material, bool) {
	return intersectTreeNode(t.triangles, t.root, ray)
}

func intersectTreeNode(triangles []triangle, tn *treeNode, ray pathtracer.Ray) (hit pathtracer.Hit, hitMaterial material, hitOk bool) {
	if tn.left != nil && tn.right != nil {
		leftTMin, _, leftOk := tn.left.box.IntersectsRay(ray)
		rightTMin, _, rightOk := tn.right.box.IntersectsRay(ray)

		if leftOk && rightOk {
			// Since both boxes are intersected, we want to start with
			// the one that's closer to the ray's origin.
			//
			// Since we know the treeNodees don't overlap, it's safe to
			// assume that tMin can be used to compare the two.
			if leftTMin < rightTMin {
				firstHit, firstHitMaterial, firstHitOk := intersectTreeNode(triangles, tn.left, ray)
				if firstHitOk {
					hit = firstHit
					hitMaterial = firstHitMaterial
					hitOk = true
					return
				}
				return intersectTreeNode(triangles, tn.right, ray)
			}
			firstHit, firstHitMaterial, firstHitOk := intersectTreeNode(triangles, tn.right, ray)
			if firstHitOk {
				hit = firstHit
				hitMaterial = firstHitMaterial
				hitOk = true
				return
			}
			return intersectTreeNode(triangles, tn.left, ray)
		}

		if leftOk {
			return intersectTreeNode(triangles, tn.left, ray)
		}

		if rightOk {
			return intersectTreeNode(triangles, tn.right, ray)
		}

		return
	}

	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64
	var closestTriangle triangle

	for _, triangleIndex := range tn.triangleIndexes {
		faceTriangle := triangles[triangleIndex]
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

	hitOk = true
	hit = pathtracer.Hit{
		From:     ray,
		Position: closestPoint,
		Normal:   closestNormal,
	}
	hitMaterial = closestTriangle.material
	return
}

type treeNode struct {
	box             scene.Box
	triangleIndexes []int
	left            *treeNode
	right           *treeNode
}

func trianglesByIndexes(triangles []triangle, indexes []int) []triangle {
	out := make([]triangle, len(indexes))
	for i, index := range indexes {
		out[i] = triangles[index]
	}
	return out
}

func boxLargestSideArea(box scene.Box) float64 {
	d := box.Max().Subtract(box.Min())
	return math.Max(math.Max(d.X*d.Y, d.X*d.Z), d.Y*d.Z)
}

func splitBoxByLongestAxis(box scene.Box) (scene.Box, scene.Box) {
	d := box.Max().Subtract(box.Min())

	if d.X >= d.Y && d.X >= d.Z {
		// X is longest.
		middle := box.Min().X + d.X/2

		a0 := box.Min()
		a1 := pathtracer.NewVector(
			middle,
			box.Max().Y,
			box.Max().Z,
		)
		a := scene.NewBox(a0, a1)

		b0 := pathtracer.NewVector(
			middle,
			box.Min().Y,
			box.Min().Z,
		)
		b1 := box.Max()
		b := scene.NewBox(b0, b1)
		return a, b
	}

	if d.Y >= d.X && d.Y >= d.Z {
		// Y is longest.
		middle := box.Min().Y + d.Y/2

		a0 := box.Min()
		a1 := pathtracer.NewVector(
			box.Max().X,
			middle,
			box.Max().Z,
		)
		a := scene.NewBox(a0, a1)

		b0 := pathtracer.NewVector(
			box.Min().X,
			middle,
			box.Min().Z,
		)
		b1 := box.Max()
		b := scene.NewBox(b0, b1)
		return a, b
	}

	// Z is longest.
	middle := box.Min().Z + d.Z/2

	a0 := box.Min()
	a1 := pathtracer.NewVector(
		box.Max().X,
		box.Max().Y,
		middle,
	)
	a := scene.NewBox(a0, a1)

	b0 := pathtracer.NewVector(
		box.Min().X,
		box.Min().Y,
		middle,
	)
	b1 := box.Max()
	b := scene.NewBox(b0, b1)
	return a, b
}
