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
    <loki-heading/>
    <loki-fable/>
    <loki-number-one/>
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
}
</script>
<style lang="scss" scoped></style>
