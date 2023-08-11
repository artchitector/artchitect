<i18n>
{
  "en": {
    "creative_process": "creative process",
    "creating": "creating",
    "seed": "seed",
    "painting": "creating"
  },
  "ru": {
    "creative_process": "творческий процесс",
    "creating": "создаю",
    "seed": "seed",
    "painting": "пишу картину"
  }
}
</i18n>
<template>
  <div class="progress-view">
    <div class="heart-heading">
      <common-entropy :entropy="entropy" v-if="entropy"/>
      <h1 class="is-size-5 has-text-success has-text-centered mb-2">
        {{ $t('creative_process') }}
      </h1>
      <div>
        <div class="tags mb-3">
          <insight-odin-progress-seed class="tag is-primary" :odin="odin"/>
          <span class="tag" v-for="word in odin.words">{{ word }}</span>
        </div>
      </div>
      <b-progress type="is-primary" :value="progress" show-value>
        {{ $t('painting') }} {{ odin.currentPaintTime }}/{{ odin.totalPaintTime }}
      </b-progress>
    </div>
  </div>
</template>

<script>
export default {
  name: "insight-odin-progress",
  props: ['entropy', 'odin'],
  computed: {
    progress() {
      if (!this.odin || (!this.odin.currentPaintTime && !this.odin.totalPaintTime)) {
        return 0
      }
      const progress = this.odin.currentPaintTime / this.odin.totalPaintTime
      return Math.floor(progress * 100)
    }
  }
}
</script>

<style lang="scss" scoped>
.progress-view {
  max-width: 800px;
  min-width: 370px;

  .tags .tag {
    font-size: 9px;
    letter-spacing: 0px;
  }
}
</style>
