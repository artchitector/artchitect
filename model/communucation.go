package model

const (
	ChanTech            = "tech" // технический канал для ненужных уведомлений
	ChanEntropy         = "entropy"
	ChanEntropyExtended = "entropy_extended"
	ChanNewArt          = "new_art"
)

// Здесь расположены модели, необходимые для передачи рантаймовых событий из Асгарда по остальным мирам

// Cargo - упаковка для событий при передаче между Асгардом и Альфхеймом (между бекэнд-серверами по Редису)
// Odin: этот груз везёт хрустальный драккар по волнам радужного моста из Асгарда в Альфхейм. Красиво?
type Cargo struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}

// Radiogram - упаковка для событий при передаче между Альфхеймом и Мидгардом (между api-gateway и клиентом, websocket)
// Odin: "Слушайте наше радио Artchitect-FM, настраивайте ваши радиостанции на канал между мирами!"
// Odin: "Радиостанция Artchitect-FM - свет маяка в тумане для заблудших душ."
type Radiogram struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}
