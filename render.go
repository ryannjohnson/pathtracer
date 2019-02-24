package pathtracer

// RenderSettings describes options related to the performance and
// output of the Render function.
type RenderSettings struct {
	BounceDepth   int
	SamplesPerRay int
}

// ImageWriter represents a 2D canvas that receives colors for each of
// its descrete pixels
//
// Whether or not WriteColor is called multiple times per pixel is up to
// the RenderSettings.
type ImageWriter interface {
	Width() int
	Height() int
	Set(x, y int, color Color)
}

// Render converts a 3D scene into a 2D image.
func Render(scene Scene, camera Camera, image ImageWriter, settings *RenderSettings) {
	width := image.Width()
	height := image.Height()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colors := make([]Color, settings.SamplesPerRay)
			xFloat := float64(x)
			yFloat := float64(y)
			for i := 0; i < settings.SamplesPerRay; i++ {
				ray := camera.Cast(xFloat, yFloat)
				colors[i] = scene.Sample(ray, settings.BounceDepth)
			}
			color := averageColors(colors)
			image.Set(x, y, color)
		}
	}
}

func averageColors(colors []Color) Color {
	output := Color{}

	total := len(colors)
	for i := 0; i < total; i++ {
		output.R += colors[i].R
		output.G += colors[i].G
		output.B += colors[i].B
	}

	totalFloat := float64(total)
	output.R = output.R / totalFloat
	output.G = output.G / totalFloat
	output.B = output.B / totalFloat

	return output
}
