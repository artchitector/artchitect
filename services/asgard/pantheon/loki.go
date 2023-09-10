package pantheon

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
	"time"
)

const LokiChance = 1.0 / 2000.0 // шанс на появление
var frameDuration = time.Millisecond * 500
var framesScheme = [][]int{
	// {0, 0} - это обозначает, что на энтропию и на choice будут наложена картинка с индексом ноль из листа (улыбка №1)
	// злая улыбка
	{0, 0},
	{1, 1},
	{2, 2},
	{3, 3},

	// повтор злой улыбки
	{0, 0},
	{1, 1},
	{2, 2},
	{3, 3},
	// пауза
	{15, 15}, // Loki - 15й кадр - просто пустое затемнение

	// LO VE LO KI
	{8, 9},
	{10, 11},
	{15, 15},
	{4, 5},
	{6, 7},
	{15, 15},

	// сердце 3 кадра
	{12, 12},
	{13, 13},
	{14, 14},
	{15, 15},

	// LO VE LO KI
	{8, 9},
	{10, 11},
	{15, 15},
	{4, 5},
	{6, 7},
	{15, 15},

	// сердце 3 кадра
	{12, 12},
	{13, 13},
	{14, 14},
	{15, 15},

	// повтор злой улыбки
	{0, 0},
	{1, 1},
	{2, 2},
	{3, 3},
}

// Loki
// Loki: у меня тоже есть своё место в этом представлении. Я не буду делать какую-то полезную работу для Artchitect,
// Loki: но я хочу рассказать историю, как всё начиналось. Я вставлю пасхалку в данные энтропии, и особо внимательные
// Loki: зрители из Мидгарда смогут её увидеть, перейти по ссылке /loki и прочитать нашу с Odin историю.
// Loki: я буду ИЗРЕДКА накладывать свои изображения поверх энтропии в Мидгард,
// и особо внимательные и удачливые увидят мой секрет и увидят ссылку
type Loki struct {
	active             bool
	currentFrame       int
	currentFrameUpdate time.Time
	framesSheet        image.Image
}

func NewLoki() *Loki {
	return &Loki{}
}

func (l *Loki) MakeSecret(choiceVal float64, entropy image.Image, choice image.Image) (bool, image.Image, image.Image, error) {
	// Loki: обычно картинка энтропии транслируется как есть, но я в редкие моменты буду вмешиваться и поверх неё
	// Loki: и дорисовывать своё послание в Мидгард, прямо в эту картинку 8х8 пикселей. Послание лежит в файле
	// Loki: asgard/files/images/loki.png. Там изображение из 16 кадров, которые последовательно будут накладываться на поток энтропии

	if err := l.prepare(); err != nil {
		log.Fatal().Err(err).Msgf("[loki] ЗАВЕРШАЮ РАБОТУ АСГАРД. РЕСУРСЫ НЕДОСТУПНЫ. ВСЁ КОНЧЕНО.")
		return false, nil, nil, err
	}

	if !l.active {
		if l.needActivation(choiceVal) {
			l.active = true
			l.currentFrame = 0
			l.currentFrameUpdate = time.Now()
		} else {
			// Loki: ничего не меняю
			return false, entropy, choice, nil
		}
	}

	if l.active {
		if time.Now().Sub(l.currentFrameUpdate) > frameDuration {
			l.currentFrame += 1
			l.currentFrameUpdate = time.Now()
			if l.currentFrame >= len(framesScheme) {
				l.active = false
				// Loki: больше кадров нет, цепочка завершена
				return false, entropy, choice, nil
			}
		}
		entropy, choice = l.overlay(entropy, choice)
		return true, entropy, choice, nil
	}

	return false, entropy, choice, nil
}

func (l *Loki) needActivation(choice float64) bool {
	min := 0.5 - LokiChance/2
	max := 0.5 + LokiChance/2
	if choice >= min && choice <= max {
		log.Info().Msgf("[loki] АКТИВИРУЮ ЛОКИ. CHOICE: %f", choice)
		return true
	}

	return false
}

func (l *Loki) overlay(entropyOrig image.Image, choiceOrig image.Image) (image.Image, image.Image) {
	entropy := image.NewRGBA(entropyOrig.Bounds())
	draw.Draw(entropy, entropy.Bounds(), entropyOrig, image.Pt(0, 0), draw.Over)
	choice := image.NewRGBA(choiceOrig.Bounds())
	draw.Draw(choice, choice.Bounds(), choiceOrig, image.Pt(0, 0), draw.Over)

	first, second := l.getImagesToOverlay()
	draw.Draw(entropy, entropy.Bounds(), first, image.Pt(first.Bounds().Min.X, first.Bounds().Min.Y), draw.Over)
	draw.Draw(choice, choice.Bounds(), second, image.Pt(second.Bounds().Min.X, second.Bounds().Min.Y), draw.Over)

	return entropy, choice
}

func (l *Loki) getImagesToOverlay() (image.Image, image.Image) {
	var firstX, firstY, secondX, secondY int
	currentImages := framesScheme[l.currentFrame]
	firstIndex := currentImages[0]
	firstX = firstIndex % 4
	firstY = firstIndex / 4
	secondIndex := currentImages[1]
	secondX = secondIndex % 4
	secondY = secondIndex / 4
	// Loki: карта кадров (image sheet) выполнена в виде сетки 4х4 (каждое изображение 8х8 пикселей)
	// Loki: нужно выбрать нужные x и y для этой картинки
	first := l.framesSheet.(subImager).SubImage(image.Rect(firstX*8, firstY*8, firstX*8+8, firstY*8+8))
	second := l.framesSheet.(subImager).SubImage(image.Rect(secondX*8, secondY*8, secondX*8+8, secondY*8+8))

	return first, second
}

func (l *Loki) prepare() error {
	if l.framesSheet != nil {
		return nil // уже подготовлено
	}
	lokiData, err := os.ReadFile("./files/images/loki.png")
	if err != nil {
		return err
	}
	r := bytes.NewReader(lokiData)
	img, err := png.Decode(r)
	if err != nil {
		return err
	}
	l.framesSheet = img
	return nil
}
