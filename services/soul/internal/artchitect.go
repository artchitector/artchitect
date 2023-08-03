package internal

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// Artchitect - главный управляющий в сервисе soul. Управляет Главным Циклом Архитектора (ГЦА)
type Artchitect struct {
	creator *Creator
}

func NewArtchitect(creator *Creator) *Artchitect {
	return &Artchitect{creator: creator}
}

func (a *Artchitect) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[artchitect] ГЦА - ВЫКЛ")
			return // stop loop

		case <-time.Tick(time.Millisecond):
			log.Debug().Msgf("[artchitect] ГЦА - ВКЛ")

			if err := a.WorkOnce(ctx); err != nil {
				log.Error().Err(err).Msgf("[artchitect] ГЦА - АВАРИЯ")
				continue
			} else {
				log.Debug().Msgf("[artchitect] ГЦА в работе...")
			}
		}
	}
}

// WorkOnce - обработка одной полезной "работы" или сон внутри ГЦА
// Работы расставлены по приоритету. Наименее приоритетное - рисование очередной картины.
// Другие операции не всегда доступны для работы и они идут первыми в списке,
// а операция рисования доступна всегда, она завершает список.
func (a *Artchitect) WorkOnce(ctx context.Context) error {
	if a.creator.IsActive() {
		if worked, art, err := a.creator.Create(ctx); err != nil {
			return errors.Wrap(err, "[artchitect] АВАРИЯ С СОЗДАТЕЛЕМ")
		} else if worked {
			log.Debug().Msgf("[artchitect] СОЗДАН АРТ #%d", art.ID)
			return nil
		}
	}

	log.Debug().Msg("[artchitect] нет работы. жду.")
	select {
	case <-ctx.Done():
		log.Debug().Msgf("[artchitect] ВЫКЛ")
		return nil
	case <-time.After(time.Second * 10):
		return nil
	}
}
