<i18n>
    {}
</i18n>
<template>
    <section v-if="locale === 'ru'" class="content">
        <h6>В описании говорится, что каждая картина уникальна и неповторима. Существует ли возможность, что машина
            может создать несколько разных картин, которые визуально будут очень похожи друг на друга?</h6>
        <p>
            Не каждое произведение Artchitect красиво и хорошо.
        </p>
        <p>
            В реальности Архитектор далеко не идеален. Есть комбинации, где он создает очень похожие работы, почти
            неотличимые друг от друга. Кроме действительно красивых
            и уникальных работ, иногда появляется значительная часть работ, которые могут быть сломанными или некрасивыми, и
            часто лишены художественного смысла.<br />Архитектор - машина, и для него не существует хороших и плохих
            картин.
        </p>
    </section>
    <section v-else class="content">
        <h6>Is there a possibility that the machine can
            create several different pictures that visually will be very similar to each other?</h6>
        <p>
            Not every art of Artchitect is pretty and good.
        </p>
        <p>
            In real life Artchitect's paintings are far from perfect. There are combinations where Artchitect creates
            very similar works. In addition to really beautiful works, there are a
            significant part that are broken and ugly, often devoid of artistic meaning.<br />
            <b>Artchitect is machine, and there are no good and bad pictures for itself.</b>
        </p>
    </section>
</template>
<script>
export default {
    name: "idea-section-uniqueness",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>