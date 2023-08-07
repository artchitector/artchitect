class WsConnection {
  constructor(url, logPrefix, events, maxReconnectAttempts) {
    this.url = url;
    this.logPrefix = logPrefix;
    this.manuallyClosed = false;
    this.maxReconnectAttempts = maxReconnectAttempts;
    this.reconnectAttempts = 0;
    this.events = events; // ['tick', 'new_card', 'creation', 'lottery', 'selection', ...]
    this.callbacks = {
      onopen: [],
      onreconnecting: [],
      onmessage: [],
      onclose: [],
      onerror: [],
    }
  }

  connect() {
    if (process.server === true) {
      return
    }
    console.log(`${this.logPrefix}: Starting connection to WebSocket Server on ${this.url}`)

    // close previous connection
    if (!!this.connection) {
      this.close()
      this.connection = null
    }

    this._connect()
  }

  _connect() {
    if (this.manuallyClosed) {
      console.log(`${this.logPrefix}: connection manually closed. not _connect`)
      return
    }

    this.reconnectAttempts += 1
    this.emit('onreconnecting', this.reconnectAttempts, this.maxReconnectAttempts)

    this.connection = new WebSocket(this.url)
    this.connection.addEventListener('close', (e) => {
      console.log(`${this.logPrefix}: websocket onclose (clean=${e.wasClean})`, e)
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        console.log(`${this.logPrefix}: websocket max reconnect attempts exceeded (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
        this.emit('onerror', new Error("max reconnect attempts exceeded"))
      } else {
        setTimeout(() => {
          this._connect()
        }, 1000)
      }
      this.emit('onclose', null)
    })
    this.connection.addEventListener('error', (e) => {
      console.log(`${this.logPrefix}: websocket onerror`, e)
    })
    this.connection.addEventListener('message', (e) => {
      const ev = JSON.parse(e.data);
      if (this.events.includes(ev.channel)) {
        let data = JSON.parse(ev.payload)
        this.emit('onmessage', ev.channel, data)
      }
    })
    this.connection.addEventListener('open', (e) => {
      console.log(`${this.logPrefix}: websocket Successfully connected to the echo websocket server ${this.url}`)
      this.reconnectAttempts = 0
      this.emit('onopen', null)

      for (let allowedEvent of this.events) {
        this.connection.send(`subscribe.${allowedEvent}`)
      }
    })
    return true
  }


  close() {
    this.manuallyClosed = true
    if (!this.connection) {
      console.log(`${this.logPrefix}: Connection already closed`)
      return
    }
    console.log(`${this.logPrefix}: Connection closing...`)
    this.connection.close()
    console.log(`${this.logPrefix}: Connection closed`)
  }

  // When WS connected successfully
  onopen(cb) {
    this.callbacks.onopen.push(cb)
  }

  // When WS connection lost and trying to reconnect
  onreconnecting(cb) {
    this.callbacks.onreconnecting.push(cb)
  }

  onmessage(cb) {
    this.callbacks.onmessage.push(cb)
  }

  onclose(cb) {
    this.callbacks.onclose.push(cb)
  }

  onerror(cb) {
    this.callbacks.onerror.push(cb)
  }

  emit(type, event) {
    let callbacks = this.callbacks[type]
    if (callbacks.length === 0) {
      return
    }
    const args = Array.prototype.slice.call(arguments, 1);
    callbacks.forEach((cb) => {
      cb(...args)
    })
  }
}

export default WsConnection
