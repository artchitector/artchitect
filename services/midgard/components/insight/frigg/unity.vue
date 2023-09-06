<i18n>
{
  "en": {
    "rest": "finished. make rest",
    "unification": "Unification of",
    "making_collage": "Making collage",
    "leads_collecting": "Selecting leads",
    "applicants": "Applicants"
  },
  "ru": {
    "rest": "завершено. отдых",
    "unification": "Объединение единства",
    "making_collage": "Создаётся коллаж",
    "leads_collecting": "Отбор лидеров",
    "applicants": "Претендентов"

  }
}
</i18n>
<template>
  <div>
    <h3 class="is-size-6 has-text-centered has-text-primary mb-0">
       {{ $t('unification') }} {{ state.unity.mask }}
    </h3>

    <div v-if="state.children.length" class="tags mb-0 has-text-centered">
      <template v-for="unity in state.children">
        <span v-if="unity.state === 'empty'"
              class="tag is-light mr-2">{{ unity.mask }}</span>
        <span v-else-if="unity.state === 'unified' || unity.state === 'pre-unified'"
              class="tag is-success mr-2">{{ unity.mask }}</span>
        <span v-if="unity.state === 'reunification'"
              class="tag is-danger mr-2">{{ unity.mask }}</span>
      </template>
    </div>

    <div v-if="state.subprocess">
      <hr class="light-hr"/>
      <insight-frigg-unity :state="state.subprocess"/>
    </div>

    <div v-else-if="!state.collageFinished && state.totalLeads > 0" class="has-text-centered">
      <template v-if="state.collageStarted">
        <div class="is-size-7" v-if="state.collageStarted">
          {{ $t('making_collage') }}
        </div>
      </template>
      <div class="is-size-7" v-else>
        {{ $t('leads_collecting') }} ({{ state.leads.length }}/{{ state.totalLeads }}).
        {{ $t('applicants') }}: {{ state.totalApplicants }}
      </div>

      <div class="is-relative">
        <b-loading :is-full-page="false" v-model="state.collageStarted"/>
        <div class="leadsbox is-relative" v-for="line in leadsLines">
          <div class="leadsline" v-for="lead in line">
            <img :class="imgClass" v-if="lead" :src="`/api/image/artchitect-${lead}-xs`"/>
            <img :class="imgClass" v-else :src="`/images/black/black-xs.jpg`"/>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="state.collageFinished" class="has-text-centered">
      <b-progress v-if="this.state.expectedEnjoyTime" :value="progress" type="is-info" show-value class="mb-1">
        {{$t('rest')}} {{state.currentEnjoyTime}}/{{state.expectedEnjoyTime}}
      </b-progress>
      <img :src="`/api/uimage/${state.unity.mask}/${state.unity.version}/m`"/>
    </div>
  </div>
</template>

<script>
export default {
  name: "insight-frigg-unity",
  props: ["state"],
  computed: {
    leadsLines() {
      let lineWidth = 0
      switch (this.state.unity.rank) {
        case 100:
          lineWidth = 4
          break
        case 1000:
          lineWidth = 5
          break
        case 10000:
          lineWidth = 6
          break
        case 100000:
          lineWidth = 7
          break
      }

      const lines = []
      for (let i = 0; i < this.state.leads.length; i += lineWidth) {
        let line = this.state.leads.slice(i, i + lineWidth)
        for (let j = line.length; j < lineWidth; j++) {
          line.push(null) // пустые картинки чёрным цветом отобразятся
        }
        lines.push(line)
      }
      return lines
    },
    imgClass() {
      switch (this.state.unity.rank) {
        case 100:
          return "lead-16"
        case 1000:
          return "lead-25"
        case 10000:
          return "lead-36"
        case 100000:
          return "lead-49"
      }
    },
    progress() {
      if (!this.state.expectedEnjoyTime) {
        return 0
      }
      const progress = this.state.currentEnjoyTime / this.state.expectedEnjoyTime
      return Math.floor(progress * 100)
    }
  }
}
</script>

<style scoped lang="scss">
.light-hr {
  margin: 0;
  margin-top: 2px;
  height: 1px;
  border: 0;
  border-top: 1px solid #5f5f5f;
}

.tags {
  justify-content: center;
}

.tag {
  padding: 0 3px;
  margin-bottom: 2px;
  font-size: 11px;
}

.leadsbox {
  display: flex;
  justify-content: center;
  gap: 7px;

  .leadsline {
    img.lead-16 {
      width: 96px;
    }

    img.lead-25 {
      width: 64px;
    }

    img.lead-36 {
      width: 48px;
    }

    img.lead-49 {
      width: 32px;
    }
  }
}
</style>

