package pathtracer

// Camera creates rays to send into a scene to generate an image.
//
// The Expose function ranges from -0.5 to 0.5 in each of its x and y
// axis regarding its intended field of view (FOV). Eg, the top left
// corner is reprented by [-0.5, 0.5], and the bottom right corner is
// represented by [0.5, -0.5].
type Camera interface {
	Cast(x, y float64) Ray
}
