package main

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/soul/external"
	"github.com/artchitector/artchitect2/services/soul/infrastructure"
	"github.com/artchitector/artchitect2/services/soul/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"image"
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

	// общие зависимости
	database := infrastructure.InitDB(ctx, config.DbDSN)
	red := infrastructure.InitRedis(config)

	// веб-камера + энтропия
	entropy := external.NewEntropy()
	webcam := infrastructure.NewWebcam(config.WebcamDeviceID, config.WebcamFrameResolution)
	webcamStream := make(chan image.Image)
	go func() {
		entropy.StartEntropyDecode(ctx, webcamStream)
	}()

	// redis-stream
	stream := external.NewStream(red)

	go func() {
		sCtx, done := context.WithTimeout(ctx, time.Second*10)
		defer done()
		ch := entropy.SubscribeEntropy(sCtx)

		for ent := range ch {
			log.Info().Msgf("[main] ПОЛУЧЕНА ЭНТРОПИИ ПО КАНАЛУ %+v", ent)
			if b, err := json.Marshal(ent); err != nil {
				log.Fatal().Msgf("JSON MARSHAL")
			} else {
				err = stream.SendCargo(ctx, model.Event{
					Channel: model.ChanEntropy,
					Payload: string(b),
				})
				if err != nil {
					log.Fatal().Msgf("SEND CARGO FAILED")
				}
			}
		}

		log.Fatal().Msgf("STOP ASGARD")
	}()

	if err := webcam.Start(ctx, webcamStream); err != nil {
		log.Fatal().Err(err).Send()
	}

	// запуск фоновых служб
	go runServices(ctx)
	// запуск Главного Цикла Архитектора (ГЦА)
	runArtchitect(ctx, config, database)
}

// runArtchitect - запуск Главного Цикла Архитектора (ГЦА)
func runArtchitect(
	ctx context.Context,
	config *infrastructure.Config,
	database *gorm.DB,
) {
	artRepo := model.NewArtRepository(database, nil)

	ai := infrastructure.NewAI(config.InvokeAIPath)
	artist := external.NewArtist(ai)
	creator := internal.NewCreator(config.CreatorActive, artist, artRepo)
	artchitect := internal.NewArtchitect(creator)
	artchitect.Run(ctx)
}

// runServices - запуск фоновых служб
func runServices(ctx context.Context) {

}
