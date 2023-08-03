package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// базовые параметры
	IsDebug      bool
	InvokeAIPath string
	DbDSN        string

	// службы
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
		// базовые параметры
		IsDebug:      env["IS_DEBUG"] == "true",
		InvokeAIPath: env["INVOKEAI_PATH"],
		DbDSN:        env["DB_DSN"],

		// службы
		CreatorActive: env["CREATOR_ACTIVE"] == "true",
	}
}
