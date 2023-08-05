package model

import "time"

const (
	ChanEntropy            = "entropy"
	ChanEntropyCalculation = "entropy_calculation"
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

// Event - упаковка для событий для передачи по каналам (внутри сервиса gate) и отправки клиентам
type Event struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}
