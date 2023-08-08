package communication

import (
	"context"
	"image"
)

// Warehouse - склад бинарных картинок
// Odin: Картинки хранятся не в Асгарде, а на серверах файловых-хранилищах.
// Odin: Warehouse инкапсулирует сложную логику сохранения на файловые серверы, скрывая её от слоя Асгарда, и от Меня pantheon.Odin.
// Loki: А SOLID можешь расшифровать?) Ты уже прокачался в программировании, как я посмотрю.
// Loki: Уже боюсь проиграть наше пари...
type Warehouse struct {
}

func (wh *Warehouse) SaveImage(ctx context.Context, artID uint, img image.Image) error {
	/*
		Odin:
		Ну чтож, дети мои, слушайте мой рассказ.
		Архитектор рисует большие
	*/
	return nil
}
