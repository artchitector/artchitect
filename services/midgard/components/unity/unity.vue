<template>
  <div class="image-container">
    <NuxtLink :to="localePath(`/unity/${unity.mask}`)">
      <img :src="`/api/uimage/${unity.mask}/${unity.version}/m`"/>
    </NuxtLink>
    <p class="unity-title">U{{ unity.mask }}</p>
  </div>
</template>

<script>
export default {
  name: "unity",
  props: ['unity']
}
</script>

<style lang="scss" scoped>
.image-container {
  position: relative;

  p.unity-title {
    position: absolute;
    margin: auto;
    background-color: rgba(0, 0, 0, 0.6);
    padding: 5px;
    top: 30px;
    width: 100%;
    text-align: center;
    color: #d4d1c3;
  }
}
</style>
