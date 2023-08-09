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
    <template v-else-if="data.Type === 'unity'">
      <unity-list :unities="data.Unities" visible-count="10" cards-in-column="2"/>
    </template>
    <template v-else-if="data.Type === 'cards'">
      <p class="has-text-centered">total: {{data.Cards.length}}</p>
      <artlist :arts="data.Cards" cards-in-column="5" card-size="s" visible-count="50" no-tags="true"/>
    </template>
  </section>
</template>

<script>
import UnityList from "@/components/unity/unity-list.vue";
import Artlist from "@/components/list/artlist.vue";

export default {
  components: {Artlist, UnityList},
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
