package goray

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"time"
)

type light struct {
	position  Vec3f
	intensity float64
}

func sceneIntersect(orig, dir Vec3f, spheres []Sphere) (bool, Material, Vec3f, Vec3f) {
	spheresDist := math.MaxFloat64
	var hit, N Vec3f
	var material Material

	for _, s := range spheres {
		sphereIntersect, distance := s.rayIntersect(orig, dir)
		if sphereIntersect && distance < spheresDist {
			spheresDist = distance
			hit = vAdd(orig, vMult(dir, distance))
			N = vSubtract(hit, s.center).normalize()
			material = s.material
		}
	}

	return (spheresDist < 1000), material, hit, N
}

var background = Vec3f{0.2, 0.7, 0.8}

func castRay(orig, dir Vec3f, spheres []Sphere, lights []light) Vec3f {
	intersect, material, point, N := sceneIntersect(orig, dir, spheres)
	if !intersect {
		return background // Background
	}

	var diffuseLightIntensity float64
	for _, l := range lights {
		lightDir := vSubtract(l.position, point).normalize()
		diffuseLightIntensity += l.intensity * math.Max(0., vDot(lightDir, N))
	}

	return vMult(material.diffuseColor, diffuseLightIntensity)
}

func render(spheres []Sphere, lights []light) {
	const width = 1024
	const height = 768
	const fov = math.Pi / 3.

	im := image.NewRGBA(image.Rect(0, 0, width, height))

	outputFile, err := os.Create("out.png")
	if err != nil {
		fmt.Printf("Can't open output file")
		return
	}

	framebuffer := [height * width]Vec3f{}

	orig := Vec3f{0, 0, 0}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			dirX := (float64(i) + 0.5) - float64(width)/2.
			dirY := -(float64(j) + 0.5) + float64(height)/2.
			dirZ := -height / (2. * math.Tan(fov/2.))
			dir := Vec3f{dirX, dirY, dirZ}.normalize()
			framebuffer[i+j*width] = castRay(orig, dir, spheres, lights)
		}
	}

	start := time.Now()

	for i, v := range framebuffer {
		off := i * 4
		im.Pix[off] = uint8(v.x * 255.)
		im.Pix[off+1] = uint8(v.y * 255.)
		im.Pix[off+2] = uint8(v.z * 255.)
		im.Pix[off+3] = 255
	}

	dataEnd := time.Now()
	dataElapsed := time.Since(start)

	png.Encode(outputFile, im)
	outputFile.Close()

	saveTime := time.Since(dataEnd)
	fmt.Printf("Data time: %s\t save time: %s\n", dataElapsed, saveTime)
}

func Scene() {
	ivory := Material{Vec3f{0.4, 0.4, 0.3}}
	redRubber := Material{Vec3f{0.3, 0.1, 0.1}}

	var spheres []Sphere
	spheres = append(spheres, Sphere{Vec3f{-3., 0., -16.}, 2., ivory})
	spheres = append(spheres, Sphere{Vec3f{-1., -1.5, -12.}, 2., redRubber})
	spheres = append(spheres, Sphere{Vec3f{1.5, -0.5, -18.}, 3., redRubber})
	spheres = append(spheres, Sphere{Vec3f{7., 5., -18.}, 4., ivory})

	var lights []light
	lights = append(lights, light{Vec3f{-20., 20., 20.}, 1.5})

	renderStart := time.Now()
	render(spheres, lights)
	renderElapsed := time.Since(renderStart)
	fmt.Printf("Render took %s\n", renderElapsed)
}
