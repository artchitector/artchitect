<template>
  <div>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="unity in line">
        <unity v-if="!!unity" :unity="unity"/>
      </div>
    </div>
  </div>
</template>

<script>
import Unity from "@/components/unity/unity.vue";

export default {
  name: "unity-list",
  components: {Unity},
  props: ["unities", "cardsInColumn", "visibleCount"],
  data() {
    return {
      currentVisible: -1
    }
  },
  computed: {
    lines() {
      let unities = []
      if (this.currentVisible === -1) {
        unities = []
      } else if (this.currentVisible === 0) {
        unities = this.unities
      } else {
        unities = this.unities.slice(0, this.currentVisible)
      }
      const chunkSize = parseInt(this.cardsInColumn)
      const chunks = [];
      for (let i = 0; i < unities.length; i += chunkSize) {
        let chunk = unities.slice(i, i + chunkSize)
        for (let j = chunk.length; j < this.cardsInColumn; j++) {
          chunk.push(null)
        }
        chunks.push(chunk)
      }
      return chunks
    }
  },
  mounted() {
    if (!!this.visibleCount) {
      this.currentVisible = parseInt(this.visibleCount)
    } else {
      this.currentVisible = 0
    }
  }
}
</script>

<style scoped>

</style>
