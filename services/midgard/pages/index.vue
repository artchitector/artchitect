<i18n>
{
  "en": {
    "title": "Artchitect - autonomous creative machine",
    "description": "Artchitect - it is an amazing autonomous creative machine capable of creating magnificent artworks inspired by the universe around us. In its continuous creativity , the machine receives inspiration from the natural entropy of the universe, represented as background light, and creates unique artworks without human participation."
  },
  "ru": {
    "title": "Artchitect - автономная творческая машина",
    "description": "Artchitect - удивительная автономная творческая машина, способная создавать великолепные картины, вдохновляясь окружающей нас Вселенной. В своем непрерывном творчестве машина черпает вдохновение из естественной энтропии Вселенной, представленной в виде фонового света, и создает уникальные произведения без участия человека. Архитектор создаёт новую картину каждые 60 секунд."
  }
}
</i18n>
<template>
  <div>
    <main-description/>
    <main-last99/>
  </div>
</template>
<script>
export default {
  head() {
    return {
      title: this.$t('title'),
      meta: [
        {hid: 'description', name: 'description', content: this.$t('description')},
        {property: 'og:title', content: this.$t('title')},
        {property: 'og:description', content: this.$t('description')},
        {property: 'og:type', content: 'image'},
        {property: 'og:image', content: `https://artchitect.space/icon196.png`}
      ]
    }
  }
}
</script>
<style lang="scss" scoped>

</style>
