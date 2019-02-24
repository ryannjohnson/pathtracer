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

	aspectRatio := float64(width) / float64(height)

	var xRatio float64 = 1
	var yRatio float64 = 1
	if aspectRatio < 1 {
		xRatio = aspectRatio
	} else {
		yRatio = 1 / aspectRatio
	}

	for yPixel := 0; yPixel < height; yPixel++ {
		for xPixel := 0; xPixel < width; xPixel++ {
			x := xRatio * (float64(xPixel)/float64(width-1) - 0.5)       // Positive is right
			y := yRatio * (float64(yPixel)/float64(height-1) - 0.5) * -1 // Positive is up

			// TODO: Implement interchangable samplers via RenderSettings.
			colors := make([]Color, settings.SamplesPerRay)
			for i := 0; i < settings.SamplesPerRay; i++ {
				ray := camera.Cast(x, y)
				colors[i] = sampleScene(scene, ray, settings.BounceDepth)
			}
			color := averageColors(colors)

			image.Set(xPixel, yPixel, color)
		}
	}
}

func sampleScene(scene Scene, ray Ray, bouncesLeft int) Color {
	if bouncesLeft <= 0 {
		return black
	}

	hit := scene.Intersect(ray)

	nextSample := func(nextRay Ray) Color {
		return sampleScene(scene, nextRay, bouncesLeft-1)
	}

	return hit.Material.Sample(hit.From, hit.Position, hit.Normal, hit.UV, nextSample)
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
