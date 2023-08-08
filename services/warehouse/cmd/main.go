package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/libraries/resizer"
	"github.com/artchitector/artchitect2/model"
	"github.com/artchitector/artchitect2/services/warehouse/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
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

	warehouse := Warehouse{
		config.ArtsPath,
		config.UnityPath,
		config.FullsizePath,
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

		r.POST("/upload_art", warehouse.HandleArt)
		r.POST("/upload_fullsize", warehouse.HandleFullsize)

		// запуск http-сервера
		log.Info().Msgf("[warehouse] HTTP-ВРАТА ВКЛ. ПОРТ:%s", config.HttpPort)
		if err := r.Run("0.0.0.0:" + config.HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msgf("[warehouse] СКЛАД ОТКЛЮЧАЕТСЯ")
}

type Warehouse struct {
	artsPath     string
	unityPath    string
	fullsizePath string
}

func (w *Warehouse) HandleArt(c *gin.Context) {
	// тут приходит jpeg файл в размере F
	data, artID, err := w.parse(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleArt] БИТЫЕ ДАННЫЕ")
		return
	}

	log.Info().Msgf("[warehouse:HandleArt] ВХОДЯЩИЙ ЗАПРОС ART_ID=%d", artID)

	b := bytes.NewReader(data)
	decoded, err := jpeg.Decode(b)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleArt] БИТЫЙ JPEG")
		return
	}

	for _, size := range []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS} {
		img := resizer.ResizeImage(decoded, size)
		var quality int
		switch size {
		case model.SizeF:
			quality = model.QualityF
		case model.SizeM:
			quality = model.QualityM
		case model.SizeS:
			quality = model.QualityS
		case model.SizeXS:
			quality = model.QualityXS
		}
		b := new(bytes.Buffer)
		err := jpeg.Encode(b, img, &jpeg.Options{Quality: quality})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			log.Error().Err(err).Msgf("[warehouse:HandleArt] НЕ СМОГ СОЗДАТЬ JPEG")
			return
		}

		filename := fmt.Sprintf("art-%d-%s.jpg", artID, size)
		if err := w.saveFile(c, w.artsPath, artID, filename, b.Bytes()); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			log.Error().Err(err).Msgf("[warehouse:HandleArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPEG-ФАЙЛ")
			return
		}
	}
}

func (w *Warehouse) HandleFullsize(c *gin.Context) {
	// а тут приходит JPEG-файл в полном размере
	data, artID, err := w.parse(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleArt] БИТЫЕ ДАННЫЕ")
		return
	}

	log.Info().Msgf("[warehouse:HandleArt] ВХОДЯЩИЙ ЗАПРОС ART_ID=%d", artID)

	b := bytes.NewReader(data)
	_, err = jpeg.Decode(b) // просто проверим, что изображение не битое
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleArt] БИТЫЙ JPEG")
		return
	}

	filename := fmt.Sprintf("art-%d.jpg", artID)
	if err := w.saveFile(c, w.fullsizePath, artID, filename, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Error().Err(err).Msgf("[warehouse:HandleArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPG-ФАЙЛ")
		return
	}
}

func (w *Warehouse) parse(c *gin.Context) (data []byte, artID int, err error) {
	var file *multipart.FileHeader
	file, err = c.FormFile("file")
	if err != nil {
		return
	}

	artIDStr := c.PostForm("art_id")
	if artIDStr == "" {
		err = errors.Errorf("[warehouse] ART_ID ОБЯЗАТЕЛЕН")
		return
	}
	artID, err = strconv.Atoi(artIDStr)
	if err != nil {
		return
	}

	var f multipart.File
	f, err = file.Open()
	if err != nil {
		return
	}

	data, err = io.ReadAll(f)
	return
}

func (w *Warehouse) saveFile(
	ctx context.Context,
	dir string,
	artID int,
	filename string,
	data []byte,
) error {
	unityMask := model.Art{ID: uint(artID)}.GetUnityMask(model.Unity1K)

	dirpath := path.Join(dir, fmt.Sprintf("U%s", unityMask))
	if err := os.MkdirAll(dirpath, 0774); err != nil {
		return err
	}
	log.Info().Msgf("DIR CREATED %s", dirpath)

	filepath := path.Join(dirpath, filename)
	err := os.WriteFile(filepath, data, 0774)
	return err
}
