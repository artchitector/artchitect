<i18n>
{
  "en": {
    "title": "Artchitect - idea"
  },
  "ru": {
    "title": "Artchitect - идея"
  }
}
</i18n>
<template>
  <endescription v-if="locale === 'en'" :entropy="entropy"/>
  <rudescription v-else :entropy="entropy"/>
</template>
<script>
import Rudescription from "@/components/idea/ru/rudescription.vue";
import Endescription from "@/components/idea/en/endescription.vue";
import WsConnection from "@/utils/ws_connection";

export default {
  components: {Endescription, Rudescription},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      connection: null,
      logPrefix: '♦️',
      entropy: null,
      status: {
        error: null,
        reconnecting: null,
      },
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, this.logPrefix, ['entropy'], 100)
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
      this.onMessage(channel, message)
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`${this.logPrefix}: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {attempt, maxAttempts}
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
      console.log(`${this.logPrefix}: connection established`)
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  },
  methods: {
    onMessage(chan, message) {
      this.entropy = message
    }
  }
}
</script>
