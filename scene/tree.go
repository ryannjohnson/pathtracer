package scene

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

// TreeShape is anything that occupies 3D space. It can be indexed in an
// AABB (axis-aligned bounding box) tree and can be intersected by rays.
type TreeShape interface {
	BoundingBox() Box
	Intersect(ray pathtracer.Ray) (distanceFromOrigin float64, position, normal pathtracer.Vector, ok bool)
	IntersectsBox(box Box) bool
}

// TreeNode indexes a scene's shapes for fast retrieval.
//
// The tree is a data structure that speeds up the retrieval of the
// first shape in a ray's path to O(log n).
//
// When generated, a tree of nodes starts with a node whose bounding box
// contains the entire scene. That node contains two child nodes, each
// split evenly along the longest axis of the parent node so that they
// fill the entire space of the parent.
//
// The ray will try to intersect each child node, favoring the closer
// intersection when both are hit. The ray continues to test
// intersections with child nodes until it reaches a node with no
// children. That node will contain indexes.
//
// The ray will then loop through all the indexes in the node, testing
// for intersections with each shape and returning the closest hit if
// any.
type TreeNode struct {
	box          Box
	shapeIndexes []int
	left         *TreeNode
	right        *TreeNode
}

// BuildTreeNode constructs a tree of TreeNodes that can be queried for
// shape indexes.
//
// The root TreeNode should contain a box that encapsulates all the
// shapes in the scene. That box gets subdivided along its longest axis
// until it either contains very few shapes or is notably smaller than
// any of the shapes it contains. This frequently happens when the box
// encapsulates a corner of 3 or more shapes.
//
// TreeNodes can either be branches or leaves. Branches have zero shape
// indexes and always have both their `left` and `right` nodes. Leaves
// always have a nonzero number of shape indexes and neither of their
// `left` or `right` nodes.
func BuildTreeNode(shapes []TreeShape, possibleShapeIndexes []int, box Box) *TreeNode {
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

		shapeBox := shape.BoundingBox()
		shapeLength := shapeBox.Max().Subtract(shapeBox.Min()).Length()
		if minShapeLength > shapeLength {
			minShapeLength = shapeLength
		}
	}
	maxBoxLength := box.Max().Subtract(box.Min()).Length() * 4
	if maxBoxLength < minShapeLength || minShapeLength < pathtracer.EPS {
		// Boxes are smaller than shapes at this point. Subdivision
		// has diminished its returns by now.
		node := &TreeNode{
			box:          box,
			shapeIndexes: shapeIndexes,
		}
		return node
	}

	boxA, boxB := splitBoxByLongestAxis(box)

	var nodeA *TreeNode
	var nodeB *TreeNode

	nodeChan := make(chan *TreeNode)
	go func() {
		nodeChan <- BuildTreeNode(shapes, shapeIndexes, boxA)
	}()
	go func() {
		nodeChan <- BuildTreeNode(shapes, shapeIndexes, boxB)
	}()
	nodeA = <-nodeChan
	nodeB = <-nodeChan

	if nodeA != nil && nodeB != nil {
		node := &TreeNode{
			box:   box,
			left:  nodeA,
			right: nodeB,
		}
		return node
	}

	if nodeA != nil {
		// Shorten the search hierarchy by returning the node directly,
		// instead of having this "middleman" node.
		return nodeA
	}

	if nodeB != nil {
		// Shorten the search hierarchy by returning the node directly,
		// instead of having this "middleman" node.
		return nodeB
	}

	// This should never happen if all the shapes are accounted for
	// in the boxes.
	panic("shapes were dropped in the tree building process")
}

// IntersectTreeNode searches the node tree, intersecting each node in
// the order according to the ray's trajectory, trying to intersect each
// shape until its first hit.
//
// NOTE: It turns out to be more performant to pass in []TreeShape
// instead of a getter interface.
func IntersectTreeNode(shapes []TreeShape, tn *TreeNode, ray pathtracer.Ray) (hit pathtracer.Hit, closestShapeIndex int, hitOk bool) {
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
