<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <img class="loader-image" :class="`loader-image-size-${size}`" src="/loader2.png" alt="loading..."/>
</template>
<script>
export default {
  name: "common-loader",
  props: ["size"]
}
</script>
<style lang="scss" scoped>
.loader-image {
  animation: spinAround 1300ms infinite linear;
  &.loader-image-size-undefined {
    height: 256px;
    width: 256px;
  }
  &.loader-image-size-l {
    height: 128px;
    width: 128px;
  }
  &.loader-image-size-s {
    height: 2rem;
    width: 2rem;
  }
}
</style>
