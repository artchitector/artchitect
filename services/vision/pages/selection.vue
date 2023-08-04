<i18n>
{
  "en": {
    "title": "Artchitect - selection",
    "selection": "selection",
    "loading": "loading",
    "total": "total",
    "ws_error": "websocket listening error",
    "ws_connecting": "websocket connecting",
    "selection_description": "Everyday Artchitect selects 10-100 chosen dreams. For this purpose used ",
    "lottery": "lottery"
  },
  "ru": {
    "title": "Artchitect - отборное",
    "selection": "отборное",
    "loading": "загрузка",
    "total": "всего",
    "ws_error": "ошибка подключения к websocket",
    "ws_connecting": "подключение к websocket",
    "selection_description": "Ежедневно Архитектор выбирает 10-100 снов из общего числа, пережитых за день (~1400 снов). Для этого используется ",
    "lottery": "лотерея"
  }
}
</i18n>
<template>
  <div>
    <section>
        {{ $t('selection_description') }}
        <NuxtLink :to="localePath('lottery')">{{ $t('lottery') }}</NuxtLink>.
    </section>
    <section>
      <h1 class="is-size-3 has-text-centered mb-5">{{ $t('selection') }}</h1>
      <div v-if="count" class="has-text-centered mb-6">{{ $t('total') }} <b>{{ count }}</b></div>
      <div class="notification is-primary" v-if="$fetchState.pending">
        {{ $t('loading') }}...
      </div>
      <div class="notification is-primary" v-else-if="$fetchState.error">
        {{ $fetchState.error.message }}
      </div>
      <div v-else>
        <div v-if="wsStatus.error" class="notification is-warning is-size-7 has-text-centered">
          {{ $t('ws_error') }}: {{ wsStatus.error.message }}
        </div>
        <div v-else-if="wsStatus.reconnecting" class="notification is-size-7 has-text-centered">
          {{ $t('ws_connecting') }} {{ wsStatus.reconnecting.attempt }}/{{ wsStatus.reconnecting.maxAttempts }}
        </div>
        <cardlist :cards="selection" cards-in-column="5" card-size="s" visible-count="50"/>
      </div>
    </section>
  </div>
</template>
<script>
import Viewer from "@/components/viewer/viewer";
import Cardlist from "@/components/list/cardlist.vue";
import WsConnection from "@/utils/ws_connection";

export default {
  components: {Cardlist, Viewer},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      connection: null,
      wsStatus: {
        error: null,
        reconnecting: null,
      },
      selection: []
    }
  },
  mounted() {
    this.connection = new WsConnection(process.env.WS_URL, '🪆', ['new_selection'], 10)
    this.connection.onmessage((channel, selection) => {
      this.wsStatus.error = null;
      this.wsStatus.reconnecting = null;
      this.selection.unshift(selection.CardID)
      console.log(`🌄: new selection (card_id=${selection.CardID})`,)
    })
    this.connection.onerror((err) => this.wsStatus.error = err)
    this.connection.onreconnecting((attempt, maxAttempts) => this.wsStatus.reconnecting = {
      attempt: attempt,
      maxAttempts: maxAttempts
    })
    this.connection.onopen(() => {
      this.wsStatus.error = null
      this.wsStatus.reconnecting = null
    })
    this.connection.connect()
  },
  beforeDestroy() {
    this.connection.close()
    this.connection = null
  },
  async fetch() {
    try {
      this.selection = await this.$axios.$get("/selection")
    } catch (e) {
      if (this.connection) {
        this.connection.close()
      }
      throw e;
    }
  },
  computed: {
    count() {
      return this.selection.length
    }
  }
}
</script>
<style lang="scss">
.image-container {
  position: relative;

}
</style>
