package obj

import (
	"io"

	"github.com/g3n/engine/loader/obj"
	"github.com/ryannjohnson/pathtracer"
)

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

	treeRoot := newTree(triangles)

	return &Scene{decoder, treeRoot}, nil
}

// Scene contains obj geometry and materials loaded from the g3n game
// engine library.
type Scene struct {
	decoder *obj.Decoder
	tree    *tree
}

// Intersect finds the first geometry a ray passes through in the scene
// and returns details about the intersection and its material.
func (s *Scene) Intersect(ray pathtracer.Ray) (hit pathtracer.Hit, material pathtracer.Material, ok bool) {
	return s.tree.intersect(ray)
}
