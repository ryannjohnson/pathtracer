package main

import (
	"fmt"
	"math"
	"os"

	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/camera"
	"github.com/ryannjohnson/pathtracer/image"
	sceneLoader "github.com/ryannjohnson/pathtracer/scene/loader"
)

func main() {
	camera := camera.NewPerspective()
	cameraMatrix := pathtracer.IdentityMatrix()
	cameraMatrix = cameraMatrix.Translate(pathtracer.NewVector(-5, 1, -5))
	cameraMatrix = cameraMatrix.Rotate(pathtracer.AxisY, math.Pi/-4)
	camera.SetTransformationMatrix(cameraMatrix)
	camera.SetFieldOfView(45)
	camera.SetDepthOfField(6, 0.25)

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

	scene, err := sceneLoader.NewOBJScene(objFile, mtlFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	img := image.NewPNG8(200, 200)

	pathtracer.Render(scene, camera, img, &pathtracer.RenderSettings{
		BounceDepth:   5,
		SamplesPerRay: 100,
	})

	if err := img.Write(os.Stdout); err != nil {
		fmt.Println(err)
	}
}
