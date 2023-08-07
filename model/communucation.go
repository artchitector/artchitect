package model

import "time"

const (
	ChanEntropy         = "entropy"
	ChanEntropyExtended = "entropy_extended"
)

// Здесь расположены модели, необходимые для передачи рантаймовых событий от soul до клиентов

// EntropyState - рассчитанная энтропия.
// Постоянно отправляется на клиент, даже когда не используется в работе. Видна повсеместно на сайте.
type EntropyState struct {
	Timestamp time.Time `json:"timestamp"`
	Entropy   Entropy   `json:"entropy"`
	Choice    Entropy   `json:"choice"`
}

// EntropyCalculationState - энтропия с подробным описанием. Видна только на странице /entropy на сайте.
// Тут еще добавляются кадр с камеры и шум для наглядности, так что событие это объёмное.
type EntropyCalculationState struct {
	Timestamp  time.Time `json:"timestamp"`
	ImageFrame string    `json:"imageFrame"` // base64-encoded (нигде не хранится)
	ImageNoise string    `json:"imageNoise"` // base64-encoded (нигде не хранится)
	Entropy    Entropy   `json:"entropy"`
	Choice     Entropy   `json:"choice"`
}

// Cargo - упаковка для событий при передаче между Асгардом и Альфхеймом (между бекэнд-серверами)
// Odin: этот груз везёт хрустальный драккар по волнам радужного моста из Асгарда в Альфхейм. Красиво?
// Odin: Жаль, что выглядит как набор каких-то символов на одном из человеческих языков
type Cargo struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}

// Radiogram - упаковка для событий при передаче между Альфхеймом и Мидгардом (между api-gateway и клиентом)
// Odin: "Слушайте наше радио Artchitect-FM, настраивайте ваши радиостанции на канал между мирами!"
// Odin: "Радиостанция Artchitect-FM - свет маяка в тумане для заблудших душ."
type Radiogram struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}
