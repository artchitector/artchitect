<i18n>
{
  "en": {
    "title": "your liked collection",
    "not_authorized": "you're not authorized",
    "login": "login page"
  },
  "ru": {
    "title": "вам нравится",
    "not_authorized": "вы не авторизованы",
    "login": "страница входа"
  }
}
</i18n>
<template>
  <section>
    <h3 class="is-size-3 has-text-centered">{{ $t('title') }}</h3>
    <div v-if="!loggedIn" class="has-text-centered">
      {{$t('not_authorized')}}
      <br/>
      <NuxtLink :to="localePath(`/login`)">
        {{ $t('login') }}
      </NuxtLink>
    </div>
    <liked-list v-else/>
  </section>
</template>
<script>
import loggedIn from "~/utils/auth";

export default {
  head() {
    return {
      title: this.$t('title')
    }
  },
  computed: {
    loggedIn() {
      return loggedIn()
    }
  }
}
</script>
<style lang="scss" scoped>

</style>
