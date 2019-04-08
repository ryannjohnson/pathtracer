package obj

import "github.com/ryannjohnson/pathtracer"

type tree struct {
	root *treeBox
}

type treeBox struct {
	positiveCorner pathtracer.Vector
	negativeCorner pathtracer.Vector
	triangles      []triangle
	left           *treeBox
	right          *treeBox
}
