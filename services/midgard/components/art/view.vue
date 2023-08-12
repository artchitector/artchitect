<i18n>
{
  "en": {
    "art": "Art"
  },
  "ru": {
    "art": "Картина"
  }
}
</i18n>
<template>
  <div class="card">
    <div class="card-content">
      <div class="media">
        <div class="columns">
          <div class="column is-one-fifth">
            <figure class="image is-64x64">
              <img :src="`data:image/png;base64, ${art.idea.seedEntropy.entropy.image}`" alt="odin's mind"/>
            </figure>
            <figure class="image is-64x64">
              <img :src="`data:image/png;base64, ${art.idea.seedEntropy.choice.image}`" alt="odin's memory"/>
            </figure>
          </div>
          <div class="column">
            <h1 class="title is-4">{{ $t('art') }} #{{ art.id }}</h1>
            <p class="subtitle is-6">{{ created }}</p>
          </div>
        </div>
      </div>
      <div class="content">
        <art-tags :art="art"/>
      </div>
    </div>
    <div class="card-image">
      <figure class="image is-2by3">
        <img :src="`/api/image/${art.id}/f`" :alt="`art_${art.id}`"/>
      </figure>
    </div>
  </div>
</template>

<script>
import {format} from "@/utils/date";

export default {
  name: "art-view",
  props: ['art'],
  computed: {
    created() {
      return format(this.art.createdAt, this.$i18n.locale)
    }
  }
}
</script>

<style scoped lang="scss">
  .card .media .columns {
    width: 100%;
    .column {
      padding-bottom: 0;
    }
    figure {
      display: inline-block;
    }
    img {
      image-rendering: pixelated;
    }
  }
</style>
