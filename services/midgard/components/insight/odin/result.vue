<i18n>
{
  "en": {
    "created": "created art",
    "rest": "make rest"
  },
  "ru": {
    "created": "написана картина",
    "rest": "отдых"
  }
}
</i18n>
<template>
  <div class="result-wrapper">
    <common-entropy v-if="entropy" :entropy="entropy"/>
    <div class="enjoy-progress">
      <p class="is-size-5 has-text-success has-text-centered">
        {{ $t('created') }}
        <NuxtLink :to="localePath(`/art/${odin.artId}`)" class="has-text-info">
          #{{ odin.artId }}
        </NuxtLink>
      </p>
      <b-progress v-if="progress" :value="progress" type="is-info" show-value>
        {{$t('rest')}} {{odin.currentEnjoyTime}}/{{odin.expectedEnjoyTime}}
      </b-progress>
    </div>
    <div class="image-container">
      <NuxtLink :to="localePath(`/art/${odin.artId}`)" class="has-text-info">
        <img :src="`/api/image/${odin.artId}/f`"/>
      </NuxtLink>
    </div>
    <div class="random-four" v-if="giving && giving.given.length > 0">
      <div v-for="give in giving.given" class="random-four-item">
        <insight-odin-rrnd :art-id="give"/>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "insight-odin-result",
  props: ['odin', 'entropy', 'giving'],
  computed: {
    progress() {
      if (!this.odin.currentEnjoyTime && !this.odin.expectedEnjoyTime) {
        return 0
      }
      const progress = this.odin.currentEnjoyTime / this.odin.expectedEnjoyTime
      return Math.floor(progress * 100)
    }
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

