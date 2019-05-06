package main

import (
	"fmt"
	"math"
	"os"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/camera"
	"github.com/ryannjohnson/pathtracer/image"
	"github.com/ryannjohnson/pathtracer/scene/obj"
)

func main() {
	camera := camera.NewPerspective()
	cameraMatrix := pathtracer.IdentityMatrix()
	cameraMatrix = cameraMatrix.Rotate(pathtracer.AxisY, math.Pi)
	cameraMatrix = cameraMatrix.Rotate(pathtracer.AxisX, math.Pi/-4.3)
	cameraMatrix = cameraMatrix.Translate(pathtracer.NewVector(0, 5, 5))
	camera.SetTransformationMatrix(cameraMatrix)
	camera.SetFieldOfView(20)
	camera.SetDepthOfField(6.666, 0.07)

	objFile, err := os.Open("./examples/triangle/scene.obj")
	if err != nil {
		fmt.Println(err)
		return
	}

	mtlFile, err := os.Open("./examples/triangle/scene.mtl")
	if err != nil {
		fmt.Println(err)
		return
	}

	scene, err := obj.NewScene(objFile, mtlFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	img := image.NewPNG8(1024, 1024)

	pathtracer.Render(scene, camera, img, &pathtracer.RenderSettings{
		BounceDepth:   5,
		SamplesPerRay: 1000,
	})

	if err := img.Write(os.Stdout); err != nil {
		fmt.Println(err)
	}
}
