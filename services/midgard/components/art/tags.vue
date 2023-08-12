<i18n>
{
  "en": {
    "seed": "seed",
    "more": "show more",
    "less": "show less"
  },
  "ru": {
    "seed": "seed",
    "more": "показать подробнее",
    "less": "показать меньше"
  }
}
</i18n>
<template>
  <div class="tags" v-if="!more">
    <span class="tag is-primary">
      {{ $t('seed') }}:{{ art.idea.seed }}
    </span>
    <span v-for="word in art.idea.words" class="tag">
      {{ word.word }}
    </span>
    <a href="#" class="tag is-info" @click="showMore">
      {{ $t('more') }}
    </a>
  </div>
  <div v-else>
    <div class="has-text-centered">
      <a href="#" @click="showLess">{{ $t('less') }}</a>
    </div>
    <table class="table is-bordered">
      <tbody>
      <tr class="is-selected">
        <td>
          seed={{ art.idea.seed }}
        </td>
        <td>
          <img class="entropy-image" :src="`data:image/png;base64, ${art.idea.seedEntropy.entropy.image}`"
               alt="odin's mind"/>
          <img class="entropy-image" :src="`data:image/png;base64, ${art.idea.seedEntropy.choice.image}`"
               alt="odin's mind"/>
        </td>
      </tr>
      <tr v-for="word in art.idea.words">
        <td>
          {{ word.word }}
        </td>
        <td>
          <img class="entropy-image" :src="`data:image/png;base64, ${word.entropy.entropy.image}`"
               alt="odin's mind"/>
          <img class="entropy-image" :src="`data:image/png;base64, ${word.entropy.choice.image}`"
               alt="odin's mind"/>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: "art-tags",
  props: ["art"],
  data() {
    return {
      more: false
    }
  },
  methods: {
    showMore() {
      this.more = true
    },
    showLess() {
      this.more = false
    }
  }
}
</script>

<style scoped lang="scss">
.table td {
  vertical-align: middle;
  padding-bottom: 0;
  padding-top: 0;
}

.entropy-image {
  margin-top: 6px;
  width: 32px;
  height: 32px;
  image-rendering: pixelated;
}
</style>
