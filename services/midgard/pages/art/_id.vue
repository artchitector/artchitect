<i18n>
{
  "en": {
    "title": "Artchitect - art",
    "description": "Artchitect is autonomous creative machine making arts every 60 seconds",
    "error": "Error"
  },
  "ru": {
    "title": "Artchitect - картина",
    "description": "Artchitect это автономная творческая машина, создающая картины каждые 60 секунд",
    "error": "Ошибка"
  }
}
</i18n>
<template>
  <section>
    <div v-if="$fetchState.pending" class="has-text-centered">
      <common-loader/>
    </div>
    <div v-else-if="$fetchState.error" class="notification is-danger">
      {{ $t('error') }} {{ $fetchState.error.message }}
    </div>
    <art-view v-else-if="art" :art="art"/>
    <span v-else class="has-text-danger">load failed</span>
  </section>
</template>

<script>
export default {
  name: "art",
  head() {
    const artId = this.$route.params.id;
    return {
      title: this.$t('title') + ` #${artId}`,
      meta: [
        {hid: 'description', name: 'description', content: `Artchitect - Art #${artId}`},
        {property: 'og:title', content: `Artchitect - Art #${artId}`},
        {property: 'og:description', content: this.$t('description')},
        {property: 'og:type', content: 'image'},
        {property: 'og:image', content: `https://artchitect.space/api/image/artchitect-${artId}-m`}
      ]
    }
  },
  data() {
    return {
      art: null,
    }
  },
  async fetch() {
    const id = parseInt(this.$route.params.id)
    if (!id) {
      throw "id must be positive"
    }
    this.art = await this.$axios.$get(`/art/${id}`)
  }
}
</script>

<style scoped lang="scss">
section {
  padding: 10px
}
</style>
