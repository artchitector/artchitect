<i18n>
{
  "en": {
    "title": "Artchitect - art",
    "card": "Art",
    "loading": "loading...",
    "error": "error: ",
    "version": "version",
    "seed": "Seed-number",
    "tags": "Keywords",
    "full_size": "view full size",
    "description": "Artchitect is autonomous creative machine making arts every 60 seconds"
  },
  "ru": {
    "title": "Artchitect - картина",
    "card": "Картина",
    "loading": "загрузка...",
    "error": "ошибка: ",
    "version": "версия",
    "seed": "Seed-номер",
    "tags": "Ключевые слова",
    "full_size": "смотреть в полном размере",
    "description": "Artchitect это автономная творческая машина, создающая картины каждые 60 секунд"
  }
}
</i18n>
<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      {{ $t('loading') }}
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $t('error') }} {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">{{ $t('card') }} #{{ card.ID }}</h1>
      <p><span class="tag is-primary is-light">{{ $t('version') }} {{ card.Version }}</span></p>
      <p>{{ created }}</p>
      <p>{{ $t('seed') }} = {{ card.Spell.Seed }}</p>
      <p class="tags">{{ $t('tags') }} = <i>{{ card.Spell.Tags }}</i></p>
      <p class="has-text-centered">
        <a :href="fullSizeUrl" target="_blank" class="is-size-7">{{ $t('full_size') }}</a>
      </p>
      <div class="image-container">
        <img :src="`/api/image/f/${card.ID}`"/>
        <liker :dream_id="card.ID" class="control-like"/>
      </div>
    </div>
  </section>

</template>
<script>
import moment from "moment"
import Liker from "@/components/utils/liker.vue";

export default {
  components: {Liker},
  head() {
    const artId = this.$route.params.id;
    return {
      title: this.$t('title') + ` #${artId}`,
      meta: [
        {hid: 'description', name: 'description', content: `Artchitect - Art #${artId}`},
        {property: 'og:title', content: `Artchitect - Art #${artId}`},
        {property: 'og:description', content: this.$t('description')},
        {property: 'og:type', content: 'image'},
        {property: 'og:image', content: `https://artchitect.space/api/image/m/${artId}`}
      ]
    }
  },
  data() {
    return {
      card: null,
    }
  },
  computed: {
    created() {
      return moment(this.card.CreatedAt).format("YYYY MMM Do HH:mm:ss")
    },
    fullSizeUrl() {
      const part = Math.floor(this.card.ID / 10000) + 1
      return `${process.env.STORAGE_URL}/art/${part}/art-${this.card.ID}.jpg`
    }
  },
  async fetch() {
    const id = parseInt(this.$route.params.id);
    if (!id) {
      throw "id must be positive integer"
    }
    this.card = await this.$axios.$get(`/card/${id}`)
  }
}
</script>
<style lang="scss" scoped>
p.tags {
  word-wrap: break-word;
  word-break: break-all;
  overflow: hidden;
}

.image-container {
  position: relative;

  .control-like {
    position: absolute;
    left: 50%;
    bottom: 20%;
    z-index: 3;
    margin-left: -20px;
    font-size: 48px;
    opacity: 0.7;
    filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
  }
}
</style>
