package external

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"image"
	"image/png"
	"strings"
)

type Artist struct {
	ai ai
}

func NewArtist(ai ai) *Artist {
	return &Artist{ai: ai}
}

func (a *Artist) MakeArt(
	ctx context.Context,
	artID uint,
	seed uint,
	tags []string,
) (image.Image, error) {
	prompt := strings.Join(tags, ",")

	// запрос картины у ИИ по seed-номеру и ключевым словам
	imgData, err := a.ai.GenerateImage(ctx, seed, prompt)
	if err != nil {
		return nil, errors.Wrap(err, "[artist] АВАРИЯ ИИ")
	}

	// картина создаётся в png ради сохранения качества
	img, err := a.decode(imgData)
	if err != nil {
		return nil, errors.Wrap(err, "[artist] ОШИБКА РАСШИФРОВКИ")
	}

	// Нанесение водяного знака (с котом и номером работы в углу картинки). Кота тоже рисовал Архитектор, но v1 (опытная версия)
	img, err = a.makeWatermark(img, artID)
	if err != nil {
		return nil, errors.Wrap(err, "[artist] ОШИБКА ВОДЯНОГО ЗНАКА")
	}
	return img, nil
}

func (a *Artist) decode(data []byte) (image.Image, error) {
	b := bytes.NewReader(data)
	img, err := png.Decode(b)
	return img, err
}

func (a *Artist) makeWatermark(img image.Image, artID uint) (image.Image, error) {
	return nil, errors.New("[artist] ВОДЯНОЙ ЗНАК - НЕ ГОТОВО")
}
