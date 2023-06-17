package image

import "strings"

func IsImage(filename string) bool {
	format := GetImageFormat(filename)

	return format == "jpg" || format == "jpeg" || format == "png" || format == "gif"
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
	default:
		return ""
	}
}
