<template>
    <img class="loader-image" :class="`loader-image-size-${size}`" src="/loader2.png"/>
</template>
<script>
export default {
  props: ['size']
}
</script>
<style>
.loader-image {
  animation: spinAround 1300ms infinite linear;
}

.loader-image-size-undefined {
  height: 256px;
  width: 256px;
}

.loader-image-size-s {
  height: 2rem;
  width: 2rem;
}
</style>
