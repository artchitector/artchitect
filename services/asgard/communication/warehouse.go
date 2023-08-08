package communication

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/libraries/resizer"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Warehouse - склад бинарных картинок
// Odin: Картинки хранятся не в Асгарде, а на серверах файловых-хранилищах.
// Odin: Warehouse инкапсулирует сложную логику сохранения на файловые серверы, скрывая её от слоя Асгарда, и от Меня pantheon.Odin.
// Loki: А SOLID можешь расшифровать?) Ты уже прокачался в программировании, как я посмотрю.
// Loki: Уже боюсь проиграть наше пари...
type Warehouse struct {
	storageURL string // хранилище для размера XF. На том конце запущен сервис warehouse для приёма файлов
	memoryURL  string // хранилище для размеров F, M, S, XS. На том конце запущен сервис warehouse для приёма файлов
}

func NewWarehouse(storageURL string, memoryURL string) *Warehouse {
	return &Warehouse{storageURL: storageURL, memoryURL: memoryURL}
}

func (wh *Warehouse) SaveImage(ctx context.Context, artID uint, img image.Image) error {
	/*
		Odin:
		Ну чтож, дети мои, слушайте мой рассказ.
		Архитектор рисует большие изображения 3328 × 5120 пикселов (PNG ~10-15Мб). Это настроено в infrastructure.AI
		Этот размер называется XF - extra full.
		Extra full - всегда максимально доступное качество. Оригинал.
		Такие размеры нужны для того, чтобы картина могла быть ОТПЕЧАТАНА на холсте 40x60 в физическом мире.

		Но на сайте такие объёмы просто не нужны, там нужны 4 других размера в формате Jpeg

		/model/utils.go:
		SizeF  = "f"  // 1024x1536 Jpeg-Quality: 90
		SizeM  = "m"  // 512x768   Jpeg-Quality: 80
		SizeS  = "s"  // 256x384   Jpeg-Quality: 75
		SizeXS = "xs" // 128x192   Jpeg-Quality: 75

		Задачи склада:
		- отправить XF на долговременный сервер-хранилище HDD, где лежат оригиналы
		- ужать до размера F и отправить сжатую картинку на быстрый memory-сервер. Оттуда их забирает Alfheimr, который передаёт картинки на Midgard.
		(это для безопасности memory-сервер скрыт за api-gateway)

		Внутри самого memory-сервера картинка будет пережата на остальные размеры, и затем будет доступна для просмотра
	*/

	if err := wh.saveExtraFull(ctx, artID, img); err != nil {
		return errors.Wrapf(err, "[warehouse] СОХРАНЕНИЕ XF-КАРТИНКИ ART=%d ПРОВАЛЕНО", artID)
	}
	if err := wh.saveToMemory(ctx, artID, img); err != nil {
		return errors.Wrapf(err, "[warehouse] СОХРАНЕНИЕ F-КАРТИНКИ ART=%d ПРОВАЛЕНО", artID)
	}

	return nil
}

func (wh *Warehouse) saveExtraFull(ctx context.Context, artID uint, img image.Image) error {
	s := time.Now()
	filename := fmt.Sprintf("art-%d.jpg", artID) // filename не имеет особого значения
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityTransfer}); err != nil {
		return errors.Wrapf(err, "[warehouse] СЖАТЬ В BIG-JPEG НЕ УДАЛОСЬ. ОТКАЗ")
	}

	err := wh.makeRequest(ctx, wh.storageURL, map[string]string{
		"art_id": fmt.Sprintf("%d", artID),
	}, filename, buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "[warehouse] ARD_ID=%d. ОШИБКА СВЯЗИ С СЕРВЕРОМ %s", artID, wh.storageURL)
	}
	log.Info().Msgf("[warehouse] XF-КАРТИНКА #%d СОХРАНЕНА НА СКЛАД. T:%s", artID, time.Now().Sub(s))
	return nil
}

func (wh *Warehouse) saveToMemory(ctx context.Context, artID uint, img image.Image) error {
	s := time.Now()
	img = wh.downscaleToF(img)

	filename := fmt.Sprintf("art-%d.jpg", artID) // filename не имеет особого значения
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityTransfer}); err != nil {
		return errors.Wrapf(err, "[warehouse] СЖАТЬ В JPEG НЕ УДАЛОСЬ. ОТКАЗ")
	}

	err := wh.makeRequest(ctx, wh.memoryURL, map[string]string{
		"art_id": fmt.Sprintf("%d", artID),
	}, filename, buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "[warehouse] ARD_ID=%d. ОШИБКА СВЯЗИ С СЕРВЕРОМ %s", artID, wh.memoryURL)
	}
	log.Info().Msgf("[warehouse] F-КАРТИНКА #%d СОХРАНЕНА НА СКЛАД. T:%s", artID, time.Now().Sub(s))
	return nil
}

func (wh *Warehouse) downscaleToF(img image.Image) image.Image {
	return resizer.ResizeImage(img, model.SizeF)
}

func (wh *Warehouse) makeRequest(
	ctx context.Context,
	url string,
	fields map[string]string,
	filename string,
	binaryData []byte,
) error {
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Поле для файла
	fileW, err := multipartWriter.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	r := bytes.NewReader(binaryData)
	_, err = io.Copy(fileW, r)
	if err != nil {
		return err
	}

	// Остальные поля
	for field, value := range fields {
		fieldW, err := multipartWriter.CreateFormField(field)
		if err != nil {
			return err
		}
		_, err = fieldW.Write([]byte(value))
		if err != nil {
			return err
		}
	}

	if err := multipartWriter.Close(); err != nil {
		return err
	}

	// Подготовка запроса
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Сам запрос
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("[warehouse] ОШИБКА ЗАПРОСА %s. Status: %d. Body: %s", url, res.StatusCode, body)
	}

	return nil
}
