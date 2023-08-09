<template>
  <div>
    <viewer ref="viewer" @liked="onLiked"/>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="art in line">
        <div v-if="!art"></div>
        <art-complex v-else-if="isComplex" :art="art" @select="select(art.ID)" :no-tags="noTags"/>
        <art-simple v-else :art-id="art" @select="select(art)"/>
      </div>
    </div>
    <div v-if="showLoadMore" class="has-text-centered">
      <button class="button" @click.prevent="showMore">show more...</button>
    </div>
  </div>
</template>

<script>
import ArtComplex from "~/components/list/art/art-complex.vue";
import ArtSimple from "~/components/list/art/art-simple.vue";

export default {
  name: "artlist",
  components: {ArtSimple, ArtComplex},
  props: [
    'arts',
    'artsInColumn',
    'artSize',
    'visibleCount', // how many arts show at start of component
    'noTags'
  ],
  data() {
    return {
      currentVisible: -1,
    }
  },
  computed: {
    lines() {
      let arts = []
      if (this.currentVisible === -1) {
        arts = []
      } else if (this.currentVisible === 0) {
        arts = this.arts
      } else {
        arts = this.arts.slice(0, this.currentVisible)
      }
      const chunkSize = parseInt(this.artsInColumn)
      const chunks = []
      for (let i = 0; i < arts.length; i += chunkSize) {
        let chunk = arts.slice(i, i + chunkSize)
        for (let j = chunk.length; j < this.artsInColumn; j++) {
          chunk.push(null)
        }
        chunks.push(chunk)
      }
      return chunks
    },
    isComplex() {
      return typeof this.arts[0] === 'object';
    },
    showLoadMore() {
      return this.currentVisible > 0 && this.currentVisible < this.arts.length;
    }
  },
  mounted() {
    if (!!this.visibleCount) {
      this.currentVisible = parseInt(this.visibleCount)
    } else {
      this.visibleCount = 0
    }
  },
  methods: {
    select(artID) {
      const ids = [];
      const isComplex = this.isComplex
      this.arts.forEach((art) => {
        if (isComplex) {
          ids.push(art.ID)
        } else {
          ids.push(art)
        }
      })
      this.$refs.viewer.show(ids, artID)
    },
    showMore() {
      this.currentVisible += parseInt(this.visibleCount)
    },
    onLiked(data) {
      this.$emit("liked", data)
    }
  }
}
</script>

<style scoped>

</style>
