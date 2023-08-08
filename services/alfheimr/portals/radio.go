package portals

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
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
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Error().Err(err).Msgf("[portal:radio] АВАРИЯ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Odin: если клиент отключится, то через этот контекст я выключу и связанные горутины
	sCtx, cancel := context.WithCancel(r.Context())
	defer cancel()

	radioCh := cp.radio.ListenRadio(sCtx)
	commands := cp.startIncomingMessageListen(r.Context(), cancel, conn)
	go func() {
		for cmd := range commands {
			log.Info().Msgf("[portal:radio] ПОЛУЧЕНА КОМАНДА %s", strings.Join(cmd, ":"))
		}
	}()

	for radiogram := range radioCh {
		// Эльф Альфхейма: Привет! В гавани мы приняли груз из Асгарда. Я записал новую радиограмму, и передаю её тут тебе, Человек.
		j, err := json.Marshal(radiogram)
		if err != nil {
			log.Error().Err(err).Msgf("[portal:radio] РАДИОГРАММА ИСПОРЧЕНА")
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, j); err != nil {
			log.Warn().Err(err).Msgf("[portal:radio] ОТПРАВКА НЕ УДАЛАСЬ")
			cancel() // Odin: если есть какие-то проблемы, то лучше разорвём канал связи совсем
			return
		}
	}

	log.Debug().Msgf("[portal:radio] РАДИОТРАНСЛЯЦИЯ ОСТАНОВЛЕНА")
}

// startIncomingMessageListen - запускается чтение сообщений, которые приходят с Мидгарда.
// Odin: Событий в Artchitect, которые отправляются в Мидгард для показа - несколько, а некоторые из них тяжёлые по объёму.
// Odin: получатель должен сначала предъявить список каналов, из которых он собирается использовать данные
// Odin: и лишь потом светлые эльфы Альфхейма будут знать, какие посылки отправлять в Мидгард, а какие нет
func (cp *RadioPortal) startIncomingMessageListen(
	ctx context.Context,
	cancel func(),
	conn *websocket.Conn,
) chan []string {
	commandCh := make(chan []string)
	go func() {
		// Как эта горутина останавливается по контексту?
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error().Err(err).Msgf("[portal:radio] ОШИБКА ЗАКРЫТИЯ ПОРТАЛА")
				} else {
					log.Debug().Msgf("[portal:radio] ПРИНИМАЮЩИЙ ОТКЛЮЧИЛСЯ")
				}
				cancel() // при ошибке соединение разрывается
				break
			}

			// Odin: мидгард запросил что-то по радио (прислал сообщение). Интересное радио, в обе стороны работает.
			log.Debug().Msgf("[portal:radio] ВХОДЯЩИЙ РАДИО-ЗАПРОС: %s", message)
			command := strings.Split(string(message), ".")
			commandCh <- command
		}
	}()

	return commandCh
}
