package pantheon

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/bits"

	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

/*
Heimdallr - страж Биврёста. Хе́ймдалль поддерживает поток передачи энтропии из LostEye в communication.Bifröst.

Heimdallr получает энтропию не в виде картинки, а в виде матрицы силы пикселя, поэтому png-картинку с
нужным цветом генерирует Heimdallr.

Odin: Нам нужно постоянно поддерживать поток энтропии, который отправляется из Асгарда через Альфхейм в Мидгард.
Odin: Мидгард должен видеть то же, что и я вижу. Heimdallr, ты страж радужного моста. Ты можешь создать этот поток.
Heimdallr: Да, Odin. Я сделаю это, но я добавлю свою маленькую деталь...
*/
type Heimdallr struct {
	huginn      *Huginn
	bifröst     bifröst
	colorScheme *colorScheme
	loki        *Loki
}

// colorScheme цветовая схема для раскрашивания энтропии
// Heimdallr: я буду менять цвет энтропии плавно, чтобы она переливалась всеми существующими RGB-цветами
type colorScheme struct {
	Red   colorSchemeColor
	Green colorSchemeColor
	Blue  colorSchemeColor
}

type colorSchemeColor struct {
	Current float64 // текущий уровень цвета
	Target  float64 // требуемый уровень цвета, к которому в скором времени придёт состояние цвета
	Step    float64 // шаг изменения Current в Target за каждый просчёт энтропии. Цвет будет плавно меняться
}

func NewHeimdallr(huginn *Huginn, bifröst bifröst, loki *Loki) *Heimdallr {
	return &Heimdallr{
		huginn:      huginn,
		bifröst:     bifröst,
		loki:        loki,
		colorScheme: nil,
	}
}

/*
StartStream - запуск потока передачи энтропии из Huginn в communication.Bifröst
Heimdallr: брррррр... возьмёмся за дело.
Odin: этот поток критически важен, его нельзя терять. При его отключении надо "ронять" весь Asgard для его перезапуска.
Heimdallr: проооще простооого... (протяжно и безразлично)
*/
func (h *Heimdallr) StartStream(ctx context.Context) {
	entropyCh := h.huginn.Subscribe(ctx, "heimdallr_transfer_stream")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[heimdallr] МММММ. НАКОНЕЦ И МОЯ ОСТАНОВКА")
			return
		case entropy := <-entropyCh:
			// log.Debug().Msgf("[heimdallr] ВИЖУ ЭНТРОПИЮ ОТ ХУГИНА %+v", entropy)
			go func(entropy model.EntropyPackExtended) {
				if err := h.transferEntropy(ctx, entropy); err != nil {
					log.Error().Err(err).Msgf("[heimdallr] ОТПРАВКА ЭНТРОПИИ НЕ СЛУЧИЛАСЬ, К СОЖАЛЕНИЮ!")
				}
			}(entropy)
		}
	}
}

func (h *Heimdallr) EnrichIdeaWithImages(idea model.Idea) model.Idea {
	idea.NumberOfWordsEntropy = h.fillEntropyPackWithImages(idea.NumberOfWordsEntropy)
	idea.SeedEntropy = h.fillEntropyPackWithImages(idea.SeedEntropy)
	for i, word := range idea.Words {
		idea.Words[i].Entropy = h.fillEntropyPackWithImages(word.Entropy)
	}
	return idea
}

// SendNewArt
// Heimdallr: Odin написал новую картину. Я отправлю её по радужному мосту, дабы нижние миры тоже могли порадоваться этому
func (h *Heimdallr) SendNewArt(ctx context.Context, art model.Art) error {
	return h.bifröst.SendDrakkar(ctx, model.ChanNewArt, art)
}

func (h *Heimdallr) SendOdinState(ctx context.Context, state model.OdinState) error {
	return h.bifröst.SendDrakkar(ctx, model.ChanOdinState, state)
}

func (h *Heimdallr) SendFriggState(ctx context.Context, state model.FriggState) error {
	return h.bifröst.SendDrakkar(ctx, model.ChanFriggState, state)
}

