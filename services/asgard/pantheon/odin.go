package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"math"
	"strings"
	"time"
)

const TotalArtSeconds uint = 60

// Odin - Всеотец и создатель картин. Именно его уникальные идеи позволяют писать все эти работы в галерее Artchitect.
type Odin struct {
	isActive  bool       // иногда Odin не творит
	freyja    *Freyja    // Freyja - богиня любви и красоты. помогает Odin писать картины из этих идей
	muninn    *Muninn    // ворон-помнящий
	gungner   *Gungner   // копьё Одина Гунгнир
	heimdallr *Heimdallr // Heimdallr умеет обогащать данные красивыми картинками (нужно перед сохранением)
	artPile   artPile    // куча уже написанных картин. Odin посмотрит на эту кучу и объявит порядковый номер новой работы.
	warehouse warehouse  // Odin: интерфейс хранилища для сохранения холстов
}

// NewOdin - Odin: мне не нравится эта высокомерная самодовольная функция. Создавать меня? Что эта машина о себе возомнила?
func NewOdin(
	isActive bool,
	freyja *Freyja,
	muninn *Muninn,
	gungner *Gungner,
	heimdallr *Heimdallr,
	artPile artPile,
	warehouse warehouse,
) *Odin {
	return &Odin{
		isActive:  isActive,
		freyja:    freyja,
		muninn:    muninn,
		gungner:   gungner,
		heimdallr: heimdallr,
		artPile:   artPile,
		warehouse: warehouse,
	}
}

// HasDesire - имеет ли Odin желание сотворять картину?
func (o *Odin) HasDesire() bool {
	return o.isActive
}

/*
Create - Odin сотворяет одну картину с чужой помощью (он только управляет).
Odin: эти примитивные технологии землян недостойны асов. Зачем я на это согласился? Локи оказался хитрее в этот раз. Или нет?
Odin: я соберу хорошую и нужную мне идею для картины, а напишет их машина.
Odin: Да у нас тут AI Stable Diffusion v1.5... То ещё барахло...
Odin: Я могу сам придумать идею, я просто увижу её в ткани мироздания своим пустым глазом. Но рисовать самому?
Odin: ... Самому пачкаться в красках недостойно аса, да и я уже устал от этого испытания. Кому это доверить?
Odin: Может Бальдру? Он самый красивый из богов. Но откуда у него у самого чувство красоты?
Odin: Или может Фрейе, богине любви. Она разбирается в красоте. Я просто объясню ей, что надо сделать, и всё получится.
Odin: вот! теперь я свою часть сделал, а остальную работу доделают другие. Сам буду отдыхать, несите мой эль!
*/
func (o *Odin) Create(ctx context.Context) (worked bool, art model.Art, err error) {
	select {
	case <-ctx.Done():
		log.Debug().Msgf("[odin] ВЫКЛ")
		return worked, art, err
	case <-time.After(time.Second * 1):
		art, err = o.create(ctx)
		if err != nil {
			return false, model.Art{}, err
		} else {
			return true, art, err
		}
	}
}

