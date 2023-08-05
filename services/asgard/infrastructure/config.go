package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// базовые параметры
	IsDebug      bool
	InvokeAIPath string

	// подключения к внешним службам
	DbDSN         string
	RedisAddr     string
	RedisPassword string

	// hardware
	WebcamDeviceID        string
	WebcamFrameResolution string

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

		// подключения к внешним службам
		DbDSN:         env["DB_DSN"],
		RedisAddr:     env["REDIS_ADDR"],
		RedisPassword: env["REDIS_PASSWORD"],

		// hardware
		WebcamDeviceID:        env["WEBCAM_DEVICE_ID"],
		WebcamFrameResolution: env["WEBCAM_FRAME_RESOLUTION"],

		// службы
		CreatorActive: env["CREATOR_ACTIVE"] == "true",
	}
}
