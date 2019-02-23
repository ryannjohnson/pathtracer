package pathtracer

// Camera describes how to send rays into a scene to return colors.
//
// The Expose function ranges from -0.5 to 0.5 in each of its x and y
// axis regarding its intended field of view (FOV). Eg, the top left
// corner is reprented by [-0.5, 0.5], and the bottom right corner is
// represented by [0.5, -0.5].
type Camera interface {
	Expose(scene Scene, x, y float64, bouncesLeft int) Color
}
