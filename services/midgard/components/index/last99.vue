<i18n>
{
  "en": {
    "last_99": "last 99 arts",
    "loading": "loading",
    "error": "Error",
    "not_loaded": "arts not loaded",
    "connecting": "connecting",
    "connected": "connected. waiting for event"
  },
  "ru": {
    "last_99": "последние 99 картин",
    "loading": "Загрузка",
    "error": "Ошибка",
    "not_loaded": "картины не загрузились",
    "connecting": "соединение",
    "connected": "соединение установлено. ожидаю событий",
    "ws_error": "websocket listening error",
    "ws_connecting": "websocket connecting"
  }
}
</i18n>
<template>
  <div>
    <h3 class="is-size-4 has-text-centered mb-4">{{$t('last_99')}}</h3>
    <div v-if="$fetchState.pending" class="notification has-text-centered">
      {{$t('loading')}}...
    </div>
    <div v-else-if="$fetchState.error" class="notification is-danger">
      {{$t('error')}} {{ $fetchState.error.message }}
    </div>
    <div v-else-if="!this.arts || !this.arts.length" class="notification is-danger">
      {{$t('not_loaded')}} :(
    </div>
    <div v-else>
      <div v-if="wsStatus.error" class="notification is-warning is-size-7 has-text-centered">
        {{$t('ws_error')}}: {{wsStatus.error.message}}
      </div>
      <div v-else-if="wsStatus.reconnecting" class="notification is-size-7 has-text-centered">
        {{$t('ws_connecting')}} {{wsStatus.reconnecting.attempt}}/{{wsStatus.reconnecting.maxAttempts}}
      </div>
      <art-list :arts="this.arts" arts-in-column="3" card-size="m" visible-count="33"/>
    </div>
  </div>
</template>

<script>
import ArtList from "@/components/list/artlist.vue";
import WsConnection from "@/utils/ws_connection";

export default {
  name: "last99",
  components: {ArtList},
  data() {
    return {
      connection: null,
      wsStatus: {
        error: null,
        reconnecting: null,
      },
      arts: []
    }
  },
  mounted() {
    this.connection = new WsConnection(process.env.WS_URL, '🌄', ['new_card'], 10)
    this.connection.onmessage((channel, newCard) => {
      this.wsStatus.error = null;
      this.wsStatus.reconnecting = null;
      const removedCard = this.arts[this.arts.length - 1];
      const arts = this.arts.slice(0, this.arts.length - 1)
      arts.unshift(newCard)
      this.arts = arts;
      console.log(`🌄: new card (id=${newCard.ID}), removed (id=${removedCard.ID})`,)
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
    this.connection = null;
  },
  async fetch() {
    try {
      this.arts = await this.$axios.$get('/arts/last/99')
      console.log('[last99] loaded arts', this.arts.length)
    } catch (e) {
      if (this.connection) {
        this.connection.close()
      }
      throw e;
    }
  }
}
</script>

<style scoped>

</style>
