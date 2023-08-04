<template>
  <div class="image-container">
    <a :href="`/card/${cardId}`" @click.prevent="select()">
      <img :src="'/api/image/m/' + cardId"/>
    </a>
  </div>
</template>

<script>
export default {
  name: "card-simple",
  props: ['cardId'],
  methods: {
    select() {
      this.$emit('select')
    }
  }
}
</script>

<style scoped>

</style>
