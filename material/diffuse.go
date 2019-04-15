package material

import (
	"math/rand"

	"github.com/ryannjohnson/pathtracer"
)

// DiffuseBounce returns a unit vector in the hemisphere of the supplied
// normal.
func DiffuseBounce(random *rand.Rand, normal pathtracer.Vector) pathtracer.Vector {
	var vector pathtracer.Vector

	for {
		vector = pathtracer.NewVector(
			random.Float64()*2-1,
			random.Float64()*2-1,
			random.Float64()*2-1,
		)

		vectorLength := vector.Length()
		if vectorLength >= 1 || vectorLength < pathtracer.EPS {
			continue
		}

		// Flips the direction of the vector if it's in the opposite
		// hemisphere.
		dotProduct := vector.DotProduct(normal)
		vector = vector.Scale(dotProduct)

		return vector.Normalize()
	}
}
