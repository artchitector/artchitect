<template>
  <div class="art-complex">
    <a :href="localePath(`/art/${art.id}`)" @click.prevent="select()">
      <img :src="`/api/image/${art.id}/m`" :alt="`art_${art.id}`"/>
    </a>
    <div class="version">
      <span class="tag is-primary is-light">{{ art.version }}</span>
    </div>
    <p class="is-size-7 words">
      {{ art.ideaWords.join(',') }}
    </p>
    <p class="is-size-7 info">
      id={{ art.id }}, <b>seed: {{ art.ideaSeed }}</b>
    </p>
  </div>
</template>

<script>
// Odin: Карточка для картины может быть:
// - simple.vue (только картинка).
// - complex.vue (с доп. данными, словами, seed-номером)
export default {
  name: "common-art-complex",
  props: ['art'],
  methods: {
    select() {
      this.$emit('select')
    }
  }
}
</script>

<style scoped lang="scss">
.art-complex {
  position: relative;

  &:hover p.words, &:hover p.info { // при наведении на блок появляются его внутренние элементы
    visibility: visible;

  }

  .version {
    position: absolute;
    right: 5px;
    bottom: 10px;
    opacity: 0.7;
  }

  a {
    display: block;
  }

  p {
    visibility: visible;
    position: absolute;
    margin: auto;
    word-break: break-all;
    background-color: rgba(0, 0, 0, 0.6);
    padding: 5px;

    &.words {
      bottom: 30px;
      color: #d4d1c3;
      visibility: hidden;
    }

    &.info {
      top: 30px;
      width: 100%;
      text-align: center;
      color: #d4d1c3;
      visibility: hidden;
    }
  }
}
</style>
