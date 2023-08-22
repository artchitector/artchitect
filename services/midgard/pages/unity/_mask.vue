<i18n>
{
  "en": {
    "title": "Artchitect - unity",
    "subtitle100": "hundred",
    "subtitle1000": "1K unity",
    "subtitle10000": "10K unity",
    "subtitle100000": "100K unity",
    "description1000": "Every unity below contains 100 arts",
    "description10000": "Every unity below contains 1000 arts",
    "description100000": "Every unity below contains 10000 arts",
    "hundred_description_left": "First art of this hundred was created",
    "hundred_description_middle": "and hundred was finished"
  },
  "ru": {
    "title": "Artchitect - единство",
    "subtitle100": "сотня картин",
    "subtitle1000": "тысячное единство",
    "subtitle10000": "10-тысячное единство",
    "subtitle100000": "100-тысячное единство",
    "description1000": "В каждом из представленных ниже единств 100 картин",
    "description10000": "В каждом из представленных ниже единств 1000 картин",
    "description100000": "В каждом из представленных ниже единств 10000 картин",
    "unity_first_part": "Первая картина в этом единстве была написана",
    "unity_no_first_art": "В этом единстве первая картина еще не написана",
    "unity_last_part": "а единство было завершено полностью",
    "unity_no_last_art": "и единство еще не завершено",
    "hundred_description_left": "Первая картина этой сотни была написана",
    "hundred_description_middle": "а вся сотня была завершена"
  }
}
</i18n>

<template>
  <section>
    <h3 class="has-text-centered is-size-4 mb-4">{{ subtitle }}</h3>
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
    <template v-else-if="data.children.length">
      <div v-if="data" class="notification has-text-centered">
        <div>{{ dateRange }}</div>
        <div>{{ $t(`description${data.unity.rank}`) }}</div>
      </div>
      <!--   TODO Unity надо переписать заново. Все страницы и компоненты-->
      <unity-list :unities="data.children" visible-count="10" cards-in-column="2"/>
    </template>
    <template v-else-if="data.arts.length">
      <div v-if="data" class="notification has-text-centered">
        <div>{{ hundredRange }}</div>
      </div>
      <common-art-list :arts="data.arts" row-size="5" initial-visible-count="100" :show-words="false"/>
    </template>
  </section>
</template>

<script>
import UnityList from "@/components/unity/unity-list.vue";
import {format, duration} from "@/utils/date"

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
      if (!this.data) {
        return "";
      }
      return this.$t(`subtitle${this.data.unity.rank}`) + ' U' + this.data.unity.mask
    },
    dateRange() {
      let first, last
      if (this.data.firstArtTime) {
        const dt = format(this.data.firstArtTime, this.$i18n.locale)
        first = `${this.$t('unity_first_part')} ${dt}`
      } else {
        return this.$t('unity_no_first_art')
      }
      if (this.data.lastArtTime) {
        const dur = duration(this.data.firstArtTime, this.data.lastArtTime, this.$i18n.locale)
        last = `${this.$t('unity_last_part')} ${dur}`
      } else {
        last = this.$t('unity_no_last_art')
      }

      return `${first}, ${last}.`
    },
    hundredRange() {
      if (!this.data.firstArtTime) {
        return ""
      }
      const dt = format(this.data.firstArtTime, this.$i18n.locale)
      const first = `${this.$t('hundred_description_left')} ${dt}`
      if (!this.data.lastArtTime) {
        return first
      }

      const dur = duration(this.data.firstArtTime, this.data.lastArtTime, this.$i18n.locale)
      return `${first}, ${this.$t('hundred_description_middle')} ${dur}.`

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
