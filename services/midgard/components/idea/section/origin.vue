<i18n>
    {}
</i18n>
<template>
    <section v-if="locale === 'ru'" class="content">
        <h6>Что является произведением в Artchitect?</h6>
        <p>Artchitect производит не картины (хотя именно их мы наблюдаем в конечном итоге). Artchitect производит цифровые
            идеи, из
            которых эти картины могут быть воссозданы на подходящем оборудовании. Эти идеи состоят из слов и чисел, которые
            были
            получены из окружающей энтропии в момент начала "рисования" очередной работы. Именно эти идеи являются
            источником
            визуальных изображений, создаваемых Архитектором. Каждая идея имеет свой уникальный номер, который отпечатан в
            углу
            каждой картины. По номеру всегда можно найти исходник на сайте <a
                href="https://artchitect.space">artchitect.space</a>
            и отследить оригинальную идею работы.
        </p>
    </section>
    <section v-else class="content">
        <h6>What is an artwork in Artchitect?</h6>
        <p>Artchitect does not produce paintings (although they are what we observe in the end). Artchitect produces digital
            ideas from which these paintings can be recreated on suitable equipment. These ideas consist of words and
            numbers
            that were obtained from the surrounding entropy at the time of the beginning of the "drawing" of the next work.
            It
            is these ideas that are the source of visual images created by the Artchitect. Each idea has its own unique
            number,
            which is printed in the corner of each painting. By number, you can always find the source code on the
            <a href="https://artchitect.space">artchitect.space</a> website and track the original idea of the work.
        </p>
    </section>
</template>
<script>
export default {
    name: "idea-section-origin",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>