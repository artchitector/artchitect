<i18n>
{
  "en": {
    "subtitle": "unification",
    "unity": "Unity",
    "state": "State",
    "collecting_children": "collect children unities",
    "unify_children": "unify children unities",
    "promote_leads": "promote lead arts",
    "prepare_thumb": "prepare collage",
    "finished": "unified",
    "children": "Children unities",
    "thumb": "Collage",
    "version": "Version"
  },
  "ru": {
    "subtitle": "объединение единств",
    "unity": "Единство",
    "state": "Статус",
    "collecting_children": "собираем дочерние единства",
    "unify_children": "объединение дочерних единств",
    "promote_leads": "выбор лидеров",
    "prepare_thumb": "подготовка коллажа",
    "finished": "едино",
    "children": "Дочерние единства",
    "thumb": "Коллаж",
    "version": "Версия"
  }
}
</i18n>
<template>
  <div>
    <h3 class="is-size-5 has-text-centered mb-2">{{$t('subtitle')}}</h3>
    <div v-for="unification in state.Unifications" class="box" :class="getBoxClass(unification)">
      <div>
        <b>{{ $t('unity') }}</b>: {{ unification.Unity.Mask }}, <b>{{ $t('state') }}</b>: {{ $t(unification.State) }}
        <template v-if="unification.CurrentProgress > 0">
          {{ unification.CurrentProgress }}/{{ unification.TotalProgress }}
        </template>
        <b>{{$t('version')}}: {{unification.Unity.Version}}</b>
      </div>
      <div v-if="unification.Version">
        {{ $t('version')}}:
      </div>
      <div v-if="unification.Children && unification.Children.length">
        {{ $t('children') }}:
        <template v-for="child in unification.Children">
          <span v-if="child.State === 'empty'" class="tag is-light mr-2">{{ child.Mask }}</span>
          <span v-else-if="child.State === 'unified'" class="tag is-success mr-2">{{ child.Mask }}</span>
          <span v-else-if="child.State === 'reunification'" class="tag is-danger mr-2">{{ child.Mask }}</span>
          <span v-else-if="child.State === 'skipped'" class="tag is-dark mr-2">{{ child.Mask }}</span>
        </template>
      </div>
      <div v-if="!!unification.Thumb" class="has-text-centered">
        {{ $t('thumb') }}:<br/>
        <img :src="`/api/image/unity/${unification.Unity.Mask}/${unification.Unity.Version}/m`"/>
      </div>
      <div v-else-if="unification.Leads && unification.Leads.length">
        <template v-for="lead in unification.Leads">
          <img class="lead-preview" :src="`/api/image/xs/${lead}`"/>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "unity",
  props: ['state'],
  methods: {
    getBoxClass(unification) {
      const cls = {}
      cls[`rank${unification.Rank}`] = true
      return cls
    }
  }
}
</script>

<style lang="scss" scoped>
img.lead-preview {
  max-width: 64px;
  margin: 3px;
}

.rank10000 {

}

.rank1000 {
  margin-left: 10px;
  margin-right: 10px;
}

.rank100 {
  margin-left: 20px;
  margin-right: 20px;
}
</style>
