<template>
  <unity-base :state="state" v-if="state"/>
  <loader v-else/>
</template>

<script>
import UnityBase from "@/components/heart/layout/unity.vue"
export default {
  name: "unity",
  components: {UnityBase},
  data() {
    return {
      state: null,
    }
  },
  methods: {
    onMessage(channelName, message) {
      this.state = message
    }
  }
}
</script>

<style scoped>

</style>
