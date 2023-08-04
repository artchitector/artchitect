<i18n>
{
  "en": {
    "title": "Artchitect - entropy"
  },
  "ru": {
    "title": "Artchitect - энтропия"
  }
}
</i18n>

<template>
  <section>
    <h1 class="is-size-3" v-if="locale === 'ru'">Датчик энтропии</h1>
    <h1 class="is-size-3" v-else>Entropy sensor</h1>
    <hr/>
    <p v-if="locale === 'ru'">
      Для рисования картин Архитектору приходится придумывать "что рисовать?" - придумывать набор ключевых слов и
      уникальный seed-номер. Для возможности выбора (принятия решений)
      Архитектор подключён к датчику энтропии (собранного из веб-камеры), который позволяет из светового шума получить
      float64-число. <b>float64-число</b> - фундамент каждого выбора в работе Artchitect, а в день их происходит десятки
      тысяч.
    </p>
    <p v-else>
      To draw arts Artchitect need know "what to draw?", it need to invent new keywords and a unique
      seed number. To be able to choose (make decisions), Artchitect is connected to an entropy sensor (made
      of a webcam), which allows Artchitect to receive a float64 number from light noise. <b>float64-number</b>>
      is the foundation of
      every choice in the work of Artchitect, and there are tens of thousands decisions per day.
    </p>
    <hr/>
    <p v-if="locale === 'ru'">
      <b>Шаг 1. Исходный кадр (Source)</b>
      <br/>
      Энтропия (в форме светового шума) получается из разницы двух соседних кадров. Текущий кадр виден ниже, его
      содержание не имеет принципиального
      значения. Красным отмечена область, с которой будет сниматься фоновый шум.
    </p>
    <p v-else>
      <b>Step 1. Original frame (Source)</b>
      <br/>
      Entropy (in the form of light noise) is obtained from the difference of two adjacent frames. The current frame is
      visible below, its content does not matter. The area from which background noise will be removed is marked in red.
    </p>
    <div class="has-text-centered">
      <img v-if="images.source !== null" :src="`data:image/jpeg;base64, ${images.source}`"
           alt="loading source stream"/>
    </div>
    <p v-if="locale === 'ru'">
      <b>Шаг 2. Шум (Noise)</b>
      <br/>
      Для экстракции светового шума используется вычитание кадров (из цветов текущего кадра вычитаются цвета
      предыдущего кадра). Разница между кадрами и есть световой шум, который дополнительно усиливается. В визуальном
      представлении шум выглядит следующим образом:
    </p>
    <p v-else>
      <b>Step 2. Noise</b>
      <br/>
      To extract light noise, frame subtraction is used (the colors of the current frame are subtracted from the colors
      of the previous frame). The difference between the frames is the light noise, which is further amplified. In the
      visual representation , the noise looks like this:
    </p>
    <div class="has-text-centered">
      <img v-if="images.noise !== null" :src="`data:image/jpeg;base64, ${images.noise}`"
           style="height: 256px; width: 256px;"
           alt="loading noise stream"/>
    </div>
    <p v-if="locale === 'ru'">
      <b>Шаг 3. Энтропия - сжатие и нормализация шума (Entropy)</b>
      <br/>
      Шум сжимается до области 8x8 пикселей и нормализуется (распределяется по шкале от 0 до
      255). Так и выглядит <b>базовая энтропия</b>, как её видит Архитектор.
    </p>
    <p v-else>
      <b>Step 3. Entropy - noise compression and normalization</b>
      <br/>
      The noise is compressed to an area of 8x8 pixels and normalized (distributed on a scale from 0 to 255). This is
      what the <b>basic entropy</b> looks like, as the Artchitect sees it.
    </p>
    <div class="has-text-centered">
      <img v-if="images.entropy !== null" :src="`data:image/png;base64, ${images.entropy}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading entropy stream"/>
      <div>{{ entropy.float }}</div>
    </div>
    <p v-if="locale === 'ru'">
      <b>Шаг 4. Инвертированная энтропия (Choice)</b>
      <br/>
      Чтобы Архитектор рисовал разнообразные картины, он при каждом выборе (ключевого слова или seed-номера) должен
      обеспечивать максимальное разнообразие (Архитектор должен выбирать очень разные слова из словаря каждый раз,
      нечасто повторяясь).
      <br/>
      Чтобы
      обеспечить такую селективность, энтропия инвертируется бинарно (каждый байт в рисунке энтропии зеркально
      отражается). Инвертированная энтропия становится очень случайной и позволяет равномерно справедливо выбирать слова
      из словаря.
    </p>
    <p v-else>
      <b>Step 4. Inverted entropy (Choice)</b>
      <br/>
      In order for Artchitect to draw a variety of pictures, it must ensure maximum diversity with each choice
      (keyword or seed number). Artchitect must choose very different words from the dictionary each time, rarely
      repeating itself.
      <br/>
      To ensure such selectivity, entropy is inverted binary (each byte in the entropy pattern is mirrored). The
      inverted entropy becomes very random and allows for evenly fair selection of words from the dictionary.
    </p>
    <div class="has-text-centered">
      <img v-if="images.choice !== null" :src="`data:image/png;base64, ${images.choice}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading choice stream"/>
      <div>{{ choice.float }}</div>
    </div>
    <p v-if="locale === 'ru'">
      <b>Шаг 5. Превращение изображения в число</b>
      <br/>
      Каждый пиксель изображения или включен, или выключен (значение его цвета ближе к красному или к чёрному). Каждый
      пиксель - это один бит из 64-битного целого числа (uint64). Картинка после бинарного преобразования становится
      uint64 числом на шкале от 0 до 18446744073709551615. Положение сгенерированного числа на этой шкале и есть точка
      выбора.
      Далее шкала превращается в float64-число, представляющее шкалу выбора - это дробное число от 0.0 до 1.0, где 0 -
      самый первый элемент в списке, 1.0 - последний элемент, 0.5 - примерно середина списка.
    </p>
    <p v-else>
      <b>Step 5. Image to number conversion</b>
      <br/>
      Each pixel of 8x8-image is either on or off (its color value is closer to red or black). Each pixel is one bit of
      a 64-bit integer (uint64). The image after binary conversion becomes a uint64 number on a scale from 0 to
      18446744073709551615. The position of generated number on the scale is the point of choice. Next uint64 number
      become a float64-a number representing the selection scale from 0.0 to 1.0, where 0 is the
      very
      first element in the list, 1.0 is the last element, 0.5 is about the middle of the list.
    </p>
    <div class="columns">
      <div class="column has-text-centered">
        <img v-if="images.entropy !== null" :src="`data:image/png;base64, ${images.entropy}`"
             style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading entropy stream"/>
        <div class="is-size-7" v-html="entropy.bytes ? entropy.bytes.match(/.{1,8}/g).join('<br/>') : '-'"></div>
        <div>uint64: {{ entropy.int }}</div>
        <div>float64: <b>{{ entropy.float }}</b></div>
      </div>
      <div class="column has-text-centered">
        <img v-if="images.choice !== null" :src="`data:image/png;base64, ${images.choice}`"
             style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading choice stream"/>
        <div class="is-size-7" v-html="choice.bytes ? choice.bytes.match(/.{1,8}/g).join('<br/>') : '-'"></div>
        <div>uint64: {{ choice.int }}</div>
        <div>float64: <b>{{ choice.float }}</b></div>
      </div>
    </div>

  </section>
