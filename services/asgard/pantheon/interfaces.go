package pantheon

import (
	"context"
	"image"

	"github.com/artchitector/artchitect/model"
)

// artPile - Куча написанных картин (репозиторий для таблицы art)
type artPile interface {
	GetNextArtID(ctx context.Context) (uint, error)
	GetMaxArtID(ctx context.Context) (uint, error)
	SaveArt(ctx context.Context, artID uint, art model.Art, idea model.Idea) (model.Art, error)
	GetLastPaintTime(ctx context.Context) (uint, error) // в миллисекундах
}

// unityPile - репозиторий единств
type unityPile interface {
	Get(ctx context.Context, mask string) (model.Unity, error)
	Create(ctx context.Context, mask, state string, rank, min, max uint) (model.Unity, error)
	Save(ctx context.Context, unity model.Unity) (model.Unity, error)
	GetNextUnityForReunification(ctx context.Context) (model.Unity, error)
	GetChildren(ctx context.Context, unity model.Unity) ([]model.Unity, error)
}

type ai interface {
	GenerateImage(ctx context.Context, seed uint, prompt string) ([]byte, error)
}

type bifröst interface {
	SendDrakkar(ctx context.Context, channel string, item interface{}) error
}

type warehouse interface {
	SaveArtImage(ctx context.Context, artID uint, img image.Image) error
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}
