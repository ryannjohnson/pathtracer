package image

import (
	"image/color"
	"math"

	"github.com/ryannjohnson/pathtracer"
)

func toNRGBA(c pathtracer.Color) color.NRGBA {
	return color.NRGBA{
		R: to8Bit(c.R),
		G: to8Bit(c.G),
		B: to8Bit(c.B),
		A: 255,
	}
}

func to8Bit(c float64) uint8 {
	if c <= 0 {
		return 0
	}
	if c >= 1 {
		return 255
	}
	return uint8(math.Floor(c * 255))
}
