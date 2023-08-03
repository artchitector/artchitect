package main

import (
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	for range time.Tick(time.Second) {
		log.Info().Msgf("[tick] hello world!")
	}
}
