package utils

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// UnityInitializer - временная служба для инициализации единств, которые уже в прошлом.
type UnityInitializer struct {
	artPile artPile
	frigg   frigg
}

func NewUnityInitializer(artPile artPile, frigg frigg) *UnityInitializer {
	return &UnityInitializer{artPile: artPile, frigg: frigg}
}

func (ui *UnityInitializer) Init(ctx context.Context) error {
	maxArtID, err := ui.artPile.GetMaxArtID(ctx)
	if err != nil {
		return errors.Wrap(err, "[unity_initializer] НЕ ПОЛУЧЕН MAX_ART_ID")
	}
	var id uint
	for id = 99; id < maxArtID; id += 100 {
		log.Info().Msgf("[unity_initializer] РАБОТАЮ С КАРТИНОЙ #%d", id)
		art, err := ui.artPile.GetArt(ctx, id)
		if err != nil {
			return errors.Wrapf(err, "[unity_initializer] ОШИБКА ДОСТУПА К КАРТИНЕ #%d", id)
		}
		err = ui.frigg.ReunifyArtUnities(ctx, art)
		if err != nil {
			return errors.Wrapf(err, "[unity_initializer] ОШИБКА СОЗДАНИЯ МНОЖНСТВ К КАРТИНЕ #%d", id)
		}
	}

	return nil
}
