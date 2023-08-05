package external

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Stream течёт сверху от soul вниз к gate, как ручей. Русло ручья - Redis pub/sub.
type Stream struct {
	red *redis.Client
}

func NewStream(red *redis.Client) *Stream {
	return &Stream{red: red}
}

// SendCargo - отправка груза вниз по ручью
func (s *Stream) SendCargo(ctx context.Context, event model.Event) error {
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
