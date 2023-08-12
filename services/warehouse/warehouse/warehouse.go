package warehouse

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/libraries/resizer"
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
)

type artRequest struct {
	Id   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"`
}

type Warehouse struct {
	ArtsPath   string
	UnityPath  string
	OriginPath string
}

func (w *Warehouse) HandleGetArt(c *gin.Context) {
	var request artRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	unityMask := model.Art{ID: request.Id}.GetUnityMask(model.Unity1K)
	var dirpath, filename, contentType string
	if request.Size == model.SizeOrigin {
		dirpath = path.Join(w.OriginPath, fmt.Sprintf("U%s", unityMask))
		filename = fmt.Sprintf("art-%d.jpg", request.Id)
		contentType = "image/jpeg"
	} else {
		dirpath = path.Join(w.ArtsPath, fmt.Sprintf("U%s", unityMask))
		filename = fmt.Sprintf("art-%d-%s.jpg", request.Id, request.Size)
		contentType = "image/jpeg"
	}
	data, err := w.readFile(dirpath, filename)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.Data(http.StatusOK, contentType, data)
}

func (w *Warehouse) HandleUploadArt(c *gin.Context) {
	// тут приходит jpeg файл в размере F
	data, artID, err := w.parse(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadArt] БИТЫЕ ДАННЫЕ")
		return
	}

	log.Info().Msgf("[warehouse:HandleUploadArt] ВХОДЯЩИЙ ЗАПРОС ART_ID=%d", artID)

	b := bytes.NewReader(data)
	decoded, err := jpeg.Decode(b)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadArt] БИТЫЙ JPEG")
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
			log.Error().Err(err).Msgf("[warehouse:HandleUploadArt] НЕ СМОГ СОЗДАТЬ JPEG")
			return
		}

		filename := fmt.Sprintf("art-%d-%s.jpg", artID, size)
		if err := w.saveFile(c, w.ArtsPath, artID, filename, b.Bytes()); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			log.Error().Err(err).Msgf("[warehouse:HandleUploadArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPEG-ФАЙЛ")
			return
		}
	}
}

func (w *Warehouse) HandleUploadOrigin(c *gin.Context) {
	// а тут приходит JPEG-файл в полном размере
	data, artID, err := w.parse(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadArt] БИТЫЕ ДАННЫЕ")
		return
	}

	log.Info().Msgf("[warehouse:HandleUploadArt] ВХОДЯЩИЙ ЗАПРОС ART_ID=%d", artID)

	b := bytes.NewReader(data)
	_, err = jpeg.Decode(b) // просто проверим, что изображение не битое
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadArt] БИТЫЙ JPEG")
		return
	}

	filename := fmt.Sprintf("art-%d.jpg", artID)
	if err := w.saveFile(c, w.OriginPath, artID, filename, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Error().Err(err).Msgf("[warehouse:HandleUploadArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPG-ФАЙЛ")
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

func (w *Warehouse) readFile(dir string, filename string) ([]byte, error) {
	b, err := os.ReadFile(path.Join(dir, filename))
	return b, err
}
