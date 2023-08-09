<template>
  <div class="image-container">
    <a :href="`/art/${artId}`" @click.prevent="select()">
      <img :src="`/api/image/${artId}/m`" :alt="`art_${artId}`"/>
    </a>
  </div>
</template>

<script>
export default {
  name: "art-simple",
  props: ['artId'],
  methods: {
    select() {
      this.$emit('select')
    }
  }
}
</script>

<style scoped>

</style>
