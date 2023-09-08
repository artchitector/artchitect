<i18n>
{
  "en": {
    "nothing": "you're not liked anything. please, like some art."
  },
  "ru": {
    "nothing": "вам еще ничего не нравилось. пожалуйста, поставьте лайк одной из картин."
  }
}
</i18n>
<template>
  <div v-if="loading" class="has-text-centered">
    <common-loader/>
  </div>
  <div v-else-if="error" class="notification is-danger is-light">
    {{ error.message }}
  </div>
  <div v-else-if="liked.length === 0" class="notification is-link is-light">
    {{ $t('nothing') }}
  </div>
  <div v-else>
    <common-art-list :arts="liked" :initial-visible-count="50" :row-size="5" :show-words="false"/>
  </div>
</template>
<script>
export default {
  name: "liked-list",
  data() {
    return {
      loading: false,
      liked: [],
      error: null,
    }
  },
  async fetch() {
    this.loading = true;
    try {
      this.liked = await this.$axios.$get("/liked")
    } catch (e) {
      this.error = e
    } finally {
      this.loading = false
    }
  }
}
</script>
<style lang="scss" scoped>

</style>
