package goray

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func Render() {
	const width = 1024
	const height = 768
	im := image.NewRGBA(image.Rect(0, 0, width, height))

	outputFile, err := os.Create("out.png")
	if err != nil {
		fmt.Printf("Can't open output file")
		return
	}

	framebuffer := [height * width]Vec3f{}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			framebuffer[i+j*width] = Vec3f{float32(j) / float32(height), float32(i) / float32(width), 0}
		}
	}

	for i, v := range framebuffer {
		off := i * 4
		im.Pix[off] = uint8(v.X * 255.)
		im.Pix[off+1] = uint8(v.Y * 255.)
		im.Pix[off+2] = uint8(v.Z * 255.)
		im.Pix[off+3] = 255
	}

	png.Encode(outputFile, im)
	outputFile.Close()
}
