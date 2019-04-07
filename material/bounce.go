package material

import (
	"math/rand"

	"github.com/ryannjohnson/pathtracer"
)

// DiffuseBounce returns a unit vector in the hemisphere of the supplied
// normal.
func DiffuseBounce(normal pathtracer.Vector) pathtracer.Vector {
	return randomDirectionWithinHemisphere(normal)
}

// SpecularBounce reflects the incident vector across the normal.
//
// https://en.wikipedia.org/wiki/Specular_reflection
func SpecularBounce(normal, incident pathtracer.Vector) pathtracer.Vector {
	// Since both values are unit vectors, this essentially results in
	// the length of the incident vector when reduced to the dimension
	// parallel to the normal.
	//
	// For example, if normal was pointing straight up the Y axis in 2D
	// space, then the dot product would represent the length of the
	// incident ray's Y value. The X axis would be zeroed out.
	//
	// Since the normal and incident will be in "opposite" directions,
	// the height will need to be flipped to stay positive.
	incidentHeight := incident.DotProduct(normal) * -1

	// If the incident ray lands at the feet of the normal, then we now
	// want the normal to be exactly twice as tall than where the
	// incident started.
	//
	// We use the normal to determine what "tall" is relative to. If
	// they're perpendicular, then tall is 0. If they're parallel, then
	// double the height of the incident is 2.
	scaledNormal := normal.Scale(2 * incidentHeight)

	// Now that we're at two-times the height of the incident vector
	// relative to the normal, we just add the incident vector again to
	// arrive at a vector at the same "height" as the original, and the
	// same distance away from the normal in the opposite direction.
	reflectedDirection := scaledNormal.Add(incident)

	return reflectedDirection
}

func randomDirectionWithinHemisphere(normal pathtracer.Vector) pathtracer.Vector {
	for {
		vector := pathtracer.NewVector(
			rand.Float64()*2-1,
			rand.Float64()*2-1,
			rand.Float64()*2-1,
		)

		vectorLength := vector.Length()
		if vectorLength >= 1 || vectorLength < pathtracer.EPS {
			continue
		}

		vector = vector.Normalize()

		if vector.DotProduct(normal) < 0 {
			return pathtracer.NewVector(vector.X*-1, vector.Y*-1, vector.Z*-1)
		}

		return vector
	}
}
