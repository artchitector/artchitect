<i18n>
{
  "en": {
    "title": "ENTROPY SENSOR"
  },
  "ru": {
    "title": "ДАТЧИК ЭНТРОПИИ"
  }
}
</i18n>
<template>
  <div class="has-text-centered entropy-sensor">
    <div class="entropy-sensor-inner">
      <span class="is-size-7">{{ $t('title') }}</span><br/>
      <NuxtLink :to="localePath('entropy')">
        <img v-if="entropy.ImagesEncoded.entropy !== null"
             :src="`data:image/png;base64, ${entropy.ImagesEncoded.entropy}`"
             style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading entropy stream"/>
      </NuxtLink>
      <NuxtLink :to="localePath('entropy')">
        <img v-if="entropy.ImagesEncoded.choice !== null"
             :src="`data:image/png;base64, ${entropy.ImagesEncoded.choice}`"
             style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading choice stream"/>
      </NuxtLink>
    </div>
  </div>
</template>

<script>
export default {
  props: ['entropy']
}
</script>

<style scoped lang="scss">
.entropy-sensor-inner {
  display: inline-block;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 0 7px;
}
</style>
