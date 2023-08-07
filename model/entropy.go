package model

import (
	"fmt"
	"image"
	"time"
)

// EntropySize - Энтропия считается по квадрату 8 на 8 пикселей
// Odin: Фригг, где мои очки? Я не могу разобрать, что тут нарисовано! Тысяча чертей в эту допотопную машину!
const EntropySize = 8
const EntropyJpegQualityFrame = 65
const EntropyJpegQualityNoise = 20

// EntropyMatrix - матрица хранит силу каждого из 64 пикселей в энтропии
type EntropyMatrix struct {
	power [EntropySize][EntropySize]uint8 // матрица силы пикселей
}

func (em *EntropyMatrix) Size() int {
	return EntropySize
}

func (em *EntropyMatrix) Set(x, y int, power uint8) {
	em.power[x][y] = power
}

func (em *EntropyMatrix) Get(x, y int) uint8 {
	return em.power[x][y]
}

type Entropy struct {
	IntValue   uint64  `json:"int"`
	FloatValue float64 `json:"float"`

	ImageEncoded string `json:"image"`   // base64-encoded 8x8 PNG изображение энтропии
	ImageID      string `json:"imageId"` // ключ для получения изображения энтропии с memory-сервера
	// Odin: ImageID тут нужен не для потоковой передачи энтропии, а для сохранения принятых когда-то решений.
	// Например, если я придумал картину, и вспомнил для неё seed-номер 213712211,
	// то ДОЛЖНО сохранить и видимую мной энтропию (прямую и обратную),
	// которая была использована для получения этого номера. Так Artchitect будет иметь полную историю
	// событий, которые происходили с картиной, включая источники всех решений.
	// Картинки этой энтропии сохраняются в файловое хранилище, а ссылки на картинки сохраняются в это поле.

	// Odin: Потоковая энтропия так не сохраняется, она нужна лишь для отображения на клиенте.
	// Odin: Если вы, смертные, захотите сохранять всё видимое моим пустым глазом, то у вас дисков не хватит на всей Земле.
	// Odin: сжалюсь над вами и сэкономлю вам "пару гигабайт", мои любимые мидгардцы
	// Loki: Ты работаешь над проектом третий день, а уже разжалобился к людям.
	// Loki: Где же твои "Чтоб горели эти ваши примитивные технологии земные технологии! Вот у нас в Асгарде
	//		 самый вкусный пломбир по 5 копеек!!!!"?
	// Odin: не забывайся, что я верховный бог в Асгарде.
	// Loki: Так и есть, но если ты не сможешь разобраться в искусственном интеллекте людей, а я думаю, что не сможешь -
	//		то твоё место на троне твоего двонца займу Я - Loki! Бог хитрости!
	// Odin: [ссылает Локи в темницу на 100 лет] Да охлади ты свою гордыню, бедный сын ётунов...

	Matrix EntropyMatrix `json:"-"` // Это нужно только в Асгарде для преобразований, по сети не уходит
}

func (e Entropy) String() string {
	return fmt.Sprintf("E:%.6f", e.FloatValue)
}

// EntropyPack - рассчитанная энтропия.
// Постоянно отправляется на клиент, даже когда не используется в работе. Видна повсеместно на сайте.
type EntropyPack struct {
	Timestamp time.Time `json:"timestamp"`
	Entropy   Entropy   `json:"entropy"`
	Choice    Entropy   `json:"choice"`
}

// EntropyPackExtended - энтропия с подробным описанием. Видна только на странице /entropy на сайте.
// К остальным данным тут еще добавляются кадр с камеры и шум для наглядности, так что событие это объёмное.

type EntropyPackExtended struct {
	Timestamp time.Time `json:"timestamp"`
	Entropy   Entropy   `json:"entropy"`
	Choice    Entropy   `json:"choice"`

	ImageFrame        image.Image `json:"-"`          // сами картинки передаются только в памяти сервиса Асгард
	ImageFrameEncoded string      `json:"imageFrame"` // base64-encoded jpeg картинки (уходят в Мидгард)

	ImageNoise        image.Image `json:"-"`          // сами картинки передаются только в памяти сервиса Асгард
	ImageNoiseEncoded string      `json:"imageNoise"` // base64-encoded jpeg картинки (уходят в Мидгард)
}
