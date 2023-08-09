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
