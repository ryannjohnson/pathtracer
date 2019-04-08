package material

import (
	"math/rand"

	"github.com/ryannjohnson/pathtracer"
)

// DiffuseBounce returns a unit vector in the hemisphere of the supplied
// normal.
func DiffuseBounce(random *rand.Rand, normal pathtracer.Vector) pathtracer.Vector {
	for {
		vector := pathtracer.NewVector(
			random.Float64()*2-1,
			random.Float64()*2-1,
			random.Float64()*2-1,
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
