package main

import (
	"fmt"
	goimg "image"
	"log"
	"os"
	"sync"

	"github.com/kalogs-c/resizenator/image"
)

func main() {
	// Open the image file
	file, err := os.Open("./destination/image.jpg")
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()
	// Decode the image
	src, err := image.ReaderToImage(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Calculate the desired width and height for resizing
	desiredSizes := make([]image.ImageSize, 30)

	wg := &sync.WaitGroup{}
	wg.Add(len(desiredSizes))

	destinationChan := make(chan goimg.Image, 10)

	for i := range desiredSizes {
		fmt.Println("Resizing image...")

		size := image.ImageSize{
			X: 100 * i,
			Y: 100 * i,
		}
		go image.ResizeImageToChan(destinationChan, wg, src, size, "BiLinear")
	}

	go image.ResizeMonitor(destinationChan, wg)

	for i := range destinationChan {
		fmt.Println("Size: ", i.Bounds().Dx(), i.Bounds().Dy())
		destinationFile, err := os.Create(
			fmt.Sprintf("./destination/dest-%dx%d.jpg", i.Bounds().Dx(), i.Bounds().Dy()),
		)
		defer destinationFile.Close()
		if err != nil {
			log.Fatalf("Failed to create destination file: %v", err)
		}
		image.ImageToWriter(i, destinationFile, image.Jpeg)
	}
}
