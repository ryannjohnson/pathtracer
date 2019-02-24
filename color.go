package pathtracer

var black = Color{0, 0, 0}

// NewColor takes its RGB values as float64 numbers.
func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

// Color is a set of RGB values, 0 being black and 1 being white. It's
// possible to have values outside of this range, which describe colors
// imperceptable to the human eye.
type Color struct {
	R, G, B float64
}
