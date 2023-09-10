<i18n>
{
  "en": {
    "title": "Fable of artchitect"
  },
  "ru": {
    "title": "Миф об artchitect"
  }
}
</i18n>
<template>
  <section class="content">
    <h1 class="is-size-4">{{ $t('title') }}</h1>
    <div class="notification" v-if="locale === 'ru'">
      Код Artchitect архитектурно и идеологически основан на сказочном рассказе по мотивам скандинавской мифологии.
      Все действующие лица (
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/odin.go"
         target="_blank">Один</a>,
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/loki.go"
         target="_blank">Локи</a>,
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/heimdallr.go"
         target="_blank">Хеймдалль</a> и
      <a href="https://github.com/artchitector/artchitect/tree/main/services/asgard/pantheon" target="_blank">другие</a>
      ) существуют в программном коде Artchitect, имеют свои роли, от их имени написаны комментарии о сути проекта и о
      технических деталях, использованных в artchitect.
      <div class="has-text-centered mt-3">
        <a href="https://github.com/artchitector/artchitect" class="button is-info is-light" target="_blank">исходный
          код artchitect на github</a>
      </div>
    </div>
    <div class="notification" v-else>
      artchitect source code is architecturally and ideologically based on a fairy tale inspired by Scandinavian
      mythology.
      All characters (
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/odin.go"
         target="_blank">Odin</a>,
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/loki.go"
         target="_blank">Loki</a>,
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/heimdallr.go"
         target="_blank">Heimdallr</a>,
      and <a href="https://github.com/artchitector/artchitect/tree/main/services/asgard/pantheon"
             target="_blank">others</a>
      ) exist within the artchitect's source code, have their roles, and
      comments about the project's essence and the technical details used in artchitect are written on their behalf.

      <div class="has-text-centered mt-3">
        <a href="https://github.com/artchitector/artchitect" class="button is-info is-light" target="_blank">artchitect
          source code on github</a>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: "loki-heading",
  props: ['locale']
}
</script>

<style lang="scss" scoped>

</style>
