<template>
  <div v-if="activeComponent === 'empty'" class="has-text-centered pt-4">
    no active component. wait for a while<br/>
    <loader class="mt-4"/>
  </div>
  <creation v-else-if="activeComponent === 'creation'" ref="creation"/>
  <unity v-else-if="activeComponent === 'unity'" ref="unity"/>
  <lottery v-else-if="activeComponent === 'lottery'" ref="lottery"/>
</template>

<script>

import Creation from "@/components/flexheart/creation/creation.vue";
import Lottery from "@/components/flexheart/lottery/lottery.vue"
import WsConnection from "@/utils/ws_connection";
import Unity from "@/components/flexheart/unity/unity.vue";

export default {
  name: "core",
  components: {Unity, Creation, Lottery},
  data() {
    return {
      logPrefix: '❤️',
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,
      activeComponent: 'empty'
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(
      process.env.WS_URL,
      this.logPrefix,
      ['creation', 'lottery', 'unity', 'heart', 'entropy_mini'],
      100,
    )
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
      // this.message = [channel, message]
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
      console.log(`${this.logPrefix}: connection established 🍏`)
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
  methods: {
    onMessage(channelName, message) {
      // В Сердце может находиться не тот компонент, по которому пришло новое сообщение.
      // Такое бывает, когда режим Архитектора переключается на иную задачу
      // (например, нарисовал и пошёл собирать множество)

      const extraEvents = ['heart', 'entropy_mini']
      if (extraEvents.indexOf(channelName) === -1 && this.activeComponent !== channelName) {
        this.activeComponent = channelName
        setTimeout(() => {
          this.$refs[channelName].onMessage(channelName, message)
        }, 100)
      }

      switch (channelName) {
        case 'entropy_mini':
        case 'creation':
        case 'heart':
          if (this.$refs.creation && this.$refs.creation.onMessage) {
            this.$refs.creation.onMessage(channelName, message)
          }
          break

        case 'unity':
          if (this.$refs.unity && this.$refs.unity.onMessage) {
            this.$refs.unity.onMessage(channelName, message)
          }
          break
        case 'lottery':
          if (this.$refs.lottery && this.$refs.lottery.onMessage) {
            this.$refs.lottery.onMessage(channelName, message)
          }
          break
      }
    },
  }
}
</script>


