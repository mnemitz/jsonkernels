package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println(fmt.Errorf("Input file name required as argument"))
		os.Exit(-1)
	}
	filePath := flag.Arg(0)
	file, err := os.Open(filePath)
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
