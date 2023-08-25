<i18n>
    {}
</i18n>
<template>
    <section class="content" v-if="locale === 'ru'">
        <h6>Каким образом определяется стиль и содержание создаваемых картин? Есть ли какая-то предустановленная база знаний
            или влияние других художественных произведений?</h6>
        <p>
            Архитектор опирается на свой внутренний словарь сущностей и стилей, который был задан авторами проекта. В этом
            словаре содержится все, что Архитектор знает о окружающем мире. Там перечислены как простые сущности, например,
            "человек", "кот", "солнце", "красный", "мрачно", "радостно", так и некоторые художественные стили или конкретные
            художники. Конечное произведение Архитектора представляет собой комбинацию сущностей и стилей, которые стали
            основой идеи картины. Рядом с каждой картиной Архитектора можно увидеть слова, составляющие идею работы.
        </p>
        <p>Искусственный интеллект Stable Diffusion получает на вход идею сюжета (seed-номер и ключевые слова), а затем
            рекомбинирует все объекты в один кадр в соответствии со своим внутренним "чувством красоты", "золотым сечением".
            Обычно большинство
            ключевых слов игнорируются. Если в ключевых слова картины будет несколько стилей - то в итоге может получиться и
            комбинация стилей, и
            один стиль из нескольких, и вообще ни один из них. Это всё будет решено в недрах "мышления" ИИ.
        </p>
        <div class="has-text-centered">
            <a href="/description_five.jpg" target="_blank">
                <img src="/description_five.jpg" />
            </a>
            <br />
            <span class="is-size-7">Картины с их параметрами</span>
        </div>
    </section>
    <section v-else class="content">
        <h6>How is the style and content of the created paintings determined? Is there some kind of pre-installed knowledge
            base or influence of other works of art?</h6>
        <p>
            Artchitect use internal dictionary of entities and styles, which was set by authors of the
            project. This dictionary contains everything, that Artchitect knows about real world. It lists both simple
            entities, for example, "man", "cat", "sun", "red", "gloomy", "joyful", as well as some artistic styles or
            specific
            artists. The final artwork of the Artchitect is a combination of entities and styles that became the basis of
            the
            idea
            of the painting. Next to each picture of the Artchitect, you can see keywords that make up the idea of the work.
        </p>
        <p>
            Stable Diffusion AI consume the general idea of plot (seed-number and keywords) and then it recombine all
            entities
            to one single frame, according to it's internal "sense of beauty" or "golden ratio". Usually <b>most of the
                keywords are ignored</b>.
            If there are several styles in the keywords of the picture, then in the end you can get a
            combination of styles, and one style from several, and none of them at all. This will all be solved in the
            bowels
            of AI "thinking".
        </p>
        <div class="has-text-centered">
            <a href="/description_five.jpg" target="_blank">
                <img src="/description_five.jpg" />
            </a>
            <br />
            <span class="is-size-7">Arts with related keywords</span>
        </div>
    </section>
</template>
<script>
export default {
    name: "idea-section-dictionary",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>