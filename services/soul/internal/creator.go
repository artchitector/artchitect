package internal

import (
	"context"
	"errors"
	"github.com/artchitector/artchitect2/model"
	"github.com/rs/zerolog/log"
	"time"
)

type Creator struct {
	isActive bool
}

func NewCreator(isActive bool) *Creator {
	return &Creator{isActive: isActive}
}

func (c *Creator) IsActive() bool {
	return c.isActive
}

func (c *Creator) Create(ctx context.Context) (bool, model.Art, error) {
	log.Info().Msgf("[creator] starting fake art generation")

	select {
	case <-ctx.Done():
		log.Debug().Msgf("[creator] ctx.Done")
		return false, model.Art{}, nil
	case <-time.After(time.Second * 5):
		return false, model.Art{}, errors.New("[creator] fake error")
	}
}
