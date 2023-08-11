<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div class="rrnd">
    <div class="rrnd-img-link">
      <NuxtLink :to="localePath(`/art/${artId}`)">
        <img :src="`/api/image/${artId}/s`"/>
      </NuxtLink>
    </div>
    <NuxtLink class="rrnd-number-link has-text-info" :to="localePath(`/art/${artId}`)">#{{ artId}}</NuxtLink>
  </div>
</template>

<script>
export default {
  name: "insignt-odin-rrnd",
  props: ["artId"]
}
</script>

<style lang="scss" scoped>
.rrnd {
  height: 100%;
  background-color: rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-items: center;
  justify-content: center;
  padding-top: 2px;

  .rrnd-img-link {
    flex-grow: 1;
    overflow: hidden;
  }

  img {
    display: block;
    max-height: 100%;
  }

  .rrnd-number-link {
    font-size: 12px;
  }
}
</style>

