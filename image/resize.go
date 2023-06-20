package image

import (
	"fmt"
	"image"
	"log"
	"sync"

	"golang.org/x/image/draw"

	"github.com/kalogs-c/resizenator/semaphore"
)

type Algorithm string

const (
	BiLinear        Algorithm = "BiLinear"
	ApproxBiLinear  Algorithm = "ApproxBiLinear"
	CatmullRom      Algorithm = "CatmullRom"
	NearestNeighbor Algorithm = "NearestNeighbor"
)

func resizeImage(
	originalImage image.Image,
	desiredSize ImageSize,
	algorithm Algorithm,
) image.Image {
	originalImageBounds := originalImage.Bounds()
	scale := desiredSize.GetScales(originalImageBounds)

	destinationImageWidth := int(float64(originalImageBounds.Dx()) / scale)
	destinationImageHeight := int(float64(originalImageBounds.Dy()) / scale)

	destinationRect := image.Rect(0, 0, destinationImageWidth, destinationImageHeight)
	destinationImage := image.NewRGBA(destinationRect)

	switch algorithm {
	case "BiLinear":
		draw.BiLinear.Scale(
			destinationImage,
			destinationImage.Bounds(),
			originalImage,
			originalImageBounds,
			draw.Over,
			nil,
		)
	case "ApproxBiLinear":
		draw.ApproxBiLinear.Scale(
			destinationImage,
			destinationImage.Bounds(),
			originalImage,
			originalImageBounds,
			draw.Over,
			nil,
		)
	case "CatmullRom":
		draw.CatmullRom.Scale(
			destinationImage,
			destinationImage.Bounds(),
			originalImage,
			originalImageBounds,
			draw.Over,
			nil,
		)
	case "NearestNeighbor":
		draw.NearestNeighbor.Scale(
			destinationImage,
			destinationImage.Bounds(),
			originalImage,
			originalImageBounds,
			draw.Over,
			nil,
		)
	default:
		log.Fatalln("Algorithm not supported")
	}

	return destinationImage
}

func ResizeImageToChan(
	imageChan chan<- image.Image,
	semaphore semaphore.Semaphore,
	wg *sync.WaitGroup,
	originalImage image.Image,
	destinationImage ImageSize,
	algorithm Algorithm,
) {
	imageChan <- resizeImage(originalImage, destinationImage, algorithm)
	semaphore.Release()
	wg.Done()
}

func ResizeMonitor(imageChan chan image.Image, wg *sync.WaitGroup) {
	wg.Wait()
	close(imageChan)
	fmt.Println("Done")
}
