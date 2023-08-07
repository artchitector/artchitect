package model

import (
	"fmt"
	"image"
	"time"
)

type Entropy struct {
	IntValue     uint64      `json:"int"`
	FloatValue   float64     `json:"float"`
	Image        image.Image `json:"-"`
	ImageEncoded string      `json:"imageEncoded"` // перед отправкой в Биврёст эту картинку надо base64-энкодить
	ImageID      string      `json:"imageId"`      // ключ для получения изображения энтропии с memory-сервера
}

func (e Entropy) String() string {
	return fmt.Sprintf("E:%.6f", e.FloatValue)
}

type EntropyPack struct {
	Timestamp time.Time `json:"timestamp"`
	Entropy   Entropy   `json:"entropy"`
	Choice    Entropy   `json:"choice"`
}
