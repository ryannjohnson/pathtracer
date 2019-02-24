package image

import (
	"image"
	"image/png"
	"io"

	"github.com/ryannjohnson/pathtracer"
)

// NewPNG8 creates a fresh 8-bit image in memory.
func NewPNG8(width, height int) *PNG8 {
	return &PNG8{
		img: image.NewNRGBA(image.Rect(0, 0, width, height)),
	}
}

// PNG8 represents an image with 8-bit channels.
type PNG8 struct {
	img *image.NRGBA
}

// Width is the number of pixels this image has along the x axis.
func (p PNG8) Width() int {
	return p.img.Rect.Dx()
}

// Height is the number of pixels this image has along the y axis.
func (p PNG8) Height() int {
	return p.img.Rect.Dy()
}

// Set updates the pixel value for a particular x, y combination.
func (p *PNG8) Set(x, y int, c pathtracer.Color) {
	p.img.Set(x, y, toNRGBA(c))
}

// Write saves the image to PNG format.
func (p PNG8) Write(w io.Writer) error {
	return png.Encode(w, p.img)
}
