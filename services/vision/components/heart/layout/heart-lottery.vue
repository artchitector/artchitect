<template>
  <div>
    <p>
      <span v-if="state.Lottery.State === 'running'">Running</span>
      <span v-else>Finished</span>
      lottery <b>{{ state.Lottery.Name }}</b>.
      <span v-if="state.EnjoyTotalTime > 0">
        Enjoy after finish: {{state.EnjoyCurrentTime}}/{{state.EnjoyTotalTime}}
      </span>
    </p>
    <div class="box has-text-centered has-background-link-light">
      <p>Winners ({{ state.Lottery.Winners.length }} of total {{ state.Lottery.TotalWinners }})</p>
      <template v-for="winnerID in winners">
        <img v-if="winnerID === 0" class="mini-preview ml-1 mr-1" src="/in-progress-lottery.jpg"/>
        <a v-else :href="`/art/${winnerID}`" @click.prevent="select(winnerID)" class="winner-link">
          <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${winnerID}`"/>
        </a>
      </template>
    </div>
    <viewer ref="viewer"/>
  </div>
</template>

<script>
export default {
  name: "heart-lottery",
  props: ['state'],
  methods: {
    select(id) {
      this.$refs.viewer.show(this.state.Lottery.Winners, id);
    }
  },
  computed: {
    winners() {
      if (!this.state.Lottery.TotalWinners) {
        return [];
      }
      const winners = new Array(this.state.Lottery.TotalWinners).fill(0);
      for (let i = 0; i < this.state.Lottery.Winners.length; i++) {
        winners[i] = this.state.Lottery.Winners[i];
      }
      return winners;
    }
  }
}
</script>

<style scoped>

</style>
