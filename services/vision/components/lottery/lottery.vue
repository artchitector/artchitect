<i18n>
{
  "en": {
    "lottery": "Lottery",
    "start_time": "Start time:",
    "collection_period": "Collection period:",
    "winners": "Winners",
    "of_total": "of total",
    "waiting": "waiting",
    "finished": "finished",
    "running": "running"
  },
  "ru": {
    "lottery": "Лотерея",
    "start_time": "Время запуска",
    "collection_period": "Период сбора",
    "winners": "Отобранные победители",
    "of_total": "из общего числа",
    "waiting": "ожидание",
    "finished": "завершена",
    "running": "идёт"
  }
}
</i18n>
<template>
  <div class="box">
    <h3 class="is-size-4">{{$t('lottery')}} "{{ lottery.Name }}" - {{ $t(lottery.State) }}</h3>
    <p>{{$t('start_time')}}: <i>{{ formatDate(lottery.StartTime) }}</i></p>
    <p>{{$t('collection_period')}}: <i>{{ formatDate(lottery.CollectPeriodStart) }} - {{
        formatDate(lottery.CollectPeriodEnd)
      }}</i></p>

    <div v-if="showWinners" class="has-text-centered box has-background-link-light">
      {{$t('winners')}} ({{ lottery.Winners.length }} {{$t('of_total')}} {{ lottery.TotalWinners }})<br/>
      <template v-for="winnerID in winners">
      <img v-if="winnerID === 0" class="mini-preview ml-1 mr-1" src="/in-progress-lottery.jpg"/>
      <a v-else :href="localePath(`/art/${winnerID}`)" @click.prevent="select(winnerID)" class="winner-link">
        <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${winnerID}`"/>
      </a>
      </template>
    </div>
    <viewer ref="viewer"/>
  </div>
</template>
<script>
import moment from "moment";

export default {
  props: ['lottery'],
  computed: {
    showWinners() {
      return this.lottery.State === "running" || this.lottery.State === "finished";
    },
    winners() {
      if (!this.lottery.TotalWinners || !this.lottery.Winners) {
        return [];
      }
      const winners = [];
      for (let i = 0; i < this.lottery.TotalWinners; i++) {
        if (i >= this.lottery.Winners.length) {
          winners.push(0);
        } else {
          winners.push(this.lottery.Winners[i]);
        }
      }
      return winners;
    }
  },
  methods: {
    formatDate(date) {
      // TODO need make global date helper and use it everywhere
      return moment(date).format("YYYY MMM Do HH:mm:ss")
    },
    select(id) {
      this.$refs.viewer.show(this.lottery.Winners, id);
    }
  }
}
</script>
<style>
img.mini-preview {
  max-height: 100px;
}

img.micro-preview {
  max-height: 60px;
}
.winner-link {
  display: inline-block;
}
</style>
