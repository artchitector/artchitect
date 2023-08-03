package internal

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// Artchitect - main mediator in artchitect workflow. Handles main work loop.
type Artchitect struct {
	creator *Creator
}

func NewArtchitect(creator *Creator) *Artchitect {
	return &Artchitect{creator: creator}
}

// Run - infinite loop of work
func (a *Artchitect) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[artchitect] ctx.Done")
			return // stop loop

		case <-time.Tick(time.Second):
			log.Debug().Msgf("[artchitect] run main loop")

			if err := a.WorkOnce(ctx); err != nil {
				log.Error().Err(err).Msgf("[artchitect] my main loop crashed")
				continue
			} else {
				log.Debug().Msgf("[artchitect] my main loop success")
			}
		}
	}
}

// WorkOnce - take one task and do it
func (a *Artchitect) WorkOnce(ctx context.Context) error {
	if a.creator.IsActive() {
		if worked, art, err := a.creator.Create(ctx); err != nil {
			return errors.Wrap(err, "[artchitect] creator.Create failed")
		} else if worked {
			log.Debug().Msgf("[artchitect] created new Art #%d", art.ID)
			return nil
		}
	}

	log.Debug().Msg("[artchitect] nothing to do. will sleep 10 seconds")
	select {
	case <-ctx.Done():
		log.Debug().Msgf("[artchitect] workOnce ctx.Done")
		return nil
	case <-time.After(time.Second * 10):
		return nil
	}
}
