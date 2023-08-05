package external

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type subscriber struct {
	ctx     context.Context
	eventCh chan model.Event
}

// Harbour - гавань, принимающая полезные грузы от верхнего мира (soul), который выше по течению.
// В гавани принимают грузы, сортируют - некоторые отправляют конечным получателям, некоторые обрабатывают сами
type Harbour struct {
	mutex       sync.Mutex
	red         *redis.Client
	subscribers []*subscriber
}

func NewListener(red *redis.Client) *Harbour {
	return &Harbour{red: red, mutex: sync.Mutex{}, subscribers: make([]*subscriber, 0)}
}

func (l *Harbour) Run(ctx context.Context) error {
	subscriber := l.red.Subscribe(
		ctx,
		model.ChanEntropy,
		model.ChanEntropyCalculation,
	)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[слушатель] ОСТАНОВ")
			return nil
		default:
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[слушатель] ОШИБКА ПРИЁМА")
				time.Sleep(time.Second)
				continue
			}
			if err := l.handle(ctx, msg); err != nil {
				log.Error().Err(err).Msgf("[слушатель] ОШИБКА ОБРАБОТКИ")
			}
		}
	}
}

// handle - обработка одного пришедшего сообщения
func (l *Harbour) handle(ctx context.Context, msg *redis.Message) error {
	broadcastChannels := []string{
		model.ChanEntropy,
		model.ChanEntropyCalculation,
	}
	idx := slices.IndexFunc(broadcastChannels, func(s string) bool { return msg.Channel == s })
	if idx != -1 {
		// нужно отправить это сообщение броадкастом всем подписанным клиентам, которые сейчас слушают websocket
		l.broadcast(ctx, msg)
	}

	return nil
}

// broadcast - веерная рассылка всем подключенным и заинтересованным клиентам
func (l *Harbour) broadcast(ctx context.Context, msg *redis.Message) {
	event := model.Event{
		Channel: msg.Channel,
		Payload: msg.Payload,
	}
	l.mutex.Lock()
	subscribers := l.subscribers[:]
	l.mutex.Unlock()

	for _, sub := range subscribers {
		// отправка сообщения каждому подписчику, если он ещё активен
		go func(s *subscriber) {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
				log.Error().Msgf("[слушатель] ОТПРАВКА ЗАВИСЛА, ГРУЗ ПОТЕРЯН. КАНАЛ:%s", msg.Channel)
			case s.eventCh <- event:
				log.Debug().Msgf("[слушатель] ГРУЗ ОТПРАВЛЕН. КАНАЛ:%s. УСПЕХ", msg.Channel)
			}
		}(sub)
	}
}

// Subscribe - подписка слушателя на канал событий. Передаются все события, подлежащие веерной рассылке (без фильтрации)
func (l *Harbour) Subscribe(subscribeCtx context.Context) chan model.Event {
	eventCh := make(chan model.Event)
	sub := subscriber{
		// с помощью этого контекста канал подписчик будет остановлен и удалён
		ctx:     subscribeCtx,
		eventCh: eventCh,
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.subscribers = append(l.subscribers, &sub)
	go func() {
		<-subscribeCtx.Done()
		l.unsubscribe(&sub)
	}()

	return eventCh
}

func (l *Harbour) unsubscribe(sub *subscriber) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	idx := slices.IndexFunc(l.subscribers, func(s *subscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[слушатель] ПОДПИСАНТ ОТСУТСТВУЕТ. НЕ УДАЛИТЬ")
		return
	}

	l.subscribers = append(l.subscribers[:idx], l.subscribers[idx+1:]...)
	log.Debug().Msgf("[слушатель] ПОДПИСАНТ #d - УДАЛЕНО", idx)
}
