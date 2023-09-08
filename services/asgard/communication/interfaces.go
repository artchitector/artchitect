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
	GetArtRecursive(ctx context.Context, ID uint) (model.Art, error)
}

type warehouse interface {
	DownloadArtImage(ctx context.Context, artID uint, size string) ([]byte, error)
}

type muninn interface {
	RememberArtNo(ctx context.Context, min uint, max uint) (uint, model.EntropyPack, error)
}

type odin interface {
	AnswerPersonalCrown(ctx context.Context, crownRequest string) (interface{}, error)
}
