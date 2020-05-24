package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	json2vecs "github.com/mnemitz/go-json2vecs-sdk"
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

	apiClient := json2vecs.Client{
		URL: "http://127.0.0.1:8000",
	}
	fromKeys := []string{"color"}
	toKeys := []string{
		"N0",
		"N1",
		"N2",
		"N3",
		"N4",
		"N5",
		"N6",
		"N7",
	}
	batch := make([]interface{}, 10000)
	i := 0
	for k := range kernels {
		verboseKernel := k.toVerboseKernel()
		batch[i] = *verboseKernel
		if i == len(batch)-1 {
			resp, err := apiClient.CountEntities(batch, fromKeys, toKeys)
			if err != nil {
				panic(err)
			}
			if resp.StatusCode >= 300 {
				fmt.Println("Error from API:", resp)
				os.Exit(-1)
			}
			i = 0
			resp.Body.Close()
			continue
		}
		i++
	}
	fmt.Printf("Send %d kernels for processing!\n", i)
	fmt.Println("Generating embeddings...")
	_, err = apiClient.Embed("color", "N0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created embeddings from color to N0, querying...")

	resp, err := apiClient.Query("color", "N0", "4294967295")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
	body := buf.Bytes()
	fmt.Println(string(body))
}
