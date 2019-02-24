package pathtracer

// Sampler converts a ray into a color.
//
// In the context of the pathtracer, the ray goes into a scene,
// intersects with geometry, then is shaded by a material.
type Sampler func(Ray) Color

// Material take a ray and returns its resulting color.
//
// Materials are encouraged to send rays back into scenes in order to
// hit a light source if necessary.
//
// Materials can also be lights, themselves.
type Material interface {
	Sample(from Ray, position Vector, normal Ray, uv UVCoordinate, nextSample Sampler) Color
}
