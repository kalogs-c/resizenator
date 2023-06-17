package image

import "image"

type ImageSize struct {
	X int
	Y int
}

func (s *ImageSize) GetScales(originalImageBounds image.Rectangle) float64 {
	xScale := float64(originalImageBounds.Dx()) / float64(s.X)
	yScale := float64(originalImageBounds.Dy()) / float64(s.Y)

	scale := xScale
	if yScale > xScale {
		scale = yScale
	}

	return scale
}
