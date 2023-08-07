<i18n>
{
  "en": {
    "version": "Version",
    "seed":"Seed",
    "tags_count": "Entities count",
    "tags":"Entities",
    "paint_progress": "Dream progress",
    "enjoy_progress":"Enjoy progress"
  },
  "ru": {
    "version": "Версия",
    "seed":"Зерно",
    "tags_count": "Кол-во сущностей",
    "tags":"Сущности",
    "paint_progress": "Прогресс сна",
    "enjoy_progress":"Прогресс наслаждения"
  }
}
</i18n>
<template>
  <div class="is-flex is-flex-direction-row">
    <viewer ref="viewer"/>
    <div class="image-container">
      <img v-if="!state.CardID" src="/in-progress.jpeg"/>
      <a v-else :href="`/art/${state.CardID}`" target="_blank" @click.prevent="viewer()">
        <img :src="`/api/image/s/${state.CardID}`"/>
      </a>
    </div>
    <div class="info-container">
      <p>{{$t('version')}}: {{ state.Version }}</p>
      <p>{{$t('seed')}}: {{ state.Seed }}</p>
      <p>
        {{$t('tags_count')}}:
        <span v-if="state.TagsCount === 0">-</span>
        <span v-else>{{ state.Tags ? state.Tags.length : 0 }}/{{ state.TagsCount }}</span>
      </p>
      <p class="is-size-7">
        {{$t('tags')}}:
        <span v-if="!state.Tags || !state.Tags.length">-</span>
        <span v-else>{{ state.Tags.join(', ') }}</span>
      </p>
      <p>
        {{$t('paint_progress')}}:
        <span v-if="state.CurrentCardPaintTime === null">
        -
      </span>
        <span v-else>
        {{ state.CurrentCardPaintTime }}/{{ state.LastCardPaintTime }}
      </span>
      </p>
      <p>
        {{$t('enjoy_progress')}}:
        <span v-if="state.CurrentEnjoyTime === null">
          -
          </span>
        <span v-else>
            {{ state.CurrentEnjoyTime }}/{{ state.EnjoyTime }}
          </span>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  name: "creation",
  props: ['state'],
  data() {
    return {}
  },
  methods: {
    viewer() {
      const ids = [this.state.CardID]
      this.$refs.viewer.show(ids, this.state.CardID);
    }
  }
}
</script>

<style lang="scss" scoped>
.heart div.image-container {
  min-width: 170px;
  width: 170px;
  padding-right: 10px;

  a {
    display: block;
  }
}

.heart hr.divider {
  margin: 0 0 0.5rem 0;
}

.heart div.info-container {

}
</style>
