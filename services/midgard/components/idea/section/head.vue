<i18n>
    {}
</i18n>
<template>
    <section class="content">
        <div class="has-text-centered">
            <img src="/jesus_anim_92.gif" alt="artchitect" style="max-height: 64px" />
        </div>

        <hr />
        <p v-if="locale == 'ru'">Ключевым компонентом Artchitect является искусственный интеллект <a
                href="https://github.com/Stability-AI/stablediffusion">Stable Diffusion</a>, доступный по лицензии MIT.
            Все картины нарисованы этим инструментом, и сама
            художественная часть картины целиком придумывается в недрах этой технологии.<br />
            Кроме того Artchitect благодарит проект <a href="https://github.com/invoke-ai/InvokeAI">InvokeAI</a> за удобную
            рабочую сборку этого ИИ.
        </p>
        <p v-else>A core-component of Artchitect is the <a href="https://github.com/Stability-AI/stablediffusion">Stable
                Diffusion</a> AI, available under the MIT license.
            All the paintings are drawn with this tool, and the artistic part of the painting itself is entirely invented in
            the depths of this technology.
            In addition, Artchitect thanks the <a href="https://github.com/invoke-ai/InvokeAI">InvokeAI</a> project for a
            convenient working assembly of this AI.
        </p>
    </section>
</template>
<script>
export default {
    name: "idea-section-head",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>