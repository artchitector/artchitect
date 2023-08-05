package model

import (
	"fmt"
	"image"
	"time"
)

type Entropy struct {
	IntValue   uint64      `json:"int"`
	FloatValue float64     `json:"float"`
	Image      image.Image `json:"-"`
	ImageID    string      `json:"imageId"` // ключ для получения изображения энтропии с memory-сервера
}

func (e Entropy) String() string {
	return fmt.Sprintf("E:%.6f", e.FloatValue)
}

type EntropyPack struct {
	Timestamp time.Time `json:"timestamp"`
	Entropy   Entropy   `json:"entropy"`
	Choice    Entropy   `json:"choice"`
}
