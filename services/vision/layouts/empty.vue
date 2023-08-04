<i18n>
{
  "en": {
  },
  "ru": {
  }
}
</i18n>
<template>
  <div>
    <div><template v-html="yandexMetrica"></template></div>
    <Nuxt/>
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
//body {
//  background-color: #2d2d2d;
//}
</style>
