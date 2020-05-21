package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	// Get input image file
	// Iterate through pixels
	// Get kernel for each pixel (corners will be incomplete)
	file, err := os.Open("test.jpeg")
	if err != nil {
		panic(err)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}
	rect := img.Bounds()
	// Empty grid where the pixels follow 8-bit RGBA color model
	// saves us from having to explicitly convert to 8-bit representation
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	kernels := make(chan kernel)
	go naiveKernels(rect, rgba, kernels)

	for k := range kernels {
		verboseJSON, err := k.toVerboseJSON()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(verboseJSON))
	}
}
