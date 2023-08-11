<i18n>
{
  "en": {
    "title": "Artchitect's insight 👁"
  },
  "ru": {
    "title": "Озарение  Artchitect 👁"
  }
}
</i18n>
<template>
        <insight-main/>
</template>

<script>
export default {
  name: "insight",
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
