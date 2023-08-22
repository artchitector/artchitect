<i18n>
{
  "en": {
    "page_title": "Artchitect - Log in",
    "title": "Log in"
  },
  "ru": {
    "page_title": "Artchitect - Вход",
    "title": "Вход"
  }
}
</i18n>
<template>
  <section class="content">
    <h3 class="is-size-5">{{ $t('title') }}</h3>
    <div v-if="loading" class="has-text-centered">
      <common-loader/>
      <br/>
    </div>
    <template v-else-if="!isLocal && !isServer && !isLoggedIn">
      <script async src="https://telegram.org/js/telegram-widget.js?21"
              data-telegram-login="ArtchitectBot" data-size="large" data-auth-url="https://artchitect.space/api/login"
              data-request-access="write"></script>
    </template>
    <template v-else-if="isLocal && !isServer && !isLoggedIn">
      <button @click.prevent="fakeLogin">Fake login</button>
    </template>

    <div v-else-if="isLoggedIn">
      <figure class="image is-128x128">
        <img class="is-rounded" :src="photoUrl"/>
      </figure>
      <p>Вы вошли как @{{ username }}</p>
      <button class="button" @click="logout()">Выйти</button>
    </div>
    <div v-else>
      <div class="has-text-centered">
        <common-loader/>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: "login",
  head() {
    return {
      title: this.$t('title')
    }
  },
  computed: {
    isServer() {
      return !process.client
    },
    isLoggedIn() {
      return process.client ? !!localStorage.getItem("token") : false
    },
    username() {
      return process.client ? localStorage.getItem("username") : null
    },
    photoUrl() {
      return process.client ? localStorage.getItem("photo_url") : null
    }
  },
  data() {
    return {
      loading: false,
      isLocal: false,
    }
  },
  mounted() {
    if (process.env.IS_LOCAL === 'true') {
      this.isLocal = true
    }
    if (this.$route.query.token) {
      this.loading = true
      localStorage.setItem("token", this.$route.query.token)
      localStorage.setItem("username", this.$route.query.username)
      localStorage.setItem("photo_url", this.$route.query.photo_url)
      if (process.client) {
        window.location.href = this.localePath(`/login`)
      }
    }
  },
  methods: {
    logout() {
      localStorage.removeItem("token")
      localStorage.removeItem("username")
      localStorage.removeItem("photo_url")
      if (process.client) {
        window.location.href = this.localePath(`/login`)
      }
    },
    fakeLogin() {
      this.loading = true
      localStorage.setItem("token", "FAKE_LOCAL_TOKEN")
      localStorage.setItem("username", "fake_artchitect")
      localStorage.setItem("photo_url", "/icon196.png")
      if (process.client) {
        window.location.href = this.localePath(`/login`)
      }
    }
  }
}
</script>
