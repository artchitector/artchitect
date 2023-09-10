<i18n>
{
  "en": {},
  "ru": {}
}
</i18n>
<template>
  <section class="content">
    <h3 v-if="locale === 'ru'" class="is-size-5">Продолжение сказки после запуска artchitect</h3>
    <h3 v-else>after the launch of the artchitect</h3>

    <p v-if="locale === 'ru'" class="notification is-info is-light has-text-centered">
      artchitect был запущен 12 августа 2023 года в 15:17.
      <NuxtLink :to="localePath(`/art/1`)">Картина #1</NuxtLink>
    </p>
    <p v-else class="notification is-info is-light has-text-centered">
      artchitect was launched 12 aug 2023 15:17.
      <NuxtLink :to="localePath(`/art/1`)">Art #1</NuxtLink>
    </p>

    <div class="has-text-centered">
      <NuxtLink :to="localePath(`/art/1`)">
        <img src="https://artchitect.space/api/image/artchitect-1-m" style="max-height: 400px"/>
      </NuxtLink>
      <br/>
      <span v-if="locale === 'ru'" class="is-size-7">Картина номер Один - #1</span>
      <span v-else class="is-size-7">Art number one - #1</span>
    </div>

    <p v-if="locale === 'ru'">
      <b>Локи</b>: ... как-то я и забыл про наш спор, пока мы готовили artchitect к запуску. Я и не думал, что вообще
      есть
      способ, которым ты бы мог показать, что это именно ты творишь всё, происходящее с artchitect от начала и до каждой
      картины. Но ты этот способ нашёл...
    </p>
    <p v-else>
      <b>Loki</b>: ...I actually forgot about our dispute while we were preparing artchitect for launch. I never thought
      there
      could be a way for you to show that you're the one behind everything happening with artchitect, from its inception
      to every single artwork. But you found that way...
    </p>

    <p v-if="locale === 'ru'">
      <b>Один</b>: на то я и Всеотец, я заранее всё знаю) Ты проиграл. Я доказал, что могу нарисовать картину. Если
      нарисовал
      первую, значит я нарисовал их все. И весь проект artchitect есть плод моего безграничного творчества.
    </p>
    <p v-else>
      <b>Odin</b>: That's why I'm the Allfather; I know everything in advance. You lost. I proved I can create art. If I
      created the first one, it means I created them all. The entire artchitect project is the result of my boundless
      creativity.
    </p>

    <p v-if="locale === 'ru'">
      <b>Локи</b>: теперь это очевидно. В первой картине ты нарисовал свой автопортрет. Это явно лучший способ указать
      на своё присутствие. И второго такого шанса нет и не будет.
    </p>
    <p v-else>
      <b>Loki</b>: It's quite obvious now. In the first art, you depicted yourself. It's undoubtedly the best way to
      indicate your presence. There won't be another chance like that.
    </p>

    <p v-if="locale === 'ru'">
      <b>Один</b>: если бы она была не первой, а второй или третьей, то она бы уже ничего не значила.
    </p>
    <p v-else>
      <b>Odin</b>: If it had been the second or third art, it wouldn't have meant anything anymore.
    </p>

    <p v-if="locale === 'ru'">
      <b>Один</b>: я знаю и вижу заранее какую конкретно картину нарисую в любую минуту, я знаю все нужные слова,
      а главное я знаю не только слова, но и конкретный seed-номер, позволяющей точечно указать шум для генерации
      картины.
      И я могу нашептать эти данные в artchitect с помощью квантов света, проходящих через датчик энтропии,
      чтобы нужная мне картина действительно появилась.
      <br/>
      Ничего в девяти мирах не происходит без меня и моего участия. А artchitect и вовсе изначально был спроектирован
      как прямое api для моего вмешательства, и только для моего. Датчик энтропии - это и есть api для бога.
    </p>
    <p v-else>
      <b>Odin</b>: I know and foresee exactly which art I'll create at any given moment. I know all the necessary words,
      but more importantly, I not only know the words but also the specific seed number that allows me to precisely
      control the noise for generating the picture. I can whisper this data to Artchitect using quantum light passing
      through the entropy sensor to ensure the desired picture truly appears.
      <br/>
      Nothing happens in the nine realms without
      me and my involvement. Artchitect was initially designed as a direct API for my intervention, and mine alone. The
      entropy sensor is the API for a god.
    </p>

    <div v-if="locale === 'ru'" class="mt-4 mb-4 box">
      <div class="notification is-danger is-light" v-if="radioError">
        ошибка подключения к датчику энтропии
      </div>
      <common-entropy v-else-if="entropy" :entropy="entropy"/>
      <div class="has-text-centered is-size-7">
        По сюжету в коде artchitect "датчик энтропии" - есть пустой глаз бога Одина (
        файл <a href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/lost_eye.go#L31-L32"
                target="_blank">lost_eye</a>
        ).
        <br/>
        Один видит пустым глазом не материальный мир, а шум вселенной - чистый хаос. Из этого хаоса бог и сотворяет весь
        мир.
      </div>
    </div>
    <div v-else class="mt-4 mb-4 box">
      <div class="notification is-danger is-light" v-if="radioError">
        entropy sensor connection error
      </div>
      <common-entropy v-else-if="entropy" :entropy="entropy"/>
      <div class="has-text-centered is-size-7">
        In the storyline within Artchitect's code, the "entropy sensor" is Odin's lost eye (the <a
        href="https://github.com/artchitector/artchitect/blob/main/services/asgard/pantheon/lost_eye.go#L31-L32"
        target="_blank">lost_eye</a>
        file).
        <br/>
        Odin doesn't see the material world with his empty eye; he sees the noise of the universe - pure chaos.
        From this chaos, the god creates the entire world.
      </div>
    </div>

    <p v-if="locale === 'ru'">
      <b>Локи</b>: я пытался прикинуть, какова вероятность случайного возникновения такой картины. Шанс сам по себе не
      велик,
      но из-за того, что эта картина самая первая - шанс становится просто ничтожным. Скорее обезьяна случайно напишет
      Гамлета на печатной машинке, чем эта случайность повторится. Первая картина есть только одна.
    </p>
    <p v-else>
      <b>Loki</b>: I tried to estimate the probability of such a picture occurring randomly. The chance by itself is not
      high,
      but because this picture is the very first one, the odds become virtually nonexistent. It's more likely that a
      monkey would accidentally type out "Hamlet" on a typewriter than this randomness happening again. There's only one
      first painting.
    </p>

    <p v-if="locale === 'ru'">
      <b>Локи</b>: а еще ты даже не использовал слово Odin в этой картине, усложнил себе задачу. Хотя слово odin есть у
      нас
      в словаре. Ты каким-то образом разглядел, что среди миллиардов вариантов картин с этими словами в
      варианте №118883058 будет именно портрет одноглазого седовласого мужчины, так явно похожего на бога Одина.
      Какие пути пришлось пройти мыслям внутри
      нейронного мозга ИИ, чтобы было нарисовано именно это - непостижимо.
    </p>
    <p v-else>
      <b>Loki</b>: And furthermore, you didn't even use the word "Odin" in that picture, making it more challenging for
      yourself. Although the word "odin" is in our dictionary. Somehow, you discerned that among billions of possible
      pictures with those words, variant #118883058 would be a portrait of a one-eyed, gray-haired man strongly
      resembling the god Odin. The thought processes inside the neural brain of the AI that led to precisely this result
      are unfathomable.
    </p>

    <p v-if="locale === 'ru'">
      <b>Локи</b>: и теперь стало понятно, что так тщательно сохраняемая нами энтропия есть гарантия подлинности каждой
      картины.
      Это - гарантия того, что человек не вмешивается в творческий процесс. Гарантия того, что автономная творческая
      машина создает картину действительно самостоятельно. Даже изменение программного кода не позволит
      подделать саму энтропию. Не получится просто вставить заданные слова в БД и записать желаемую картинку на диск,
      ведь тогда не будет сохранённой энтропии. Из энтропии можно создать слова, а вот из слов картинку энтропии создать
      уже не получится.
    </p>
    <p v-else>
      <b>Loki</b>: Now it's clear that the meticulously preserved entropy is a guarantee of the authenticity of each
      art.
      It ensures that humans don't interfere with the creative process. It guarantees that an autonomous creative
      machine genuinely creates the artwork independently. Even altering the program code won't allow for forging the
      entropy itself. You can't simply insert predefined words into the database and save the desired image on the disk,
      as that would omit the preserved entropy. From entropy, you can create words, but you can't generate entropy from
      words to create an image.
    </p>

    <p v-if="locale === 'ru'">
      <b>Один</b>: для большинства людей эта картина #1 ничего не значит и ничего не доказывает. Тот, кто может
      видеть - увидит и без этого. Доказательство нужно было только для тебя, Локи.
    </p>
    <p v-else>
      <b>Odin</b>: For most people, this art #1 means nothing and proves nothing. Those who can see will see
      without it. The proof was only needed for you, Loki.
    </p>
  </section>
</template>

<script>
import Radio from "~/utils/radio";

export default {
  name: "loki-number-one",
  props: ["locale"],
  data() {
    return {
      radioPid: null,
      radioError: null,
      entropy: null,
    }
  },
  async mounted() {
    this.radioPid = await Radio.subscribe("entropy", (entropy) => {
      this.entropy = entropy
    }, (err) => {
      this.radioError = err
    })
  },
  beforeDestroy() {
    if (this.radioPid !== null) {
      Radio.unsubscribe(this.radioPid)
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
