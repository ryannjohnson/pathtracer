# Outline

## Render(renderSettings, scene, camera, imageWriter)

* Iterate through image dimensions and get rays from camera
* Cast ray into scene => Color
* Save color to image

## Cast(scene, ray, bouncesLeft) Color

* Find the closest geometry in front of the ray (Hit)
* Material emits hit to a limited depth

## Hit{Ray, Geometry, Material, Position, Normal, UVCoordinate}

## Material.Emit(scene, hit, bouncesLeft) => Color

# Structs

* RenderSettings
  - bounceDepth
  - samplesPerPixel
* Camera
  - Expose(scene, x, y, bouncesLeft) Color
* Scene
  - Intersect(ray) *Hit
* Hit
  - Emit(scene, bouncesLeft) Color
* Geometry
  - Intersect(ray) *Hit
* Material
  - Emit(scene, ray, position, normal, UVCoordinate, bouncesLeft) Color
* Color [3]float64
* Ray [3]float64
* Matrix [16]float64
