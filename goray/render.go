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
			hit = orig.add(dir.mult(distance))
			N = hit.subtract(s.center).normalize()
			material = s.material
		}
	}

	return (spheresDist < 1000), material, hit, N
}

var background = Vec3f{0.2, 0.7, 0.8}

func castRay(orig, dir Vec3f, spheres []Sphere, lights []light, depth uint) Vec3f {
	if depth > 4 {
		return background // Background
	}

	intersect, material, point, N := sceneIntersect(orig, dir, spheres)
	if !intersect {
		return background // Background
	}

	Nscaled := N.mult(1e-3)
	reflectDir := dir.reflect(N).normalize()
	var reflectOrig Vec3f
	if reflectDir.dot(N) < 0 {
		reflectOrig = point.subtract(Nscaled)
	} else {
		reflectOrig = point.add(Nscaled)
	}
	reflectColor := castRay(reflectOrig, reflectDir, spheres, lights, depth+1)

	refractDir := dir.refract(N, material.refractiveIndex).normalize()
	var refractOrig Vec3f
	if refractDir.dot(N) < 0 {
		refractOrig = point.subtract(Nscaled)
	} else {
		refractOrig = point.add(Nscaled)
	}
	refractColor := castRay(refractOrig, refractDir, spheres, lights, depth+1)

	var diffuseLightIntensity, specularLightIntensity float64
	for _, l := range lights {
		lightDir := l.position.subtract(point).normalize()
		lightDistance := l.position.subtract(point).norm()

		var shadowOrig Vec3f
		if lightDir.dot(N) < 0 {
			shadowOrig = point.subtract(Nscaled)
		} else {
			shadowOrig = point.add(Nscaled)
		}

		shadowIntersect, shadMaterial, shadowPoint, shadowN := sceneIntersect(shadowOrig, lightDir, spheres)
		if shadowIntersect && shadowPoint.subtract(shadowOrig).norm() < lightDistance {
			continue
		}

		_ = shadMaterial
		_ = shadowPoint
		_ = shadowN

		diffuseLightIntensity += l.intensity * math.Max(0., lightDir.dot(N))
		reflectedLight := lightDir.mult(-1.).reflect(N).mult(-1.).dot(dir)
		intensity := math.Pow(math.Max(0., reflectedLight), material.specularExponent)
		specularLightIntensity += intensity * l.intensity
	}

	diffuseColor := material.diffuseColor.mult(diffuseLightIntensity).mult(material.albedo[0])

	onesVec := Vec3f{1., 1., 1.}
	specularColorComponent := onesVec.mult(specularLightIntensity).mult(material.albedo[1])
	reflectedColorComponent := reflectColor.mult(material.albedo[2])
	refractedColorComponent := refractColor.mult(material.albedo[3])

	return diffuseColor.add(specularColorComponent).add(reflectedColorComponent).add(refractedColorComponent)
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
			framebuffer[i+j*width] = castRay(orig, dir, spheres, lights, 0)
		}
	}

	start := time.Now()

	for i, v := range framebuffer {
		off := i * 4
		max := math.Max(v.x, math.Max(v.y, v.z))
		if max > 1. {
			v = v.mult(1. / max)
		}
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
	ivory := Material{1.0, [4]float64{0.6, 0.3, 0.1, 0.0}, Vec3f{0.4, 0.4, 0.3}, 50.}
	glass := Material{1.5, [4]float64{0.0, 0.5, 0.1, 0.8}, Vec3f{0.6, 0.7, 0.8}, 125.}
	redRubber := Material{1.0, [4]float64{0.9, 0.1, 0.0, 0.0}, Vec3f{0.3, 0.1, 0.1}, 10.}
	mirror := Material{1.0, [4]float64{0.0, 10.0, 0.8, 0.0}, Vec3f{1.0, 1.0, 1.0}, 1425.}

	var spheres []Sphere
	spheres = append(spheres, Sphere{Vec3f{-3., 0., -16.}, 2., ivory})
	spheres = append(spheres, Sphere{Vec3f{-1., -1.5, -12.}, 2., glass})
	spheres = append(spheres, Sphere{Vec3f{1.5, -0.5, -18.}, 3., redRubber})
	spheres = append(spheres, Sphere{Vec3f{7., 5., -18.}, 4., mirror})

	var lights []light
	lights = append(lights, light{Vec3f{-20., 20., 20.}, 1.5})
	lights = append(lights, light{Vec3f{30., 50., -25.}, 1.8})
	lights = append(lights, light{Vec3f{30., 20., 30.}, 1.7})

	renderStart := time.Now()
	render(spheres, lights)
	renderElapsed := time.Since(renderStart)
	fmt.Printf("Render took %s\n", renderElapsed)
}
