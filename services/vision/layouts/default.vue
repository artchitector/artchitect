<i18n>
{
  "en": {
    "main": "main",
    "lottery": "lottery",
    "selection": "selection",
    "liked": "liked",
    "unity": "unity",
    "idea": "idea",
    "entropy": "entropy"
  },
  "ru": {
    "main": "главная",
    "lottery": "лотерея",
    "selection": "отборное",
    "liked": "вам нравится",
    "unity": "единства",
    "idea": "идея",
    "entropy": "энтропия"
  }
}
</i18n>
<template>
  <div>
    <div><template v-html="yandexMetrica"></template></div>
    <div class="has-text-centered is-header pt-1 pb-2">
      <div>
        <NuxtLink :to="localePath('index')">
          <h1 class="is-size-3 has-text-weight-bold">
            artchitect.space
          </h1>
        </NuxtLink>
      </div>
      <div>
        <NuxtLink :to="localePath('heart')">❤️</NuxtLink>
        <span> - </span>
        <NuxtLink :to="localePath('index')">{{$t('main')}}</NuxtLink>
        <span> - </span>
        <NuxtLink :to="localePath('/idea')">{{$t('idea')}}</NuxtLink>
        <span> - </span>
        <NuxtLink :to="localePath('/entropy')">{{$t('entropy')}}</NuxtLink>
        <span> - </span>
        <NuxtLink :to="localePath('/unity')">{{$t('unity')}}</NuxtLink>
        <span> - </span>
<!--        <NuxtLink :to="localePath('lottery')">{{$t('lottery')}}</NuxtLink>-->
<!--        <span> - </span>-->
<!--        <NuxtLink :to="localePath('selection')">{{$t('selection')}}</NuxtLink>-->
<!--        <span> - </span>-->
        <template v-if="loggedIn">
          <NuxtLink :to="localePath('/liked')">{{$t('liked')}}</NuxtLink>
          <span> - </span>
        </template>
        <NuxtLink :to="switchLocalePath('en')">en</NuxtLink>
        <NuxtLink :to="switchLocalePath('ru')">ru</NuxtLink>
      </div>
      <div class="security-block">
        <NuxtLink :to="localePath('login')">
          <figure v-if="avatar" class="image is-32x32">
            <img :src="avatar" class="is-rounded"/>
          </figure>
          <span v-else>🔑</span>
        </NuxtLink>
      </div>
    </div>

    <div class="container is-max-desktop">
      <Nuxt/>
    </div>

    <footer class="footer">
      <p class="has-text-centered">
        github: <a href="https://github.com/artchitector/artchitect">https://github.com/artchitector/artchitect</a><br/>
        twitter (en): <a href="https://twitter.com/artchitectspace">@artchitectspace</a><br/>
        telegram (ru): <a href="https://t.me/artchitect_space">t.me/artchitect_space</a><br/>
        contacts: <a href="mailto:i@nkuhta.ru">i@nkuhta.ru</a> or <a href="https://t.me/artchitector">@artchitector</a>
      </p>
    </footer>
  </div>
</template>
<script>
import loggedIn from '@/utils/auth'
export default {
  computed: {
    avatar() {
      return process.client ? localStorage.getItem("photo_url") : null;
    },
    loggedIn() {
      return loggedIn()
    },
    yandexMetrica() {
      return `
      <!-- Yandex.Metrika counter -->
    <script type="text/javascript" >
      (function(m,e,t,r,i,k,a){m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
        m[i].l=1*new Date();
        for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
        k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)})
      (window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");

      ym(92022377, "init", {
        clickmap:true,
        trackLinks:true,
        accurateTrackBounce:true
      });
    <\/script>
<noscript><div><img src="https://mc.yandex.ru/watch/92022377" style="position:absolute; left:-9999px;" alt="" /></div></noscript>
<!-- /Yandex.Metrika counter -->
      `
    }
  }
}
</script>
<style lang="scss">
.is-header {
  background-color: #2d2d2d;
  position: relative;

  h1 {
    color: #d4d1c3;
  }
}

footer.footer {
  background-color: #2d2d2d;
  color: hsl(0, 0%, 96%);
}

.security-block {
  position: absolute;
  top: 10px;
  right: 10px;

  img {
    box-shadow: 0px 0px 5px #d4d1c3;
  }
}
</style>
