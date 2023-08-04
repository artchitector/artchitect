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

    <div id="wrapper">
      <header ref="header" class="has-text-centered">
        <NuxtLink :to="localePath('index')">
          <img src="/jesus_anim_92.gif" alt="artchitect"/>
        </NuxtLink>
        <h1 class="is-size-3 has-text-weight-bold">
          <NuxtLink :to="localePath('index')">
            artchitect.space
          </NuxtLink>
        </h1>
        <div class="subtitle">autonomous creative machine</div>
      </header>
      <div id="main">
        <article :style="{'height': `${articleHeight}px`}">
          <Nuxt/>
        </article>
      </div>
    </div>


  </div>
</template>
<script>
import loggedIn from '@/utils/auth'
export default {
  data() {
    return {
      resizeTimeout: null,
      articleHeight: 0
    }
  },
  mounted() {
    window.addEventListener('resize', this.onResize)
    this.resizeTimeout = setTimeout(this.onResize, 200)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize)
  },
  methods: {
    onResize() {
      clearTimeout(this.resizeTimeout)
      this.resizeTimeout = setTimeout(() => {
        const headerHeight = this.$refs.header.clientHeight
        this.articleHeight = window.innerHeight - headerHeight
        console.log('article height', this.articleHeight)
      }, 50)
    }
  },
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
//html, body {margin: 0; height: 100%; overflow: hidden}

#wrapper {
  background-color: #2d2d2d;
  display: flex;
  min-height: 100vh;
  flex-direction: column;
  margin: 0;
  color: #d4d1c3;
}
#main {
  display: flex;
  flex: 1;
}
header {
  padding: 1em 1em 0.3em 1em;
  background-color: #d4d1c3;
  //background-image: url("/heart/back1.jpg");
  background-size:100% 100%;
  color: #2d2d2d;
  img {
    margin: 0;
    box-shadow: 0 0 10px #2d2d2d;
    max-width: 64px;
  }
  h1 {
    margin: 0;
    margin-top: -35px;
    padding: 0;
    text-shadow: 0 0 3px #d4d1c3;
    a {
      color: #2d2d2d;
      &:hover {
        color: #2d2d2d;
      }
    }
  }
  .subtitle {
    margin-top: -12px;
    font-size: 14px;
  }
}
#main > article {
  flex: 1;
  order: 1;
}
#main > nav,
#main > aside {
  flex: 0 0 20vw;
}
#main > nav {
  background: #D7E8D4;
  order: 3;
}
#main > aside {
  background: beige;
  order: 2;
}
</style>