// transferEntropy
// Heimdallr: так. Моя задача из матрицы силы пикселей сделать видимую картинку, раскрасить её как мне нравится
// Heimdallr: и далее отправить в нижние миры по радужному мосту, до наших любимый людей из Мидгарда.
func (h *Heimdallr) transferEntropy(ctx context.Context, entropy model.EntropyPackExtended) error {
	var err error
	h.updateColorScheme(entropy)
	entropy.Entropy.ImageEncoded = h.encodeEntropyImageFromMatrix(entropy.Entropy.Matrix)
	entropy.Choice.ImageEncoded = h.encodeEntropyImageFromMatrix(entropy.Choice.Matrix)

	entropy.ImageFrameEncoded, err = h.encodeJpeg(entropy.ImageFrame, model.EntropyJpegQualityFrame)
	if err != nil {
		log.Error().Msgf("[heimdallr] НЕ СМОГ СДЕЛАТЬ JPEG ЭНТРОПИИ")
	}
	entropy.ImageNoiseEncoded, err = h.encodeJpeg(entropy.ImageNoise, model.EntropyJpegQualityNoise)
	if err != nil {
		log.Error().Msgf("[heimdallr] НЕ СМОГ СДЕЛАТЬ JPEG ОБРАТНОЙ ЭНТРОПИИ")
	}

	if err = h.sendDrakkar(ctx, model.ChanEntropyExtended, entropy); err != nil {
		return errors.Wrap(err, "[heimdallr] BIFRÖST СЛОМАН, ДРАККАР С ГРУЗОМ ЭНТРОПИИ УТЕРЯН В ТКАНИ ПРОСТРАНСТВА")
	}

	// Odin: в Альфхейм и Мидгард отправятся две посылки с энтропией
	// Odin: model.EntropyPackExtended содержит еще и кадр с шумом, большие jpeg-картинки. Это тяжёлая модель,
	// и она нужна лишь в одном месте - на странице entropy
	// Odin: model.EntropyPack содержит только минимальные картинки 8х8 (весом по 200байт), это лёгкая модель, и
	// она используется везде в Artchitect для отображения текущей энтропии
	// Heimdallr: я отправлю два разные пакета в разные каналы. Читатели разберутся, кому какой нужен.
	miniEntropy := model.EntropyPack{
		Timestamp: entropy.Timestamp,
		Entropy:   entropy.Entropy,
		Choice:    entropy.Choice,
	}
	miniEntropy, err = h.encodeEntropyImagesForTransfer(miniEntropy)
	if err = h.sendDrakkar(ctx, model.ChanEntropy, miniEntropy); err != nil {
		return errors.Wrap(err, "[heimdallr] BIFRÖST СЛОМАН, ДРАККАР С ГРУЗОМ ЭНТРОПИИ УТЕРЯН В ТКАНИ ПРОСТРАНСТВА")
	}

	return nil
}

func (h *Heimdallr) fillEntropyPackWithImages(pack model.EntropyPack) model.EntropyPack {
	if pack.Entropy.ImageEncoded == "" {
		pack.Entropy.ImageEncoded = h.encodeEntropyImageFromMatrix(pack.Entropy.Matrix)
	}
	if pack.Choice.ImageEncoded == "" {
		pack.Choice.ImageEncoded = h.encodeEntropyImageFromMatrix(pack.Choice.Matrix)
	}
	return pack
}

func (h *Heimdallr) sendDrakkar(ctx context.Context, channel string, pack interface{}) error {
	// Heimdallr: теперь отправлю этот ценный драккар по волнам Биврёста...
	return h.bifröst.SendDrakkar(ctx, channel, pack)
}

// makeEntropyImage - матрица сил пикселей в энтропии становится 8х8 PNG картинкой
func (h *Heimdallr) makeEntropyImage(matrix model.EntropyMatrix) image.Image {
	bounds := image.Rect(0, 0, matrix.Size(), matrix.Size())
	img := image.NewRGBA(bounds)
	for x := 0; x < matrix.Size(); x++ {
		for y := 0; y < matrix.Size(); y++ {
			colr := h.makeColor(matrix.Get(x, y))
			img.Set(x, y, colr)
		}
	}
	return img
}

