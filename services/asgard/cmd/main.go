package main

import (
	"context"
	"github.com/artchitector/artchitect2/libraries/warehouse"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/asgard/communication"
	"github.com/artchitector/artchitect2/services/asgard/infrastructure"
	"github.com/artchitector/artchitect2/services/asgard/pantheon"
	frigg2 "github.com/artchitector/artchitect2/services/asgard/pantheon/frigg"
	"github.com/artchitector/artchitect2/services/asgard/utils"
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
	//warehouse := communication.NewWarehouse(config.WarehouseOriginUrl, config.WarehouseArtUrls)
	wh := warehouse.NewWarehouse(config.WarehouseArtUrls, config.WarehouseOriginUrl)
	keyhole := communication.NewKeyhole(config.HttpPort, webcam)

	// хранилища сущностей
	artPile := model.NewArtPile(database)
	unityPile := model.NewUnityPile(database)

	giving := communication.NewGiving(artPile, muninn, bifröst)

	heimdall := pantheon.NewHeimdallr(huginn, bifröst)
	ai := infrastructure.NewAI(config.UseFakeAI, config.InvokeAIPath)
	freyja := pantheon.NewFreyja(ai)
	gungner := pantheon.NewGungner()

	friggCollage := frigg2.NewFriggCollage(wh, gungner)
	frigg := pantheon.NewFrigg(friggCollage, muninn, unityPile, artPile)

	odin := pantheon.NewOdin(
		config.CreatorActive,
		config.CreateTotalTimeSec,
		frigg,
		freyja,
		muninn,
		gungner,
		heimdall,
		artPile,
		wh,
	)

	// запуск фоновых служб
	go runServices(ctx, lostEye, huginn, heimdall, webcam, keyhole, giving, odin, bifröst)

	// запуск временных сервисов, которые обрабатывают задачи. Artchitect будет запущен лишь после их выполнения
	runTemporary(ctx, artPile, frigg)

	// запуск Главного Цикла Архитектора (ГЦА)
	runArtchitect(ctx, odin, frigg)
}

// runArtchitect - запуск Главного Цикла Архитектора (ГЦА)
func runArtchitect(
	ctx context.Context,
	odin *pantheon.Odin,
	frigg *pantheon.Frigg,
) {
	artchitect := pantheon.NewIntention(odin, frigg)
	artchitect.Run(ctx)
}

// runServices - запуск фоновых служб
func runServices(
	ctx context.Context,
	eye *pantheon.LostEye,
	huginn *pantheon.Huginn,
	heimdall *pantheon.Heimdallr,
	webcam *infrastructure.Webcam,
	keyhole *communication.Keyhole,
	giving *communication.Giving,
	odin *pantheon.Odin,
	bifröst *communication.Bifröst,
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
	go func() {
		if err := keyhole.StartHttpServer(ctx); err != nil {
			log.Fatal().Err(err).Msgf("[keyhole] ПАДЕНИЕ KEYHOLE. СТОП-МАШИНА, АСГАРД!")
		}
	}()
	go func() {
		giving.StartGiving(ctx)
		log.Debug().Msgf("[main] GIVING ОСТАНОВЛЕН")
	}()
	go func() {
		bifröst.ListenPrivateOdinRequests(ctx, odin)
		log.Debug().Msgf("[main] ОСТАНОВЛЕНО ЧТЕНИЕ ЛИЧНЫХ ПРОШЕНИЙ К ОДИНУ")
	}()
}

func runTemporary(ctx context.Context, artPile *model.ArtPile, frigg *pantheon.Frigg) {
	ui := utils.NewUnityInitializer(artPile, frigg)
	if err := ui.Init(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
}
