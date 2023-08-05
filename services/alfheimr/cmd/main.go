package main

import (
	"context"
	"github.com/artchitector/artchitect2/services/alfheimr/external"
	"github.com/artchitector/artchitect2/services/alfheimr/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	_ = infrastructure.InitDB(ctx, config.DbDSN)
	// СЛУШАТЕЛЬ РЕДИСА И СОБЫТИЙ. ФОНОВЫЙ ЗАПУСК ОБРАБОТКИ
	red := infrastructure.InitRedis(config)
	lis := external.NewListener(red)
	go func() {
		if err := lis.Run(ctx); err != nil {
			log.Fatal().Msgf("[ВРАТА] СЛУШАТЕЛЬ - АВАРИЯ")
		}
	}()
	// ТЕСТОВЫЙ ЗАПУСК СЛУШАТЕЛЯ КАНАЛА
	go func() {
		time.Sleep(time.Second)

		tCtx, tDone := context.WithTimeout(ctx, time.Second*5)
		defer tDone()

		ch := lis.Subscribe(tCtx)
		for e := range ch {
			log.Info().Msgf("[ТЕСТ] ПОЛУЧЕНО СООБЩЕНИЕ: %+v", e)
		}
	}()

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

		// http-ручки
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
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