// makeColor - сделать цвет по силе пикселя
// Odin: сила пикселя лежит от 0 (полностью не светится) до 255 (полностью светится), но цвет не определён
// Heimdallr: Я охраняю РАДУЖНЫЙ мост, и знаю о цветах многое.
// Heimdallr: Я сделаю свой цвет. Это и есть моя упомянутая выше деталь.
func (h *Heimdallr) makeColor(power uint8) color.Color {
	// greenToRedProportion := 165.0 / 255.0

	// Цвета могут быть почти чёрные (по значениям Current, но иметь какой-то цвет). Надо нормализовать цвета так,
	// чтобы цвет был яркий

	return color.RGBA{
		R: uint8(float64(power) * h.colorScheme.Red.Current),
		G: uint8(float64(power) * h.colorScheme.Green.Current),
		B: uint8(float64(power) * h.colorScheme.Blue.Current),
		// R: power,
		// G: uint8(greenToRedProportion * float64(power)),
		// B: 0,
		A: math.MaxUint8,
	}
}

func (h *Heimdallr) updateColorScheme(entropy model.EntropyPackExtended) {
	const Chance = 0.90
	// Heimdallr: я собираюсь плавно менять цвет от одного к другому, чтобы он переливался
	if h.colorScheme == nil {
		// начинается всё с белого цвета
		h.colorScheme = &colorScheme{
			// шкала цвета 1.0 во float64 = 255 в uint8. 1.0 - полный красный, 0.5 - половина красного, 0.0 - нет красного
			Red: colorSchemeColor{
				Current: 1.0, // если Current != Target - цвет находится в изменении
				Target:  1.0,
				Step:    0.0,
			},
			Green: colorSchemeColor{
				Current: 1.0,
				Target:  1.0,
				Step:    0.0,
			},
			Blue: colorSchemeColor{
				Current: 1.0,
				Target:  1.0,
				Step:    0.0,
			},
		}
	}

	colors := map[string]*colorSchemeColor{
		"rød":   &h.colorScheme.Red,   // Odin: задействуем милый глазу норвежский язык!
		"grønn": &h.colorScheme.Green, // Odin: задействуем милый глазу норвежский язык!
		"blå":   &h.colorScheme.Blue,  // Odin: задействуем милый глазу норвежский язык!
	}

	// Heimdallr: менять цвета я буду тоже опираясь на энтропию, она послужит зерном к этим изменениям
	changedColor := false
	for farge, col := range colors { // Odin: farge - "цвет" по-норвежски. там лежит норвежское наименование
		if col.Current != col.Target {
			// Heimdallr: цвет уже меняется, оставим его
			col.Current += col.Step // цвет сдвинут на один шаг
			// Heimdallr: после сдвига цвет уже может быть успешно установлен
			if col.Step > 0 {
				// Яркость цвета росла
				if col.Current >= col.Target {
					col.Current = col.Target
					col.Step = 0.0 // цвет окончательно установлен
					// log.Debug().Msgf("%s DONE", farge)
				}
				// цвет еще в процессе перехода
			} else {
				// Яркость цвета уменьшалась на этом шаге
				if col.Current <= col.Target {
					col.Current = col.Target
					col.Step = 0.0 // цвет окончательно установлен
					// log.Debug().Msgf("%s DONE", farge)
				}
			}
		} else if !changedColor {
			// Heimdallr: Этот цвет сейчас не находится в процессе изменения. Мне надо придумать,
			// менять его или оставить пока. Чтобы иметь выбор, на что опереться, я возьму uint64 от уже инвертированной
			// энтропии (от choice) и еще раз это число инвертирую, только уже полностью (LostEye до этого инвертировал попиксельно).
			// Это будет моим случайным числом, по которому я определю необходимость менять цвет и уровень+скорость изменения
			invertedChoiceI := bits.Reverse64(entropy.Choice.IntValue)
			// Heimdallr: Huginn любезно согласился оставить это знание публичным, я воспользуюсь
			chance := h.huginn.UintToFloat(invertedChoiceI)
			if chance > Chance {
				// цвет надо менять
				col.Target = entropy.Choice.FloatValue
				// поменяем цвет за 300 шагов или меньше
				steps := 300.0 * entropy.Choice.FloatValue
				col.Step = (col.Target - col.Current) / steps
				log.Info().Msgf(
					"[Heimdallr] МЕНЯЮ ЦВЕТ %s БЫЛО:%.3f БУДЕТ:%.3f ШАГ:%.3f ШАГОВ:%d",
					farge,
					col.Current,
					col.Target,
					col.Step,
					int(steps),
				)
				// Odin: ловко ты это придумал. Иногда бывает тёмный цвет, но это терпимо.
				changedColor = true
			}
		}
	}
}

