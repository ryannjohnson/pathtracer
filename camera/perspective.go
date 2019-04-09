package camera

import (
	"math"
	"math/rand"

	"github.com/ryannjohnson/pathtracer"
)

// NewPerspective creates a camera with some defaults attached.
func NewPerspective() *Perspective {
	return &Perspective{
		fieldOfView:          30,
		transformationMatrix: pathtracer.IdentityMatrix(),
	}
}

// Perspective is a camera that simulates how the human eye works,
// casting rays from a point behind the camera in the world.
//
// Its default orientation is facing towards the positive Z axis.
//
// The depth of field distance refers to how far from the camera after
// transformations will be the focal target in the scene.
//
// The depth of field radius describes how wide the circle will be that
// new rays are cast from. This circle will be placed 1 distance unit
// away from the focal target.
type Perspective struct {
	depthOfFieldDistance float64
	depthOfFieldRadius   float64 // From lens
	fieldOfView          float64 // Degrees
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
func (c Perspective) Cast(random *rand.Rand, x, y float64) pathtracer.Ray {
	// Find origin vector (relative to 0,0,0) from the fieldOfView.
	fieldOfViewRadians := c.fieldOfView * math.Pi / 180

	// Find the direction for the x, y coordinate by using FOV as 100%.
	m := pathtracer.IdentityMatrix()
	m = m.Rotate(pathtracer.AxisX, y*fieldOfViewRadians)
	m = m.Rotate(pathtracer.AxisY, x*fieldOfViewRadians)

	direction := pathtracer.AxisZ.Transform(m)

	// Find the final position for the ray's origin based on the circle
	// created by the camera's field of view.
	focalLength := 1 / fieldOfViewRadians
	center := pathtracer.Vector{X: 0, Y: 0, Z: focalLength * -1}

	// Where the pixel rests in 3D space on the lens of the camera.
	origin := center.Add(direction.Scale(focalLength))

	ray := pathtracer.Ray{
		Origin:    origin,
		Direction: direction,
	}

	// Includes translation, rotation, and scale of the camera.
	ray = ray.Transform(c.transformationMatrix)

	if c.depthOfFieldRadius >= pathtracer.EPS {
		// Slide the origin up to the part of the environment we want to be
		// in complete focus. Do this now while we have the
		focalOrigin := ray.Origin.Add(ray.Direction.Scale(c.depthOfFieldDistance))

		// Create a vector perpendicular to the ray direction.
		perpendicularAxis := arbitraryOrthogonal(ray.Direction).Normalize()

		// Rotate the perpendicular vector around the ray by some random
		// amount. This is important to ensure that samples are taken from
		// all directions surrounding the focal target.
		m = pathtracer.IdentityMatrix()
		m = m.Rotate(ray.Direction, 2*math.Pi*random.Float64())
		perpendicularAxis = perpendicularAxis.Transform(m)

		// Adjust the direction of the ray as if it were cast from a circle
		// around its original place on the lens.
		m = pathtracer.IdentityMatrix()
		m = m.Rotate(perpendicularAxis, math.Atan2(c.depthOfFieldRadius, c.depthOfFieldDistance))
		direction = ray.Direction.Transform(m)

		// Slide the origin back down to the lens, using the new focal
		// direction.
		origin = focalOrigin.Subtract(direction.Scale(c.depthOfFieldDistance))

		ray = pathtracer.Ray{
			Origin:    origin,
			Direction: direction,
		}
	}

	return ray
}

// Return an arbitrary vector perpendicular to the unit vector supplied.
//
// https://stackoverflow.com/a/43454629/5307109
func arbitraryOrthogonal(v pathtracer.Vector) pathtracer.Vector {
	w := pathtracer.NewVector(0, 0, 0)

	if v.X < v.Y && v.X < v.Z {
		w.X = 1
	}
	if v.Y <= v.X && v.Y < v.Z {
		w.Y = 1
	}
	if v.Z <= v.X && v.Z <= v.Y {
		w.Z = 1
	}

	return v.CrossProduct(w)
}

// SetDepthOfField sets the distance and radius of the depth of field.
//
// The distance is not scaled according to the camera's transformation
// matrix.
//
// The radius is calculated as if the circle is positioned at the lens.
func (c *Perspective) SetDepthOfField(distance, radius float64) {
	c.depthOfFieldDistance = distance
	c.depthOfFieldRadius = radius
}

// SetFieldOfView expects an angle in degrees.
//
// https://en.wikipedia.org/wiki/Field_of_view
func (c *Perspective) SetFieldOfView(fov float64) {
	c.fieldOfView = fov
}

// SetTransformationMatrix sets this camera's transformation matrix,
// which is applied to every ray cast from the camera.
func (c *Perspective) SetTransformationMatrix(m pathtracer.Matrix) {
	c.transformationMatrix = m
}
