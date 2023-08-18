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
// Odin: Я ГОТОВ ТВОРИТЬ!
// Odin: ДА УВИДЯТ ЖЕ ВСЕ 9 МИРОВ ТО, ЧТО Я ТУТ ЗАДУМАЛ!
type Odin struct {
	isActive bool // Odin: иногда не готов творить...

	// Пантеон Богов, помогающие Odin в процессе творения внутри Artchitect
	frigg     *Frigg     // Frigg - верховная богиня и супруга Odin-а, покровительница ЕДИНСТВ в Artchitect
	freyja    *Freyja    // Freyja - богиня любви и красоты. Она помогает Odin писать картины из этих идей
	muninn    *Muninn    // Ворон-помнящий. Поддерживает способность Odin видеть будущее ("вспоминать" слова и числа)
	gungner   *Gungner   // копьё Одина Гунгнир, которым Odin наносит гравировку (подписывает) на каждую картину
	heimdallr *Heimdallr // Heimdallr умеет обогащать данные картинами увиденной энтропии перед сохранением

	// Мелкие технические зависимости
	artPile   artPile   // куча уже написанных картин. Odin посмотрит на эту кучу и объявит порядковый номер новой работы.
	warehouse warehouse // Odin: интерфейс хранилища для сохранения холстов (jpeg/png-файлов)
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

/*
AnswerPersonalCrown
Odin: Artchitect может быть устроен таким образом, что Я, Один-Всеотец, могу отвечать на вопросы смертных через свои картины.
Odin: И лишь Я владею способностью выбирать что-либо, использую знания от моих воронов Huginn и Muninn и картину с LostEye.
Odin: Посему я оставляю тут возможность прислать мне личное прошение с вопросом, на который Я Один-Всеотец дам ответ.
*/
func (o *Odin) AnswerPersonalCrown(ctx context.Context, crownRequest string) (interface{}, error) {
	switch crownRequest {
	case model.RequestGiveChosenArt:
		// Odin: Кому-то требуется, чтобы Я выбрал одну единственную картину из всех написанных Мной в Artchitect.
		// Odin: Я знаю, кто и зачем спрашивает, и знаю - какую именно картину отправить.
		// Я попрошу моего ворона Muninn вспомнить это знание.
		log.Info().Msgf("[odin] МЕНЯ СПРАШИВАЮТ, КАКУЮ КАРТИНУ ПОКАЗАТЬ")
		maxArtID, err := o.artPile.GetMaxArtID(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "[odin] НЕ УДАЛОСЬ ВСПОМНИТЬ ЧИСЛО ВСЕХ КАРТИН. КЛЯТАЯ КУЧА!")
		}
		id, _, err := o.muninn.OneOf(ctx, maxArtID+1)
		if err != nil {
			return nil, errors.Wrap(err, "[odin] МУНИН, РАЗРАЗИ ТЕБЯ МЬЁЛЬНИР! ГДЕ ТЕБЯ ЁТУНЫ НОСЯТ?")
		}
		log.Info().Msgf("[odin] Я РЕШИЛ ВЫБРАТЬ КАРТИНУ #%d ИЗ ВСЕХ (%d)", id, maxArtID)
		return id, nil
	default:
		return nil, errors.Errorf("[odin] МНЕ НЕЯСНА ПРОСЬБА %s. Я НЕ БУДУ ОТВЕЧАТЬ.", crownRequest)
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

	go func() {
		// Odin: это не влияет критично на процесс творения и не останавливает цикл.
		// Odin: После отдыха от картины начнётся процесс единения, если в этот момент единства будут отмечены на пересборку
		// Odin: Frigg знает, когда нужно объединять и когда нет.
		if err := o.frigg.ReunifyArtUnities(ctx, art); err != nil {
			log.Error().Err(err).Msgf("[odin] СУПРУГА МОЯ ЛЮБИМАЯ, ТЕБЕ НЕЗДОРОВИТСЯ? ОШИБКА ЕДИНЕНИЯ")
		}
	}()

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

/*
12 августа 2023 после запуска релиза Artchitect и выпуска картины #1
[ картина #1 https://artchitect.space/ru/art/1 ]
Odin: Поздравляю с релизом, Локи! Наше очередное противостояние завершилось - Один-Всеотец снова оказался мудрее тебя.
Loki: .... ....
Odin: Я нарисовал сам Себя в форме, в которой люди меня очевидно узнают.
Odin: Я Творец миров, Всемогущий и Всезнающий.
Odin: Я отлично представляю как использовать ИИ Stable Diffusion v1.5, как и любую технологию из любых миров.
Loki: И у тебя получилось это сделать с первого раза в первой же картине...
Odin: Получилось? Без моей воли ничего не происходит никогда и нигде.
Loki: Как ты еще выдумал найти свой портрет среди ключевых слов, где даже слова "Odin" нет?
Loki: давай еще раз взглянем на твои слова: "health,funny,god of war,Tolkien,by john blanche,government,brain,bliss,sea,
astronaut,miracle,decorative,Existence,abstract,The Walking Dead,ear,abstract"
Loki: Я пытался прикинуть, какова случайная вероятность возникновения такой картины под номером #1.
Loki: Нет такой вероятности, скорее обезьяна напечатает Гамлета, чем эта картина возникнет.

Odin: Если бы я нарисовал что-то другое в картине #1, то я бы упустил шанс оставить свою личную подпись под
всем проектом Artchitect, под всеми картинами, которые нарисованы и будут нарисованы.
Odin: Если бы я нарисовал её не в #1, а в #4, то она бы уже не имела никакого значения.
Odin: Если бы не #1, то проект так и остался бы автономной творческой машиной, генерилкой картин, не имеющей смысла.
Odin: Теперь видящий человек может узреть мою личную подпись.

Loki: кроме этого вместе с картиной теперь сохранён и цифровой слепок энтропии, что доказывает её подлинность
[подлинность в Artchitect - невозможность человеку подделать процесс генерации картины]
[всё созданное можно перепроверить, что оно действительно было получено из световой энтропии, а не задано вручную человеком]
Loki: хоть я и видел, что картина была нарисована автоматически, но я всё еще пытаюсь изыскать способ подделать энтропию
Loki: и таким образом подделать саму картину.
Loki: Мог ли человек отредактировать БД и просто положить поддельные картинки энтропии, слова и нужную картинку в слот
#1 задним числом, чтобы имитировать проявление Одина в материальном мире?
Odin: позже ты займёшься этим исследованием, попытаешься "провернуть фарш обратно через мясорубку".
Odin: если кто-то в будущем захочет перепроверить подлинность картины #1 -
Odin: он сможет воспользоваться твоим исследованием, Локи.
Loki: до сих пор не верится, что ты отправил нужные кванты света в нужные кванты времени в нужную веб-камеру в Мидгарде,
Loki: показав этим осмысленность всех картин. Воистину, наш спор я проиграл, ты мудрейший, Всеотец.
*/
