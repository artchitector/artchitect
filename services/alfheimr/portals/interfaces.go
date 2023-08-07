package portals

import (
	"context"
	"github.com/artchitector/artchitect2/model"
)

type radio interface {
	ListenRadio(subscribeCtx context.Context) chan model.Radiogram
}
