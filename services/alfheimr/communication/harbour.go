package communication

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type subscriber struct {
	ctx     context.Context
	eventCh chan model.Radiogram
}

/*
Harbour - гавань, принимающая полезные грузы от верхнего мира Asgard, который выше по течению радужного моста.
В гавани принимают грузы, сортируют - некоторые отправляют конечным получателям, некоторые обрабатывают (внутри Альфхейма)

Odin: светлые эльфы Альфхейма придумали поставить радиомачту в гавани,
Odin: а все грузы оцифровать и передавать в видео РАДИОГРАММ.
Loki: а до этого момента грузы на деревянных драккарах плыли? Зачем цифровать цифровое? Я не вижу драккара, вижу redis-пакет.
Odin: Локи, ты хоть и бог, но сейчас решил мыслить совсем узколобо.
Odin: Мир бесконечен и имеет много уровней проявления одного явления.
Есть технический аспект (алгоритмы, код, данные, redis-message), и этот аспект, кстати, тоже нефизический, а абстрактный. Код - просто буквы.
Есть физический аспект (как двигаются электроны в интегральных схемах и проводах, зажигаются пиксели на экране).
А есть аспект на уровне абстрактной идеи Artchitect, сказочной в нашем случае, и тут эльфы отправляют радиограммы.
Какой-то светлый эльф сидит в радио-студии с чашечкой кофе, и рассказывает слушателям, как выглядит очередная картина,
сотворённая в Асгарде. На уровне сказочной идеи Artchitect я могу как угодно называть эти явления, хоть драккарами,
хоть пакетами, грузами, радиограммами, радужными мостами, и вешать на них любые ярлыки из мира относительного ума.
Я придумываю, как это выглядит на уровне идеи, а код подстроится под эту картинку.
Loki: твоя мудрость не знает границ, Всеотец. Но как же программисты, которые потом это будут поддерживать? Ты о них подумал?
Odin: Вряд ли им будет легко, но я постараюсь всю свою задумку комментировать
Odin: поймут содержание и суть сказки - тогда поймут и всю архитектуру проекта.
Odin: представь себя на месте программиста из будущего, который вместо сухой документации может прочитать забавный рассказ с
картинками за 10-20 минут, и после этого он уже поймёт и суть проекта, и как в нём ориентироваться - где искать все эти мосты и драккары.

Odin: оцифрованные грузы эльфы через радио рассылают всем
Odin: это нужно для направления потока данных внутри этой программы
Odin: golang настолько примитивен по сравнению с технологиями Асгарда...
Loki: Всеотец, с тобой довольно душно в этой программе находиться)
Odin: с тобой мы еще не закончили наш спор.

Внутренний подписчик может настроиться на радио, и так слушать данные о грузах в гавани
Odin: некоторые грузы не будут переданы в мидгард, а будут использованы локально внутри Альфхейма
Odin: некоторые грузы будут ретранслированы по внешнему радио (websocket), которое доступно из Мидгарда (браузера)
*/
type Harbour struct {
	mutex    sync.Mutex
	red      *redis.Client
	listener []*subscriber // слушатели радио
}

func NewHarbour(red *redis.Client) *Harbour {
	return &Harbour{red: red, mutex: sync.Mutex{}, listener: make([]*subscriber, 0)}
}

// Run - запуск процесса получения грузов из Асгарда по радужному мосту (redis-у) и ретрансляции их по радио в виде радиограмм
// Odin: запускайте это в горутине, или оно остановит всё остальное
func (l *Harbour) Run(ctx context.Context) error {
	subscriber := l.red.Subscribe(
		ctx,
		model.ChanEntropy,
		model.ChanEntropyExtended,
		model.ChanNewArt,
	)
	log.Info().Msgf("[harbour] ГАВАНЬ НАЧИНАЕТ ПРИЁМ ГРУЗОВ")
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[harbour] ОСТАНОВ РАБОТЫ ГАВАНИ")
			return nil
		default:
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[harbour] ОШИБКА ПРИЁМА ГРУЗА")
				time.Sleep(time.Second)
				continue
			}
			if err := l.handle(ctx, msg); err != nil {
				log.Error().Err(err).Msgf("[harbour] ОШИБКА ОБРАБОТКИ ГРУЗА")
			}
		}
	}
}

