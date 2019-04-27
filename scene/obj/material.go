package obj

import (
	"math/rand"

	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/ryannjohnson/pathtracer"
	"github.com/ryannjohnson/pathtracer/material"
)

type objMaterial struct {
	source *obj.Material
}

func (m objMaterial) Sample(random *rand.Rand, hit pathtracer.Hit, nextSample pathtracer.Sampler) pathtracer.Color {
	color := pathtracer.NewColor(0, 0, 0)

	if m.source.Diffuse.R >= pathtracer.EPS || m.source.Diffuse.G >= pathtracer.EPS || m.source.Diffuse.B >= pathtracer.EPS {
		ray := pathtracer.Ray{
			Origin:    hit.Position,
			Direction: material.DiffuseBounce(random, hit.Normal),
		}

		colorFromScene := nextSample(ray)
		colorToCamera := colorFromScene.Multiply(math32ToColor(m.source.Diffuse))
		color = color.Add(colorToCamera)
	}

	if m.source.Specular.R >= pathtracer.EPS || m.source.Specular.G >= pathtracer.EPS || m.source.Specular.B >= pathtracer.EPS {
		// TODO: Account for specular glossiness.
		ray := pathtracer.Ray{
			Origin:    hit.Position,
			Direction: material.SpecularBounce(hit.Normal, hit.From.Direction),
		}

		colorFromScene := nextSample(ray)
		colorToCamera := colorFromScene.Multiply(math32ToColor(m.source.Diffuse))
		color = color.Add(colorToCamera)
	}

	color = color.Add(math32ToColor(m.source.Emissive))

	return color
}

func math32ToColor(m32 math32.Color) pathtracer.Color {
	return pathtracer.NewColor(
		float64(m32.R),
		float64(m32.G),
		float64(m32.B),
	)
}
