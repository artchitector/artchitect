package pantheon

import (
	"context"
)

// artPile - Куча написанных картин (репозиторий для таблицы art)
type artPile interface {
	GetNextArtID(ctx context.Context) (uint, error)
}

type ai interface {
	GenerateImage(ctx context.Context, seed uint, prompt string) ([]byte, error)
}
