package camera

import (
	"math"

	"github.com/ryannjohnson/pathtracer"
)

// Perspective is a camera that simulates how the human eye works,
// casting rays from a point behind the camera in the world.
type Perspective struct {
	fieldOfView          float64 // Degrees [0, 180)
	transformationMatrix pathtracer.Matrix
}

// Cast converts the x and y coordinates into a Ray that can be cast
// from that point on the 2D plane.
//
// The perspective camera plots ray origins along the surface of a
// sphere. The size of the sphere is dictated by the magnitude of
// fieldOfView.
//
// When depth of field is applied, the origin of the ray will no longer
// originate at the x, y coordinate of the original plane but will
// instead be cast to intersect the focal point in front of the camera.
func (c Perspective) Cast(x, y float64) pathtracer.Ray {
	// Find origin vector (relative to 0,0,0) from the fieldOfView.
	radians := c.fieldOfView * math.Pi / 180

	// Find the direction for the x, y coordinate by using FOV as 100%.
	m := pathtracer.NewMatrix()
	m = m.Rotate(pathtracer.AxisX, y*radians)
	m = m.Rotate(pathtracer.AxisY, x*radians)

	direction := pathtracer.AxisZ.Transform(m)

	// Find the final position for the ray's origin based on the circle
	// created by the camera's field of view.
	radius := 1 / radians
	center := pathtracer.Vector{X: 0, Y: 0, Z: radius * -1}

	origin := center.Add(direction.Scale(radius))

	ray := pathtracer.Ray{
		Origin:    origin,
		Direction: direction,
	}

	return ray.Transform(c.transformationMatrix)
}
