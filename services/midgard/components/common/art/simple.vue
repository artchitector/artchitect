<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div>
    <a :href="localePath(`/art/${this.artId}`)" @click.prevent="select">
      <img :src="`/api/image/artchitect-${artId}-m`" :alt="`art_${artId}`"/>
    </a>
  </div>
</template>
<script>
export default {
  name: "common-art-simple",
  props: ["artId"],
  methods: {
    select() {
      this.$emit('select')
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