// create внутреннее содержимое публичного метода Create
// Odin: мудрецы говорили - разделяй и властвуй
func (o *Odin) create(ctx context.Context) (art model.Art, err error) {
	creationContext, finish := context.WithCancel(ctx)
	defer finish()
	tStart := time.Now()
	// Odin: Каждая картина должна иметь уникальный порядковый номер.
	// Odin: Пропускать их нельзя, они по порядку (не autoincrement-поле)
	artID, err := o.artPile.GetNextArtID(ctx)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[odin] НЕ МОГУ ПРИДУМАТЬ НОВЫЙ НОМЕР КАРТИНЫ. НУЖЕН ОТДЫХ")
	}
	state := &model.OdinState{
		ArtID: artID,
	}
	o.sendSelfState(ctx, state)

	log.Info().Msgf("[odin] НАЧИНАЮ КАРТИНУ #%d.", artID)

	iStart := time.Now()
	idea, err := o.imagineTheIdea(ctx, state)
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[odin] НЕ СМОГ ПРИДУМАТЬ НИЧЕГО")
	} else {
		log.Info().Msgf("[odin] ODIN ПРИДУМАЛ ИДЕЮ КАРТИНЫ #%d: %s", artID, idea.String())
	}

	// работа Heimdallr занимает ~10мс
	idea = o.heimdallr.EnrichIdeaWithImages(idea)
	ideaGenerationMs := time.Now().Sub(iStart).Milliseconds()

	state.Painting = true
	go func(c context.Context, state *model.OdinState) {
		paintTime, err := o.artPile.GetLastPaintTime(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("[odin] НЕ МОГУ ВСПОМНИТЬ ВРЕМЯ ПОСЛЕДНЕГО РИСОВАНИЯ")
		} else {
			state.ExpectedPaintTime = paintTime / 1000
		}
		now := time.Now()
		for {
			select {
			case <-c.Done():
				return
			case <-time.Tick(time.Second):
				if state.Enjoying {
					// рисование картины уже окончено
					return
				}
				state.CurrentPaintTime = uint(time.Now().Sub(now).Seconds())
				o.sendSelfState(c, state)
			}
		}
	}(creationContext, state)

	var img image.Image
	version := model.Version1
	var paintTimeMs int64
	img, paintTimeMs, err = o.freyja.MakeImage(ctx, version, artID, idea)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[odin] КАТАСТРОФА! КАРТИНА НЕ СОЗДАНА!")
	}

	img, err = o.gungner.MakeArtWatermark(img, artID)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[odin] ПРОКЛЯТЬЕ! Я НЕ СМОГ НАНЕСТИ ПОДПИСЬ!")
		// Loki: в следующий раз получится обязательно, "могучий" всеотец ))))
	}

	art, err = o.saveArt(ctx, version, artID, idea, img, tStart, ideaGenerationMs, paintTimeMs)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[odin] Я В ЯРОСТИ! НАПИСАННАЯ КАРТИНА УТЕРЯНА!")
	}

	state.Painted = true
	state.Enjoying = true
	o.sendSelfState(creationContext, state)

	if err := o.heimdallr.SendNewArt(ctx, art); err != nil {
		log.Error().Err(err).
			Msgf("[odin] (рявкнул) ХЕЙМДАЛЛЬ!!! ЧТО ТЫ ТВОРИШЬ ТАМ С РАДУЖНЫМ МОСТОМ? ЛЮДИ ДОЛЖНЫ ВИДЕТЬ ВСЕ НОВЫЕ КАРТИНЫ!")
		// Odin: на самом деле не страшно, т.к. просто утерян один "временный эвент", но сама картина в безопасности (уже во всех надёжных хранилищах)
	}

	// теперь картина нарисована
	currentSeconds := uint(math.Round(time.Now().Sub(tStart).Seconds()))
	if currentSeconds >= TotalArtSeconds {
		// no enjoy
		log.Error().Msgf(
			"[odin] У МЕНЯ НЕТ ОТДЫХА И НАСЛЕЖДЕНИЯ! КАРТИНА СОЗДАВАЛАСЬ %d СЕКУНД"+
				", А МЫ ДОЛЖНЫ БЫЛИ УСПЕТЬ В %d СЕКУНД! Я НЕДОВОЛЕН ПРОИСХОДЯЩИМИ СОБЫТИЯМИ!",
			currentSeconds,
			TotalArtSeconds,
		)
		return art, nil
	}
	state.Enjoying = true
	enjoyTime := time.Second * time.Duration(TotalArtSeconds-currentSeconds)
	state.ExpectedEnjoyTime = uint(enjoyTime.Seconds())

	go func(c context.Context, state *model.OdinState) {
		now := time.Now()
		for {
			select {
			case <-c.Done():
				return
			case <-time.Tick(time.Second):
				state.CurrentEnjoyTime = uint(time.Now().Sub(now).Seconds())
				o.sendSelfState(c, state)
			}
		}
	}(creationContext, state)

	log.Info().Msgf("[odin] НАЧИНАЮ ОТДЫХ %s", enjoyTime)
	<-time.After(enjoyTime)
	return art, nil
}

