<template>
  <div class="image-container">
    <a :href="`/art/${art.id}`" @click.prevent="select()">
      <img :src="`/api/image/${art.id}/m`" :alt="`art_${artId}`"/>
    </a>
    <div class="likes">
      <span>
        <font-awesome-icon v-if="!art.Liked" icon="fa-solid fa-heart" class="has-color-base"/>
        <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
        {{art.Likes}}
      </span>
    </div>
    <div class="version">
      <span class="tag is-primary is-light">{{ art.Version }}</span>
    </div>
    <p v-if="!noTags" class="is-size-7 tags">{{ art.ideaWords.join(",") }}</p>
    <p class="is-size-7 info">id={{ art.id }}, <b>seed: {{ art.ideaSeed }}</b></p>
  </div>
</template>

<script>
export default {
  name: "art-complex",
  props: ['art', 'noTags'],
  methods: {
    select() {
      this.$emit('select')
    }
  }
}
</script>

<style lang="scss" scoped>
.image-container {
  position: relative;

  .version {
    position: absolute;
    right: 5px;
    bottom: 10px;
    opacity: 0.7;
  }

  .likes {
    position: absolute;
    left: 5px;
    bottom: 10px;
    opacity: 0.7;
    color: #d4d1c3;
    background-color: rgba(0, 0, 0, 0.6);
    padding: 0 2px 0 4px;
    border-radius: 3px;
  }

  a {
    display: block;
  }
}

.image-container p {
  visibility: visible;
  position: absolute;
  margin: auto;
  word-break: break-all;
  background-color: rgba(0, 0, 0, 0.6);
  padding: 5px;
}

.image-container p.tags {
  bottom: 30px;
  color: #d4d1c3;
  visibility: hidden;
}

.image-container:hover p.tags {
  visibility: visible;
}

.image-container p.info {
  top: 30px;
  width: 100%;
  text-align: center;
  color: #d4d1c3;
  visibility: hidden;
}

.image-container:hover p.info {
  visibility: visible;
}

</style>
