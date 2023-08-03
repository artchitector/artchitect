package internal

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

// Creator - создаёт новый арт.
type Creator struct {
	isActive bool
	artist   artist
	artRepo  artRepo
}

func NewCreator(
	isActive bool,
	artist artist,
	artRepo artRepo,
) *Creator {
	return &Creator{
		isActive: isActive,
		artist:   artist,
		artRepo:  artRepo,
	}
}

// IsActive - ВКЛ/ВЫКЛ
func (c *Creator) IsActive() bool {
	return c.isActive
}

// Create - создание одного арта с нуля. Creator соберёт все нужные параметры и запросит AI нарисовать картинку
func (c *Creator) Create(ctx context.Context) (worked bool, art model.Art, err error) {
	select {
	case <-ctx.Done():
		log.Debug().Msgf("[creator] ВЫКЛ")
		return worked, art, err
	case <-time.After(time.Second * 1):
		art, err = c.create(ctx)
		if err != nil {
			return false, model.Art{}, err
		} else {
			return true, art, err
		}
	}
}

func (c *Creator) create(ctx context.Context) (art model.Art, err error) {
	artID, err := c.artRepo.GetNextArtID(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("[creator] НЕВОЗМОЖНО ПОЛУЧИТЬ НОМЕР НОВОЙ КАРТИНЫ")
	}

	log.Info().Msgf("[creator] КАРТИНА #%d. НАЧАЛО.", artID)

	i := rand.Intn(10000)
	_, err = c.artist.MakeArt(ctx, artID, uint(i), []string{"purple", "by kinkade", "elemental", "sattelite", "sauron", "age"})
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[creator] ТРЕВОГА! КАРТИНА НЕ СОЗДАНА!")
	}

	return model.Art{}, nil // TODO продолжить тут
}
