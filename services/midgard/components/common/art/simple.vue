<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div>
    <a :href="localePath(`/art/${this.artId}`)" @click.prevent="select">
      <img :src="`/api/image/${artId}/m`" :alt="`art_${artId}`"/>
    </a>
  </div>
</template>
<script>
export default {
  name: "common-art-simple",
  props: ["artId"],
  methods: {
    select() {
      console.warn(`need viewer art#${this.artId}`)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
