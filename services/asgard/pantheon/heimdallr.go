package pantheon

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
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
	huginn  *Huginn
	bifröst bifröst
}

func NewHeimdallr(huginn *Huginn, bifröst bifröst) *Heimdallr {
	return &Heimdallr{huginn: huginn, bifröst: bifröst}
}

/*
StartStream - запуск потока передачи энтропии из Huginn в communication.Bifröst
Heimdallr: брррррр... возьмёмся за дело.
Odin: этот поток критически важен, его нельзя терять. При его отключении надо "ронять" весь Asgard для его перезапуска.
Heimdallr: проооще простооого... (протяжно и безразлично)
*/
func (h *Heimdallr) StartStream(ctx context.Context) {
	entropyCh := h.huginn.Subscribe(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[heimdallr] МММММ. НАКОНЕЦ И МОЯ ОСТАНОВКА")
			return
		case entropy := <-entropyCh:
			//log.Debug().Msgf("[heimdallr] ВИЖУ ЭНТРОПИЮ ОТ ХУГИНА %+v", entropy)

			// Heimdallr: так. Моя задача из матрицы силы пикселей сделать видимую картинку
			var err error
			entropy.Entropy.ImageEncoded = h.encodeEntropyImage(entropy.Entropy.Matrix)
			entropy.Choice.ImageEncoded = h.encodeEntropyImage(entropy.Choice.Matrix)
			entropy.ImageFrameEncoded, err = h.encodeJpeg(entropy.ImageFrame, model.EntropyJpegQualityFrame)
			if err != nil {
				log.Error().Msgf("[heimdallr] НЕ СМОГ СДЕЛАТЬ JPEG ЭНТРОПИИ")
			}
			entropy.ImageNoiseEncoded, err = h.encodeJpeg(entropy.ImageNoise, model.EntropyJpegQualityNoise)
			if err != nil {
				log.Error().Msgf("[heimdallr] НЕ СМОГ СДЕЛАТЬ JPEG ОБРАТНОЙ ЭНТРОПИИ")
			}

			if err = h.sendDrakkar(ctx, model.ChanEntropyExtended, entropy); err != nil {
				log.Error().Msgf("[heimdallr] BIFRÖST СЛОМАН, ДРАККАР С ГРУЗОМ ЭНТРОПИИ УТЕРЯН В ТКАНИ ПРОСТРАНСТВА")
				// Heimdallr: поток не прерываю. Проблема может быть сетевая из-за Redis
				continue
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
			if err = h.sendDrakkar(ctx, model.ChanEntropy, miniEntropy); err != nil {
				log.Error().Msgf("[heimdallr] BIFRÖST СЛОМАН, ДРАККАР С ГРУЗОМ ЭНТРОПИИ УТЕРЯН В ТКАНИ ПРОСТРАНСТВА")
				// Heimdallr: поток не прерываю. Проблема может быть сетевая из-за Redis
				continue
			}
		}
	}
}

func (h *Heimdallr) sendDrakkar(ctx context.Context, channel string, pack interface{}) error {
	// Heimdallr: теперь отправлю этот ценный драккар по волнам Биврёста...
	var b []byte
	var err error
	if b, err = json.Marshal(&pack); err != nil {
		return errors.Wrap(err, "[heimdallr] ЭНТРОПИЯ ИСПОРЧЕНА. БЛЭКАУТ!")
	}
	err = h.bifröst.SendDrakkar(ctx, model.Cargo{
		Channel: channel,
		Payload: string(b),
	})
	return err
}

// MakeEntropyImage - матрица сил пикселей в энтропии становится 8х8 PNG картинкой
func (h *Heimdallr) MakeEntropyImage(matrix model.EntropyMatrix) image.Image {
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
// Heimdallr: я сделаю золотой
func (h *Heimdallr) makeColor(power uint8) color.Color {
	greenToRedProportion := 165.0 / 255.0 // для золотого цвета
	return color.RGBA{
		R: power,
		G: uint8(greenToRedProportion * float64(power)),
		B: 0,
		A: math.MaxUint8,
	}
}

// encodeEntropyImage - картинка энтропии в base64-encode виде для передаче по json до Мидгарда
func (h *Heimdallr) encodeEntropyImage(matrix model.EntropyMatrix) string {
	img := h.MakeEntropyImage(matrix)
	b := bytes.Buffer{}
	if err := png.Encode(&b, img); err != nil {
		log.Fatal().Err(err).Msgf("[heimdallr] Я СОЗДАЛ ИСПОРЧЕННУЮ КАРТИНКУ. МОЙ ДОЗОР ОКОНЧЕН, АСГАРД ЗАСЫПАЕТ НА ВРЕМЯ!")
		// Odin: ой, это такой краш...
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (h *Heimdallr) encodeJpeg(img image.Image, quality int) (string, error) {
	b := bytes.Buffer{}
	err := jpeg.Encode(&b, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return "", err
	}
	data := b.Bytes()
	log.Info().Msgf("[heimdallr] %d. IMAGE SIZE = %dкБ", quality, len(data)/1024)
	return base64.StdEncoding.EncodeToString(data), nil

}
