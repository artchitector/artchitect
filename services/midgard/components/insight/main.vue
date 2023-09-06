<i18n>
{
  "en": {
    "loading": "...waiting for asgard-transmission..."
  },
  "ru": {
    "loading": "...ожидаем вестей с Асгарда..."
  }
}
</i18n>
<template>
  <div v-if="activeComponent === 'empty'" class="has-text-centered pt-4">
    {{ $t('loading') }}
    <br/>
    <common-loader class="mt-4"/>
  </div>
  <insight-odin v-else-if="activeComponent === 'odin'"
                :odin="state.odin"
                :entropy="state.entropy"
                :giving="state.giving"
                ref="odin"/>
  <insight-frigg v-else-if="activeComponent === 'frigg'"
                 :frigg="state.frigg"
                 :entropy="state.entropy"
                 ref="frigg"/>

</template>

<script>
import Radio from "@/utils/radio";

export default {
  name: "insight-main",
  data() {
    return {
      activeComponent: 'empty',
      radioPid: {
        entropy: null,
        odin: null,
        frigg: null,
        giving: null, // giving - выдаёт 4 случайные картины из всех написанных, меняя одну картину в выборке раз в 3 секунды
      },
      state: {
        entropy: null,
        odin: null,
        frigg: null,
        giving: null,
      }
    }
  },
  async mounted() {
    this.radioPid.entropy = Radio.subscribe("entropy", (entropy) => {
      this.state.entropy = entropy
    }, (err) => {
      console.error("[RADIO-ENTROPY] ОШИБКА ПОДКЛЮЧЕНИЯ К РАДИО", err)
    })
    this.radioPid.odin = Radio.subscribe("odin_state", (odinState) => {
      this.onMessage("odin", odinState)
    }, (err) => {
      console.error("[RADIO-ODIN] ОШИБКА ПОДКЛЮЧЕНИЯ К РАДИО", err)
    })
    this.radioPid.frigg = Radio.subscribe("frigg_state", (friggState) => {
      this.onMessage("frigg", friggState)
    }, (err) => {
      console.error("[RADIO-FRIGG] ОШИБКА ПОДКЛЮЧЕНИЯ К РАДИО", err)
    })
    this.radioPid.giving = Radio.subscribe("giving", (giving) => {
      this.state.giving = giving
    }, (err) => {
      console.error("[RADIO-GIVING] ОШИБКА ПОДКЛЮЧЕНИЯ К РАДИО", err)
    })
  },
  beforeDestroy() {
    Radio.unsubscribe(this.radioPid.entropy)
    Radio.unsubscribe(this.radioPid.odin)
    Radio.unsubscribe(this.radioPid.frigg)
    Radio.unsubscribe(this.radioPid.giving)
  },
  methods: {
    onMessage(stateType, state) {
      if (stateType === 'odin') {
        if (this.activeComponent !== 'odin') {
          this.activeComponent = 'odin'
        }
        this.state.odin = state
      } else if (stateType === "frigg") {
        if (this.activeComponent !== "frigg") {
          this.activeComponent = "frigg"
        }
        this.state.frigg = state
      }
    }
  }
}
</script>

<style scoped lang="scss">

</style>
