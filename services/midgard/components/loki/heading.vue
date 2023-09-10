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
    <div class="notification">
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
  </section>
</template>

<script>
export default {
  name: "loki-heading"
}
</script>

<style lang="scss" scoped>

</style>
