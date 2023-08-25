<i18n>
    {}
</i18n>
<template>
    <section v-if="locale === 'ru'" class="content">
        <h6>Как машина черпает вдохновение из естественной энтропии Вселенной? Как она интерпретирует и использует эту
            информацию?</h6>
        <p>
            Для процесса творения Архитектор использует случайные числа-зерна (seeds), которые являются исходными точками
            для
            создания картины.<br />
            Числа-зерна считываются веб-камерой из случайных кадров. При этом, не содержимое кадра, а только
            световой шум используется для генерации случайных чисел.<br />
            Световой шум настолько различается на соседних кадрах, даже в темноте, что обеспечивает отличную случайную
            генерацию чисел. Архитектор принимает десятки тысяч самостоятельных решений каждый день и считывает десятки
            тысяч
            кадров, что поддерживает непрерывный процесс творения.
            <br />
            <NuxtLink :to="localePath('entropy')">Алгоритм преобразования энтропии в выбор в реальном времени</NuxtLink>
        </p>

        <div class="has-text-centered">
            <img src="/one_frame.jpeg" height="480" width="640" />
            <br />
            <span class="is-size-7">Кадр пространства, на которое Архитектор смотрит каждый день. <br />Из подобных кадров,
                где
                изображение не меняется месяцами, он и рисует все свои картины</span>
        </div>
    </section>
    <section class="content" v-else>
        <h6>How does a machine gets inspiration from the natural entropy of the universe? How does she interpret and use
            this information?</h6>
        <p>
            <b>Artchitect uses entropy of surrounding lightnoise catched by webcam.</b><br />
            In the beginning of creative process Artchitect
            collects seed-numbers, which are the source of decisions needed
            to construct idea of plot. Seed-numbers are read with the webcam from random frames. <br />
            At the same time, not the content of the frame, but
            only light noise is used to generate random numbers. The light noise is so different on adjacent frames, even in
            the dark, which provides excellent random number generation.
            <br />Artchitect makes thousands of independent decisions every day and reads thousands
            of frames, which supports the continuous creative process.
            <br />
            <NuxtLink :to="localePath('entropy')">Entropy conversion algorithm in realtime</NuxtLink>
        </p>
        <div class="has-text-centered">
            <img src="/one_frame.jpeg" height="480" width="640" />
            <br />
            <span class="is-size-7">Usual Artchitect's frame - source of int64 number .
                From such frames, where the image does not change for months, Artchitect draws all his arts</span>
        </div>
    </section>
</template>
<script>

export default {
    name: "idea-section-entropy",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>