<i18n>
{
  "en": {
    "loading": "loading",
    "error": "error",
    "ws_error": "websocket listening error",
    "ws_connecting": "websocket connecting",
    "only_10": "only last 10 lotteries shows"
  },
  "ru": {
    "loading": "загрузка",
    "error": "ошибка",
    "ws_error": "ошибка подключения к websocket",
    "ws_connecting": "подключение к websocket",
    "only_10": "только 10 последних лотерей показывается"
  }
}
</i18n>
<template>
  <div>
    <div class="notification is-primary" v-if="$fetchState.pending">
      {{$t('loading')}}...
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{$t('error')}}: {{ $fetchState.error.message }}
    </div>
    <div v-else>
      <div v-if="wsStatus.error" class="notification is-warning is-size-7 has-text-centered">
        {{$t('ws_error')}}: {{wsStatus.error.message}}
      </div>
      <div v-else-if="wsStatus.reconnecting" class="notification is-size-7 has-text-centered">
        {{$t('ws_connecting')}} {{wsStatus.reconnecting.attempt}}/{{wsStatus.reconnecting.maxAttempts}}
      </div>
      <lottery v-for="lottery in lotteries" v-bind:key="lottery.ID" :lottery="lottery"/>
      <div class="is-size-7 has-text-centered" v-if="lotteries.length >= 10">
        {{$t('only_10')}}
      </div>
    </div>
  </div>
</template>

<script>
import WsConnection from "@/utils/ws_connection";

export default {
  name: "lottery-list",
  data() {
    return {
      connection: null,
      wsStatus: {
        error: null,
        reconnecting: null,
      },
      lotteries: [],
    }
  },
  mounted() {
    this.connection = new WsConnection(process.env.WS_URL, '🏆', ['lottery'], 10)
    this.connection.onmessage((channel, lottery) => {
      this.updateLottery(lottery.Lottery)
    });
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
      this.lotteries = await this.$axios.$get('/lottery/10')
    } catch (e) {
      if (this.connection) {
        this.connection.close();
      }
      throw e;
    }
  },
  methods: {
    updateLottery(lottery) {
      if (!this.lotteries || !this.lotteries.length) {
        return
      }
      console.log(`🏆: update lottery id=${lottery.ID}, winners count: ${lottery.Winners.length}`)
      for (let i = 0; i <= this.lotteries.length; i++) {
        const l = this.lotteries[i];
        if (l.ID === lottery.ID) {
          this.$set(this.lotteries, i, lottery);
          break;
        }
      }
    }
  }
}
</script>

<style scoped>

</style>
