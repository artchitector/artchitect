package main

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/asgard/communication"
	"github.com/artchitector/artchitect2/services/asgard/infrastructure"
	"github.com/artchitector/artchitect2/services/asgard/pantheon"
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
	red := infrastructure.InitRedis(config)

	// веб-камера -> глаз Одина (pantheon.LostEye) -> Хугин (pantheon.Huginn)
	webcam := infrastructure.NewWebcam(config.WebcamDeviceID, config.WebcamFrameResolution)
	lostEye := pantheon.NewLostEye()
	huginn := pantheon.NewHuginn(lostEye)
	muninn := pantheon.NewMuninn(huginn)

	// redis-bifröst
	bifröst := communication.NewBifröst(red)
	heimdall := pantheon.NewHeimdallr(huginn, bifröst)

	// запуск фоновых служб
	go runServices(ctx, lostEye, huginn, heimdall, webcam)
	// запуск Главного Цикла Архитектора (ГЦА)
	runArtchitect(ctx, config, database, muninn)
}

// runArtchitect - запуск Главного Цикла Архитектора (ГЦА)
func runArtchitect(
	ctx context.Context,
	config *infrastructure.Config,
	database *gorm.DB,
	muninn *pantheon.Muninn,
) {
	artPile := model.NewArtPile(database)
	warehouse := communication.NewWarehouse(config.WarehouseFullsizeUrl, config.WarehouseArtUrls)

	ai := infrastructure.NewAI(config.InvokeAIPath)
	freyja := pantheon.NewFreyja(ai)
	creator := pantheon.NewOdin(config.CreatorActive, freyja, muninn, artPile, warehouse)
	artchitect := pantheon.NewArtchitect(creator)
	artchitect.Run(ctx)
}

// runServices - запуск фоновых служб
func runServices(
	ctx context.Context,
	eye *pantheon.LostEye,
	huginn *pantheon.Huginn,
	heimdall *pantheon.Heimdallr,
	webcam *infrastructure.Webcam,
) {
	webcamStream := make(chan image.Image)
	go func() {
		// Глаз начинает смотреть в ткань мироздания
		eye.StartEntropyDecode(ctx, webcamStream)
	}()
	go func() {
		// Хугин смотрит в глаз и осмысляет поступающую энтропию
		huginn.StartEntropyRealize(ctx)
	}()
	go func() {
		// pantheon.Heimdallr запускает ретрансляцию энтропии из pantheon.LostEye в communication.Bifröst
		heimdall.StartStream(ctx)
	}()
	go func() {
		// Запуск чтения кадров с веб-камеры
		if err := webcam.Start(ctx, webcamStream); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()
}
