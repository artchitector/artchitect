<i18n>
{
  "en": {
    "title": "Artchitect's heart ❤️"
  },
  "ru": {
    "title": "Сердце  Artchitect ❤️"
  }
}
</i18n>
<template>
        <flexheart-core/>
</template>

<script>
export default {
  name: "heart",
  layout: "simple",
  head() {
    return {
      title: this.$t('title')
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
