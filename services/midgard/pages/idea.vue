<i18n>
{
  "en": {
    "title": "Artchitect - idea"
  },
  "ru": {
    "title": "Artchitect - идея"
  }
}
</i18n>
<template>
  <idea-main/>
</template>
<script>


export default {
  name: "idea",
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      entropy: null
    }
  },
  mounted() {
    
  },
  beforeDestroy() {

  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  }
}
</script>
