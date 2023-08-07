<i18n>
{
  "en": {
    "last_art": "last art"
  },
  "ru": {
    "last_art": "последняя работа"
  }
}
</i18n>
<template>
  <div class="lastdream-container">

    <div class="left-column">
      <rnd :dream_id="four && four.length > 0 ? four[0] : 0"/>
      <rnd :dream_id="four && four.length > 1 ? four[1] : 0"/>
    </div>

    <div class="center-column">
      <div class="column-image">
        <div class="link-container">
          {{$t('last_art')}}
          <NuxtLink :to="localePath(`/art/${last}`)" class="has-text-info">#{{ last }}</NuxtLink>
        </div>
        <div class="image-container">
          <NuxtLink :to="localePath(`/art/${last}`)" class="has-text-info">
            <img v-if="last" :src="`/api/image/m/${last}`"/>
          </NuxtLink>
          <div class="control-like">
            <liker :dream_id="last"/>
          </div>
        </div>
      </div>

    </div>

    <div class="right-column">
      <rnd :dream_id="four && four.length > 2 ? four[2] : 0"/>
      <rnd :dream_id="four && four.length > 3 ? four[3] : 0"/>
    </div>

  </div>
</template>

<script>
import Rnd from "@/components/flexheart/creation/rnd.vue";
import Liker from "@/components/utils/liker.vue";

export default {
  name: "lastdream",
  components: {Liker, Rnd},
  props: ["last", "four"],
}
</script>

<style lang="scss" scoped>
.lastdream-container {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: stretch;
  background-color: rgba(0, 0, 0, 0.1);
  max-width: 800px;

  .left-column, .right-column {
    flex-grow: 1;
    justify-content: center;
  }

  .center-column {
    flex-grow: 3;
    margin: 0 5px;

    .column-image {
      text-align: center;
      margin: 5px 5px;
      position: relative;

      .image-container {
        position: relative;

        .control-like {
          position: absolute;
          left: 50%;
          bottom: 10%;
          z-index: 3;
          transform: translate(-50%, -10%);
        }
      }

      img {
        display: block;
        width: 100%;
      }

      .link-container {
        font-size: 12px;
        background-color: rgba(0, 0, 0, 0.5);
        padding: 0 3px;
      }
    }
  }

  .left-column, .right-column, .center-column {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: stretch;
  }
}
</style>
