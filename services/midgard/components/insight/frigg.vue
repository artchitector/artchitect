<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <div class="wrapper">
    <insight-frigg-unity class="wrapper-cell" :state="frigg"/>
  </div>
</template>

<script>
export default {
  name: "insight-frigg",
  props: ["frigg", "entropy"]
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
