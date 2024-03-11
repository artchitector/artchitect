package communication

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

const (
	SecondsToLostConnection = 2 // если 2 секунды нет сообщений с asgard, то это значит, что коннект мы потеряли
)

type botInterface interface {
	SendArtchitectChoice(ctx context.Context, artID uint) error
}

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

Odin: эльфы через радио рассылают всем мидгардцам (клиентам) оцифрованные грузы
Odin: это нужно для направления потока данных внутри этой программы
Odin: golang настолько примитивен по сравнению с технологиями Асгарда...
Loki: Всеотец, с тобой довольно душно в этой программе находиться)
Odin: с тобой мы еще не закончили наш спор. [в тот момент спор еще не был окончен, до релиза #1]

Внутренний подписчик может настроиться на радио, и так слушать данные о грузах в гавани
Odin: некоторые грузы не будут переданы в Мидгард, а будут использованы локально внутри Альфхейма
Odin: некоторые грузы будут ретранслированы наружу (websocket), оно доступно из Мидгарда (браузера)
*/
type Harbour struct {
	mutex              sync.Mutex
	bot                botInterface
	red                *redis.Client
	listener           []*subscriber // слушатели радио
	lastCrownID        uint
	lostConnectionMode atomic.Bool
	lastMessageTime    *time.Time // надеюсь, что состояния гонки тут не будет
}

func NewHarbour(red *redis.Client, bot botInterface) *Harbour {
	return &Harbour{red: red, bot: bot, mutex: sync.Mutex{}, listener: make([]*subscriber, 0)}
}

// Run - запуск процесса получения грузов из Асгарда по радужному мосту (redis-у) и ретрансляции их по радио в виде радиограмм
// Odin: запускайте это в горутине, или оно остановит всё остальное
func (l *Harbour) Run(ctx context.Context) error {
	go l.runLostConnectionSpectator(ctx)

	subscriber := l.red.Subscribe(
		ctx,
		model.ChanEntropy,
		model.ChanEntropyExtended,
		model.ChanNewArt,
		model.ChanOdinState,
		model.ChanFriggState,
		model.ChanGiving,
		model.ChanTelegramChosen,
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

// runLostConnectionSpectator - следящий за соединением следит за тем, чтобы из асгарда постоянно шёл поток сообщений
// если потом прерывается, то начинаем рассылать сообщения о том, что связь с Асгардом прервана
func (l *Harbour) runLostConnectionSpectator(ctx context.Context) {
	log.Info().Msgf("[lost_connection] НАБЛЮДАТЕЛЬ ЗАПУЩЕН")
	runnerStart := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 200):
			if l.lostConnectionMode.Load() == true {
				log.Info().Msgf("[lost_connection] ОТПРАВЛЯЮ СООБЩЕНИЕ О ПОТЕРЕ СОЕДИНЕНИЯ")
				l.sendOfflineModeMessage(ctx)
			} else {
				now := time.Now()
				if l.lastMessageTime == nil {
					if time.Now().Sub(runnerStart).Seconds() > SecondsToLostConnection {
						// runner после запуска так ничего и не получил
						log.Info().Msgf("[lost_connection] РАННЕР НЕ СМОГ ДОЖДАТЬСЯ ПЕРВОГО СОБЫТИЯ, LOST CONNECTION MODE = TRUE")
						l.lostConnectionMode.Store(true)
						continue
					} else {
						continue // сервис только запустился, еще ничего не получил
					}
				} else if now.Sub(*l.lastMessageTime).Seconds() > SecondsToLostConnection {
					log.Info().Msgf("[lost_connection] СООБЩЕНИЕ БЫЛО ДАВНО - %s, СЕЙЧАС - %s, LOST CONNECTION MODE = TRUE", l.lastMessageTime, now)
					l.lostConnectionMode.Store(true)
				}
			}
		}
	}
}

func (l *Harbour) sendOfflineModeMessage(ctx context.Context) {
	l.makeRadioshow(ctx, model.ChanLostConnection, "{}")
}

