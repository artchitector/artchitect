package external

import (
	"context"
	"github.com/pkg/errors"
	"image"
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

	// Нанесение водяного знака (с котом и номером работы в углу картинки)
	img, err = a.makeWatermark(img, artID)
	return img, nil
}

func (a *Artist) decode(data []byte) (image.Image, error) {
	return nil, errors.New("[artist] РАСШИФРОВКА - НЕ ГОТОВО")
}

func (a *Artist) makeWatermark(img image.Image, artID uint) (image.Image, error) {
	return nil, errors.New("[artist] ВОДЯНОЙ ЗНАК - НЕ ГОТОВО")
}
