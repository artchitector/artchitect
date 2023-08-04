<i18n>
{
  "en": {
    "heart": "creative process",
    "maintenance": "Artchitect need rest. Maintenance. No creative process.",
    "error": "Error",
    "connecting": "connecting",
    "connected": "connected. waiting for event"
  },
  "ru": {
    "heart": "творческий процесс",
    "maintenance": "Artchitect нуждается в отдыхе. Техобслуживание. Творческий процесс не запущен.",
    "error": "Ошибка",
    "connecting": "соединение",
    "connected": "соединение установлено. ожидаю событий"
  }
}
</i18n>
<template>
  <section class="heart">
    <div v-if="maintenance" class="notification is-warning">
      {{ $t('maintenance') }}
    </div>
    <template v-else>
      <h3 class="has-text-centered is-size-4">
        <NuxtLink :to="localePath('heart')" class="">
          {{ $t('heart') }}
        </NuxtLink>
      </h3>
      <hr class="divider"/>
      <div v-if="status.error" class="notification is-danger has-text-centered">
        {{ $t('error') }}: {{ status.error.message }}
      </div>
      <div v-else-if="status.reconnecting" class="notification has-text-centered">
        <loader size="s"/>
        <br/>
        {{ $t('connecting') }} {{ status.reconnecting.attempt }}/{{ status.reconnecting.maxAttempts }}
      </div>
      <p v-else-if="!stateChannel" class="notification has-text-centered">
        <loader size="s"/>
        <br/>
        {{ $t('connected') }}
      </p>
      <creation v-else-if="stateChannel === 'creation'" :state="state"/>
      <heart-lottery v-else-if="stateChannel === 'lottery'" :state="state"/>
      <unity v-else-if="stateChannel === 'unity'" :state="state"/>
      <p v-else>
        unknown state {{ stateChannel }}
      </p>
    </template>
  </section>
</template>
<script>
import WsConnection from '../../utils/ws_connection'
import Creation from "@/components/heart/layout/creation.vue";
import Unity from "@/components/heart/layout/unity.vue";

export default {
  components: {Unity, Creation},
  data() {
    return {
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,
      stateChannel: null,
      state: null,
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, '🧡', ['creation', 'lottery', 'unity'], 10)

    this.connection.onmessage((channel, state) => {
      this.status.error = null
      this.status.reconnecting = null
      this.stateChannel = channel
      this.state = state
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`🧡: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {
        attempt: attempt,
        maxAttempts: maxAttempts,
      }
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
}
</script>
<style lang="scss">

</style>
