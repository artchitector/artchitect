package communication

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
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
func (s *Bifröst) SendDrakkar(ctx context.Context, cargo model.Cargo) error {
	var err error
	if err = s.red.Publish(ctx, cargo.Channel, cargo.Payload).Err(); err != nil {
		return errors.Wrap(err, "[ПОТОК] ГРУЗ УТОНУЛ")
	}
	//log.Debug().Msgf("[bifröst] ГРУЗ ОТПРАВЛЕН %s", string(cargoData))
	return nil
}

func (s *Bifröst) SendDrakkarWithPack(ctx context.Context, channel string, item interface{}) error {
	j, err := json.Marshal(&item)
	if err != nil {
		return errors.Wrap(err, "[bifröst] ОШИБКА УПАКОВКИ ГРУЗА ПЕРЕД ОТПРАВКОЙ ПО РАДУЖНОМУ МОСТУ")
	}
	cargo := model.Cargo{
		Channel: channel,
		Payload: string(j),
	}
	return s.SendDrakkar(ctx, cargo)
}
