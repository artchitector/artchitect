package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// Odin - Всеотец и создатель картин. Именно его уникальные идеи позволяют писать все эти работы в галерее Artchitect.
type Odin struct {
	isActive bool    // иногда Odin не творит
	freyja   *Freyja // Freyja - богиня любви и красоты. помогает Odin рисовать картины из этих идей
	muninn   *Muninn // ворон-помнящий
	artPile  artPile // куча уже написанных картин. Odin посмотрит на эту кучу и объявит порядковый номер новой работы.
}

// NewOdin - Odin: мне не нравится эта высокомерная самодовольная функция. Создавать меня? Что эта машина о себе возомнила?
func NewOdin(
	isActive bool,
	freyja *Freyja,
	muninn *Muninn,
	artRepo artPile,
) *Odin {
	return &Odin{
		isActive: isActive,
		freyja:   freyja,
		muninn:   muninn,
		artPile:  artRepo,
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
	// Odin: Каждая картина должна иметь уникальный порядковый номер.
	// Odin: Пропускать их нельзя, они по порядку (не autoincrement-поле)
	artID, err := o.artPile.GetNextArtID(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("[odin] НЕ МОГУ ПРИДУМАТЬ НОВЫЙ НОМЕР КАРТИНЫ. НУЖЕН ОТДЫХ")
	}

	log.Info().Msgf("[odin] НАЧИНАЮ КАРТИНУ #%d.", artID)

	idea, err := o.imagineTheIdea(ctx)
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[odin] НЕ СМОГ ПРИДУМАТЬ НИЧЕГО")
	}

	_, err = o.freyja.MakeArt(ctx, artID, idea)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[odin] ТРЕВОГА! КАРТИНА НЕ СОЗДАНА!")
	}

	return model.Art{}, nil // TODO продолжить тут
}

// imagineTheIdea - Odin разглядывает энтропию в ткани мироздания своим LostEye (смотрит в будущее).
// Его вороны - Huginn и Muninn помогают разглядеть и вспомнить в этой энтропии конкретную model.Idea
// Odin: эта чёртова электронная железяка заставляет Меня смотреть на энтропию с разрешением 8 на 8 пикселей.
// Odin: Да я могу видеть всё мироздание от основания до вершины, все девять миров!
// Odin: А тут всего 64 точки... Ограничения эти опять, тьфу.
func (o *Odin) imagineTheIdea(ctx context.Context) (model.Idea, error) {
	seed, seedEntropy, err := o.muninn.rememberSeed(ctx)
	if err != nil {
		return model.Idea{}, errors.Wrap(err, "[odin] Я ЗАБЫЛ SEED У ЭТОЙ КАРТИНЫ. МУНИН, ТЫ ЗАБОЛЕЛ?")
	}
	numberOfWords, numberEntropy, err := o.muninn.rememberNumberOfWords(ctx)
	if err != nil {
		return model.Idea{}, errors.Wrap(err, "[odin] ЗАБЫЛ КОЛИЧЕСТВО СЛОВ У ЭТОЙ КАРТИНЫ. ЭТОЙ КАРТИНЫ УЖЕ НЕ БУДЕТ")
	}
	words := make([]model.Word, 0, numberOfWords)
	for i := 0; i < int(numberOfWords); i++ {
		word, err := o.muninn.rememberWord(ctx)
		if err != nil {
			return model.Idea{}, errors.Wrapf(err, "[odin] ЗАБЫЛ НУЖНОЕ %d-е СЛОВО У ЭТОЙ КАРТИНЫ. СТАРОСТЬ?", i+1)
		}
		words = append(words, word)
	}

	idea := model.Idea{
		Seed:                 seed,
		SeedEntropy:          seedEntropy,
		NumberOfWordsEntropy: numberEntropy,
		Words:                words,
	}

	return idea, nil
}
