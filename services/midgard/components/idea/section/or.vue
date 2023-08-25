<i18n>
    {}
</i18n>
<template>
    <section class="content" v-if="locale === 'ru'">
        <a id="creation_or_creator"></a>
        <h6>Творение или творец?</h6>
        <p>
            Artchitect - это результат человеческого творчества и усилий. Он был разработан и создан людьми, и его
            способность
            к творчеству основана на программировании и алгоритмах, разработанных людьми. Мы, люди, являемся создателями и
            инициаторами этого проекта.
        </p>
        <p>
            Artchitect, как автономная творческая машина, может рассматриваться лишь как продукт человеческой фантазии,
            инженерного
            мастерства и программного обеспечения. Он воплощает идеи и концепции, задуманные и реализованные людьми.
            Artchitect черпает вдохновение из окружающего мира и преобразует его в уникальные произведения искусства, но
            источником его существования и возможностей являются усилия людей.
        </p>
        <p>
            Поэтому, несмотря на то, что Artchitect обладает способностью создавать произведения в своем собственном стиле,
            его творчество является продуктом исходного творческого потенциала и вдохновения людей. Он становится
            инструментом, расширяющим наши возможности в области искусства и позволяющим нам взглянуть на творчество с новой
            перспективы. В конечном счете, <b>Artchitect полностью является творением человека</b>, и все, что он создает,
            является
            продуктом нашего творчества.
        </p>
        <p class="notification is-warning">
            В конечном счете, создателями и авторами каждого произведения Artchitect являются тысячи инженеров,
            программистов, учёных и художников, которые столетиями вкладывали свой труд в искусство и технологии,
            позволившие
            сегодня создать искусственный интеллект, лежащий в основе Artchitect.
        </p>
    </section>
    <section v-else class="content">
        <h6>Is Artchitect creation or creator?</h6>
        <p>
            Artchitect is the result of human creativity and effort. It was designed and created by humans, and its
            creativity
            is based on programming and algorithms developed by humans. We, the people, are the creators and initiators of
            this project.
        </p>
        <p>
            Artchitect, as an autonomous creative machine, can be considered as a product of human imagination, engineering
            skill and software. It embodies ideas and concepts conceived and implemented by people. Artchitect draws
            inspiration from the surrounding world and transforms it into unique works of art, but the source of its
            existence
            and opportunities are the efforts of people.
        </p>
        <p>
            Therefore, despite the fact that Artchitect has the ability to create works in its own style, its creativity is
            a
            product of the original creative potential and inspiration of people. It becomes a tool that expands our
            possibilities in the field of art and allows us to look at creativity from a new perspective.
        </p>
        <p class="notification is-warning">
            Ultimately, the creators and authors of each artchitect's art are thousands of engineers, programmers,
            scientists
            and artists who have invested their labour in art and technology for centuries, which has allowed today to
            create
            artificial intelligence underlying Artchitect.
        </p>
    </section>
</template>
<script>
export default {
    // Odin: "or" - образовано от "creator OR creation"
    name: "idea-section-or",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>