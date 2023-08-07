<i18n>
{
  "en": {
    "created_art_is": "created art is"
  },
  "ru": {
    "created_art_is": "создана картина"
  }
}
</i18n>
<template>
  <div class="result-wrapper">
    <entropy v-if="entropy" :entropy="entropy"/>
    <div class="enjoy-progress">
      <p class="is-size-5 has-text-success has-text-centered">{{$t('created_art_is')}}
        <NuxtLink :to="localePath(`/art/${dream_id}`)" class="has-text-info">#{{ dream_id }}</NuxtLink>
      </p>
      <progress class="progress is-warning" :value="progress" max="100">-</progress>
    </div>
    <div class="image-container">
      <NuxtLink :to="localePath(`/art/${dream_id}`)" class="has-text-info">
        <img :src="`/api/image/f/${dream_id}`"/>
      </NuxtLink>
      <div class="control-like">
        <liker :dream_id="dream_id"/>
      </div>
    </div>
    <div class="random-four">
      <div v-for="f in rndFour" class="random-four-item">
        <rrnd class="is-max-width" :f="f"/>
      </div>
    </div>
  </div>
</template>

<script>
import Liker from "@/components/utils/liker.vue";
import Rnd from "@/components/flexheart/creation/rnd.vue";
import Rrnd from "@/components/flexheart/creation/rrnd.vue";
import Entropy from "@/components/entropy/entropy.vue";

export default {
  name: "result",
  components: {Entropy, Rrnd, Rnd, Liker},
  props: ["dream_id", "totalEnjoyTime", "currentEnjoyTime", "rndFour", "entropy"],
  computed: {
    progress() {
      if (!this.totalEnjoyTime || !this.currentEnjoyTime) {
        return 0;
      }
      const progress = this.currentEnjoyTime / this.totalEnjoyTime;
      return Math.floor(progress * 100);
    },
  }

}
</script>

<style lang="scss" scoped>
.result-wrapper {
  text-align: center;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;

  .image-container {
    position: relative;
    flex-grow: 2;
    overflow: hidden;

    img {
      max-height: 100%
    }

    .control-like {
      position: absolute;
      left: 50%;
      bottom: 10%;
      z-index: 3;
      transform: translate(-50%, -10%);
    }
  }

  .random-four {
    flex-grow: 1;
    flex-shrink: 0.3;
    overflow: hidden;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: stretch;
    gap: 10px;
    .random-four-item {
      img {
        max-height: 100%;
      }
    }
  }

  .is-max-width {
    max-width: 150px;
  }
}
</style>
