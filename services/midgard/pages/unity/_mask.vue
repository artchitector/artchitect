<i18n>
{
  "en": {
    "title": "Artchitect - unity",
    "subtitle": "unity"
  },
  "ru": {
    "title": "Artchitect - единство",
    "subtitle": "единство"
  }
}
</i18n>

<template>
  <section>
    <h3 class="has-text-centered is-size-4">{{ subtitle }}</h3>
    <template v-if="$fetchState.pending">
      <div class="has-text-centered">
        <loader/>
      </div>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">
        {{ $fetchState.error.message }}
      </div>
    </template>
    <template v-else-if="data.children.length">
      <!--   TODO Unity надо переписать заново. Все страницы и компоненты-->
      <unity-list :unities="data.children" visible-count="10" cards-in-column="2"/>
    </template>
    <template v-else-if="data.arts.length">
      <p class="has-text-centered">total: {{ data.arts.length }}</p>
      <common-art-list :arts="data.arts" row-size="5" initial-visible-count="100" :show-words="false"/>
    </template>
  </section>
</template>

<script>
import UnityList from "@/components/unity/unity-list.vue";

export default {
  components: {UnityList},
  head() {
    let mask = this.$route.params.mask
    return {
      title: `${this.$t('title')} U${mask}`
    }
  },
  data() {
    return {
      data: null
    };
  },
  computed: {
    subtitle() {
      let mask = this.$route.params.mask
      return `${this.$t('subtitle')} U${mask}`
    }
  },
  async fetch() {
    let mask = this.$route.params.mask
    this.data = await this.$axios.$get(`/unity/${mask}`)
  }
}
</script>

<style scoped>

</style>
