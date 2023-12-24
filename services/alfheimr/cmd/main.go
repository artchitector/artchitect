package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/artchitector/artchitect/libraries/warehouse"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/services/alfheimr/communication"
	"github.com/artchitector/artchitect/services/alfheimr/infrastructure"
	"github.com/artchitector/artchitect/services/alfheimr/portals"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	unityPile := model.NewUnityPile(db)
	likePile := model.NewLikePile(db)

	// WAREHOUSE
	wh := warehouse.NewWarehouse(config.ArtWarehouseURL, config.OriginWarehouseURL)

	// AUTH
	authService := infrastructure.NewAuthService(config.AllowFakeAuth, []byte(config.JwtSecret))

	// СБОРКА ПОРТАЛОВ (ХЕНДЛЕРОВ)
	radioPortal := portals.NewRadioPortal(harbour)
	artPortal := portals.NewArtPortal(artPile, harbour)
	imPortal := portals.NewImagePortal(wh)
	unityPortal := portals.NewUnityPortal(unityPile, artPile)
	likePortal := portals.NewLikePortal(authService, likePile, harbour, config.ArtchitectorID)
	authPortal := portals.NewAuthPortal(
		authService,
		config.JwtSecret,
		config.ArtchitectHost,
		config.BotToken,
	)

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

		r.GET("/art/chosen", artPortal.HandleChosen)
		r.GET("/art/max", artPortal.HandleMaxArt)
		r.GET("/art/:id", artPortal.HandleArt)
		r.GET("/art/:id/flat", artPortal.HandleArtFlat)
		r.GET("/arts/last/:last", artPortal.HandleLast)
		r.GET("/image/:name", imPortal.HandleArtImage)
		r.GET("/uimage/:mask/:version/:size", imPortal.HandleUnityImage)
		// http://localhost/api/ uimage 0XXXXX 0 f
		r.GET("/unity", unityPortal.HandleMain)
		r.GET("/unity/:mask", unityPortal.HandleUnity)
		r.GET("/me", authPortal.HandleMe)
		r.GET("/login", authPortal.HandleLogin)
		r.GET("/liked", likePortal.HandleLikedList)
		r.GET("/liked/:art_id", likePortal.HandleLikedArt)
		r.POST("/like", likePortal.HandleLike)

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