// imagineTheIdea - Odin разглядывает энтропию в ткани мироздания своим LostEye (смотрит в будущее).
// Его вороны - Huginn и Muninn помогают разглядеть и вспомнить в этой энтропии конкретную model.Idea
// Odin: эта чёртова электронная железяка заставляет Меня смотреть на энтропию с разрешением 8 на 8 пикселей.
// Odin: Да я могу видеть всё мироздание от основания до вершины, все девять миров!
// Odin: А тут всего 64 точки... Ограничения эти опять, тьфу.
func (o *Odin) imagineTheIdea(ctx context.Context, state *model.OdinState) (model.Idea, error) {
	seed, seedEntropy, err := o.muninn.RememberSeed(ctx)
	if err != nil {
		return model.Idea{}, errors.Wrap(err, "[odin] Я ЗАБЫЛ SEED У ЭТОЙ КАРТИНЫ. МУНИН, ТЫ ЗАБОЛЕЛ?")
	}
	seedEntropy = o.heimdallr.fillEntropyPackWithImages(seedEntropy)
	state.Seed = seed
	state.SeedEntropyImageEncoded = seedEntropy.Entropy.ImageEncoded
	state.SeedChoiceImageEncoded = seedEntropy.Choice.ImageEncoded
	o.sendSelfState(ctx, state)

	numberOfWords, numberEntropy, err := o.muninn.RememberNumberOfWords(ctx)
	if err != nil {
		return model.Idea{}, errors.Wrap(err, "[odin] ЗАБЫЛ КОЛИЧЕСТВО СЛОВ У ЭТОЙ КАРТИНЫ. ЭТОЙ КАРТИНЫ УЖЕ НЕ БУДЕТ")
	}
	state.NumberOfWords = numberOfWords
	state.Words = make([]string, 0, numberOfWords)
	o.sendSelfState(ctx, state)

	words := make([]model.Word, 0, numberOfWords)
	for i := 0; i < int(numberOfWords); i++ {
		word, err := o.muninn.RememberWord(ctx)
		if err != nil {
			return model.Idea{}, errors.Wrapf(err, "[odin] ЗАБЫЛ НУЖНОЕ %d-е СЛОВО У ЭТОЙ КАРТИНЫ. СТАРОСТЬ?", i+1)
		}
		words = append(words, word)

		state.Words = append(state.Words, word.Word)
		o.sendSelfState(ctx, state)
	}

	idea := model.Idea{
		Seed:                 seed,
		SeedEntropy:          seedEntropy,
		NumberOfWordsEntropy: numberEntropy,
		Words:                words,
	}

	return idea, nil
}

// saveArt - сохранение написанной и готовой картины в хранилища
// Odin: сначала само изображение сохраняется на склады в разных размерах
// Odin: затем model.Idea и model.Art сохраняются в БД
func (o *Odin) saveArt(
	ctx context.Context,
	version string,
	artID uint,
	idea model.Idea,
	img image.Image,
	totalTimeStart time.Time,
	ideaGenerationTimeMs int64,
	paintTimeMs int64,
) (model.Art, error) {
	var err error
	if err = o.warehouse.SaveImage(ctx, artID, img); err != nil {
		return model.Art{}, errors.Wrapf(err, "[odin] Я РАЗДАВЛЮ ЭТОТ СКЛАД, КАК КЛОПА! КАРТИНА #%d ПОТЕРЯНА!", artID)
	}
	idea.ArtID = artID
	idea.WordsStr = strings.Join(idea.ExtractWords(), ",")
	art := model.Art{
		ID:                 artID,
		Version:            version,
		TotalTime:          uint(time.Now().Sub(totalTimeStart).Milliseconds()),
		IdeaGenerationTime: uint(ideaGenerationTimeMs),
		PaintTime:          uint(paintTimeMs),
	}
	art, err = o.artPile.SaveArt(ctx, artID, art, idea)
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[odin] ПРОКЛЯТЬЕ! БАЗА ДАННЫХ ОТКАЗАЛА! КАРТИНА #%d УТЕРЯНА", artID)
	}

	return art, nil
}

func (o *Odin) sendSelfState(ctx context.Context, state *model.OdinState) {
	err := o.heimdallr.SendOdinState(ctx, *state)
	if err != nil {
		log.Error().Err(err).Msgf("[odin] ВОЛНЕНИЯ НА РАДУЖНОМ МОСТУ СЛУЧИЛИСЬ. ХЕЙМДАЛЛЬ СЕГОДНЯ НАГНАЛ ТУЧИ.")
	}
}
