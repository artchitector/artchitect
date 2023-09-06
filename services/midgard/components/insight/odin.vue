<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div class="wrapper">
    <template v-if="odin.painted">
      <insight-odin-result :odin="odin" :giving="giving" :entropy="entropy" @show="$emit('show', $event)"/>
    </template>
    <template v-else>
      <div class="wrapper-cell">
        <insight-odin-progress :odin="odin" :entropy="entropy" @show="$emit('show', $event)"/>
      </div>
      <div class="wrapper-cell">
        <insight-odin-arts v-if="giving" :giving="giving" @show="$emit('show', $event)"/>
      </div>
    </template>
  </div>
</template>

<script>
import ViewerViewer from "~/components/viewer/viewer.vue";

export default {
  name: "insight-odin",
  components: {ViewerViewer},
  props: ["odin", "entropy", "giving"]
}
</script>

<style scoped lang="scss">
.wrapper {
  width: 100%;
  height: 100%;
  padding: 2px 10px 10px 10px;

  letter-spacing: 1px;
  text-transform: uppercase;

  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;

  .wrapper-cell {
    align-self: center;
  }
}
</style>
