package warehouse

// Warehouse - склад бинарных картинок
// Odin: Картинки хранятся не в Асгарде, а на серверах файловых-хранилищах.
// Odin: Warehouse инкапсулирует сложную логику сохранения на файловые серверы, скрывая её от слоя Асгарда, и от Меня pantheon.Odin.
// Loki: А SOLID можешь расшифровать?) Ты уже прокачался в программировании, как я посмотрю.
// Loki: Уже боюсь проиграть наше пари... (UPD: уже проиграл)
type Warehouse struct {
	artWarehouseURL    string // Odin: На этом сервере хранятся операционные картинки в обычных разрешениях для сайта
	originWarehouseURL string // Odin: На этом сервере хранятся большие JPEG-файлы (оригиналы) для печати
}

func NewWarehouse(artWarehouseURL string, originWarehouseURL string) *Warehouse {
	return &Warehouse{artWarehouseURL: artWarehouseURL, originWarehouseURL: originWarehouseURL}
}
