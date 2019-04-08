package obj

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/scene"
)

type tree struct {
	root *treeBox
}

func (t *tree) intersect(ray pathtracer.Ray) (pathtracer.Hit, material, bool) {
	return intersectTreeBox(t.root, ray)
}

func intersectTreeBox(tb *treeBox, ray pathtracer.Ray) (hit pathtracer.Hit, hitMaterial material, hitOk bool) {
	if tb.left != nil && tb.right != nil {
		leftTMin, _, leftOk := tb.left.box.IntersectsRay(ray)
		rightTMin, _, rightOk := tb.right.box.IntersectsRay(ray)

		if leftOk && rightOk {
			// Since both boxes are intersected, we want to start with
			// the one that's closer to the ray's origin.
			//
			// Since we know the treeBoxes don't overlap, it's safe to
			// assume that tMin can be used to compare the two.
			if leftTMin < rightTMin {
				firstHit, firstHitMaterial, firstHitOk := intersectTreeBox(tb.left, ray)
				if firstHitOk {
					hit = firstHit
					hitMaterial = firstHitMaterial
					hitOk = true
					return
				}
				return intersectTreeBox(tb.right, ray)
			}
			firstHit, firstHitMaterial, firstHitOk := intersectTreeBox(tb.right, ray)
			if firstHitOk {
				hit = firstHit
				hitMaterial = firstHitMaterial
				hitOk = true
				return
			}
			return intersectTreeBox(tb.left, ray)
		}

		if leftOk {
			return intersectTreeBox(tb.left, ray)
		}

		if rightOk {
			return intersectTreeBox(tb.right, ray)
		}

		return
	}

	var closestPoint pathtracer.Vector
	var closestNormal pathtracer.Vector
	var closestDistance = math.MaxFloat64
	var closestTriangle triangle

	for _, faceTriangle := range tb.triangles {
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

type treeBox struct {
	box       scene.Box
	triangles []triangle
	left      *treeBox
	right     *treeBox
}
