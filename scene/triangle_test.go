package scene

import (
	"testing"

	"github.com/ryannjohnson/pathtracer"
)

type testTriangle struct {
	v0, v1, v2 pathtracer.Vector
}

func (t testTriangle) Vertex0() pathtracer.Vector { return t.v0 }
func (t testTriangle) Vertex1() pathtracer.Vector { return t.v1 }
func (t testTriangle) Vertex2() pathtracer.Vector { return t.v2 }

func TestIntersectTriangle(t *testing.T) {
	testCases := []struct {
		name                       string
		ray                        pathtracer.Ray
		triangle                   Triangle
		ok                         bool
		intersectionPoint          pathtracer.Vector
		intersectionNormal         pathtracer.Vector
		distanceAlongRayFromOrigin float64
	}{
		{
			name: "should pass if the ray passes through the middle of the triangle",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, -1, 0),
				pathtracer.NewVector(1, 1, 1),
				pathtracer.NewVector(1, 1, -1),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(1, 0, 0),
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: 1,
		},
		{
			name: "should pass if the ray passes through the middle of the triangle if the verticies are counter-clockwise",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, -1, 0),
				pathtracer.NewVector(1, 1, -1),
				pathtracer.NewVector(1, 1, 1),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(1, 0, 0),
			intersectionNormal:         pathtracer.NewVector(1, 0, 0),
			distanceAlongRayFromOrigin: 1,
		},
		{
			name: "should pass if a diagonal ray passes through the corner of the triangle",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 1, 1).Normalize(),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, -1, 0),
				pathtracer.NewVector(1, 1, 1),
				pathtracer.NewVector(1, 1, -1),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(1, 1, 1),
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: 1.732050807568877,
		},
		{
			name: "should pass if a diagonal ray passes through the corner of the triangle in the negative direction",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(-1, -1, -1).Normalize(),
			},
			triangle: testTriangle{
				pathtracer.NewVector(-1, 1, 0),
				pathtracer.NewVector(-1, -1, -1),
				pathtracer.NewVector(-1, -1, 1),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(-1, -1, -1),
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: 1.732050807568877,
		},
		{
			name: "should pass if the ray passes through the middle of the triangle at an extreme angle",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, -1, 0),
				pathtracer.NewVector(1000000, 1, 1),
				pathtracer.NewVector(1000000, 1, -1),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(500000.4999999999, 0, 0),
			intersectionNormal:         pathtracer.NewVector(-2.000001999998e-06, 0.9999999999979999, 0),
			distanceAlongRayFromOrigin: 500000.4999999999,
		},
		{
			name: "should pass if the ray passes through an edge of the triangle",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, 0, -1),
				pathtracer.NewVector(1, 0, 1),
				pathtracer.NewVector(1, 1, 0),
			},
			ok:                         true,
			intersectionPoint:          pathtracer.NewVector(1, 0, 0),
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: 1,
		},
		{
			name: "should fail if the ray misses the triangle",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, 2, -1),
				pathtracer.NewVector(1, 1, 0),
				pathtracer.NewVector(1, 2, 1),
			},
			ok:                         false,
			intersectionPoint:          pathtracer.NewVector(1, 0, 0),
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: 1,
		},
		{
			name: "should fail if the triangle is behind the ray's origin",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(-1, -1, 0),
				pathtracer.NewVector(-1, 1, 1),
				pathtracer.NewVector(-1, 1, -1),
			},
			ok:                         false,
			intersectionNormal:         pathtracer.NewVector(-1, 0, 0),
			distanceAlongRayFromOrigin: -1,
		},
		{
			name: "should fail if the triangle and ray are parallel",
			ray: pathtracer.Ray{
				Origin:    pathtracer.NewVector(0, 0, 0),
				Direction: pathtracer.NewVector(1, 0, 0),
			},
			triangle: testTriangle{
				pathtracer.NewVector(1, 0, 0),
				pathtracer.NewVector(2, 1, 0),
				pathtracer.NewVector(1, 1, 0),
			},
			ok:                 false,
			intersectionNormal: pathtracer.NewVector(0, 0, 1),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			distance, point, normal, ok := IntersectTriangle(testCase.ray, testCase.triangle)
			if ok != testCase.ok {
				t.Fatal("ok", ok, "didn't match expected", testCase.ok)
			}

			if point != testCase.intersectionPoint {
				t.Error("intersectionPoint", point, "doesn't match expected", testCase.intersectionPoint)
			}
			if normal != testCase.intersectionNormal {
				t.Error("intersectionNormal", normal, "doesn't match expected", testCase.intersectionNormal)
			}
			if distance != testCase.distanceAlongRayFromOrigin {
				t.Error("distanceAlongRayFromOrigin", distance, "doesn't match expected", testCase.distanceAlongRayFromOrigin)
			}
		})
	}
}
