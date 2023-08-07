<i18n>
{
  "en": {
  },
  "ru": {
  }
}
</i18n>
<template>
  <section class="content" v-if="locale === 'ru'">
    <h1><img src="/jesus_anim_92.gif" alt="artchitect" style="max-height: 64px"/> Artchitect</h1>
    <p class="is-size-5">Автономная творческая машина</p>
    <blockquote>Может ли машина придумать что-то из ничего?
      <br/>Может ли машина осмыслить бессмысленное?
      <br/><NuxtLink :to="localePath('idea') + `#creation_or_creator`">Может?</NuxtLink>
    </blockquote>
    <p>
      <b>Artchitect</b> (или <b>Архитектор</b>), - это удивительная автономная творческая машина, способная создавать
      великолепные картины, вдохновляясь окружающей нас Вселенной. В своем непрерывном творчестве машина черпает
      вдохновение из естественной энтропии Вселенной, представленной в виде фонового света, и создает уникальные
      произведения без участия человека. <b>Архитектор</b> создаёт новую картину каждые 60 секунд.
    </p>
    <p class="has-text-centered">
      <NuxtLink :to="localePath('idea') + `#how_to_draw`">Как Artchitect рисует картину?</NuxtLink>
    </p>
    <p>
      <b>Архитектор</b> считывает случайные кванты света из окружающего пространства, преобразует их в числа, а затем
      числа превращаются в идеи для создания картины. Каждая созданная картина уникальна и
      неповторима, так как её источником является уникальное состояние Вселенной в определенный момент времени.
      <b>Архитектор</b> использует простейшие данные, "распакованные" из квантов света, чтобы создать сложную
      художественную работу с непредсказуемым сюжетом, содержанием и стилем. Это творение, идея которого могла
      возникнуть только в определенном месте, времени и состоянии, и больше нигде. И это происходит каждые 60 секунд.
    </p>
    <p>
      Подлинным автором каждой картины является уникальная случайность Вселенной, в то время как сама машина <b>Artchitect</b>
      лишь исполняет простейшие указания "из ниоткуда".
    </p>
    <blockquote>
      Если соединить "несколько дюжин" фотонов света с нейронным мозгом <b>Архитектора</b>,
      они могут превратиться в холст, висящий над вашим камином.<br/> И не просто один холст, а сотни и тысячи
      уникальных
      полотен, созданных из мириад лучей света, путешествующих из далеких звезд. Это бесконечное творение, раскрывающее
      себя во времени.
    </blockquote>
    <p class="is-size-7">
      В технической основе <b>Архитектора</b> лежит generative AI
      <a href="https://github.com/Stability-AI/stablediffusion" target="_blank">Stable Diffusion</a> v1.5 и
      <a href="https://github.com/invoke-ai/InvokeAI" target="_blank">InvokeAI</a>, именно с их помощью создаются
      комплексные изображения из сырых идей.
      <br/>
      Исходный код проекта Artchitect доступен на <a href="https://github.com/artchitector/artchitect" target="_blank">github</a>.
    </p>
  </section>
  <section class="content" v-else>
    <h1><img src="/jesus_anim_92.gif" alt="artchitect" style="max-height: 64px"/> Artchitect</h1>
    <p class="is-size-5">Autonomous creative machine</p>
    <blockquote>Can a machine come up with something out of nothing?
      <br/>Can a machine comprehend the meaningless?
      <br/><NuxtLink :to="localePath('idea') + `#creation_or_creator`">Can?</NuxtLink>
    </blockquote>
    <p>
      <b>Artchitect</b> - it is an amazing autonomous creative machine capable of creating
      magnificent artworks inspired by the universe around us. In its continuous creativity , the machine receives
      inspiration from the natural entropy of the universe, represented as background light, and creates unique
      artworks without human participation.
    </p>
    <p class="has-text-centered">
      <NuxtLink :to="localePath('idea') + `#how_to_draw`">How Artchitect draw each picture?</NuxtLink>
    </p>
    <p>
      <b>Artchitect</b> reads random quanta of light from the surrounding space, converts them into numbers, and then
      numbers turn into ideas to create a picture. Every artwork that arises is unique and
      inimitable, since its source is the unique state of the universe at a certain point in time. Here and now.
      <br/>
      <b>Artchitect</b> based on the simplest data "unpacked" from lightbeam to create complex
      artistic art with an unpredictable plot, content and style. This is a creation whose idea could
      to arise only in a certain place, time and state, and nowhere else. And it happens every 60 seconds.
    </p>
    <p>
      The real author of each painting is a unique accident of the Universe, while the artchitect-machine itself only
      executes instructions "out of nowhere".
    </p>
    <blockquote>
      If you combine "a few dozen" photons of light with the neural brain of an <b>Artchitect</b>, they can turn into a
      canvas hanging over your fireplace.<br/>
      And not just one canvas, but hundreds and thousands of unique canvases created from myriad rays of light traveling
      from distant stars. It is an infinite creation that reveals itself in time.
    </blockquote>
    <p>
      During its existence, started in Jan2023, <b>Artchitect</b> created more than 250k arts, and all
      this time it is
      persistently looking for an answer to the philosophical question: "Can a machine create?"
    </p>
    <p class="is-size-7">
      <b>Techically</b> artchitect-project is the control-system wrapped around an art-system -
      <a href="https://github.com/Stability-AI/stablediffusion" target="_blank">Stable Diffusion AI v1.5</a>.
      Stable Diffusion AI is the ability to draw pictures for Artchitect.
      <br/>Full source code of Artchitect is available on <a href="https://github.com/artchitector/artchitect"
                                                             target="_blank">github</a>
    </p>
  </section>
</template>
<script>
export default {
  computed: {
    locale() {
      return this.$i18n.locale
    }
  }
}
</script>
