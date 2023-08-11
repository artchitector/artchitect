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
        giving: null // givin - выдаёт 4 случайные картины из всех написанных, меняя одну картину в выборке раз в 3 секунды
      },
      state: {
        entropy: null,
        odin: null,
        giving: null,
      }
    }
  },
  async mounted() {
    this.radioPid.entropy = await Radio.subscribe("entropy", (entropy) => {
      this.state.entropy = entropy
    })
    this.radioPid.odin = await Radio.subscribe("odin_state", (odinState) => {
      this.onMessage("odin", odinState)
    })
    this.radioPid.giving = await Radio.subscribe("giving", (giving) => {
      this.state.giving = giving
    })
  },
  beforeDestroy() {
    Radio.unsubscribe(this.radioPid.entropy)
    Radio.unsubscribe(this.radioPid.odin)
    Radio.unsubscribe(this.radioPid.giving)
  },
  methods: {
    onMessage(stateType, state) {
      if (stateType === 'odin') {
        if (this.activeComponent !== 'odin') {
          this.activeComponent = 'odin'
        }
        this.state.odin = state
      }
    }
  }
}
</script>

<style scoped lang="scss">

</style>
