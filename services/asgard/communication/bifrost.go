package communication

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Bifröst - радужный мост, соединяющий миры.
// Биврёст как ручей, течёт сверху по мирам от Asgard вниз к Alfheimr и далее к Midgard.
// Технически русло ручья - Redis pub/sub.
// По ручью отправляются драккары с грузами.

type Bifröst struct {
	red *redis.Client
}

func NewStream(red *redis.Client) *Bifröst {
	return &Bifröst{red: red}
}

// SendDrakkar - отправка хрустального драккара с грузом вниз по ручью
func (s *Bifröst) SendDrakkar(ctx context.Context, cargo model.Cargo) error {
	var cargoData []byte
	var err error
	if err = s.red.Publish(ctx, cargo.Channel, cargo.Payload).Err(); err != nil {
		return errors.Wrap(err, "[ПОТОК] ГРУЗ УТОНУЛ")
	}
	log.Debug().Msgf("[bifröst] ГРУЗ ОТПРАВЛЕН %s", string(cargoData))
	return nil
}
