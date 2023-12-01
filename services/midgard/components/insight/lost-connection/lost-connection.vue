<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div class="p-6">
    <div class="notification is-danger">
      Соединение утеряно
    </div>
    <div class="notification is-danger is-light">
      В данный момент нет связи с Асгардом - core-сервером Artchitect, который пишет картины.
      <br/>
      Он отключился или потерял интернет, причины пока неизвестны. Новые картины сейчас не пишутся.
    </div>
    <div class="notification is-primary is-light">
      Зато написанные картины вы всегда можете найти в разделе "
      <NuxtLink :to="localePath('unity')">Единства</NuxtLink>
      ".
    </div>
    <div v-if="lastArtID">
      <div class="has-text-centered">
        ПОСЛЕДНЯЯ КАРТИНА БЫЛА
        <NuxtLink :to="localePath(`/art/${lastArtID}`)" class="has-text-info">#{{ lastArtID }}</NuxtLink>
      </div>
      <div class="has-text-centered">
      <NuxtLink :to="localePath(`/art/${lastArtID}`)" class="has-text-info">
        <img :src="`/api/image/artchitect-${lastArtID}-f`"/>
      </NuxtLink>
    </div>
  </div>
  </div>
</template>

<script>
export default {
  name: "insight-lost-connection",
  data() {
    return {
      lastArtID: null,
    }
  },
  async fetch() {
    try {
      this.lastArtID = await this.$axios.$get("/art/max")
    } catch (err) {
      console.error("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОСЛЕДНЮЮ КАРТИНУ", err)
    }
  },
}
</script>

<style scoped lang="scss">

</style>
