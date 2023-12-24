package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/artchitector/artchitect/services/metricsbridge/infrastructure"
	"github.com/artchitector/artchitect/services/metricsbridge/service"
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

	redis := infrastructure.InitRedis(config)
	mb := service.NewMetricsbridge(redis, config.ScrapeURL, config.CurrentHostName)

	// ЗАПУСК HTTP-СЕРВЕРА
	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}

		r.GET("/:hostname/metrics", mb.GetMetricsHandler)

		// запуск http-сервера
		log.Info().Msgf("[warehouse] HTTP-ВРАТА ВКЛ. ПОРТ:%s", config.HttpPort)
		if err := r.Run("0.0.0.0:" + config.HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	go func() {
		if !config.EnableTransfer {
			log.Info().Msgf("[main] metrics transfer disabled")
			return
		}
		err := mb.RunMetricsTransfer(ctx)
		if err != nil {
			log.Warn().Err(err).Msgf("[main] metrics transfer stop")
		} else {
			log.Info().Msgf("[main] metrics transfer context.stop")
		}
	}()

	<-ctx.Done()
	log.Info().Msgf("[warehouse] МОСТ МЕТРИК ОТКЛЮЧАЕТСЯ")
}
