package pathtracer

// Material take a ray and returns its resulting color.
//
// Materials are encouraged to send rays back into scenes in order to
// hit a light source if necessary.
//
// Materials can also be lights, themselves.
type Material interface {
	Emit(scene Scene, from Ray, point Point, normal Ray, uvCoordinate Point, bouncesLeft int) Color
}
