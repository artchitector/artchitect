package external

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Bifrost - радужный мост, соединяющий миры.
// Биврёст как ручей, течёт сверху по мирам от Asgard вниз к Alfheimr и далее к Midgard.
// Техническое русло ручья - Redis pub/sub.

type Bifrost struct {
	red *redis.Client
}

func NewStream(red *redis.Client) *Bifrost {
	return &Bifrost{red: red}
}

// SendCargo - отправка груза вниз по ручью
func (s *Bifrost) SendCargo(ctx context.Context, event model.Event) error {
	var cargo []byte
	var err error
	if cargo, err = json.Marshal(event); err != nil {
		return errors.Wrap(err, "[ПОТОК] ГРУЗ НЕРАЗУМЕН")
	}
	if err = s.red.Publish(ctx, event.Channel, cargo).Err(); err != nil {
		return errors.Wrap(err, "[ПОТОК] ГРУЗ УТОНУЛ")
	}
	log.Debug().Msgf("[ПОТОК] ГРУЗ ОТПРАВЛЕН")
	return nil
}
