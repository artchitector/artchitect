<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <span>
    {{ $t('seed') }}={{ odin.seed }}
    &nbsp;
    <b-tooltip position="is-bottom" multilined>
      <img v-if="odin.seedEntropyImageEncoded"
           :src="`data:image/png;base64, ${odin.seedEntropyImageEncoded}`"
           alt="odin's mind"/>
      <template v-slot:content>
        <img v-if="odin.seedEntropyImageEncoded"
             :src="`data:image/png;base64, ${odin.seedEntropyImageEncoded}`"
             style="width:64px;height:64px;image-rendering: pixelated;"
             alt="odin's mind"/>
      </template>
    </b-tooltip>
    &nbsp;
    <b-tooltip position="is-bottom" multilined>
      <img v-if="odin.seedChoiceImageEncoded"
           :src="`data:image/png;base64, ${odin.seedChoiceImageEncoded}`"
           alt="odin's mind"/>
      <template v-slot:content>
        <img v-if="odin.seedChoiceImageEncoded"
             :src="`data:image/png;base64, ${odin.seedChoiceImageEncoded}`"
             style="width:64px;height:64px;image-rendering: pixelated;"
             alt="odin's choice"/>
      </template>
    </b-tooltip>
  </span>
</template>

<script>
export default {
  name: "insight-odin-progress-seed",
  props: ['odin']
}
</script>

<style scoped lang="scss">

</style>
