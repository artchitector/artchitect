<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="art in line">
        <common-art-simple :art-id="art.id"/>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "common-art-list",
  props: [
    'arts',
    'rowSize',
    'initialVisibleCount'
  ],
  data() {
    return {
      visibleCount: this.initialVisibleCount || -1,
    }
  },
  computed: {
    lines() {
      let arts = []
      if (this.visibleCount === -1) {
        arts = []
      } else if (this.visibleCount === 0) {
        arts = this.arts
      } else {
        arts = this.arts.slice(0, this.visibleCount)
      }

      const chunkSize = parseInt(this.rowSize)
      const chunks = []
      for (let i = 0; i < arts.length; i += chunkSize) {
        let chunk = arts.slice(i, i + chunkSize)
        for (let j = chunk.length; j < this.rowSize; j++) {
          chunk.push(null)
        }
        chunks.push(chunk)
      }
      return chunks
    },
    isComplex() {
      return !!this.arts.length && typeof this.arts[0] === 'object'
    },
    showLoadMore() {
      return this.visibleCount > 0 && this.visibleCount < this.arts.length
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
