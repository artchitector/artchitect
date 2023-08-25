<i18n>
    {}
</i18n>
<template>
    <section v-if="locale === 'ru'" class="content">
        <h6>Есть ли какая-либо форма регулировки или контроля со стороны человека в процессе создания картин?</h6>
        <p>
            Человек не может вмешиваться в рабочий процесс Архитектора (такого управления просто нет), но человек задает
            изначальные условия - алгоритмы преобразования, и, что самое главное, определяет словарь знаний Архитектора.
            Архитектор может рисовать картины только в рамках своих знаний о мире, в рамках словаря. Словарь Архитектора
            включает более 1200 различных сущностей и стилей, из которых составляются картины.
            <br />
            Текущий словарь можно
            <a target="_blank"
                href="https://github.com/artchitector/artchitect/blob/main/services/asgard/files/dictionary/v1.yaml">посмотреть</a>
            в исходном коде Архитектора.
        </p>
    </section>
    <section class="content" v-else>
        <h6>Is there any human control in the process of creating art?</h6>
      <p>
        Human cannot interfere in the Artchitect's workflow (there is no such controls), but human sets the
        initial conditions - transformation algorithms, and, most importantly, defines the Artchitect's dictionary of
        knowledge. Artchitect can generate ideas only within the set of his knowledge about the real world. Artchitect's
        dictionary includes more than 1200 different entities and styles from
        which arts are made.
        <br/>
        Current dictionary can be found in
        <a target="_blank"
                href="https://github.com/artchitector/artchitect/blob/main/services/asgard/files/dictionary/v1.yaml">sourcecode</a>. There
        are multiple versions of dictionary (smaller, larger). Dictionary can be expanded in future to make more various
        arts.
      </p>
    </section>
</template>
<script>
export default {
    name: "idea-section-regulatory",
    computed: {
        locale() {
            return this.$i18n.locale;
        }
    }
}
</script>
<style lang="scss" scoped></style>