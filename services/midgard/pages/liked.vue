<i18n>
{
  "en": {
    "title": "Artchitect - liked",
    "subtitle": "you liked"
  },
  "ru": {
    "title": "Artchitect - нравится",
    "subtitle": "вам нравится"
  }
}
</i18n>
<template>
  <section>
    <h1 class="has-text-centered is-size-4">{{ $t('subtitle') }}</h1>
    <template v-if="$fetchState.pending">
      <div class="has-text-centered">
        <loader/>
      </div>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">{{ $fetchState.error.message }}</div>
    </template>
    <template v-else>
      <cardlist :cards="liked" cards-in-column="5" card-size="m" visible-count="30" @liked="onLiked"/>
    </template>
  </section>
</template>
<script>
import Cardlist from "@/components/list/cardlist.vue";

export default {
  name: "liked",
  components: {Cardlist},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      liked: [],
    }
  },
  async fetch() {
    if (process.client) {
      this.liked = await this.$axios.$get('/liked')
    }
  },
  methods: {
    onLiked(dt) {
      if (dt.Liked) {
        // if like added, then prepend card to list
        console.log('added like', dt.CardID)
        this.liked.unshift(dt.CardID)
      } else {
        // if like removed, then remove card from list
        let idx = this.liked.findIndex(x => x === dt.CardID)
        this.liked.splice(idx, 1)
      }
    }
  }
}
</script>

<style scoped>

</style>
