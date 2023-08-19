package utils

import (
	"context"
	"github.com/artchitector/artchitect2/model"
)

// artPile - репозиторий картин
type artPile interface {
	GetMaxArtID(ctx context.Context) (uint, error)
	GetArt(ctx context.Context, ID uint) (model.Art, error)
}

type frigg interface {
	ReunifyArtUnities(ctx context.Context, art model.Art) error
}
