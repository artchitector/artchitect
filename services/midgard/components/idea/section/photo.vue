<i18n>
    {}
</i18n>
<template>
    <section class="content">
        <div class="has-text-centered">
            <a href="/artchitect_in_real_world.jpg">
                <img src="/artchitect_in_real_world.jpg" alt="Artchitect in physical world" style="max-height: 300px" />
            </a>
            <br />
            <span v-if="locale === 'ru'" class="is-size-7">Нейронный мозгом Архитектора является RTX 3060 12Gb.
                <br />Каждую минуту видеокарта производит свой труд, греясь и разгоняя кулеры. Затем утихает,
                отдыхает.</span>
            <span v-else class="is-size-7">The neural brain of the Artchitect is RTX 3060 12Gb.
                <br />Every minute the video card does its work, warming up and overclocking the coolers. Then it calms
                down, rests.</span>
        </div>
    </section>
</template>
<script>
export default {
    name: "idea-section-photo",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>