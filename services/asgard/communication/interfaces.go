package communication

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"image"
)

type webcam interface {
	GetNextFrame(ctx context.Context) (image.Image, error)
}

type artPile interface {
	GetMaxArtID(ctx context.Context) (uint, error)
}

type muninn interface {
	OneOf(ctx context.Context, maxval uint) (uint, model.EntropyPack, error)
}
