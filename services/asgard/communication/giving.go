package communication

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// TODO ОТКОММЕНТИРОВАТЬ ЭТО БЕЗОБРАЗИЕ

const (
	StepInterval   = 1 // отправка каждую 1 секунду
	ChangeInterval = 3 // смена одной из картинок каждые 3 секунды
	GivenArtsLen   = 4 // 4 картинки меняются по кругу (каждые ChangeInterval меняет одну картинку)
)

type Giving struct {
	artPile          artPile
	muninn           muninn
	bifröst          *Bifröst
	lastChange       time.Time
	lastChangedIndex int
}

func NewGiving(artPile artPile, muninn muninn, bifröst *Bifröst) *Giving {
	return &Giving{artPile: artPile, muninn: muninn, bifröst: bifröst}
}

func (g *Giving) StartGiving(ctx context.Context) {
	state := &model.GivingState{Given: make([]uint, GivenArtsLen)}
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[giving] ОСТАНОВ")
			return
		case <-time.Tick(time.Second * time.Duration(StepInterval)):
			if err := g.step(ctx, state); err != nil {
				log.Error().Err(err).Msgf("[giving] ОШИБКА ШАГА")
			}
		}
	}
}

func (g *Giving) step(ctx context.Context, state *model.GivingState) error {
	if state.LastArtID == 0 || len(state.Given) == 0 {
		if err := g.initState(ctx, state); err != nil {
			return errors.Wrap(err, "[giving] ОШИБКА НАПОЛНЕНИЯ СОСТОЯНИЯ")
		}

		err := g.bifröst.SendDrakkarWithPack(ctx, model.ChanGiving, state)
		return err
	}

	if err := g.updateState(ctx, state); err != nil {
		return errors.Wrap(err, "[giving] ОШИБКА ОБНОВЛЕНИЯ СОСТОЯНИЯ")
	}
	err := g.bifröst.SendDrakkarWithPack(ctx, model.ChanGiving, state)
	return err
}

func (g *Giving) initState(ctx context.Context, state *model.GivingState) error {
	maxArtID, err := g.artPile.GetMaxArtID(ctx)
	if err != nil {
		return errors.Wrap(err, "[giving] ОШИБКА ЗНАНИЯ ПОСЛЕДНЕЙ КАРТИНЫ")
	}
	state.LastArtID = maxArtID
	for i := 0; i < GivenArtsLen; i++ {
		oneOf, _, err := g.muninn.OneOf(ctx, maxArtID)
		if err != nil {
			return errors.Wrap(err, "[giving] МУНИН НЕ ДАЛ ЯСНОГО ОТВЕТА")
		}
		state.Given[i] = oneOf + 1
		g.lastChangedIndex = i
	}
	g.lastChange = time.Now()
	return nil
}

func (g *Giving) updateState(ctx context.Context, state *model.GivingState) error {
	maxArtID, err := g.artPile.GetMaxArtID(ctx)
	if err != nil {
		return errors.Wrap(err, "[giving] ОШИБКА ЗНАНИЯ ПОСЛЕДНЕЙ КАРТИНЫ")
	}
	state.LastArtID = maxArtID
	changeInterval := time.Second * time.Duration(ChangeInterval)
	if g.lastChange.Add(changeInterval).Before(time.Now()) {
		// пора менять картинку
		currentIndex := g.lastChangedIndex + 1
		if currentIndex >= GivenArtsLen {
			currentIndex = 0
		}
		oneOf, _, err := g.muninn.OneOf(ctx, maxArtID)
		if err != nil {
			return errors.Wrap(err, "[giving] МУНИН НЕ ДАЛ ЯСНОГО ОТВЕТА")
		}
		state.Given[currentIndex] = oneOf + 1
		g.lastChangedIndex = currentIndex
		g.lastChange = time.Now()
	}

	return nil
}
