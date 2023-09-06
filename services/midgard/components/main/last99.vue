<i18n>
{
  "en": {
    "last99": "last 99 arts",
    "not_loaded": "arts not loaded"
  },
  "ru": {
    "last99": "последние 99 картин"
  }
}
</i18n>
<template>
  <section>
    <div>
      <h3 class="is-size-4 has-text-centered mb-4">{{ $t('last99') }}</h3>
      <div v-if="$fetchState.pending" class="notification has-text-centered">
        <common-loader size="l"/>
      </div>
      <div v-else-if="$fetchState.error" class="notification is-danger">
        {{ $fetchState.error.message }}
      </div>
      <div v-else-if="!this.arts.length" class="notification is-danger">
        {{ $t('not_loaded') }}
      </div>
      <common-art-list :arts="this.arts" row-size="3" initial-visible-count="33"/>
    </div>
  </section>
</template>
<script>
import Radio from "@/utils/radio";

export default {
  name: "main-last99",
  data() {
    return {
      radioPid: null,
      arts: [],
    }
  },
  async fetch() {
    this.arts = await this.$axios.$get("/arts/last/99")
  },
  async mounted() {
    this.radioPid = Radio.subscribe("new_art", (art) => {
      console.log(`[LAST99] НОВАЯ КАРТИНА #${art.id}`)
      this.arts.unshift(art)
      if (this.arts.length > 99) {
        this.arts.pop()
      }
    }, (error) => {
      alert(error)
    })
  },
  beforeDestroy() {
    if (this.radioPid) {
      Radio.unsubscribe(this.radioPid)
    }
  }
}
</script>
<style lang="scss" scoped>

</style>
