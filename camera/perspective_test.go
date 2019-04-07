package camera

import (
	"testing"

	"github.com/ryannjohnson/pathtracer"
)

func TestPerspective(t *testing.T) {
	testCases := []struct {
		name                 string
		fov                  float64
		transformationMatrix pathtracer.Matrix
		x, y                 float64
		ray                  pathtracer.Ray
	}{
		{
			name:                 "should cast a ray at zero degrees when x = 0 and y = 0",
			fov:                  30,
			transformationMatrix: pathtracer.IdentityMatrix(),
			x:                    0,
			y:                    0,
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(0, 0, 1),
			},
		},
		{
			name:                 "should cast a ray from above origin when translated straight up",
			fov:                  30,
			transformationMatrix: pathtracer.IdentityMatrix().Translate(pathtracer.NewVector(0, 1, 0)),
			x:                    0,
			y:                    0,
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 1, 0),
				Direction: pathtracer.NewVector(0, 0, 1),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			camera := NewPerspective()
			camera.SetFieldOfView(testCase.fov)
			camera.SetTransformationMatrix(testCase.transformationMatrix)
			ray := camera.Cast(testCase.x, testCase.y)
			if ray != testCase.ray {
				t.Fatal("ray", ray, "doesn't equal expected", testCase.ray)
			}
		})
	}
}
