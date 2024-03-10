# Artchitect

![artchitect_logo](files/images/logo_anim_92.gif)

#### https://artchitect.space

### Documentation and code comments are completely in Russian. Please, use translator if you need to understand.

> Artchitect (или Архитектор) - это удивительная автономная творческая машина, способная создавать великолепные
> картины,
> вдохновляясь окружающей нас Вселенной. В своем непрерывном творчестве машина черпает вдохновение из естественной
> энтропии
> Вселенной, представленной в виде фонового света, и создает уникальные произведения без участия человека.

### Алгоритм преобразования энтропии (света) в число

![Алгоритм преобразования энтропии (света) в число](files/images/entropy_extraction_algorithm.png)

### Архитектура программного кода

> Если код не читается, как приключенческая книга - его никто не захочет читать.
> 
> На сайте проекта есть [сказочная история](https://artchitect.space/loki) об идее возниковения artchitect.

Код artchitect написан в духе скандинавской мифологии. Картины сотворяет верховный бог
[Асгарда](services/asgard/) -
[Odin](services/asgard/pantheon/odin.go)
, ему помогают его
верные [Huginn](services/asgard/pantheon/huginn.go), [Muninn](services/asgard/pantheon/muninn.go),
и [многие другие](services/asgard/pantheon).

#### Главный цикл творения

1. процесс творения запускается через [Стремление (Intention)](services/asgard/pantheon/intention.go). Может
   быть [призван](services/asgard/pantheon/intention.go#L62) Odin для сотворения,
   или [призвана](services/asgard/pantheon/intention.go#L53) Frigg для объединения.
2. [Odin](services/asgard/pantheon/odin.go) начинает сотворение картины. Первым делом
   он [вспоминает](services/asgard/pantheon/odin.go#L183) номер для неё.
3. Далее Один [придумывает идею картины](services/asgard/pantheon/odin.go#L307). Ему нужно придумать её в таком виде, в
   котором поймёт ИИ Stable Diffusion - в форме [seed-числа](services/asgard/pantheon/odin.go#339)
   и [набора ключевых слов](services/asgard/pantheon/odin.go#L342).
4. Чтобы придумать идею для новой картины, Один смотрит в хаос мироздания
   своим [пустым глазом](services/asgard/pantheon/lost_eye.go#L183), а
   ворон [Хугин](services/asgard/pantheon/huginn.go) ("думающий")
   помогает
   этот хаос [интерпретировать](services/asgard/pantheon/huginn.go#L166) в виде цифрового
   слепка [энтропии](model/entropy.go#L35). В виде int64- и float64-числа.
5. Интерперитованная энтропия есть лишь число, но для окончательного формирования идеи Один вспоминает конкретные
   [понятия](services/asgard/pantheon/odin.go#L318)
   и [seed-номер](services/asgard/pantheon/odin.go#L308). Вспоминать помогает второй ворон
   Одина - [Мунин](services/asgard/pantheon/muninn.go) ("помнящий").
   Он [превращает](services/asgard/pantheon/muninn.go#L48)
   хаотичную энтропию в реальные "предметы",
   выбирая нужный [предмет из списка всех](services/asgard/pantheon/muninn.go:109), которые помнит. Один ему говорит,
   какой предмет выбирать через [float64-число](services/asgard/pantheon/muninn.go#L114)
   из энтропии, которую ранее [осмыслил](services/asgard/pantheon/huginn.go#L166) Хугин. float64 число работает как
   указатель на шкале, "например, выбери 100й
   предмет из 1250"
6. Когда идея готова, она наполнена всеми словами и seed-номером, то
   Один [передаёт](services/asgard/pantheon/odin.go#L233) эту идею [Фрейе](services/asgard/pantheon/freyja.go), которая
   уже напрямую [работает](services/asgard/pantheon/freyja.go#L43) с искусственным интеллектом,
   формируя картину в виде [цифрового изображения](services/asgard/pantheon/freyja.go#L60).
7. На готовое изображение Один с помощью своего копья Gungner [наносит](services/asgard/pantheon/odin.go#L238) свою
   подпись - порядковый номер картины в углу.
8. Далее всё сотворённое [сохраняется в хранилище](services/asgard/pantheon/odin.go#L244) цифровой галереи Artchitect.
9. Раз в некоторое время Frigg [берёт на себя управление](services/asgard/pantheon/frigg.go#L87), чтобы собрать
   очередное единство. Это случается 1 раз после 10 написанных картин.
10. Все данные в реальном времени, включая [расшифрованные изображения энтропии](services/asgard/pantheon/heimdallr.go#L164), транслируются в [Midgard](services/midgard) (фронтэнд artchitect)
    через [Alfheimr](services/alfheimr) (api-gateway artchitect). В процессе
    трансляции [Heimdallr](services/asgard/pantheon/heimdallr.go) отправляет пакеты
    драккарами
    через [радужный мост Bifröst](/home/artchitector/artchitect/artchitect/services/asgard/communication/bifrost.go), на
    другом конце которого светлые эльфы [получат груз](services/alfheimr/communication/harbour.go)
    и [передадут](services/alfheimr/portals/radio.go) его дальше
    в [Мидгард](services/midgard/components/insight/main.vue#L66)
    людям, чтобы они тоже смогли увидеть картины своими глазами.

![artchitecture](files/images/artchitecture.jpg)

#### Рабочий стенд Artchitect

Stable Diffusion AI v1.5 работает на Nvidia RTX 3060 12Gb.
![artchitect_installation](files/images/artchitect_hardware.jpg)