</template>

<script>

import WsConnection from "~/utils/ws_connection";

export default {
  name: "entropy",
  data() {
    return {
      player: null,
      logPrefix: "🎆",
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,

      images: {
        source: null,
        noise: null,
        entropy: null,
        choice: null
      },
      entropy: {
        bytes: null,
        int: null,
        float: null,
      },
      choice: {
        bytes: null,
        int: null,
        float: null,
      }
    }
  },
  head() {
    return {
      title: this.$t('title')
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, this.logPrefix, ['entropy'], 100)
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
      this.onMessage(channel, message)
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`${this.logPrefix}: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {attempt, maxAttempts}
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
      console.log(`${this.logPrefix}: connection established`)
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  },
  methods: {
    onMessage(chan, msg) {
      if (msg.Entropy) {
        this.entropy.bytes = msg.Entropy.Binary
        this.entropy.float = msg.Entropy.Float64
        this.entropy.int = msg.Entropy.Uint64
      } else {
        this.entropy.bytes = null
        this.entropy.float = null
        this.entropy.int = null
      }

      if (msg.Choice) {
        this.choice.bytes = msg.Choice.Binary
        this.choice.float = msg.Choice.Float64
        this.choice.int = msg.Choice.Uint64
      } else {
        this.choice.bytes = 0
        this.choice.float = 0
        this.choice.int = 0
      }

      for (const imageName in msg.ImagesEncoded) {
        this.images[imageName] = msg.ImagesEncoded[imageName]
      }
    }
  }
}
</script>


<style scoped lang="scss">

</style>
