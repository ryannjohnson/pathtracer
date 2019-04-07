package material

import (
	"math/rand"

	"github.com/ryannjohnson/pathtracer"
)

// DiffuseBounce returns a unit vector in the hemisphere of the supplied
// normal.
func DiffuseBounce(normal pathtracer.Vector) pathtracer.Vector {
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
