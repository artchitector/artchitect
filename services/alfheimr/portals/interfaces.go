package portals

import (
	"context"
	"github.com/artchitector/artchitect2/model"
)

type radio interface {
	ListenRadio(subscribeCtx context.Context) chan model.Radiogram
}

type artPile interface {
	GetArtRecursive(ctx context.Context, ID uint) (model.Art, error)
	GetLastArts(ctx context.Context, last uint) ([]model.Art, error)
}

type warehouse interface {
	GetArtImage(ctx context.Context, artID uint, size string) ([]byte, error)
	GetArtOrigin(ctx context.Context, artID uint) ([]byte, error)
}
