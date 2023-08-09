package main

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/alfheimr/communication"
	"github.com/artchitector/artchitect2/services/alfheimr/infrastructure"
	"github.com/artchitector/artchitect2/services/alfheimr/portals"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// подготовка служебных golang-вещей
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// загрузка конфигов из env
	config := infrastructure.InitEnv()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	if config.IsDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	// СБОРКА ЗАВИСИМОСТЕЙ
	db := infrastructure.InitDB(ctx, config.DbDSN)
	// СЛУШАТЕЛЬ РЕДИСА И СОБЫТИЙ. ФОНОВЫЙ ЗАПУСК ОБРАБОТКИ
	red := infrastructure.InitRedis(config)

	harbour := communication.NewHarbour(red)
	go func() {
		if err := harbour.Run(ctx); err != nil {
			log.Fatal().Msgf("[ГЛАВНЫЙ] ГАВАНЬ ПРЕКРАТИЛА РАБОТУ. ODIN НЕДОВОЛЕН")
		}
	}()

	// КУЧИ ДАННЫХ
	artPile := model.NewArtPile(db)

	// СБОРКА ПОРТАЛОВ (ХЕНДЛЕРОВ)
	radioPortal := portals.NewRadioPortal(harbour)
	artPortal := portals.NewArtPortal(artPile)

	// ЗАПУСК HTTP-СЕРВЕРА
	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost",
				"https://artchitect.space",
				"https://ru.artchitect.space",
				"https://eu.artchitect.space",
			},
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}

		// Odin: это порталы, через которые Мидгард может связаться с Альфхеймом
		// Odin: есть порталы, которые закрываются сразу (http-api), а есть длительно действующий (websocket)
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		r.GET("/art/:id", artPortal.HandleArt)
		r.GET("/arts/last/:last", artPortal.HandleLast)

		// connection - Портал с постоянной связью c Мидгардом (вебсокете)
		r.GET("/radio", func(c *gin.Context) {
			radioPortal.Handle(c.Writer, c.Request)
		})

		// запуск http-сервера
		log.Info().Msgf("[gate] HTTP-ВРАТА ВКЛ. ПОРТ:%s", config.HttpPort)
		if err := r.Run("0.0.0.0:" + config.HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("[gate] ВРАТА ВЫКЛ")
}
