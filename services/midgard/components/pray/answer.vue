<i18n>
{
  "en": {
    "loading": "Your answer is loading",
    "next_in_queue": "Your request is next in queue. Time to wait - less than 120 seconds.",
    "before": "Requests in queue before your:",
    "time_to_wait": "Time to wait",
    "seconds": "seconds",
    "concentrate": "Close your eyes and concentrate on your request!",
    "running": "Your answer is in work. Time to wait - less than 60 seconds",
    "answer": "Your answer is:"
  },
  "ru": {
    "loading": "Ваш ответ загружается",
    "next_in_queue": "Ваш ответ следующий в очереди. Времени на ожидание - менее 120 секунд.",
    "before": "Запросов в очереди перед вами:",
    "time_to_wait": "Времени ждать",
    "seconds": "секунд",
    "concentrate": "Закройте глаза и сконцентрируйтесь на вашем запросе!",
    "running": "Ваш ответ в работе. Времени ждать - менее 60 секунд",
    "answer": "Ваш ответ:"
  }
}
</i18n>
<template>
  <div class="has-text-centered">
    <template v-if="loading">
      <loader v-if="loading"/>
      <p></p>
      <template v-if="state">
        <template v-if="state.State === 'waiting'">
          <p v-if="state.Queue === 0">{{$t('next_in_queue')}}</p>
          <p v-else-if="state.Queue > 0">{{$t('before')}} {{ state.Queue }}. {{$t('time_to_wait')}} -
            {{ (state.Queue + 1) * 60 + 60 }} {{$t('seconds')}}.</p>
          <p><b>{{$t('concentrate')}}</b></p>
        </template>
        <template v-else-if="state.State === 'running'">
          <p>{{$t('running')}}</p>
        </template>
      </template>
    </template>
    <div v-else-if="error" class="notification is-danger">{{ error }}</div>
    <div v-else-if="state && state.State === 'answered' && state.Answer > 0">
      <p>{{$t('answer')}}</p>
      <NuxtLink :to="localePath(`/card/${state.Answer}`)">
        <img :src="`/api/image/f/${state.Answer}`"/>
      </NuxtLink>
    </div>
  </div>
</template>

<script>

export default {

  name: "answer",
  props: ['id'],
  data() {
    return {
      loading: true,
      error: null,
      state: null,
      interval: null,
    }
  },
  async mounted() {
    const id = parseInt(this.id)
    if (!id) {
      this.error = "ID must be type uint"
      return
    }
    const password = localStorage.getItem("last_pray_password")
    if (!password) {
      this.error = "Access to your request was lost. Please, make new request. Sorry :("
      return
    }
    this.startListening(id, password)
  },
  methods: {
    async startListening(id, password) {
      this.loading = true;
      this.interval = setInterval(async () => {
        try {
          this.state = await this.$axios.$post(`/pray/answer`, {
            id: id,
            password: password,
          })
          console.log(`🙏: current pray state ${id} ${this.state.State}`)
          if (this.state.State === "answered") {
            this.loading = false;
            clearInterval(this.interval)
          }
        } catch (e) {
          this.error = e.message;
          this.loading = false;
          clearInterval(this.interval)
        }
      }, 1000)
    }
  },
  beforeDestroy() {
    clearInterval(this.interval)
  }
}
</script>

<style scoped>

</style>