// handle - обработка одного пришедшего груза
func (l *Harbour) handle(ctx context.Context, msg *redis.Message) error {
	broadcastChannels := []string{
		model.ChanEntropy,
		model.ChanEntropyExtended,
	}
	idx := slices.IndexFunc(broadcastChannels, func(s string) bool { return msg.Channel == s })
	if idx != -1 {
		// нужно отправить это сообщение броадкастом всем подписанным клиентам, которые сейчас слушают websocket
		l.makeRadioshow(ctx, msg.Channel, msg.Payload)
	}

	switch msg.Channel {
	case model.ChanNewArt:
		if err := l.handleNewArt(ctx, msg); err != nil {
			return errors.Wrapf(err, "[harbour] ОШИБКА ОБРАБОТКИ СОБЫТИЯ %s", msg.Channel)
		}

	}

	return nil
}

/*
handleNewArt - гавань получила груз new_art
// Odin: структура model.Art со всеми вложениями очень большая, и отправлять её в чистом виде в Мидгард накладно из-за трафика
// Odin: тут гавань переупакует груз в более лёгкую структуру portals.FlatArt
*/
func (l *Harbour) handleNewArt(ctx context.Context, msg *redis.Message) error {
	var art model.Art
	if err := json.Unmarshal([]byte(msg.Payload), &art); err != nil {
		return errors.Wrap(err, "[harbour] ОШИБКА JSON-РАСПАКОВКИ ГРУЗА. ГРУЗ ART ПОВРЕЖДЁН.")
	}

	flat := model.MakeFlatArt(art)
	j, err := json.Marshal(flat)
	if err != nil {
		return errors.Wrap(err, "[harbour] ОШИБКА УПАКОВКИ РАДИОГРАММЫ ИЗ FLAT-ART СТРУКТУРЫ")
	}

	log.Debug().Msgf("[harbour] ПЕРЕОТПРАВЛЯЮ FLAT-ART #%d В КАНАЛ %s", flat.ID, model.ChanNewArt)
	l.makeRadioshow(ctx, model.ChanNewArt, string(j))
	return nil
}

// makeRadioshow - светлые эльфы в гавани передают груз по радио в виде радиограммы
// Odin: в этой радиостанции есть лишь одна частота, через которую переваливается всё. Разделение тут излишне
// Odin: отдельные типы грузов будут отброшены самим слушателем, если мидгардец не выразил намеренье эти грузы получать.

func (l *Harbour) makeRadioshow(ctx context.Context, channel string, payload string) {
	event := model.Radiogram{
		Channel: channel,
		Payload: payload,
	}
	l.mutex.Lock()
	subscribers := l.listener[:]
	l.mutex.Unlock()

	for _, sub := range subscribers {
		// отправка сообщения каждому подписчику, если он ещё активен
		go func(s *subscriber) {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
				log.Error().Msgf("[radio] РАДИОГРАММА ПОТЕРЯНА. КАНАЛ:%s", channel)
			case s.eventCh <- event:
				//log.Debug().Msgf("[radio] РАДИОГРАММА ОТПРАВЛЕНА. КАНАЛ:%s. УСПЕХ", msg.Channel)
			}
		}(sub)
	}
}

// ListenRadio - подписка радиослушателя.
// Передаются все события, подлежащие веерной трансляции (без фильтрации)
// Слушатель сам должен отфильтровать те радиограммы, которые ему не нужны
func (l *Harbour) ListenRadio(subscribeCtx context.Context) chan model.Radiogram {
	eventCh := make(chan model.Radiogram)
	sub := subscriber{
		// с помощью этого контекста канал подписчик будет остановлен и удалён
		ctx:     subscribeCtx,
		eventCh: eventCh,
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.listener = append(l.listener, &sub)
	go func() {
		<-subscribeCtx.Done()
		l.unlisten(&sub)
	}()

	return eventCh
}

func (l *Harbour) unlisten(sub *subscriber) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	idx := slices.IndexFunc(l.listener, func(s *subscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[radio] СЛУШАТЕЛЬ ОТСУТСТВУЕТ. НЕ УДАЛИТЬ")
		return
	}

	l.listener = append(l.listener[:idx], l.listener[idx+1:]...)
	log.Debug().Msgf("[radio] СЛУШАТЕЛЬ #%d - УДАЛЕНО", idx)
}
