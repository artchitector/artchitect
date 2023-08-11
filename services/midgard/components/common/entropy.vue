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
      <span class="is-size-7">{{ $t('title') }}</span>
      <div>
        <NuxtLink :to="localePath('entropy')">
          <img v-if="entropy.entropy.image"
               :src="`data:image/png;base64, ${entropy.entropy.image}`"
               alt="odin's mind"/>
        </NuxtLink>
        <NuxtLink :to="localePath('entropy')">
          <img v-if="entropy.choice.image"
               :src="`data:image/png;base64, ${entropy.choice.image}`"
               alt="odin's choice"/>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "common-entropy",
  props: ['entropy']
}
</script>

<style scoped lang="scss">
.entropy-sensor-inner {
  display: inline-block;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 0 7px;

  img {
    width: 64px;
    height: 64px;
    image-rendering: pixelated;
  }
}
</style>
