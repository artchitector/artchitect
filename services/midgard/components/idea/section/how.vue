<i18n>
    {}
</i18n>
<template>
    <section class="content" v-if="locale === 'ru'">
        <a id="how"></a>
        <h6>Как Artchitect рисует картину сам, без помощи людей?</h6>
        <p>
            В самом начале рисования Artchitect должен принять несколько решений, определяющих суть картины.
        </p>
        <blockquote class="mb-0">
            Машина не знает, что ей рисовать, у неё нет ни подсознания, ни чувств. Машина ни к чему не стремится, у неё нет
            своих идей о
            том, что вообще делать.<br />
            Идея проекта Artchitect - дать машине воображение, подсознание и творческую случайность.
        </blockquote>
        <p class="mb-0">
            <br />
            Для непрерывной генерации идей Artchitect подключен к датчику энтропии, из которого он получает float64-числа.
            Эти
            числа специальными алгоритмами выбора преобразуются в идею картины.
        </p>
        <div v-if="entropy" class="has-text-centered">
            <common-entropy :entropy="entropy" />
        </div>
        <div class="has-text-centered">
            <NuxtLink :to="localePath('entropy')">О датчике энтропии</NuxtLink>
        </div>
        <p class="mb-0">
            Artchitect непрерывно общается с фоновым пространством (энтропией), получая ответы в виде чисел.
            <br />
            <br />
            На языке человеческого общения построение каждой картины выглядит так:
        </p>
        <p class="pl-4 pt-2">
            <i>Artchitect:</i> Привет! Я начинаю рисовать картину номер #16269.
            <br /> <i>Энтропия:</i> Привет! Что тебе нужно знать для этого?
            <br /> <i>Artchitect:</i> Сколько взять ключевых слов для новой идеи, от 1 до 28 слов?
            <br /> <i>Энтропия:</i> Возьми <b>4</b>. <span class="has-text-grey is-size-7">[В этот момент считывается кадр
                с
                веб-камеры, по которому опреляется число на
                шкале от 1 до 28]</span>
            <br /> <i>Artchitect:</i> Принято 4 слова. У меня есть словарь из 1250 сущностей, и я хочу выбрать первое
            ключевое слово из двенадцати. Какое слово мне выбрать из своего словаря?
            <br /> <i>Энтропия:</i> Возьми слово <b>№175</b> из своих 1250 слов. <span
                class="has-text-grey is-size-7">[числа 175 и 1250 выдуманы, совпадение не точное]</span>
            <br /> <i>Artchitect:</i> Отлично. Это слово "architect". Давай второе слово...
            <br /> <span class="has-text-grey is-size-7">... так Artchitect получает все 4 слова</span>
            <br /> <i>Artchitect:</i> Хорошо, у меня теперь есть 4 слова: "architect,rose,romanticism,abstract".
            <br /> <i>Artchitect:</i> Я могу нарисовать 4294967295 различных вариантов картины с такими словами. Подскажи,
            какой вариант мне рисовать?
            <br /> <i>Энтропия:</i> Выбери вариант <b>№3202070715</b>.
            <br /> <i>Artchitect:</i> Спасибо! Это всё, что мне нужно знать. Начинаю рисовать.
            <br />
            <br />
            Далее Архитектор отправляет полученные решения на вход Stable Diffusion AI и получает изображение.<br />Так
            Архитектор создал картину
            <NuxtLink :to="localePath(`/art/16269`)">#16269</NuxtLink>
        </p>

        <div class="has-text-centered">
            <NuxtLink :to="localePath(`/art/16269`)">
                <img src="https://artchitect.space/api/image/16269/m" style="max-height: 400px" />
            </NuxtLink>
            <br />
            <span class="is-size-7">Картина с номером #16269. Ключевые слова - "architect,rose,romanticism,abstract".
                Seed = 3202070715</span>
        </div>
    </section>
    <section v-else class="content">
        <a id="how_to_draw"></a>
        <h6>How does Artchitect create a picture without human help?</h6>
        <p>
            At the beginning of creative process, Artchitect must make several decisions that determine the essence of
            the
            picture.
        </p>
        <blockquote class="mb-0">
            The machine does not know what to draw, it has no subconscious, no feelings. The machine does not strive for
            anything, it does not have its own ideas about what to do at all.<br />
            The idea of the Artchitect project is to give imagination, subconsciousness and creative randomness to machine.
        </blockquote>
        <p class="mb-0">
            For continuous creation Architect is connected to an entropy sensor, from which it receives float64
            numbers. These numbers are transformed by special selection algorithms into the idea of the picture.
        </p>
        <div v-if="entropy" class="has-text-centered">
            <common-entropy :entropy="entropy" />
        </div>
        <div class="has-text-centered">
            <NuxtLink :to="localePath('entropy')">About entropy sensor</NuxtLink>
        </div>
        <p class="mb-1">
            Artchitect constantly communicates with the background space (entropy), asking questions and getting answers.
            <br />
            <br />
            If we imagine a dialogue between Artchitect and Entropy, the creation of the picture will look like this:
        </p>
        <p class="pl-4">
            <i>Artchitect:</i> Hello! I'm starting to draw a picture #16269.
            <br /> <i>Entropy:</i> Hi! What do you need to know for this?
            <br /> <i>Artchitect:</i> How many keywords I need to take for a new idea, from 1 to 28 words?
            <br /> <i>Entropy:</i> Take <b>4</b>.
            <span class="has-text-grey is-size-7">[At this moment, a frame is read from the webcam, which determines the
                number on a scale from 1 to 28]</span>
            <br /> <i>Artchitect:</i> 4 words accepted. I have a dictionary of 1250 entities and I want to choose the first
            keyword out of three. Which word should I choose from my dictionary?
            <br /> <i>Entropy:</i> Take word <b>#175</b> out of your 1,250 words. <span
                class="has-text-grey is-size-7">[numbers 175 and 1250 are made up, the match is not exact]</span>
            <br /> <i>Artchitect:</i> Great. This is the word "architect". Give me the second word...
            <br /> <span class="has-text-grey is-size-7">... so Artchitect gets all 4 words</span>
            <br /> <i>Artchitect:</i> OK, I now have 4 words: "architect,rose,romanticism,abstract".
            <br /> <i>Artchitect:</i> I can draw 4294967295 different versions of the picture with these words. Tell me,
            which
            version should I draw?
            <br /> <i>Entropy:</i> Choose version <b>#№3202070715</b>.
            <br /> <i>Artchitect:</i> Thank you! That's all I need to know. I'm starting to draw.
            <br />
            <br />
            Next, Artchitect sends the received idea to the Stable Diffusion AI input (as prompt+seed) and receives an
            image.<br />So Artchitect created the painting
            <NuxtLink :to="localePath(`/art/16269`)">#16269</NuxtLink>
        </p>

        <div class="has-text-centered">
            <NuxtLink :to="localePath(`/art/16269`)">
                <img src="https://artchitect.space/api/image/16269/m" style="max-height: 400px" />
            </NuxtLink>
            <br />
            <span class="is-size-7">Art #16269. Keywords - "emotional,Big Bang,death of civilizations". Seed =
                3201685114</span>
        </div>
    </section>
</template>
<script>
import Radio from "@/utils/radio";
export default {
    name: "idea-section-how",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    },
    data() {
        return {
            radioPid: null,
            entropy: null,
        }
    },
    async mounted() {
        this.radioPid = await Radio.subscribe("entropy", (entropy) => {
            this.entropy = entropy
        })
    },
    beforeDestroy() {
        if (this.radioPid !== null) {
            Radio.unsubscribe(this.radioPid)
        }
    }
}
</script>
<style lang="scss" scoped></style>