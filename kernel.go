package main

import (
	"encoding/json"
	"image"
	"image/color"
)

type kernel struct {
	Value      color.RGBA
	Neighbours [8]*color.RGBA
}

func (k kernel) TopNeighbours() [3]*color.RGBA {
	return [3]*color.RGBA{
		k.Neighbours[0],
		k.Neighbours[1],
		k.Neighbours[2],
	}
}

func (k kernel) BottomNeighbours() [3]*color.RGBA {
	return [3]*color.RGBA{
		k.Neighbours[5],
		k.Neighbours[6],
		k.Neighbours[7],
	}
}

func (k kernel) LeftNeighbours() [3]*color.RGBA {
	return [3]*color.RGBA{
		k.Neighbours[0],
		k.Neighbours[3],
		k.Neighbours[5],
	}
}

func (k kernel) RightNeighbours() [3]*color.RGBA {
	return [3]*color.RGBA{
		k.Neighbours[2],
		k.Neighbours[4],
		k.Neighbours[7],
	}
}

func (k kernel) toVerboseJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value color.RGBA
		N0,
		N1,
		N2,
		N3,
		N4,
		N5,
		N6,
		N7 *color.RGBA `json:",omitempty"`
	}{
		k.Value,
		k.Neighbours[0],
		k.Neighbours[1],
		k.Neighbours[2],
		k.Neighbours[3],
		k.Neighbours[4],
		k.Neighbours[5],
		k.Neighbours[6],
		k.Neighbours[7],
	})
}

func getNthNeighbourPosition(x, y, n int) (int, int) {
	offsX := kernelOffsetOrder[n][0]
	offsY := kernelOffsetOrder[n][1]
	return x + offsX, y + offsY
}

var kernelOffsetOrder = [8][2]int{
	[2]int{-1, -1},
	[2]int{0, -1},
	[2]int{1, -1},
	[2]int{-1, 0},
	[2]int{1, 0},
	[2]int{-1, 1},
	[2]int{0, 1},
	[2]int{1, 1},
}

func naiveKernels(rect image.Rectangle, rgba *image.RGBA, kernels chan kernel) {
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			k := kernel{
				Value:      rgba.RGBAAt(i, j),
				Neighbours: [8]*color.RGBA{},
			}
			for n := range k.Neighbours {
				xOffs, yOffs := getNthNeighbourPosition(i, j, n)
				// if out of bounds, leave it nil
				inBounds := xOffs >= rect.Min.X &&
					xOffs < rect.Max.X &&
					yOffs >= rect.Min.Y &&
					yOffs < rect.Max.Y
				if !inBounds {
					continue
				}
				nVal := rgba.RGBAAt(xOffs, yOffs)
				k.Neighbours[n] = &nVal
			}
			kernels <- k
		}
	}
	close(kernels)
}
