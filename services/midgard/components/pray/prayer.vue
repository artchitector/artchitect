<i18n>
{
  "en": {
    "title": "Request",
    "ucan": "Here you can receive personal painting from Artchitect. It will randomly create one next artwork definitely for you.",
    "submit": "Get personal card"

  },
  "ru": {
    "title": "Запрос",
    "ucan": "Здесь вы можете осуществить запрос к Архитектору и получить ответ в виде персональной картины",
    "submit": "Получить карточку"
  }
}
</i18n>
<template>
  <div class="content">
    <h1 class="is-size-4">{{$t('title')}}</h1>
    <p>
      {{$t('ucan')}}
    </p>
    <p>
      {{$t('try')}}
    </p>
    <hr/>
    <div class="notification">
      <span v-html="$t('important1')"></span>
      <NuxtLink to="/">artchitect.space</NuxtLink>.
      <span v-html="$t('important2')"></span>
    </div>
    <p class="has-text-centered">
      <button class="button is-primary" :disabled="loading" @click.prevent="submit()">{{$t('submit')}}
      </button>
    </p>
    <p v-if="error" class="notification is-danger">{{ error }}</p>
    <p v-if="loading" class="has-text-centered">
      <loader/>
    </p>
  </div>
</template>

<script>
export default {
  name: "prayer",
  data() {
    const randomPassword = (Math.random() + 1).toString(36).substring(7);
    return {
      password: randomPassword,
      loading: false,
      error: null,
    }
  },
  methods: {
    async submit() {
      try {
        this.loading = true
        localStorage.setItem("last_pray_password", this.password);  // only emitter of pray can see answer (temporary password needed)
        const result = await this.$axios.$post("/pray", {
          password: this.password
        })
        let url = this.localePath(`/prayer/${result}`)
        await this.$router.push(url)
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>

</style>
