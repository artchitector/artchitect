package infrastructure

import (
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// базовые параметры
	IsDebug            bool
	HttpPort           string
	ArtWarehouseURL    string
	OriginWarehouseURL string

	// подключения к внешним службам
	DbDSN         string
	RedisAddr     string
	RedisPassword string

	// auth
	JwtSecret      string
	AllowFakeAuth  bool
	ArtchitectHost string
	ArtchitectorID uint

	// telegram
	BotToken                string
	ChatArtchitectChoice    int64
	ChatArtchitectorChoice  int64
	ArtchitectChoiceEnabled bool
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

	artchitectorID, err := strconv.ParseUint(env["ARTCHITECTOR_ID"], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	chatArtchitectChoice, err := strconv.ParseInt(env["CHAT_ARTCHITECT_CHOICE"], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msgf("[config] ОШИБКА ЧТЕНИЯ CHAT_ARTCHITECT_CHOICE")
	}
	chatArtchitectorChoice, err := strconv.ParseInt(env["CHAT_ARTCHITECTOR_CHOICE"], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msgf("[config] ОШИБКА ЧТЕНИЯ CHAT_ARTCHITECTOR_CHOICE")
	}

	return &Config{
		// базовые параметры
		IsDebug:            env["IS_DEBUG"] == "true",
		HttpPort:           env["HTTP_PORT"],
		ArtWarehouseURL:    env["ART_WAREHOUSE_URL"],
		OriginWarehouseURL: env["ORIGIN_WAREHOUSE_URL"],

		// подключения к внешним службам
		DbDSN:         env["DB_DSN"],
		RedisAddr:     env["REDIS_ADDR"],
		RedisPassword: env["REDIS_PASSWORD"],

		// auth
		JwtSecret:      env["JWT_SECRET"],
		AllowFakeAuth:  env["ALLOW_FAKE_AUTH"] == "true",
		ArtchitectHost: env["ARTCHITECT_HOST"],
		ArtchitectorID: uint(artchitectorID),

		BotToken:               env["BOT_TOKEN"],
		ChatArtchitectChoice:   chatArtchitectChoice,
		ChatArtchitectorChoice: chatArtchitectorChoice,
	}
}
