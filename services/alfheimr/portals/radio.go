package portals

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
	"sync"
)

const (
	CommandSubscribe   = "subscribe"
	CommandUnsubscribe = "unsubscribe"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*
RadioPortal
Odin: Дам имя этой http-ручке - /radio
Odin: Мидгардцы подключаются к этому радио из мира Альфхейма, чтобы обрести длительную рантаймовую связь
с Artchitect и всеми его данными.
Odin: Выполнено на примитивной технологии вебсокетов, но свою работу делает. Мидгардцы используют свои веб-браузеры для подключения.
Odin: через этот портал человек на другом конце будет видеть состояние всего Artchitect до глубоких уровней творения
Odin: Я хочу, чтобы Artchitect был понятен и прозрачен со стороны, вот и придумал такой путь.
*/
type RadioPortal struct {
	radio radio
}

func NewRadioPortal(radio radio) *RadioPortal {
	return &RadioPortal{radio: radio}
}

func (cp *RadioPortal) Handle(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusBadRequest)
	//w.Write([]byte("STOP"))
	//return

	connID := makeRadioConnectionID(3)
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Error().Err(err).Msgf("[portal:radio:%s] АВАРИЯ", connID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Odin: если клиент отключится, то через этот контекст я выключу и связанные горутины
	sCtx, cancel := context.WithCancel(r.Context())
	defer cancel()

	radioCh := cp.radio.ListenRadio(sCtx)
	commands := cp.startIncomingMessageListen(r.Context(), cancel, connID, conn)
	mutex := sync.Mutex{}
	activeChannels := make([]string, 0, 0)

	// ОБРАБОТКА ВХОДЯЩИХ КОМАНД (ПОДКЛЮЧЕНИЕ НА КАНАЛЫ)
	go func() {
		for cmd := range commands {
			if len(cmd) != 2 || (cmd[0] != CommandSubscribe && cmd[0] != CommandUnsubscribe) {
				err = errors.Errorf("[portal:radio:%s] КОМАНДА НЕ РАСПОЗНАНА: %s", connID, strings.Join(cmd, ","))
				log.Error().Err(err).Send()
			}
			if cmd[0] == CommandSubscribe {
				chIdx := slices.IndexFunc(activeChannels, func(s string) bool { return s == cmd[1] })
				if chIdx == -1 {
					mutex.Lock()
					activeChannels = append(activeChannels, cmd[1])
					mutex.Unlock()
					log.Info().Msgf("[portal:radio:%s] ПОДПИСКА НА КАНАЛ %s", connID, cmd[1])
					cp.techMessage(sCtx, conn, connID, fmt.Sprintf("%s ПОДПИСАН НА КАНАЛ %s", connID, cmd[1]))
				}
			} else if cmd[0] == CommandUnsubscribe {
				log.Info().Msgf("[portal:radio:%s] ОТПИСКА ОТ КАНАЛА %s", connID, cmd[1])
				chIdx := slices.IndexFunc(activeChannels, func(s string) bool { return s == cmd[1] })
				if chIdx != -1 {
					mutex.Lock()
					activeChannels = append(activeChannels[:chIdx], activeChannels[chIdx+1:]...)
					mutex.Unlock()
					cp.techMessage(sCtx, conn, connID, fmt.Sprintf("%s ОТПИСКА ОТ %s", connID, cmd[1]))
				}
			}
		}
		log.Info().Err(err).Msgf("[portal:radio:%s] ОБРАБОТКА КОМАНДА ЗАВЕРШЕНА", connID)
	}()

	for radiogram := range radioCh {
		// Эльф Альфхейма: Привет! В гавани мы приняли груз из Асгарда. Я записал новую радиограмму, и передаю её тут тебе, Человек.

		// Если радиослушатель подписался на этот канал, то только тогда надо ему отправить сообщение
		mutex.Lock()
		chIdx := slices.IndexFunc(activeChannels, func(s string) bool { return s == radiogram.Channel })
		mutex.Unlock()
		if chIdx == -1 {
			// Слушатель не подписан, ему сообщение отправляеть не будет
			continue
		}

		j, err := json.Marshal(radiogram)
		if err != nil {
			log.Error().Err(err).Msgf("[portal:radio:%s] РАДИОГРАММА ИСПОРЧЕНА", connID)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, j); err != nil {
			log.Warn().Err(err).Msgf("[portal:radio:%s] ОТПРАВКА НЕ УДАЛАСЬ", connID)
			cancel() // Odin: если есть какие-то проблемы, то лучше разорвём канал связи совсем
			return
		}
	}

	log.Debug().Msgf("[portal:radio:%s] РАДИОТРАНСЛЯЦИЯ ОСТАНОВЛЕНА", connID)
}

// startIncomingMessageListen - запускается чтение сообщений, которые приходят с Мидгарда.
// Odin: Событий в Artchitect, которые отправляются в Мидгард для показа - несколько, а некоторые из них тяжёлые по объёму.
// Odin: получатель должен сначала предъявить список каналов, из которых он собирается использовать данные
// Odin: и лишь потом светлые эльфы Альфхейма будут знать, какие посылки отправлять в Мидгард, а какие нет
func (cp *RadioPortal) startIncomingMessageListen(
	ctx context.Context,
	cancel func(),
	connID string,
	conn *websocket.Conn,
) chan []string {
	commandCh := make(chan []string)
	go func() {
		// Как эта горутина останавливается по контексту?
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error().Err(err).Msgf("[portal:radio:%s] ОШИБКА ЗАКРЫТИЯ ПОРТАЛА", connID)
				} else {
					log.Debug().Msgf("[portal:radio:%s] ПРИНИМАЮЩИЙ ОТКЛЮЧИЛСЯ", connID)
				}
				cancel() // при ошибке соединение разрывается
				break
			}

			// Odin: мидгард запросил что-то по радио (прислал сообщение). Интересное радио, в обе стороны работает.
			log.Debug().Msgf("[portal:radio:%s] ВХОДЯЩИЙ РАДИО-ЗАПРОС: %s", connID, message)
			command := strings.Split(string(message), ".")
			commandCh <- command
		}
	}()

	return commandCh
}

func (cp *RadioPortal) techMessage(
	ctx context.Context,
	conn *websocket.Conn,
	connID string,
	text string,
) {
	radiogram := model.Radiogram{
		Channel: model.ChanTech,
		Payload: text,
	}
	j, err := json.Marshal(radiogram)
	if err != nil {
		log.Error().Err(err).Msgf("[portal:radio:%s] ОШИБКА MARSHAL-ИНГА TECH-СООБЩЕНИЯ", connID)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, j); err != nil {
		log.Warn().Err(err).Msgf("[portal:radio:%s] ОТПРАВКА TECH-СООБЩЕНИЯ НЕ УДАЛАСЬ", connID)
		return
	}
}
