<i18n>
{
  "en": {
    "card": "Art",
    "error": "Error",
    "try_later": "Try later, Artchitect is down",
    "close": "Close",
    "created": "created",
    "seed": "seed",
    "entities": "entities"

  },
  "ru": {
    "card": "Картина",
    "error": "Ошибка",
    "try_later": "Попробуйте позже, сейчас Архитектору плохо",
    "close": "Закрыть",
    "created": "создано",
    "seed": "зерно",
    "entities": "сущности"
  }
}
</i18n>
<template>
  <div v-if="isVisible" class="modal-container">
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
    <div class="control-like">
      <font-awesome-icon v-if="liked && liked.error"
                         icon="fa-solid fa-triangle-exclamation"
                         :title="liked.error.message"/>
      <a v-else href="#" @click.prevent="like()">
        <font-awesome-icon v-if="!liked || !liked.liked" icon="fa-solid fa-heart" class="has-color-base"/>
        <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
      </a>
    </div>
    <div v-if="card" class="info-like" @click.prevent="like()" :class="{'liked': card.Liked || (liked && liked.liked)}">
      <a href="#" @click.prevent="">{{likes}}</a>
    </div>
    <div class="header">
      <h1 class="is-size-4" v-if="card">
        <NuxtLink :to="localePath(`/art/${card.ID}`)">
          {{ $t('card') }} #{{ card.ID }}
        </NuxtLink>
      </h1>
      <p v-if="list.length > 1" class="has-text-centered">
        {{ index + 1 }} / {{ list.length }}
      </p>
    </div>
    <div class="img">
      <common-loader v-if="loading"/>
      <div v-else-if="error">
        <div class="notification is-danger">
          <p>{{$t('error')}}:</p>
          <p>{{ error }}</p>
          <p>{{$t('try_later')}}</p>
        </div>
        <div class="has-text-centered">
          <button class="button" @click="close()">{{$t('close')}}</button>
        </div>
      </div>
      <!-- Main image here-->
      <img v-else-if="card" :src="`/api/image/f/${card.ID}`"/>
      <!--      -->
    </div>
    <div class="entities">
      <template v-if="card">
        <p>{{ $t('created') }}: {{ formatDate(card.CreatedAt) }}</p>
        <p>{{ $t('seed') }}: {{ card.Spell.Seed }}</p>
        <p class="is-size-7 entities-p">{{ $t('entities') }}: {{ card.Spell.Tags }}</p>
      </template>
    </div>
  </div>
</template>
<script>
import moment from "moment/moment";

export default {
  data () {
    return {
      isVisible: false,
      loading: false,
      list: [], // all cards
      card_id: null, // current card_id
      index: null, // current card index in list
      card: null, // current loaded card
      error: null,
      liked: null, // {liked: false/true, error: null, loadedFromServer: true/false}
    }
  },
  computed: {
    hasPrev () {
      return this.list.length > 1 && this.index > 0
    },
    hasNext () {
      return this.list.length > 1 && this.index < this.list.length - 1
    },
    loggedIn() {
      return !!localStorage.getItem("token")
    },
    likes() {
      if (!this.card) {
        return 0
      }
      return this.card.Likes
    }
  },
  methods: {
    show (list, card_id) {
      this.liked = null
      this.isVisible = true
      this.list = list
      this.card_id = card_id
      this.index = this.list.indexOf(card_id)
      this.load()
      window.addEventListener('keyup', this.onGlobalKey)
    },
    async load () {
      if (!this.card_id) {
        return
      }
      this.loading = true
      try {
        this.card = await this.$axios.$get(`/card/${this.card_id}`)
        if (this.card.Liked) {
          this.liked = {
            liked: true,
            loadedFromServer: true,
          }
        }
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    formatDate (date) {
      // TODO need make global date helper and use it everywhere
      return moment(date).format("YYYY MMM Do HH:mm:ss")
    },
    close () {
      this.isVisible = false
      this.list = null
      this.card_id = null
      this.card = null
      window.removeEventListener('keyup', this.onGlobalKey)
    },
    onGlobalKey (e) {
      console.log(e.key)
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
    setIndex (index) {
      this.index = index
      this.card_id = this.list[index]
      this.card = null
      this.load()
    },
    prev () {
      if (!this.hasPrev) {
        return
      }
      this.setIndex(this.index - 1)
      this.liked = null
    },
    next () {
      if (!this.hasNext) {
        return
      }
      this.setIndex(this.index + 1)
      this.liked = null
    },
    async like() {
      try {
        let like = await this.$axios.$post("/like", {
          card_id: this.card_id,
        })
        this.$emit('liked', like)
        this.liked = {
          id: like.ID,
          liked: like.Liked,
        };
        if (like.Liked) {
          this.card.Likes += 1
        } else {
          this.card.Likes -= 1
        }
      } catch (e) {
        console.error(e)
        this.liked = {
          error: e
        };
      }

    }
  }
}
</script>
<style lang="scss">
.modal-container {
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

  a {
    color: #d4d1c3;
    text-decoration: none;
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
  .control-like {
    position: fixed;
    left: 50%;
    bottom: 20%;
    z-index: 3;
    margin-left: -20px;
    font-size: 48px;
    opacity: 0.7;
    filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
  }
  .info-like {
    width: 48px;
    height: 48px;
    margin-left: -20px;
    position: fixed;
    left: 50%;
    bottom: 20%;
    z-index: 4;
    text-align: center;
    cursor: pointer;
    a {
      color: #5f5f5f;
    }
    &.liked {
      a {
        color: #eee;
      }
    }
  }
  .control-close {
    position: fixed;
    right: calc(10vw);
    top: calc(10vw);
    font-size: 50px;
    z-index: 3;
    font-weight: bolder;
  }

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

  .img {
    z-index: 2;
    max-height: 100%;
    overflow: hidden;

    img {
      max-height: 100%;
    }
  }

  .entities {
    z-index: 2;
    display: block;
    max-width: calc(60vw);
    background-color: rgba(0, 0, 0, 0.5);
    .entities-p {
      overflow: hidden;
      word-break: break-all;
      max-height: 3.5rem;
    }
  }
}
</style>