// encodeEntropyImageFromMatrix - картинка энтропии в base64-encode виде для передаче по json до Мидгарда
func (h *Heimdallr) encodeEntropyImageFromMatrix(matrix model.EntropyMatrix) string {
	img := h.makeEntropyImage(matrix)
	if encoded, err := h.encodeEntropyImage(img); err != nil {
		log.Fatal().Err(err).Msgf("[heimdallr] Я СОЗДАЛ ИСПОРЧЕННУЮ КАРТИНКУ. МОЙ ДОЗОР ОКОНЧЕН, АСГАРД ЗАСЫПАЕТ НА ВРЕМЯ!")
		// Odin: ой, это такой краш...
		return encoded
	} else {
		return encoded
	}
}

// encodeEntropyImageFromMatrix - картинка энтропии в base64-encode виде для передаче по json до Мидгарда
func (h *Heimdallr) encodeEntropyImage(img image.Image) (string, error) {
	b := bytes.Buffer{}
	if err := png.Encode(&b, img); err != nil {
		return "", errors.Wrap(err, "[heimdallr] ОТКАЗ PNG")
		// Odin: ой, это такой краш...
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func (h *Heimdallr) encodeJpeg(img image.Image, quality int) (string, error) {
	b := bytes.Buffer{}
	err := jpeg.Encode(&b, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return "", err
	}
	data := b.Bytes()
	// log.Debug().Msgf("[heimdallr] %d. IMAGE SIZE = %dкБ", quality, len(data)/1024)
	// Heimdallr: Картинка frame - 10-20Кб, картинка шума около 40Кб. И так несколько раз в секунду, такой траффик.
	return base64.StdEncoding.EncodeToString(data), nil
}

func (h *Heimdallr) encodeEntropyImagesForTransfer(entropy model.EntropyPack) (model.EntropyPack, error) {
	var loki bool
	var err error
	entropyImg := h.makeEntropyImage(entropy.Entropy.Matrix)
	choiceImg := h.makeEntropyImage(entropy.Choice.Matrix)

	// Loki: здесь я вмешаюсь в работу и оставлю себе бекдор, через который я буду отправлять в Мидгард свои послания
	// Loki: Хеймдалль, надеюсь на твоё молчание. Пусть это останется между нами.
	// Heimdallr: Какое дело?...

	// Heimdallr: Локи добавляет свои знания об Artchitect для людей из Мидгарда.
	loki, entropyImg, choiceImg, err = h.loki.MakeSecret(entropy.Choice.FloatValue, entropyImg, choiceImg)
	if err != nil {
		// Heimdallr: Loki не справился со своей задачей
		log.Error().Err(err).Msgf("[heimdallr] ЛОКИ НЕ СПРАВИЛСЯ СО СВОЕЙ ТАЙНОЙ МИССИЕЙ")
		// Heimdallr: но тогда отдаю оригинальные картинки
	} else if loki {
		entropy.Loki = true // Loki: значит, что я добавил свой секрет на картинку
	}

	var entropyImageEncoded, choiceImageEncoded string
	entropyImageEncoded, err = h.encodeEntropyImage(entropyImg)
	if err != nil {
		return entropy, errors.Wrap(err, "[heimdallr] ПРОБЛЕМА С УПАКОВКОЙ КАРТИНКИ ЭНТРОПИИ")
	}
	choiceImageEncoded, err = h.encodeEntropyImage(choiceImg)
	if err != nil {
		return entropy, errors.Wrap(err, "[heimdallr] ПРОБЛЕМА С УПАКОВКОЙ КАРТИНКИ ИНВЕРТИРОВАННОЙ ЭНТРОПИИ")
	}
	entropy.Entropy.ImageEncoded = entropyImageEncoded
	entropy.Choice.ImageEncoded = choiceImageEncoded
	return entropy, nil
}
