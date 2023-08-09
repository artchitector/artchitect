package main

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/asgard/communication"
	"github.com/artchitector/artchitect2/services/asgard/infrastructure"
	"github.com/artchitector/artchitect2/services/asgard/pantheon"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// внешние связи
	bifröst := communication.NewBifröst(red)
	warehouse := communication.NewWarehouse(config.WarehouseFullsizeUrl, config.WarehouseArtUrls)

	// хранилища сущностей
	artPile := model.NewArtPile(database)

	heimdall := pantheon.NewHeimdallr(huginn, bifröst)
	ai := infrastructure.NewAI(config.InvokeAIPath)
	freyja := pantheon.NewFreyja(ai)
	odin := pantheon.NewOdin(config.CreatorActive, freyja, muninn, heimdall, artPile, warehouse)

	// запуск фоновых служб
	go runServices(ctx, lostEye, huginn, heimdall, webcam)
	// запуск Главного Цикла Архитектора (ГЦА)
	runArtchitect(ctx, odin)
}

// runArtchitect - запуск Главного Цикла Архитектора (ГЦА)
func runArtchitect(
	ctx context.Context,
	odin *pantheon.Odin,
) {
	artchitect := pantheon.NewArtchitect(odin)
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
