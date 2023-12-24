package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/artchitector/artchitect/services/warehouse/infrastructure"
	warehouse2 "github.com/artchitector/artchitect/services/warehouse/warehouse"
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

	warehouse := warehouse2.Warehouse{
		config.ArtsPath,
		config.UnityPath,
		config.OriginPath,
	}

	// ЗАПУСК HTTP-СЕРВЕРА
	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}

		r.GET("/art/:id/:size", warehouse.HandleGetArt)
		r.GET("/unity/:mask/:version/:size", warehouse.HandleGetUnity)

		r.POST("/upload/art", warehouse.HandleUploadArt)
		r.POST("/upload/unity", warehouse.HandleUploadUnity)
		r.POST("/upload/origin", warehouse.HandleUploadOrigin)

		// запуск http-сервера
		log.Info().Msgf("[warehouse] HTTP-ВРАТА ВКЛ. ПОРТ:%s", config.HttpPort)
		if err := r.Run("0.0.0.0:" + config.HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msgf("[warehouse] СКЛАД ОТКЛЮЧАЕТСЯ")
}
