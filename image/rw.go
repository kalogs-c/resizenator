package image

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

type ImageFormat string

const (
	Jpeg ImageFormat = "jpeg"
	Png  ImageFormat = "png"
	Gif  ImageFormat = "gif"
)

func ReaderToImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ImageToWriter(img image.Image, w io.Writer, format ImageFormat) error {
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
