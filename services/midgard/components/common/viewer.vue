<i18n>
{
  "en": {
    "art": "Art",
    "error": "error",
    "try_later": "try later...",
    "close": "close",
    "created": "created",
    "seed": "seed-number",
    "words": "words"
  },
  "ru": {
    "art": "Картина",
    "error": "ошибка",
    "try_later": "попробуйте позже...",
    "close": "закрыть",
    "created": "создано",
    "seed": "seed-номер",
    "words": "ключевые слова"
  }
}
</i18n>
<template>
  <div v-if="visible" class="viewer-container">
    <!-- Фон и элементы управления (предыдущий, следующий, закрыть) -->
    <div class="background" @click="close()"></div>
    <div class="control-close">
      <a href="#" @click.prevent="close()">x</a>
    </div>
    <div class="control-prev" v-if="hasPrev">
      <a href="#" @click.prevent="prev()"><</a>
    </div>
    <div class="control-next" v-if="hasNext">
      <a href="#" @click.prevent="next()">></a>
    </div>

    <!-- Заголовок окна -->
    <div class="header">
      <h1 class="is-size-4" v-if="current.art">
        <NuxtLink :to="localePath(`/art/${current.art.id}`)">
          {{ $t('art') }} #{{ current.art.id }}
        </NuxtLink>
      </h1>
      <p v-if="list.length > 1" class="has-text-centered">
        {{ current.index + 1 }} / {{ list.length }}
      </p>
    </div>

    <div class="image-container">
      <common-loader v-if="loading"/>
      <div v-else-if="error">
        <div class="notification is-danger">
          <p>{{ $t('error') }}</p>
          <p>{{ error }}</p>
          <p>{{ $t('try_later') }}</p>
        </div>
        <div class="has-text-centered">
          <button class="button" @click="close()">{{ $t('close') }}</button>
        </div>
      </div>
      <img v-else-if="current.art" :src="`/api/image/${current.art.id}/f`"/>
    </div>

    <div class="words">
      <template v-if="current.art">
        <p>{{ $t('created') }}: {{ formatDate(current.art.CreatedAt) }}</p>
        <p>{{ $t('seed') }}: {{ current.art.ideaSeed }}</p>
        <p class="is-size-7 words-p">{{ $t('words') }}: {{ current.art.ideaWords.join(', ') }}</p>
      </template>
    </div>
  </div>
</template>

<script>
import {format} from "@/utils/date";
export default {
  // Просмотрщик картин
  name: "common-viewer",
  data() {
    return {
      visible: false,
      loading: false,
      error: null,

      list: [], // arts
      current: { // текущий выбранный арт в viewer-е
        index: null,
        artID: null,
        art: null,
      }
    }
  },
  computed: {
    hasPrev() {
      return this.list.length > 1 && this.current.index > 0
    },
    hasNext() {
      return this.list.length > 1 && this.current.index < this.list.length - 1
    },
  },
  methods: {
    show(list, artID) {
      this.list = list // в виде массива ID - [1,2,3,4,5,6...]
      this.current.artID = artID
      this.current.index = this.list.indexOf(artID)
      this.visible = true
      this.load()
      window.addEventListener('keyup', this.onGlobalKey)
    },
    async load() {
      if (!this.current.artID) {
        return
      }
      this.loading = true
      try {
        this.current.art = await this.$axios.$get(`/art/${this.current.artID}/flat`)
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    formatDate(d) {
      return format(d)
    },
    setIndex (index) {
      this.current.index = index
      this.current.artID = this.list[index]
      this.current.art = null
      this.load()
    },
    next() {
      return this.hasNext && this.setIndex(this.current.index + 1)
    },
    prev() {
      return this.hasPrev && this.setIndex(this.current.index - 1)
    },
    close () {
      this.visible = false
      this.list = []
      this.current.artID = null
      this.current.art = null
      this.current.index = 0
      window.removeEventListener('keyup', this.onGlobalKey)
    },
    onGlobalKey (e) {
      if (e.key === 'Escape') {
        this.close()
      } else if (e.key === 'ArrowLeft') {
        this.prev()
      } else if (e.key === 'ArrowRight') {
        this.next()
      } else if (e.key === '+' || e.key === '=') {
        this.like()
      }
    },
  }
}
</script>

<style scoped lang="scss">
  .viewer-container {
    padding: 20px;
    position: fixed;
    z-index: 1;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #d4d1c3;
    gap: 10px;

    .background {
      bottom: 0;
      left: 0;
      position: absolute;
      right: 0;
      top: 0;
      background-color: rgba(0, 0, 0, 0.8);
    }
    .header {
      z-index: 2;

      h1 {
        background-color: rgba(0, 0, 0, 0.5);
      }
    }
    a {
      color: #d4d1c3;
      text-decoration: none;
    }

    .image-container {
      z-index: 2;
      max-height: 100%;
      overflow: hidden;

      img {
        max-height: 100%;
      }
    }

    .words {
      z-index: 2;
      display: block;
      max-width: calc(60vw);
      background-color: rgba(0, 0, 0, 0.5);
      .words-p {
        overflow: hidden;
        word-break: break-all;
        max-height: 3.5rem;
      }
    }

    .control-prev {
      position: fixed;
      left: calc(10vw);
      top: 50%;
      font-size: 50px;
      z-index: 3;
      font-weight: bolder;
    }
    .control-next {
      position: fixed;
      right: calc(10vw);;
      top: 50%;
      font-size: 50px;
      z-index: 3;
      font-weight: bolder;
    }
    .control-close {
      position: fixed;
      right: calc(10vw);
      top: calc(10vw);
      font-size: 50px;
      z-index: 3;
      font-weight: bolder;
    }
  }
</style>
