package pantheon

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/png"
	"strings"
	"time"
)

// Freyja - Я богиня красоты и любви, так что я лучше остальных знаю толк в красоте. Потому Odin и доверил мне эту работу.
// Freyja - Я знаю, как заставить примитивный искусственный интеллект мидгардцев работать на нашу галерею Artchitect.
type Freyja struct {
	// Odin: ох уж эти кустарные технологии, недостойные асов. Я мог бы сотворить сотни миров лишь по щелчку пальцев,
	// 		но в этот раз вынужден сидеть внутри чёртова механизма, который сделали смертные из Мидгарда!
	//		Локи провёл меня на этот раз. Но... отступать уже поздно.
	ai ai
}

func NewFreyja(ai ai) *Freyja {
	return &Freyja{ai: ai}
}

func (a *Freyja) MakeImage(
	ctx context.Context,
	version string, // Freyja: тут заложена возможность использовать разные версии Stable Diffusion
	artID uint,
	idea model.Idea, // Thor: Держи идею, Фрейя! В ней Один "зашептал" свою идею.
) (image.Image, int64, error) {
	// Freyja: не такая уж сложная работа, всего строчку собрать.
	prompt := strings.Join(idea.ExtractWords(), ",")

	// Freyja: тут земной примитивный ИИ рисует картинку
	pStart := time.Now()
	imgData, err := a.ai.GenerateImage(ctx, idea.Seed, prompt)
	if err != nil {
		return nil, 0, errors.Wrap(err, "[freyja] ЭТОТ ИИ СЛОМАН, НЕСИТЕ ДРУГОГО")
	}
	log.Debug().Msgf("[freyja] ПРИМИТИВНЫЙ ИИ НАРИСОВАЛ ЧТО-ТО")

	img, err := a.decode(imgData)
	if err != nil {
		return nil, 0, errors.Wrap(err, "[freyja] КАКОЕ-ТО НЕДОРАЗУМЕНИЕ")
	}

	// Odin: ХОЧУ, чтобы на каждой картине напечатался водяной знак с номером картины, а рядом с ним был КОТ!
	// Freyja: Есть подходящий кот на картине "Есть ли кошачий Бог?" за авторством Artchitect (опытной первой версии).
	// Odin: других вариантов я не вижу, так что ок.
	/*
		[artchitector]: для любопытствующих смотреть файл services/asgard/files/images/is_there_cat_god.jpg.
		Это была первая отладка Artchitect, рисовалось всё без энтропии. Это было самое начало проекта Artchitect.
		Card #1294.
		Created: 2023 Jan 6 18:31
		Seed: 4091966908
		Words: intricate,cat,Sun,galactic,nuclear,symmetrical,Allah,girl,stunning beautiful,europe,dynamic lighting,
			greek,darkblue,art,sadness,light,fantastically beautiful,red,Gothic,train,john constable,textured,yellow,
			tribal patterns,hyper,high details,electricity
	*/

	// Фрейя наносит водяной знак (с котом и номером работы в углу картинки).
	img, err = a.makeWatermark(img, artID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "[freyja] ВОДЯНОЙ ЗНАК НЕ НАНЕСЁН. КАРТИНА ОТПРАВЛЯЕТСЯ В УТИЛЬ")
	}

	return img, time.Now().Sub(pStart).Milliseconds(), nil
}

func (a *Freyja) decode(data []byte) (image.Image, error) {
	b := bytes.NewReader(data)
	// Freyja: Красота этой картины не должна пострадать от мерзкого сжатия, поэтому я использую PNG.
	img, err := png.Decode(b)
	return img, err
}

func (a *Freyja) makeWatermark(img image.Image, artID uint) (image.Image, error) {
	return img, nil
	//return nil, errors.New("[freyja] ВОДЯНЫЕ ЗНАКИ МЫ ЕЩЕ НЕ ИЗГОТОВИЛИ")
}
