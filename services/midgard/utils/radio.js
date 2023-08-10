// Odin: На Alfheimr установлена радиостанция, которая ретранслирует все события из мира Асгарда
// Odin: радио зиждется на технологии websocket (поэтому оно еще и двухстороннее)
// Odin: Ниже коннектор к этому радио, чтобы мир Мидгарда мог подслушивать поток.
// Odin: Получение радиограмм работает по-подписке.
// ...Для профессионалов: Всего 9.99Р в месяц, и вы получаете PRO-подписку с событиями из раздела PRO.
// ...Для наших самых требовательных и особенных клиентов: всего за 39.99Р вы получите ARCHI-подписку с особенными событиями,
//  которые не будут нужны ни только никому, но и не будут нужны даже вам.
// Loki: Ты отвлёкся. Хочешь стриминговый сервис сделать? Можем ИИ и для музыки запустить.
// Odin: хорошая идея на следующий проект) Будем музыку писать. Но, вернёмся к подпискам:
// Odin: ! ДЛЯ ВСЕХ ОСТАЛЬНЫХ, не таких требовательных клиентов - все события Artchitect всегда бесплатны навсегда.
// ...Все каналы открыты!
//
// Odin: как бы это ни обзывалось, но одной конкретной сессии в браузере не нужны вообще все события, ИМ НЕТ ЧИСЛА!
// Odin: с помощью подписки можно переключить вашу радиостанцию на нужную частоту и слушать ТОЛЬКО НУЖНЫЕ события
// Odin: кто-то будет смотреть страницу энтропии, и будет получать только энтропию.
// Odin: те, кто на главной, будут получать события new_art, чтобы показывать новые картинки в списке, и так далее.
// Odin: Каждый компонент может указывать, что именно он хочет получать от Radio
// Odin: при первой подписке радио активируется и работает глобально в сессии браузера (Artchitect работает как SPA в браузере)
// Odin: vue-компоненты подписываются на нужные канали, и отписываются (тоже надо предусмотреть!)
// Odin: если компонент перестаёт слушать определённую частоту, то надо отправить запрос на прекращение вещания.
class Radio {
  constructor(url) {
    this.url = url;
    this.connection = null;
    this.reconnectAttempts = 0;
    this.shutdown = false;

    this.activeChannels = [];
    this.pidCounter = 1;
    this.subscribers = {};
  }

  async subscribe(channel, cb) {
    if (!this.connection) {
      this.connect()
      await this.waitConnection()
    }

    if (this.activeChannels.indexOf(channel) !== -1) {
      // Слушатель уже подписан
    } else {
      this.connection.send(`subscribe.${channel}`)
      this.activeChannels.push(channel)
    }
    const pid = this.pidCounter++
    this.subscribers[pid] = {
      channel: channel,
      callback: cb,
    }

    return pid
  }

  unsubscribe(pid) {
    const subscriber = this.subscribers[pid]
    if (!subscriber) {
      throw `[RADIO] НЕТ ПОДПИСЧИКА С PID=${pid}`
    }
    delete this.subscribers[pid]

    const needUnsubscribe = !this.isChannelSubscribed(subscriber.channel)
    if (!needUnsubscribe) {
      return
    }
    setTimeout(async () => {
      this.activeChannels.splice(this.activeChannels.indexOf(subscriber.channel), 1)
      await this.waitConnection()
      this.connection.send(`unsubscribe.${subscriber.channel}`)
    })
  }

  connect(cb) {
    if (process.server === true) {
      return
    }
    if (this.reconnectAttempts > 10) {
      this.shutdown = true;
      console.log('[RADIO] ПОДКЛЮЧЕНИЕ ЗАКРЫТО НАВСЕГДА. БОЛЕЕ 10 ПОПЫТОК ПОДКЛЮЧЕНИЯ ПРОВАЛЕНО')
      return
    } else
      console.log(`[RADIO] ${this.reconnectAttempts}/10 ПОДКЛЮЧАЮСЬ К РАДИОСТАНЦИИ ${this.url}`)
    this.connection = new WebSocket(this.url)

    this.connection.addEventListener("close", () => {
      console.log('[RADIO] ПОДКЛЮЧЕНИЕ ЗАКРЫТО')
      this.reconnectAttempts += 1
      setTimeout(() => {
        this.connect(cb)
      }, 1000)

    })

    this.connection.addEventListener("error", (e) => {
      console.log('[RADIO] ОШИБКА: ', e)
    })

    this.connection.addEventListener("message", (e) => {
      const ev = JSON.parse(e.data)
      if (ev.channel === "tech") {
        // ТЕХНИЧЕСКОЕ СООБЩЕНИЕ, ПРОСТО ЛОГИРУЕТСЯ
        console.log(`[RADIO] TECH: ${ev.payload}`)
        return
      }
      this.onMessage(ev)
    })

    this.connection.addEventListener("open", () => {
      console.log('[RADIO] ПОДКЛЮЧЕНО')
      this.reconnectAttempts = 0
    })
  }

  async waitConnection() {
    return await new Promise((resolve, reject) => {
      if (this.shutdown) {
        reject("[RADIO] РАДИОСТАНЦИЯ ПОГАШЕНА. МНЕ НЕ ПОДКЛЮЧИТЬСЯ")
        return
      }
      if (this.connection.readyState === this.connection.OPEN) {
        resolve()
        return
      }
      let attempts = 0;
      const interval = setInterval(() => {
        if (attempts > 10) {
          clearInterval(interval)
          reject("[RADIO] СЛИШКОМ МНОГО ПОПЫТОК ОЖИДАНИЯ КОННЕКТА. СТОП")
        }
        if (this.connection.readyState === this.connection.OPEN) {
          resolve()
          clearInterval(interval)
        } else {
          attempts++
        }
      }, 100)
    })
  }

  onMessage(ev) {
    const payload = JSON.parse(ev.payload)
    for (let pid in this.subscribers) {
      const s = this.subscribers[pid]
      if (s.channel === ev.channel) {
        s.callback(payload)
      }
    }
  }

  isChannelSubscribed(channel) {
    for (let i = 0; i < this.subscribers.length; i++) {
      if (this.subscribers[i].channel === channel) {
        return true
      }
    }
    return false
  }
}

const r = new Radio(process.env.WS_URL)
export default r
