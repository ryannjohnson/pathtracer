package pathtracer

// Material take a ray and returns its resulting color.
//
// Materials are encouraged to send rays back into scenes in order to
// hit a light source if necessary.
//
// Materials can also be lights, themselves.
type Material interface {
	Sample(from Ray, position Vector, normal Ray, uv UVCoordinate, nextSample func(ray Ray) Color) Color
}
