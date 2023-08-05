package model

import (
	"fmt"
	"image"
)

const (
	EntropyTypeDirect = "entropy"
	EntropyTypeChoice = "inverted"
)

type Entropy struct {
	Type       string      `json:"type"`
	IntValue   uint64      `json:"int"`
	FloatValue float64     `json:"float"`
	Image      image.Image `json:"-"`
	ImageID    string      `json:"imageId"` // ключ для получения изображения энтропии с memory-сервера
}

func (e Entropy) String() string {
	return fmt.Sprintf("E:%.6f", e.FloatValue)
}
