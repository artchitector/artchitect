<i18n>
{
  "en": {
    "page_title": "Artchitect - fable"
  },
  "ru": {
    "page_title": "Artchitect - сказка"
  }
}
</i18n>
<template>
  <div>
    <loki-heading :locale="locale"/>
    <loki-fable :locale="locale"/>
    <loki-number-one :locale="locale"/>
  </div>
</template>
<script>
export default {
  name: "loki",
  head() {
    return {
      title: this.$t('page_title')
    }
  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  }
}
</script>
<style lang="scss" scoped></style>
