package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"golang.org/x/image/draw"
)

func main() {
	// Open the image file
	file, err := os.Open("./destination/download.jpeg")
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Decode the image
	src, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Calculate the desired width and height for resizing
	width := 100
	height := 100

	// Resize the image
	resizedImage := ResizeImage(src, width, height)

	// Create the output file
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the resized image and save it to the output file
	err = WriteImage(resizedImage, outFile, "png")
	if err != nil {
		log.Fatalf("Failed to encode and save image: %v", err)
	}

	fmt.Println("Image resized and saved successfully!")
}

func ReaderToImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func WriteImage(img image.Image, w io.Writer, format string) error {
	switch format {
	case "jpeg":
		return jpeg.Encode(w, img, nil)
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func ResizeImage(originalImage image.Image, desiredWidth int, desiredHeight int) image.Image {
	originalImageBounds := originalImage.Bounds()
	xScale := float64(originalImageBounds.Dx()) / float64(desiredWidth)
	yScale := float64(originalImageBounds.Dy()) / float64(desiredHeight)

	scale := xScale
	if yScale > xScale {
		scale = yScale
	}

	destinationImageWidth := int(float64(originalImageBounds.Dx()) / scale)
	destinationImageHeight := int(float64(originalImageBounds.Dy()) / scale)
	destinationRect := image.Rect(0, 0, destinationImageWidth, destinationImageHeight)
	destinationImage := image.NewRGBA(destinationRect)

	draw.NearestNeighbor.Scale(
		destinationImage,
		destinationImage.Bounds(),
		originalImage,
		originalImageBounds,
		draw.Over,
		nil,
	)

	return destinationImage
}
