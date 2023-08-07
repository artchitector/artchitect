<template>
  <div class="column-image">
    <NuxtLink v-if="dream_id > 0" :to="localePath(`/art/${dream_id}`)" class="has-text-info">
      <img :src="`/api/image/s/${dream_id}`"/>
    </NuxtLink>
    <img v-else src="/in-progress.jpeg"/>

    <div v-if="dream_id > 0" class="link-container">
      <NuxtLink v-if="dream_id > 0" :to="localePath(`/art/${dream_id}`)" class="has-text-info">#{{ dream_id }}</NuxtLink>
    </div>
  </div>
</template>

<script>
export default {
  name: "rnd",
  props: ["dream_id"]
}
</script>

<style lang="scss" scoped>
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
