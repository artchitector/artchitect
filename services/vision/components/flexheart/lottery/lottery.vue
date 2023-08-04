<template>
  <heart-lottery :state="state" v-if="state"/>
  <div class="has-text-centered pt-3" v-else>
    <loader/>
  </div>
</template>

<script>
import HeartLottery from "@/components/heart/layout/heart-lottery.vue";

export default {
  name: "lottery",
  components: {HeartLottery},
  data() {
    return {
      state: null,
    }
  },
  methods: {
    onMessage(channelName, state) {
      console.log(state)
      this.state = state
    }
  }
}
</script>

<style scoped>

</style>
