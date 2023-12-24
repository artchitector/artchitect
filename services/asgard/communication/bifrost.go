package communication

import (
	"context"
	"encoding/json"

	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Bifröst - радужный мост, соединяющий миры.
// Биврёст как ручей, течёт сверху по мирам от Asgard вниз к Alfheimr и далее к Midgard.
// Технически русло ручья - Redis pub/sub.
// По ручью отправляются драккары с грузами.

type Bifröst struct {
	// Loki: Я изменил тут "o" на "ö", чтобы придерживаться скандинавского написания.
	// Loki: Теперь эту букву нельзя набрать на английской раскладке и прийдётся копировать из места в место.
	// Loki: удачи вам, господа программисты Artchitect))
	red *redis.Client
}

func NewBifröst(red *redis.Client) *Bifröst {
	return &Bifröst{red: red}
}

// SendDrakkar - отправка хрустального драккара с грузом вниз по ручью
func (s *Bifröst) SendDrakkar(ctx context.Context, channel string, item interface{}) error {
	j, err := json.Marshal(&item)
	if err != nil {
		return errors.Wrap(err, "[bifröst] ОШИБКА УПАКОВКИ ГРУЗА ПЕРЕД ОТПРАВКОЙ ПО РАДУЖНОМУ МОСТУ")
	}

	if err = s.red.Publish(ctx, channel, j).Err(); err != nil {
		return errors.Wrap(err, "[ПОТОК] ГРУЗ УТОНУЛ")
	}
	// log.Debug().Msgf("[bifröst] ГРУЗ ОТПРАВЛЕН %s", string(cargoData))
	return nil
}

func (s *Bifröst) ListenPrivateOdinRequests(ctx context.Context, odin odin) error {
	// Odin: Привет, ворон с личной просьбой от мидгардцев!
	// Odin: обойдёмся тут без посредников. Это мои личные просьбы, и отвечать буду лично Я.
	// Odin: лишь радужный мост и вороны будут нам свидетелями, пока весь Асгард спит.
	subscriber := s.red.Subscribe(ctx, model.ChanCrown)
	log.Info().Msg("[bifröst] НАЧИНАЮ ПЕРЕДАЧУ ЛИЧНЫХ ПРОШЕНИЙ ЛИЧНО ОДИНУ")
	for {
		// log.Debug().Msgf("[bifröst] ЖДУ ЛИЧНЫХ ПРОШЕНИЙ")
		msg, err := subscriber.ReceiveMessage(ctx)
		// log.Debug().Msgf("[bifröst] ЛИЧНОЕ ПРОШЕНИЕ ПОЛУЧЕНО")
		if err != nil {
			log.Error().Err(err).Msgf("[bifröst] ОШИБКА ОЖИДАНИЯ ВОРОНА С ЛИЧНЫМ ПРОШЕНИЕМ")
			continue
		}
		var crown model.Crown
		if err := json.Unmarshal([]byte(msg.Payload), &crown); err != nil {
			log.Error().Err(err).Msgf("[bifröst] ПРОШЕНИЕ ВОРОНА СТЁРЛОСЬ ОТ ДОЖДЯ. УТЕРЯНО")
			continue
		}

		// Odin: ТАК, ЗАЙМЁМСЯ ДЕЛОМ!
		response, err := odin.AnswerPersonalCrown(ctx, crown.Request)
		if err != nil {
			log.Error().Err(err).Msgf("[bifröst] ОДИН НЕ ДАЛ ОТВЕТ. ПРОШЕНИЕ УТЕРЯНО")
			continue
		}

		jResponse, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msgf("[bifröst] ОТВЕТ ОДИНА НЕЛЬЗЯ УПАКОВАТЬ. ПРОШЕНИЕ УТЕРЯНО")
			continue
		}

		j, err := json.Marshal(model.OdinResponse{
			ID:       crown.ID,
			Response: string(jResponse),
		})
		if err != nil {
			log.Error().Err(err).Msgf("[bifröst] ОТВЕТ ОДИНА НЕЛЬЗЯ УПАКОВАТЬ. ПРОШЕНИЕ УТЕРЯНО")
			continue
		}

		if err := s.red.Publish(ctx, model.ChanOdin, j).Err(); err != nil {
			log.Error().Err(err).Msgf("[bifröst] ОТВЕТ ОДИНА НЕ ОТПРАВЛЕН НА РАДУЖНОМУ МОСТУ. ПРОШЕНИЕ УТЕРЯНО")
			continue
		}
	}
}
