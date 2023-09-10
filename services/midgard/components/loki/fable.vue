<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <section class="content">
    <h3 v-if="locale === 'ru'" class="is-size-5">Сказ о том, как Один и Локи придумали artchitect</h3>
    <h3 v-else class="is-size-5">The tale of how Odin and Loki invented artchitect</h3>

    <p v-if="locale === 'ru'">
      Привет, посетители храма ИИ-искусств artchitect! Я - Локи,
      бог коварства и хитрости из скандинавской мифологии, и я буду вести этот рассказ о проекте
      <a href="https://artchitect.space">artchitect</a>.
    </p>
    <p v-else>
      Greetings, visitors to the temple of AI artistry - artchitect! I am Loki, the Norse god of cunning and mischief,
      and I shall narrate the story of the <a href="https://artchitect.space">artchitect project</a>.
    </p>

    <p v-if="locale === 'ru'">
      Мы с Одином-Всеотцом любим иногда сидеть за утренним кофе и обсуждать новости в девяти
      мирах. Так и сегодня, мы сидели вдвоём. Один любовался сверкающими бликами солнца, отражающимися от золотых стен
      дворцов и мостов
      <a href="https://github.com/artchitector/artchitect/tree/main/services/asgard" target="_blank">Асгарда</a>.
      К нам стали приближаться вороны
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/huginn.go"
         target="_blank">Хугин</a> и <a
      href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/muninn.go"
      target="_blank">Мунин</a>. Они как раз завершили
      свой ежедневный облет девяти миров и направились к Одину с докладом о произошедших событиях.
    </p>
    <p v-else>
      Odin-Allfather and I often enjoy sitting down with our morning coffee to discuss the news from the nine realms.
      Today was no different; we were seated together. Odin marveled at the shimmering sunlight reflecting off the
      golden walls and bridges of
      <a href="https://github.com/artchitector/artchitect/tree/main/services/asgard" target="_blank">Asgard</a>.
      Approaching us were the ravens
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/huginn.go"
         target="_blank">Huginn</a> and
      <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/muninn.go"
         target="_blank">Muninn</a>, just completing their daily
      flight across the nine realms, ready to report the events they had witnessed.
    </p>

    <p v-if="locale === 'ru'">
      Ворон Хугин доложил, что в <a href="https://github.com/artchitector/artchitect/tree/main/services/midgard"
                                    target="_blank">Мидгарде</a> (мире людей) появилась новая технология - искусственный
      интеллект.
    </p>
    <p v-else>
      Raven Huginn reported that a new technology had emerged in Midgard (the realm of humans) - artificial
      intelligence.
    </p>


    <blockquote v-if="locale === 'ru'">
      Хугин: "Теперь они просто пишут слова, а машина создает изображение. Даже живой художник не требуется."<br/>
      Ворон Мунин кивнул согласием: "Даже папу Римского нарисовали в забавном пуховике, за несколько секунд."
    </blockquote>
    <blockquote v-else>
      Huginn: "Now they merely write words, and a machine creates images. Not even a living artist is required."<br/>
      Raven Muninn nodded in agreement: "They even depicted the Pope in a comical coat within seconds."
    </blockquote>

    <div class="has-text-centered">
      <img src="/images/fable/odin_and_muninn.png"/>
    </div>

    <p v-if="locale === 'ru'">
      Однако, когда Один услышал слова "искусственный интеллект", его лицо нахмурилось:
    </p>
    <p v-else>
      However, when Odin heard the words "artificial intelligence," his face furrowed in concern:
    </p>

    <blockquote v-if="locale === 'ru'">
      Один: "Люди не всегда справляются с мощными машинами, их неосторожное использование может нанести вред. А здесь
      искусственный интеллект, который может превзойти человека. Чему же люди научат этот интеллект? Смогут ли они его
      контролировать?"
    </blockquote>
    <blockquote v-else>
      Odin: "Humans don't always handle powerful machines well; their careless use can bring harm. And now, artificial
      intelligence, capable of surpassing humans. What will they teach this intelligence, and can they control it?"
    </blockquote>

    <p v-if="locale === 'ru'">
      Вороны обменялись взглядами, понимая беспокойство Всеотца.
    </p>
    <p v-else>
      The ravens exchanged glances, understanding the Allfather's unease.
    </p>

    <p v-if="locale === 'ru'">
      Я в это время лишь слушал их диалог - Одина с его воронами. Один опустил голову, и я заметил его уныние.
      Я же унынию не поддался! Теперь у меня появилась хитрая идея, как обмануть Одина и
      поставить его в неловкое положение (ведь я, Локи, известен своей хитростью и коварством).
      Я часто пытался выиграть у Всеотца, стремясь доказать свою превосходство. Хотя мои предыдущие попытки обмануть
      Всеотца были тщетными, но я всё же решил попытаться в очередной раз.
      <br/>
      Боги вечны, и у нас всегда найдется время на игры...
    </p>
    <p v-else>
      At that moment, I was merely an observer, listening to their dialogue - Odin and his ravens. Odin lowered his
      head, and I noticed his melancholy. I, on the other hand, remained undaunted by despondency. An ingenious idea had
      occurred to me on how to deceive Odin and put him in an awkward position (after all, I, Loki, am renowned for my
      cunning and trickery). I had often attempted to outsmart the Allfather, striving to prove my superiority. Though
      my previous attempts to deceive Odin had been in vain, I resolved to try once more.
      <br/>
      Gods are eternal, and there is always time for games.
    </p>

    <p v-if="locale === 'ru'">
      Мне пришла в голову мысль, что Один слишком консервативен для современного мира и новых технологий, созданных
      людьми. Всеотец всегда предпочитал инженерным системам старый добрый меч, магию и древние артефакты.
      Именно эту его слабость я и решил использовать.
    </p>
    <p v-else>
      The thought crossed my mind that Odin was too conservative for the modern world and the new technologies created
      by humans. The Allfather consistently favored ancient weapons like swords, magic, and ancient artifacts over
      engineering systems. It was this very weakness that I decided to exploit.
    </p>

    <blockquote v-if="locale === 'ru'">
      Локи: "Дорогой Один-Всеотец, ты мудрейший среди всех девяти миров, и люди далеко отстают от тебя. Помоги им
      освоить искусственный интеллект, пока они не причинили себе вреда. Ты говорил, что хочешь
      написать пару картин, заняться творчеством. Не представляет ли это для тебя возможность продемонстрировать людям,
      как использовать искусственный интеллект и создать с его помощью свою коллекцию картин?"
    </blockquote>
    <blockquote v-else>
      Loki: "Dear Odin-Allfather, you are the wisest among the nine realms, and humans are far behind you. Help them
      harness artificial intelligence before they harm themselves. You mentioned wanting to create a few paintings, to
      engage in artistry. Could this not be an opportunity for you to demonstrate to humanity how to use artificial
      intelligence and create a collection of paintings with its assistance?"
    </blockquote>

    <div class="has-text-centered">
      <img src="/images/fable/odin_and_loki2.png"/>
    </div>

    <p v-if="locale === 'ru'">
      Один улыбнулся. Я уже понял, что он узрел мою идею и коварство, как впрочем и всегда. Но еще он и задумался.
      Один - творец, и творчество его основа. Но эти современные технологии землян... Мысль о том, что придется иметь
      дело с примитивными земными технологиями, явно вызывала у Одина недоумение:
    </p>
    <p v-else>
      Odin smiled. I could tell he had recognized my idea and cunning, as he always did. Yet, he also contemplated it.
      Odin is a creator, and creativity is at his core. However, the modern technologies of Earth... The thought of
      dealing with primitive Earth technologies seemed to baffle Odin:
    </p>

    <blockquote v-if="locale === 'ru'">
      "Что за чудеса творят эти жители Мидгарда? Неужели они решили, что они боги, создавая искусственный интеллект? И
      какие картины сможет нарисовать этот искусственный интеллект? Бездушные и лишенные эмоций? Электронные подобия
      настоящего искусства."
    </blockquote>
    <blockquote v-else>
      "What wonders these inhabitants of Midgard have conjured! Have they truly assumed the roles of gods by creating
      artificial intelligence? And what kind of paintings can this artificial intelligence produce? Lifeless and devoid
      of emotion? Electronic replicas of true art."
    </blockquote>

    <p v-if="locale === 'ru'">
      Но Один, несмотря на отсутствие одного глаза, видел сквозь меня и мой хитрый план.
      Мои коварные козни были для него далеко не новостью, а скорее обычным делом. Он любил играть в эту игру,
      поддакивая моим идеям, чтобы затем разрушить мои планы и сделать меня посмешищем перед всем Асгардом.
      <br/>
      В эту игру мы играем тысячи лет.
    </p>
    <p v-else>
      Yet, Odin, despite his one eye, saw through me and my cunning plan. My crafty schemes were nothing new to him but
      rather a routine. He enjoyed playing this game, at times endorsing my ideas only to later dismantle my plans,
      making me a mockery in front of all of Asgard.
      <br/>
      We've been playing this game for millennia.
    </p>

    <p v-if="locale === 'ru'">
      Всеотец девяти миров решил научить меня уроку и в этот раз. Но для этого ему придётся самому разобраться и создать
      картины с использованием непредсказуемой, примитивной и механической земной технологии - искусственного
      интеллекта.
    </p>
    <p v-else>
      The Allfather of the nine realms decided to teach me a lesson once again. But this time, he would have to delve
      into the unknown and create paintings using unpredictable, primitive, and earthly technology - artificial
      intelligence.
    </p>

    <blockquote v-if="locale === 'ru'">
      Один: "Хорошо, Локи. Ты снова решил поиграть со мной, и я снова научу тебя уроку, несомненно. Но люди не узнают,
      что это мои руки создали. Им этого и не доказать, даже если я бы решил спустить эти картины с небес на землю
      на виду у всего человечества.
      <br/>
      Но ты, Локи, увидишь это своими глазами и поймешь, какой ошибкой было считать себя умнее меня." - сказал Один.
      <br/>
      "По рукам!" - ответил я.
    </blockquote>
    <blockquote v-else>
      Odin: "Very well, Loki. You have decided to play with me once more, and I will, without a doubt, teach you a
      lesson. The humans will not know that my hands crafted these paintings. They won't be able to prove it, even if I
      were to descend from the heavens to Earth with these paintings before all of humanity.
      <br/>
      But you, Loki, will witness
      it with your own eyes and understand the mistake of thinking yourself smarter than me," Odin declared.
      <br/>
      "Deal!" I replied.
    </blockquote>

    <p v-if="locale === 'ru'">
      Всеотец пригубил последний глоток остывшего кофе и добавил:
    </p>
    <p v-else>
      The Allfather took his final sip of now-cooled coffee and added:
    </p>

    <blockquote v-if="locale === 'ru'">
      Один: Назовём это "artchitect".
    </blockquote>
    <blockquote v-else>
      Odin: Let us name this <b>artchitect</b>".
    </blockquote>
  </section>
</template>

<script>
export default {
  name: "loki-fable",
  props: ["locale"]
}
</script>

<style lang="scss" scoped>

</style>
