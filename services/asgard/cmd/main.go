package main

import (
	bytes2 "bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/asgard/communication"
	"github.com/artchitector/artchitect2/services/asgard/infrastructure"
	"github.com/artchitector/artchitect2/services/asgard/pantheon"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"image"
	"image/png"
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

	// redis-stream
	stream := communication.NewStream(red)

	// TODO ЭТО НАДО УБРАТЬ ИЗ MAIN
	go func() {
		sCtx, done := context.WithTimeout(ctx, time.Second*10000)
		defer done()
		ch := huginn.Subscribe(sCtx)

		for ent := range ch {
			log.Info().Msgf("[main] ПОЛУЧЕНА ЭНТРОПИИ ПО КАНАЛУ %+v", ent)

			if ent.Entropy.Image != nil {
				b := bytes2.Buffer{}
				if err := png.Encode(&b, ent.Entropy.Image); err != nil {
					log.Error().Err(err).Msgf("[main] PNG НЕ УДАЛСЯ. ЖАЛЬ")
				} else {
					bbb := b.Bytes()
					log.Info().Msgf("SIZE: %db", len(bbb))
					ent.Entropy.ImageEncoded = base64.StdEncoding.EncodeToString(bbb)
				}
			}
			if ent.Choice.Image != nil {
				b := bytes2.Buffer{}
				if err := png.Encode(&b, ent.Choice.Image); err != nil {
					log.Error().Err(err).Msgf("[main] PNG НЕ УДАЛСЯ. ЖАЛЬ")
				} else {
					bbb := b.Bytes()
					log.Info().Msgf("SIZE: %db", len(bbb))
					ent.Choice.ImageEncoded = base64.StdEncoding.EncodeToString(bbb)
				}
			}

			if b, err := json.Marshal(ent); err != nil {
				log.Fatal().Msgf("JSON MARSHAL")
			} else {
				err = stream.SendDrakkar(ctx, model.Cargo{
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
