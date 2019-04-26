package scene

import (
	"testing"

	"github.com/ryannjohnson/pathtracer"
)

func TestBoxIntersectsTriangle(t *testing.T) {
	testCases := []struct {
		name     string
		box      Box
		triangle testTriangle
		result   bool
	}{
		{
			name: "should not intersect when too far negative on the x axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(-2, 0, 0),
				pathtracer.NewVector(-2, 1, 0),
				pathtracer.NewVector(-2, 0, 1),
			},
			result: false,
		},
		{
			name: "should not intersect when too far positive on the x axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(2, 0, 0),
				pathtracer.NewVector(2, 1, 0),
				pathtracer.NewVector(2, 0, 1),
			},
			result: false,
		},
		{
			name: "should not intersect when too far negative on the y axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(1, -2, 0),
				pathtracer.NewVector(0, -2, 0),
				pathtracer.NewVector(0, -2, 1),
			},
			result: false,
		},
		{
			name: "should not intersect when too far positive on the y axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(1, 2, 0),
				pathtracer.NewVector(0, 2, 0),
				pathtracer.NewVector(0, 2, 1),
			},
			result: false,
		},
		{
			name: "should not intersect when too far negative on the z axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(1, 0, -2),
				pathtracer.NewVector(0, 1, -2),
				pathtracer.NewVector(0, 0, -2),
			},
			result: false,
		},
		{
			name: "should not intersect when too far positive on the z axis",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(1, 0, 2),
				pathtracer.NewVector(0, 1, 2),
				pathtracer.NewVector(0, 0, 2),
			},
			result: false,
		},
		{
			name: "should intersect when the triangle plane intersects with a box corner",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 3, 0),
				pathtracer.NewVector(3, 0, 0),
				pathtracer.NewVector(0, 0, 3),
			},
			result: true,
		},
		{
			name: "should not intersect when the separating plane is along the triangle normal",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 3, 0),
				pathtracer.NewVector(3, 0, 0),
				pathtracer.NewVector(0, 0, 3.5),
			},
			result: false,
		},
		{
			name: "should not intersect when the triangle plane intersects with the box but the triangle doesn't",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 3, 0),
				pathtracer.NewVector(3, 0, 0),
				pathtracer.NewVector(3, 3, 0),
			},
			result: false,
		},
		{
			name: "should intersect when the triangle is entirely inside the box",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 0.5, 0),
				pathtracer.NewVector(0.5, 0, 0),
				pathtracer.NewVector(0, 0, 0.5),
			},
			result: true,
		},
		{
			name: "should intersect when the triangle is partially inside the box",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 0.5, 0),
				pathtracer.NewVector(8, 0, 0),
				pathtracer.NewVector(0, 0, 0.5),
			},
			result: true,
		},
		{
			name: "should intersect when the box is in the middle of the triangle",
			box: NewBox(
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(1, 1, 1),
			),
			triangle: testTriangle{
				pathtracer.NewVector(0, 8, 0),
				pathtracer.NewVector(8, -4, 0),
				pathtracer.NewVector(-8, -4, 0),
			},
			result: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.box.IntersectsTriangle(testCase.triangle)
			if result != testCase.result {
				t.Error("result", result, "did not match expected", testCase.result)
			}
		})
	}
}
