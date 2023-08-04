package main

import (
	"context"
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

	// веб-камера + энтропия
	entropy := external.NewEntropy()
	webcam := infrastructure.NewWebcam(config.WebcamDeviceID, config.WebcamFrameResolution)
	webcamStream := make(chan image.Image)
	go func() {
		entropy.StartEntropyDecode(ctx, webcamStream)
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
