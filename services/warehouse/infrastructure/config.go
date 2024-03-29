package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	IsDebug    bool
	HttpPort   string
	ArtsPath   string
	UnityPath  string
	OriginPath string
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
		IsDebug:    env["IS_DEBUG"] == "true",
		HttpPort:   env["HTTP_PORT"],
		ArtsPath:   env["ARTS_PATH"],
		UnityPath:  env["UNITY_PATH"],
		OriginPath: env["ORIGIN_PATH"],
	}
}
