<i18n>
{
  "en": {
    "show_more": "show more..."
  },
  "ru": {
    "show_more": "показать больше..."
  }
}
</i18n>
<template>
  <div>
    <common-viewer ref="viewer"/>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="art in line">
        <common-art-complex v-if="!!art && typeof art === 'object'"
                            :art="art"
                            @select="showViewer(art.id)"
                            :show-words="showWords"/>

        <common-art-simple v-else-if="art"
                           :art-id="art"
                           @select="showViewer(art)"/>
      </div>
    </div>
    <div v-if="showLoadMore" class="has-text-centered">
      <button class="button" @click.prevent="showMore">{{ $t('show_more') }}</button>
    </div>
  </div>
</template>
<script>
export default {
  name: "common-art-list",
  props: [
    'arts',
    'rowSize',
    'initialVisibleCount',
    'showWords',
  ],
  data() {
    return {
      visibleCount: parseInt(this.initialVisibleCount) || -1,
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
          chunk.push(null) // TODO Odin: тут какой-то баг. эти элементы отображаются как art_null
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
  },
  methods: {
    showMore() {
      this.visibleCount += parseInt(this.initialVisibleCount)
      console.log(this.visibleCount)
    },
    showViewer(artID) {
      const ids = [];
      this.arts.forEach((art) => {
        if (typeof art === 'object') {
          ids.push(art.id)
        } else {
          ids.push(art)
        }
      })
      this.$refs.viewer.show(ids, artID)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
