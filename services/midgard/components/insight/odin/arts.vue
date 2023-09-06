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
  <div class="lastart-container">
    <div class="left-column">
      <insight-odin-rnd v-if="giving.given.length > 0" :art-id="giving.given[0]" @show="$emit('show', $event)"/>
      <insight-odin-rnd v-if="giving.given.length > 1" :art-id="giving.given[1]" @show="$emit('show', $event)"/>
    </div>

    <div class="center-column">
      <div class="column-image" v-if="giving.lastArtID > 0">
        <div class="link-container">
          {{ $t('last_art') }}
          <a :href="localePath(`/art/${giving.lastArtID}`)" class="has-text-info" @click.prevent="$emit('show', giving.lastArtID)">
            #{{ giving.lastArtID }}
          </a>
          <div class="image-container">
            <a :href="localePath(`/art/${giving.lastArtID}`)" class="has-text-info" @click.prevent="$emit('show', giving.lastArtID)">
              <img :src="`/api/image/${giving.lastArtID}/m`" :alt="`art_${giving.lastArtID}`"/>
            </a>
          </div>
        </div>
      </div>
    </div>

    <div class="right-column">
      <insight-odin-rnd v-if="giving.given.length > 2" :art-id="giving.given[2]"  @show="$emit('show', $event)"/>
      <insight-odin-rnd v-if="giving.given.length > 3" :art-id="giving.given[3]"  @show="$emit('show', $event)"/>
    </div>
  </div>
</template>

<script>
export default {
  name: "insight-odin-arts",
  props: ['giving']
}
</script>

<style scoped lang="scss">
.lastart-container {
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
