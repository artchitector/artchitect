<i18n>
{
  "en": {
    "already_created": "already created",
    "currently_creating": "currently creating",
    "seed_number": "seed",
    "creating": "creating"
  },
  "ru": {
    "already_created": "создано",
    "currently_creating": "творческий процесс",
    "seed_number": "зерно",
    "creating": "создаётся"
  }
}
</i18n>

<template>
  <div class="progress-view">
    <div v-if="message.CardID">
      {{$t('already_created')}} #{{message.CardID}}
    </div>
    <div v-else class="heart-heading">
      <entropy v-if="entropy" :entropy="entropy"/>
      <h1 class="is-size-5 has-text-success has-text-centered mb-2">{{$t('currently_creating')}}</h1>
      <div>
        <div class="tags mb-3">
          <span class="tag is-primary">{{$t('seed')}}={{message.Seed}}</span>
          <span class="tag" v-for="tag in message.Tags">{{ tag }}</span>
        </div>
      </div>
      <div class="is-size-7 has-text-centered">
        {{$t('creating')}}
        <span v-if="progress">({{message.CurrentCardPaintTime}}/{{message.LastCardPaintTime}})</span>
      </div>
      <progress class="progress is-primary" :value="progress" max="100">-</progress>
    </div>
  </div>
</template>

<script>
import Entropy from "@/components/entropy/entropy.vue";
export default {
  name: "progress-view",
  components: {Entropy},
  props: ["message", "entropy"],
  computed: {
    progress() {
      if (!this.message) {
        return 0;
      }
      if (!this.message.LastCardPaintTime || !this.message.CurrentCardPaintTime) {
        return 0;
      }
      const progress = this.message.CurrentCardPaintTime / this.message.LastCardPaintTime;
      return Math.floor(progress * 100);
    },
  }
}
</script>

<style lang="scss" scoped>
  .progress-view {
    max-width: 800px;
    min-width: 370px;
    .tags .tag {
      font-size: 9px;
      letter-spacing: 0px;
    }
  }
</style>
