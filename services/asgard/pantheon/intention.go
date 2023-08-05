package pantheon

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

/*
Intention или Стремление - это намерение Одина что-то сделать в Асгарде.

	Он может нарисовать картину, или потребовать произвести какую-то другую работу.
	Стремление запускается Главный Цикл Творения (ГЦТ), который по кругу заставляет творить картины
	и выполнять другие события из жизни Artchitect.
*/
type Intention struct {
	odin *Odin // Один-Всеотец рисует картины

}

func NewArtchitect(odin *Odin) *Intention {
	return &Intention{odin: odin}
}

func (a *Intention) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[СТРЕМЛЕНИЕ] ГЦТ - ОСТАНОВ")
			return // stop loop

		case <-time.Tick(time.Millisecond):
			if err := a.WorkOnce(ctx); err != nil {
				log.Error().Err(err).Msgf("[СТРЕМЛЕНИЕ] ГЦТ - СБОЙ")
				continue
			}
		}
	}
}

// WorkOnce - обработка одной полезной "работы" или сон внутри ГЦА
// Работы расставлены по приоритету. Наименее приоритетное - рисование очередной картины.
// Другие операции не всегда доступны для работы и они идут первыми в списке,
// а операция рисования доступна всегда, она завершает список.
func (a *Intention) WorkOnce(ctx context.Context) error {
	if a.odin.HasDesire() {
		if worked, art, err := a.odin.Create(ctx); err != nil {
			return errors.Wrap(err, "[СТРЕМЛЕНИЕ] ОДИН НЕ СОЗДАЛ КАРТИНУ. ОДИН В ЯРОСТИ")
		} else if worked {
			log.Debug().Msgf("[СТРЕМЛЕНИЕ] ОДИН СОЗДАЛ КАРТИНУ #%d", art.ID)
			return nil
		}
	}

	log.Debug().Msg("[СТРЕМЛЕНИЕ] НЕТ РАБОТЫ. АСГАРД ОТДЫХАЕТ.")
	select {
	case <-ctx.Done():
		log.Debug().Msgf("[СТРЕМЛЕНИЕ] ОСТАНОВ")
		return nil
	case <-time.After(time.Second * 10):
		// пропускает 10 секунд, если работы нет (такое бывает локалько или при сбоях, Odin любит рисовать без перерывов)
		return nil
	}
}
