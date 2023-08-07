package main

import (
	"context"
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
	webcamStream := make(chan image.Image)
	lostEye := pantheon.NewLostEye()
	go func() {
		// Глаз начинает смотреть в ткань мироздания
		lostEye.StartEntropyDecode(ctx, webcamStream)
	}()

	huginn := pantheon.NewHuginn(lostEye)
	go func() {
		// Хугин смотрит в глаз и осмысляет поступающую энтропию
		huginn.StartEntropyRealize(ctx)
	}()

	// redis-bifröst
	bifröst := communication.NewBifröst(red)
	heimdall := pantheon.NewHeimdallr(huginn, bifröst)
	go func() {
		heimdall.StartStream(ctx)
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
	//artRepo := model.NewArtRepository(database, nil)
	//
	//ai := infrastructure.NewAI(config.InvokeAIPath)
	//freyja := pantheon.NewFreyja(ai)
	//creator := pantheon.NewOdin(config.CreatorActive, freyja, artRepo)
	//artchitect := pantheon.NewArtchitect(creator)
	//artchitect.Run(ctx)
}

// runServices - запуск фоновых служб
func runServices(ctx context.Context) {

}
