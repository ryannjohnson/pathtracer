package main

import (
	"fmt"
	"math"
	"os"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/camera"
	"github.com/ryannjohnson/pathtracer/image"
)

type dummyUVCoordinate struct{}

func (c dummyUVCoordinate) U() float64 {
	return 0
}
func (c dummyUVCoordinate) V() float64 {
	return 0
}

type dummyMaterial struct {
	Color pathtracer.Color
}

func (m dummyMaterial) Sample(from pathtracer.Ray, position pathtracer.Vector, normal pathtracer.Ray, uv pathtracer.UVCoordinate, nextSample pathtracer.Sampler) pathtracer.Color {
	return m.Color
}

type dummyScene struct {
	Color pathtracer.Color
}

func (s dummyScene) Intersect(ray pathtracer.Ray) *pathtracer.Hit {
	return &pathtracer.Hit{Material: dummyMaterial{Color: s.Color}}
}

func main() {
	camera := camera.NewPerspective()
	cameraMatrix := pathtracer.IdentityMatrix()
	cameraMatrix = cameraMatrix.Scale(5)
	cameraMatrix = cameraMatrix.Rotate(pathtracer.AxisX, math.Pi/4)
	cameraMatrix = cameraMatrix.Rotate(pathtracer.AxisY, math.Pi/4)
	cameraMatrix = cameraMatrix.Translate(pathtracer.NewVector(-5, 5, -5))
	camera.SetTransformationMatrix(cameraMatrix)

	scene := dummyScene{Color: pathtracer.NewColor(0.5, 0.5, 0.5)}

	img := image.NewPNG8(64, 64)

	pathtracer.Render(scene, camera, img, &pathtracer.RenderSettings{
		BounceDepth:   5,
		SamplesPerRay: 5,
	})

	if err := img.Write(os.Stdout); err != nil {
		fmt.Println(err)
	}
}
