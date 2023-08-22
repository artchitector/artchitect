<i18n>
{
  "en": {
    "title": "Artchitect - unity",
    "subtitle": "100K unities",
    "description": "Hundred of arts have become a 100-unity, thousands of arts have become a 1000-unity, and the same goes for 10,000 and 100,000 unities. Unities is a useful browser to view Architect's arts from a global perspective.",
    "subdescription": "Every unity below contains 100 000 arts"
  },
  "ru": {
    "title": "Artchitect - единство",
    "subtitle": "100-тысячные единства",
    "description": "Сотня картин становится 100-единством, тысяча картин становится 1000-единством, и то же самое касается 10 000 и 100 000 единств. Единства - это полезный браузер, позволяющий увидеть картины Архитектора с глобальной точки зрения.",
    "subdescription": "В каждом представленном ниже единстве 100 тысяч картин"
  }
}
</i18n>

<template>
  <div>
    <section>
      {{ $t('description') }}
    </section>
  <section>
    <h3 class="has-text-centered is-size-4 mb-4">{{ $t('subtitle') }}</h3>
    <template v-if="$fetchState.pending">
      <div class="has-text-centered">
        <common-loader/>
      </div>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">
        {{ $fetchState.error.message }}
      </div>
    </template>
    <template v-else>
      <div class="notification mt-3">
        {{ $t('subdescription')}}
      </div>
      <unity-list :unities="unities" visible-count="50" cards-in-column="2"/>
    </template>
  </section>
  </div>
</template>

<script>
import UnityList from "@/components/unity/unity-list.vue";

export default {
  name: "unity",
  components: {UnityList},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      unities: []
    };
  },
  async fetch() {
    this.unities = await this.$axios.$get("/unity")
  }
}
</script>

<style scoped>

</style>
