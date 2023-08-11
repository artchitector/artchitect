<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div class="column-image">
    <NuxtLink :to="localePath(`/art/${artId}`)" class="has-text-info">
      <img :src="`/api/image/${artId}/s`"/>
    </NuxtLink>
    <div class="link-container">
      <NuxtLink :to="localePath(`/art/${artId}`)" class="has-test-info">#{{ artId }}</NuxtLink>
    </div>
  </div>
</template>

<script>
export default {
  name: "insight-odin-rnd",
  props: ['artId']
}
</script>

<style scoped lang="scss">
.column-image {
  text-align: center;
  margin: 5px 5px;
  position: relative;

  img {
    display: block;
    width: 100%;
  }

  .link-container {
    font-size: 12px;
    background-color: rgba(0, 0, 0, 0.5);
    padding: 0 3px;
  }
}
</style>
