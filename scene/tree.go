package scene

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

type Shape interface {
	Intersect(ray pathtracer.Ray) (distanceAlongRayFromOrigin float64, intersectionPoint, triangleNormal pathtracer.Vector, ok bool)
	IntersectsBox(box Box) bool
	Length() float64
}

type TreeNode struct {
	box          Box
	shapeIndexes []int
	left         *TreeNode
	right        *TreeNode
}

func BuildTreeNode(shapes []Shape, possibleShapeIndexes []int, box Box) *TreeNode {
	shapeIndexes := make([]int, 0)
	for _, shapeIndex := range possibleShapeIndexes {
		shape := shapes[shapeIndex]
		if shape.IntersectsBox(box) {
			shapeIndexes = append(shapeIndexes, shapeIndex)
		}
	}

	hasNoShapes := len(shapeIndexes) == 0
	if hasNoShapes {
		return nil
	}

	if len(shapeIndexes) <= 3 {
		// TODO: Test using 1 and 2 as the thresholds and measure
		// performance between identical renders.
		node := &TreeNode{
			box:          box,
			shapeIndexes: shapeIndexes,
		}
		return node
	}

	// No need to subdivide the containing box anymore if the smallest
	// shape is bigger than the largest side of the box.
	var minShapeLength = math.MaxFloat64
	for _, shapeIndex := range shapeIndexes {
		shape := shapes[shapeIndex]

		shapeLength := shape.Length()
		if minShapeLength > shapeLength {
			minShapeLength = shapeLength
		}
	}

	maxBoxLength := boxLongestEdge(box)
	if maxBoxLength < minShapeLength {
		// Boxes are smaller than shapes at this point. Subdivision
		// has diminished its returns by now.
		node := &TreeNode{
			box:          box,
			shapeIndexes: shapeIndexes,
		}
		return node
	}

	boxA, boxB := splitBoxByLongestAxis(box)

	nodeA := BuildTreeNode(shapes, shapeIndexes, boxA)
	nodeB := BuildTreeNode(shapes, shapeIndexes, boxB)

	if nodeA != nil && nodeB != nil {
		return &TreeNode{
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

	// This should never happen if all the shapes are accounted for
	// in the boxes.
	panic("shapes were dropped in the tree building process")
}

func IntersectTreeNode(shapes []Shape, tn *TreeNode, ray pathtracer.Ray) (hit pathtracer.Hit, closestShapeIndex int, hitOk bool) {
	if tn.left != nil && tn.right != nil {
		leftTMin, _, leftOk := tn.left.box.IntersectsRay(ray)
		rightTMin, _, rightOk := tn.right.box.IntersectsRay(ray)

		if leftOk && rightOk {
			// Since both boxes are intersected, we want to start with
			// the one that's closer to the ray's origin.
			//
			// Since we know the TreeNodees don't overlap, it's safe to
			// assume that tMin can be used to compare the two.
			if leftTMin < rightTMin {
				hit, closestShapeIndex, hitOk = IntersectTreeNode(shapes, tn.left, ray)
				if hitOk {
					return
				}
				return IntersectTreeNode(shapes, tn.right, ray)
			}
			hit, closestShapeIndex, hitOk = IntersectTreeNode(shapes, tn.right, ray)
			if hitOk {
				return
			}
			return IntersectTreeNode(shapes, tn.left, ray)
		}

		if leftOk {
			return IntersectTreeNode(shapes, tn.left, ray)
		}

		if rightOk {
			return IntersectTreeNode(shapes, tn.right, ray)
		}

		return
	}

	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64

	// Use these distances to determine if the shape intersections
	// happen within the volume of the box.
	//
	// For example, if the current box contains a ground plane that
	// eventually would be intersected but not inside the current
	// box, then it cannot be counted as valid _yet_.
	tMin, tMax, _ := tn.box.IntersectsRay(ray)

	for _, shapeIndex := range tn.shapeIndexes {
		shape := shapes[shapeIndex]
		distance, point, normal, didIntersect := shape.Intersect(ray)
		if !didIntersect {
			continue
		}

		if distance < tMin || distance > tMax {
			// Intersected outside of the box's volume.
			continue
		}

		if closestDistance < distance {
			continue
		}

		closestPoint = point
		closestNormal = normal
		closestDistance = distance
		closestShapeIndex = shapeIndex
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
	return
}

func boxLongestEdge(box Box) float64 {
	d := box.Max().Subtract(box.Min())
	return math.Max(math.Max(d.X, d.Y), d.Z)
}

func splitBoxByLongestAxis(box Box) (Box, Box) {
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
		a := NewBox(a0, a1)

		b0 := pathtracer.NewVector(
			middle,
			box.Min().Y,
			box.Min().Z,
		)
		b1 := box.Max()
		b := NewBox(b0, b1)
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
		a := NewBox(a0, a1)

		b0 := pathtracer.NewVector(
			box.Min().X,
			middle,
			box.Min().Z,
		)
		b1 := box.Max()
		b := NewBox(b0, b1)
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
	a := NewBox(a0, a1)

	b0 := pathtracer.NewVector(
		box.Min().X,
		box.Min().Y,
		middle,
	)
	b1 := box.Max()
	b := NewBox(b0, b1)
	return a, b
}
