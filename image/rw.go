package image

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
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
	case "webp":
		return webp.Encode(w, img, &webp.Options{Quality: 100, Lossless: true})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
