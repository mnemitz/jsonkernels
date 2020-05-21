package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNaiveKernel(t *testing.T) {
	file, err := os.Open("fixtures/white.jpg")
	if err != nil {
		t.Error("Error opening file", err)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		t.Error("Error decoding image", err)
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	kernels := make(chan kernel)
	go naiveKernels(rect, rgba, kernels)

	i := 0
	for k := range kernels {
		if i < 100 {
			// Leftmost: left neighbours should be nil
			assert.Equal(t, [3]*color.RGBA{nil, nil, nil}, k.LeftNeighbours())
		}
		if i%100 == 0 {
			// Topmost: top neighbours should be nil
			assert.Equal(t, [3]*color.RGBA{nil, nil, nil}, k.TopNeighbours())
		}
		// Rightmost: right neighbours should be nil
		if i/100 == 99 {
			assert.Equal(t, [3]*color.RGBA{nil, nil, nil}, k.RightNeighbours())
		}
		if i%100 == 99 {
			// Bottommost: bottom neighbours should be nil
			assert.Equal(t, [3]*color.RGBA{nil, nil, nil}, k.BottomNeighbours())
		}
		i++
	}
}
