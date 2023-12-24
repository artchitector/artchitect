package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	IsDebug  bool
	HttpPort string

	RedisAddr     string
	RedisPassword string
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
		IsDebug:       env["IS_DEBUG"] == "true",
		HttpPort:      env["HTTP_PORT"],
		RedisAddr:     env["REDIS_ADDR"],
		RedisPassword: env["REDIS_PASSWORD"],
	}
}
