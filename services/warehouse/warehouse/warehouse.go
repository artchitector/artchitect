package warehouse

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/artchitector/artchitect/libraries/resizer"
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type artRequest struct {
	Id   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"`
}

type unityRequest struct {
	Mask    string `uri:"mask" binding:"required"`
	Version uint   `uri:"version" binding:"numeric"`
	Size    string `uri:"size" binding:"required"`
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
	var dirpath, filename string
	if request.Size == model.SizeOrigin {
		dirpath = path.Join(w.OriginPath, fmt.Sprintf("U%s", unityMask))
		filename = fmt.Sprintf("art-%d.jpg", request.Id)
	} else {
		dirpath = path.Join(w.ArtsPath, fmt.Sprintf("U%s", unityMask))
		filename = fmt.Sprintf("art-%d-%s.jpg", request.Id, request.Size)
	}
	data, err := w.readFile(dirpath, filename)
	if err != nil {
		log.Error().Err(err).Msgf("[warehouse] ОШИБКА ЧТЕНИЯ ФАЙЛА %s/%s", dirpath, filename)
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.Data(http.StatusOK, "image/jpeg", data)
}

func (w *Warehouse) HandleGetUnity(c *gin.Context) {
	var request unityRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	parentMask := fmt.Sprintf("U%sXXXXX", string(request.Mask[0]))
	dirpath := path.Join(w.UnityPath, parentMask)
	filename := fmt.Sprintf("U%s-%d-%s.jpg", request.Mask, request.Version, request.Size)

	data, err := w.readFile(dirpath, filename)
	if err != nil {
		log.Error().Err(err).Msgf("[warehouse] ОШИБКА ЧТЕНИЯ ФАЙЛА %s/%s", dirpath, filename)
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.Data(http.StatusOK, "image/jpeg", data)
}

func (w *Warehouse) HandleUploadArt(c *gin.Context) {
	// тут приходит jpeg файл в размере F
	data, artID, err := w.parseArt(c)
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

	unityMask := model.Art{ID: uint(artID)}.GetUnityMask(model.Unity1K)
	subdir := fmt.Sprintf("U%s", unityMask)
	filename := fmt.Sprintf("art-%d-%s.jpg", artID, "%s")
	if err := w.resizeAndSave(c, decoded, w.ArtsPath, subdir, filename); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Error().Err(err).Msgf("[warehouse:HandleUploadArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPEG-ФАЙЛ")
		return
	} else {
		c.String(http.StatusOK, "ok")
	}
}

func (w *Warehouse) HandleUploadUnity(c *gin.Context) {
	fileData, mask, version, err := w.parseUnity(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadUnity] БИТЫЕ ДАННЫЕ")
		return
	}

	log.Info().Msgf("[warehouse:HandleUploadArt] ВХОДЯЩИЙ UNITY-ЗАПРОС U%s/%d", mask, version)

	b := bytes.NewReader(fileData)
	decoded, err := jpeg.Decode(b)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Warn().Err(err).Msgf("[warehouse:HandleUploadUnity] БИТЫЙ JPEG. U%s/%d", mask, version)
		return
	}

	// Frigg: все единства делятся на папки по ведущему знаку, Каждое 100К-единство окажется в папке вместе со всеми его детьми
	subdir := fmt.Sprintf("U%sXXXXX", string(mask[0]))
	filename := fmt.Sprintf("U%s-%s-%s.jpg", mask, version, "%s")
	if err := w.resizeAndSave(c, decoded, w.UnityPath, subdir, filename); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Error().Err(err).Msgf("[warehouse:HandleUploadUnity] НЕ УДАЛОСЬ СОХРАНИТЬ JPEG-ФАЙЛ")
		return
	} else {
		c.String(http.StatusOK, "ok")
	}
}

func (w *Warehouse) HandleUploadOrigin(c *gin.Context) {
	// а тут приходит JPEG-файл в полном размере
	data, artID, err := w.parseArt(c)
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

	unityMask := model.Art{ID: uint(artID)}.GetUnityMask(model.Unity1K)
	subDir := fmt.Sprintf("U%s", unityMask)
	filename := fmt.Sprintf("art-%d.jpg", artID)

	if err := w.saveFile(c, w.OriginPath, subDir, filename, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Error().Err(err).Msgf("[warehouse:HandleUploadArt] НЕ УДАЛОСЬ СОХРАНИТЬ JPG-ФАЙЛ")
		return
	}

	c.String(http.StatusOK, "ok")
}

func (w *Warehouse) parseArt(c *gin.Context) (data []byte, artID int, err error) {
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
	defer f.Close()
	if err != nil {
		return
	}

	data, err = io.ReadAll(f)
	return
}

func (w *Warehouse) resizeAndSave(
	c *gin.Context,
	decoded image.Image,
	path string,
	subdir string,
	filemask string,
) error {
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
			return err
		}

		filename := fmt.Sprintf(filemask, size)
		if err := w.saveFile(c, path, subdir, filename, b.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func (w *Warehouse) parseUnity(c *gin.Context) (data []byte, unityMask string, unityVersion string, err error) {
	var file *multipart.FileHeader
	file, err = c.FormFile("file")
	if err != nil {
		return
	}

	unityMask = c.PostForm("mask")
	if unityMask == "" {
		err = errors.Errorf("[warehouse] MASK ОБЯЗАТЕЛЕН")
		return
	}

	unityVersion = c.PostForm("version")
	if unityVersion == "" {
		err = errors.Errorf("[warehouse] VERSION ОБЯЗАТЕЛЕН")
		return
	}

	var f multipart.File
	f, err = file.Open()
	defer f.Close()
	if err != nil {
		return
	}

	data, err = io.ReadAll(f)
	return
}

func (w *Warehouse) saveFile(
	ctx context.Context,
	baseDir string,
	subDir string,
	filename string,
	data []byte,
) error {
	dirpath := path.Join(baseDir, subDir)
	if err := os.MkdirAll(dirpath, 0o774); err != nil {
		return err
	}
	log.Info().Msgf("DIR CREATED %s", dirpath)

	filepath := path.Join(dirpath, filename)
	err := os.WriteFile(filepath, data, 0o774)
	return err
}

func (w *Warehouse) readFile(dir string, filename string) ([]byte, error) {
	b, err := os.ReadFile(path.Join(dir, filename))
	return b, err
}
