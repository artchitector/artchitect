package resizer

import (
	"image"
	"math"

	"github.com/artchitector/artchitect/model"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog/log"
)

func ResizeImage(img image.Image, size string) image.Image {
	var height, width uint

	switch size {
	case model.SizeF:
		width = model.WidthF
	case model.SizeM:
		width = model.WidthM
	case model.SizeS:
		width = model.WidthS
	case model.SizeXS:
		width = model.WidthXS
	default:
		log.Fatal().Msgf("[resizer] wrong size. crash")
	}

	height = uint(math.Round(float64(width) * model.HeightToWidth))
	img = resize.Resize(width, height, img, resize.Lanczos3)
	return img
}
