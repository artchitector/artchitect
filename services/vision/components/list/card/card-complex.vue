<template>
  <div class="image-container">
    <a :href="`/art/${card.ID}`" @click.prevent="select()">
      <img :src="'/api/image/m/' + card.ID"/>
    </a>
    <div class="likes">
      <span>
        <font-awesome-icon v-if="!card.Liked" icon="fa-solid fa-heart" class="has-color-base"/>
        <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
        {{card.Likes}}
      </span>
    </div>
    <div class="version">
      <span class="tag is-primary is-light">{{ card.Version }}</span>
    </div>
    <p v-if="!noTags" class="is-size-7 tags">{{ card.Spell.Tags }}</p>
    <p class="is-size-7 info">id={{ card.ID }}, <b>seed: {{ card.Spell.Seed }}</b></p>
  </div>
</template>

<script>
export default {
  name: "card-complex",
  props: ['card', 'noTags'],
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
