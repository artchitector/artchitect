<i18n>
{
  "en": {
    "title": "Artchitect - request"
  },
  "ru": {
    "title": "Artchitect - запрос"
  }
}
</i18n>
<template>
  <div>
    <section>
      <prayer/>
    </section>
  </div>
</template>
<script>
import Prayer from "@/components/pray/prayer.vue";

export default {
  components: {Prayer},
  head() {
    return {
      title: this.$t('title')
    }
  },
}
</script>
