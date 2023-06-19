package image

import (
	"strings"
)

type ImageFormat string

const (
	Jpeg ImageFormat = "jpeg"
	Png  ImageFormat = "png"
	Gif  ImageFormat = "gif"
	Webp ImageFormat = "webp"
)

func IsImage(filename string) bool {
	format := GetImageFormat(filename)

	return format == "jpg" || format == "jpeg" || format == "png" || format == "gif" ||
		format == "webp"
}

func GetImageFormat(filename string) ImageFormat {
	format := filename[strings.LastIndex(filename, ".")+1:]

	switch format {
	case "jpg":
		return Jpeg
	case "jpeg":
		return Jpeg
	case "png":
		return Png
	case "gif":
		return Gif
	case "webp":
		return Webp
	default:
		return ""
	}
}
