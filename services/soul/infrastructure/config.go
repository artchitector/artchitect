package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// main
	IsDebug bool

	// services
	CreatorActive bool
}

func InitEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	env, err := godotenv.Read()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Config{
		// main
		IsDebug: env["IS_DEBUG"] == "true",

		// services
		CreatorActive: env["CREATOR_ACTIVE"] == "true",
	}
}
