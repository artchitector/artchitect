<template>
  <div class="control-like">
    <font-awesome-icon v-if="liked && liked.error"
                       icon="fa-solid fa-triangle-exclamation"
                       :title="liked.error.message"/>
    <a v-else href="#" @click.prevent="like()">
      <font-awesome-icon v-if="!liked || !liked.liked" icon="fa-solid fa-heart" class="has-color-base"/>
      <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
    </a>
  </div>
</template>

<script>
export default {
  name: "liker",
  props: ["dream_id"],
  data() {
    return {
      liked: {
        liked: false,
        error: null,
      }
    }
  },
  mounted() {
    this.initLiked()
  },
  methods: {
    async like() {
      try {
        let like = await this.$axios.$post("/like", {
          card_id: this.dream_id,
        })
        this.$emit('liked', like)
        this.liked = {
          id: like.ID,
          liked: like.Liked,
        };
      } catch (e) {
        console.error(e)
        this.liked = {
          error: e
        };
      }
    },
    async initLiked() {
      try {
        let like = await this.$axios.$get(`/liked/${this.dream_id}`)
        this.liked.liked = like.Liked
      } catch (e) {
        console.error(e)
        this.liked.error = e
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.control-like {
  font-size: 48px;
  opacity: 0.7;
  filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
}
</style>