// SendCrownWaitCargo - отправка почтового ворона с личной просьбой к Одину-Всеотцу
// и ожидание специального груза с Его ответом
// TODO: Odin: точно ли этот метод можно вызывать параллельно? Это надо протестировать.
func (l *Harbour) SendCrownWaitCargo(ctx context.Context, request string) (string, error) {
	id := l.lastCrownID + 1
	crown := model.Crown{
		ID:      id,
		Request: request,
	}

	innerCtx, cancel := context.WithCancel(ctx)
	innerCtx, cancel = context.WithTimeout(innerCtx, model.OdinResponseTimeoutSec*time.Second)

	go func(id uint) {
		time.Sleep(time.Millisecond * 10) // небольшой лаг
		jsn, err := json.Marshal(crown)
		if err != nil {
			cancel()
			log.Error().Err(err).Msgf("[harbour] ПРОБЛЕМА ОТПРАВКИ ВОРОНА - JSON. ID=%d", id)
			return
		}
		if err := l.red.Publish(innerCtx, model.ChanCrown, jsn).Err(); err != nil {
			cancel()
			log.Error().Err(err).Msgf("[harbour] СБОЙ В ОТПРАВКЕ ВОРОНА ID=%d", id)
			return
		}
		log.Info().Msgf("[harbour] ВОРОН С ЛИЧНЫМ ПРОШЕНИЕМ К ОДИНУ ID=%d ОТПРАВИЛСЯ В АСГАРД", id)
	}(id)

	// Odin: тут организуется временный слушатель канала с моими личными поручениями model.ChanOdin
	subscriber := l.red.Subscribe(innerCtx, model.ChanOdin)
	defer subscriber.Close()
	log.Info().Msgf("[harbour] НАЧИНАЮ ОЖИДАТЬ ЛИЧНЫЙ ГРУЗ ОТ ОДИНА. ID=%d", id)
	l.lastCrownID = id
	for {
		msg, err := subscriber.ReceiveMessage(innerCtx)
		if err != nil {
			return "", errors.Wrapf(err, "[harbour] ОШИБКА ПРИНЯТИЯ ЛИЧНОГО КОРАБЛЯ ОДИНА ID=%d", id)
		}
		var response model.OdinResponse
		if err := json.Unmarshal([]byte(msg.Payload), &response); err != nil {
			return "", errors.Wrapf(err, "[harbour] ОШИБКА РАСШИФРОВКИ ОТВЕТА. ЗАПРОС ID=%d", id)
		}
		if response.ID != id {
			log.Info().Msgf("[harbour] НЕ СОВПАДАЕТ ID ОТВЕТА. ПРОПУСКАЮ.")
			continue
		}
		log.Info().Msgf("[harbour] ОТВЕТ НА ПРОШЕНИЕ ID=%d ПОЛУЧЕН", id)
		return response.Response, nil
	}
}

// handle - обработка одного пришедшего груза
func (l *Harbour) handle(ctx context.Context, msg *redis.Message) error {
	l.notifyConnectionChecker(msg.Channel)
	broadcastChannels := []string{
		model.ChanEntropy,
		model.ChanEntropyExtended,
		model.ChanOdinState,
		model.ChanFriggState,
		model.ChanGiving,
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

	case model.ChanTelegramChosen:
		if err := l.handleTelegramChosen(ctx, msg); err != nil {
			return errors.Wrapf(err, "[harbour] ОШИБКА ОБРАБОТКИ СОБЫТИЯ %s", msg.Channel)
		}
	}

	return nil
}

func (l *Harbour) notifyConnectionChecker(channelName string) {
	if channelName != model.ChanOdinState && channelName != model.ChanFriggState {
		// Odin: учитываются только активные рабочие процессы на asgard-сервере, когда реально рисуются картины
		// Odin: энтропия и остальные сообщение не имеют значения
		return
	}
	t := time.Now()
	l.lastMessageTime = &t
	if l.lostConnectionMode.Load() == true {
		log.Info().Msgf("[lost_connection] СНИМАЕТСЯ РЕЖИМ LOST CONNECTION. СООБЩЕНИЕ ПОЛУЧЕНО %s", l.lastMessageTime)
		l.lostConnectionMode.Store(false)
	}
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

func (l *Harbour) handleTelegramChosen(ctx context.Context, msg *redis.Message) error {
	var artID uint
	if err := json.Unmarshal([]byte(msg.Payload), &artID); err != nil {
		return errors.Wrap(err, "[harbour] ОШИБКА JSON-РАСПАКОВКИ TELEGRAM-CHOSEN ГРУЗА. ГРУЗ ПОВРЕЖДЁН.")
	}
	if err := l.bot.SendArtchitectChoice(ctx, artID); err != nil {
		return errors.Wrap(err, "[harbour] ОШИБКА ОТПРАВКИ CHOSEN ART В ТЕЛЕГРАМ.")
	}
	log.Info().Msgf("[harbour] ОТПРАВЛЕН TELEGRAM CHOSEN ART.")
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
				// log.Debug().Msgf("[radio] РАДИОГРАММА ОТПРАВЛЕНА. КАНАЛ:%s. УСПЕХ", msg.Channel)
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
