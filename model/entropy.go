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
	Type       string
	IntValue   uint64
	FloatValue float64
	Image      image.Image
}

func (e Entropy) String() string {
	return fmt.Sprintf("E:%.6f", e.FloatValue)
}
