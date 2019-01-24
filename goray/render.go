package goray

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

func castRay(orig, dir Vec3f, sphere Sphere) Vec3f {
	intersect := sphere.rayIntersect(orig, dir)
	if !intersect {
		return Vec3f{0.2, 0.7, 0.8} // Background
	}

	return Vec3f{0.4, 0.4, 0.3}
}

func render(sphere Sphere) {
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
			framebuffer[i+j*width] = castRay(orig, dir, sphere)
		}
	}

	for i, v := range framebuffer {
		off := i * 4
		im.Pix[off] = uint8(v.x * 255.)
		im.Pix[off+1] = uint8(v.y * 255.)
		im.Pix[off+2] = uint8(v.z * 255.)
		im.Pix[off+3] = 255
	}

	png.Encode(outputFile, im)
	outputFile.Close()
}

func Scene() {
	s := Sphere{Vec3f{-3., 0., -16.}, 2.}
	render(s)
}
