package external

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/bits"
	"os"
	"sync"
	"time"
)

const (
	SquareSize              = 64 * 7
	NoiseAmplifierRatio int = 30
	EntropyImageSide        = 8 // размер изображения в пикселях (энтропия 8х8 пикселей, всего 64 байта, это под uint64 число)
)

type EntropyPack struct {
	Entropy model.Entropy
	Choice  model.Entropy
}

type subscriber struct {
	ctx context.Context
	ch  chan EntropyPack
}

type Entropy struct {
	previousFrame        image.Image
	currentFrame         image.Image
	mutex                sync.Mutex
	lastExtractedEntropy *EntropyPack
	sMutex               sync.Mutex
	subscribers          []*subscriber
}

func NewEntropy() *Entropy {
	return &Entropy{
		previousFrame:        nil,
		currentFrame:         nil,
		mutex:                sync.Mutex{},
		lastExtractedEntropy: nil,
		sMutex:               sync.Mutex{},
		subscribers:          make([]*subscriber, 0),
	}
}

func (e *Entropy) StartEntropyDecode(ctx context.Context, stream chan image.Image) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[entropy] ОСТАНОВ")
			return
		case img := <-stream:
			saved, err := e.handleEvent(ctx, img)
			if err != nil {
				log.Error().Err(err).Msgf("[entropy] ОШИБКА ОБРАБОТКИ")
			} else if saved {
				// ok?
			}
		}
	}
}

func (e *Entropy) SubscribeEntropy(subscriberCtx context.Context) chan EntropyPack {
	ch := make(chan EntropyPack)
	sub := subscriber{
		ctx: subscriberCtx,
		ch:  ch,
	}
	e.sMutex.Lock()
	defer e.sMutex.Unlock()

	e.subscribers = append(e.subscribers, &sub)
	go func() {
		<-subscriberCtx.Done()
		e.unsubscribe(&sub)
	}()
	return ch
}

func (e *Entropy) unsubscribe(sub *subscriber) {
	e.sMutex.Lock()
	defer e.sMutex.Unlock()

	idx := slices.IndexFunc(e.subscribers, func(s *subscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[entropy] КАНАЛ ОТСУТСТВУЕТ. ПРОБЛЕМА")
		return
	}

	e.subscribers = append(e.subscribers[:idx], e.subscribers[idx+1:]...)
	log.Debug().Msgf("[entropy] ПОДПИСЧИК %d - УДАЛЕНО. УСПЕХ.", idx)
}

func (e *Entropy) notifyListeners(ctx context.Context, entropy EntropyPack) {
	e.sMutex.Lock()
	subscribers := e.subscribers[:]
	e.sMutex.Unlock()

	for _, sub := range subscribers {
		// отправка энтропии всем слушателям
		go func(s *subscriber) {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
				log.Error().Msgf("[entropy] ОТПРАВКА ЗАВИСЛА, ГРУЗ ПОТЕРЯН")
			case s.ch <- entropy:
				log.Debug().Msgf("[entropy] ГРУЗ ЭНТРОПИИ ОТПРАВЛЕН")
			}
		}(sub)
	}
}

func (e *Entropy) handleEvent(ctx context.Context, img image.Image) (bool, error) {
	if err := e.saveFrame(ctx, img); err != nil {
		return false, errors.Wrap(err, "[entropy] СОХРАНЕНИЕ КАДРА. АВАРИЯ.")
	}

	entropy, choice, ready, err := e.extractEntropy(ctx)
	if err != nil {
		return false, err
	} else if !ready {
		log.Debug().Msgf("[entropy] ЭНТРОПИЯ НЕ ГОТОВА")
		return false, nil
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.lastExtractedEntropy = &EntropyPack{
		Entropy: entropy,
		Choice:  choice,
	}
	log.Debug().Msgf("[entropy] ЭНТРОПИЯ E:%s C:%s", entropy, choice)
	e.notifyListeners(ctx, *e.lastExtractedEntropy)
	return true, nil
}

// saveFrame - сохраняет кадр в текущее состояние, чтобы воспользоваться им далее
func (e *Entropy) saveFrame(ctx context.Context, img image.Image) error {
	var err error
	if img, err = e.extractSquare(img); err != nil {
		return err
	}

	if e.currentFrame != nil {
		e.previousFrame = e.currentFrame
	}
	e.currentFrame = img

	return nil
}

func (e *Entropy) extractSquare(frame image.Image) (image.Image, error) {
	oldBounds := frame.Bounds()
	if oldBounds.Dx() < SquareSize || oldBounds.Dy() < SquareSize {
		return nil, errors.Errorf("[entropy] КВАДРАТ МАЛ ДЛЯ ВЫРЕЗА. РАЗМЕР %dх%d", oldBounds.Dx(), oldBounds.Dy())
	}
	squareRect := image.Rect(0, 0, SquareSize, SquareSize)
	squareImg := image.NewRGBA(squareRect)

	leftOffset := (oldBounds.Dx() - squareRect.Dx()) / 2
	topOffset := (oldBounds.Dy() - squareRect.Dy()) / 2

	for x := 0; x < SquareSize; x++ {
		for y := 0; y < SquareSize; y++ {
			squareImg.Set(x, y, frame.At(x+leftOffset, y+topOffset))
		}
	}

	return squareImg, nil
}

func (e *Entropy) extractEntropy(ctx context.Context) (model.Entropy, model.Entropy, bool, error) {
	noise, ready, err := e.extractNoise()
	if err != nil {
		return model.Entropy{}, model.Entropy{}, false, err
	}
	if !ready {
		log.Debug().Msgf("[entropy] ШУМ НЕ ГОТОВ")
		return model.Entropy{}, model.Entropy{}, false, nil
	}

	entropyImage, entropyVal := e.noiseToEntropy(noise)
	entropyObj := e.makeEntropyStruct(model.EntropyTypeDirect, entropyImage, entropyVal)

	choiceImage, choiceVal := e.invertEntropy(entropyImage)
	choiceObj := e.makeEntropyStruct(model.EntropyTypeDirect, choiceImage, choiceVal)

	return entropyObj, choiceObj, true, nil
}

// extractNoise - вычитание цветов двух соседних кадров с целью изъять шум из этой разницы
func (e *Entropy) extractNoise() (image.Image, bool, error) {
	if e.previousFrame == nil || e.currentFrame == nil {
		return nil, false, nil
	}

	bounds := e.currentFrame.Bounds()
	noiseImage := image.NewRGBA(bounds)

	for x := 0; x <= bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			oldColor := e.previousFrame.At(x, y)
			newColor := e.currentFrame.At(x, y)
			if _, ok := oldColor.(color.RGBA); !ok {
				return nil, false, errors.New("[entropy] СТАРЫЙ ЦВЕТ НЕ RGBA")
			}
			if _, ok := newColor.(color.RGBA); !ok {
				return nil, false, errors.New("[entropy] НОВЫЙ ЦВЕТ НЕ RGBA")
			}

			// новый цвет, его RBG каналы
			var newR, newG, newB int
			newR = int(newColor.(color.RGBA).R) - int(oldColor.(color.RGBA).R)
			if newR < 0 {
				newR *= -1
			}
			newG = int(newColor.(color.RGBA).G) - int(oldColor.(color.RGBA).G)
			if newG < 0 {
				newG *= -1
			}
			newB = int(newColor.(color.RGBA).B) - int(oldColor.(color.RGBA).B)
			if newB < 0 {
				newB *= -1
			}

			// ВАЖНО! Усиление и изменение цветов на шумовой картине не влияет на результат. Дальше шум проходит нормализацию,
			// Абсолютные значение не так важны, всё строится на относительной светимости пикселей.
			// Цвет можно выбирать по своему вкусу и дизайну

			/*
				УСИЛЕНИЕ ЦВЕТА
				Это требуется лишь для наглядности, что шум есть. Будет видно пользователю.
				На следующем шаге весь шум будет перенормализован, поэтому усиление уже не будет играть роль.
				Относительность пикселей сохранится.
			*/
			noiseColor := color.RGBA{
				R: uint8(255) - uint8(newR*NoiseAmplifierRatio),
				G: uint8(255) - uint8(newG*NoiseAmplifierRatio),
				B: uint8(255) - uint8(newB*NoiseAmplifierRatio),
				A: 255,
			}
			noiseImage.SetRGBA(x, y, noiseColor)
		}
	}

	return noiseImage, true, nil
}

func (e *Entropy) noiseToEntropy(noise image.Image) (image.Image, uint64) {
	noiseBounds := noise.Bounds()
	entropyBounds := image.Rect(0, 0, EntropyImageSide, EntropyImageSide)
	entropyImage := image.NewRGBA(entropyBounds)

	proportion := noiseBounds.Dx() / entropyBounds.Dx()
	// power - интенсивность пикселя, а min и max - это шкала для нормализации всего шума внутри этой шкалы.
	// Все 64 пикселя энтропии будут нормализованы по шкале 0 до 255 (так усиление с предыдущего шага и исчезнет)
	var minPixelPower int64 = math.MaxInt64
	var maxPixelPower int64 = math.MinInt64
	powers := make([][]int64, 0, 8)

	// собираем силу каждого пикселя.
	// размер шума - 448х448, а энтропии 8х8. В каждом пикселе энтропии находится 56х56 пикселей исходного шума.
	// Среднее значение зоны 56х56 пикселей будет сжато в 1 пиксель энтропии
	for x := 0; x < entropyBounds.Dx(); x++ {
		powers = append(powers, make([]int64, 8))
		for y := 0; y < entropyBounds.Dy(); y++ {
			var powerOfPixel int64
			for nx := x * proportion; nx < x*proportion+proportion; nx++ {
				for ny := y * proportion; ny < y*proportion+proportion; ny++ {
					clr := noise.At(nx, ny)
					powerOfPixel += int64(clr.(color.RGBA).R)
					powerOfPixel += int64(clr.(color.RGBA).G)
					powerOfPixel += int64(clr.(color.RGBA).B)
				}
			}
			if powerOfPixel < minPixelPower {
				minPixelPower = powerOfPixel
			}
			if powerOfPixel > maxPixelPower {
				maxPixelPower = powerOfPixel
			}
			powers[x][y] = powerOfPixel
		}
	}

	var entropyValue uint64
	scale := maxPixelPower - minPixelPower
	for x := 0; x < entropyBounds.Dx(); x++ {
		for y := 0; y < entropyBounds.Dy(); y++ {
			// нормализация - из энергии текущего пикселя вычитаем энергию самого слабого пикселя,
			//так лишняя энергия будет погашена во всех пикселях
			powerOfPixel := powers[x][y] - minPixelPower

			redPower := math.Round(float64(powerOfPixel) / float64(scale) * math.MaxUint8)

			if redPower >= float64(math.MaxUint8)/2.0 {
				// картинка 8х8 пикселей тут превращается в uint64 число. Каждый пиксель - это бит числа.
				// если пиксель ближе к красному цвету по шкале - то это включенный бит uint64-числа - 1
				// если пиксель ближе к чёрному цвету по шкале - то это выключенный бит uint64-числа - 0
				// 64 пикселя = 64 бита = uint64 число
				byteIndex := x*EntropyImageSide + y
				entropyValue = entropyValue | 1<<(63-byteIndex)
			}

			entropyImage.SetRGBA(x, y, color.RGBA{
				R: uint8(redPower),
				G: 0,
				B: 0,
				A: 255,
			})
		}
	}

	return entropyImage, entropyValue
}

// invertEntropy - шумовая энтропия есть плавно изменяющая величина
// Но для работы Архитектора нужно иметь очень разнообразный выбор
// (выбирать разнообразные слова из словаря для одной картины, то есть даже 2 соседних кадра должны давать большой разброс значений)
// Чтобы превратить "плавную" энтропию в "резкий" выбор - картинка энтропии бинарно инвертируется (попиксельно)
func (e *Entropy) invertEntropy(entropy image.Image) (image.Image, uint64) {
	bounds := entropy.Bounds()
	choiceImage := image.NewRGBA(bounds)
	var choiceValue uint64

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			col := entropy.At(x, y)
			power := col.(color.RGBA).R                         // сила пикселя от 0 до 255
			power = bits.Reverse8(power)                        // теперь сила этого пикселя инвертирована бинарно
			choiceImage.Set(x, y, color.RGBA{R: power, A: 255}) // на картинке ставим инвертированный пиксель для наглядности

			// считаем включенные биты с картинки и готовим uint64-число
			isEnabledPixel := power >= math.MaxUint8/2
			if isEnabledPixel {
				byteIndex := x*EntropyImageSide + y
				choiceValue = choiceValue | 1<<(63-byteIndex)
			}
		}
	}

	return choiceImage, choiceValue
}

func (e *Entropy) makeEntropyStruct(entropyType string, img image.Image, val uint64) model.Entropy {
	return model.Entropy{
		Type:       entropyType,
		IntValue:   val,
		FloatValue: float64(val) / float64(math.MaxUint64),
		Image:      img,
	}
}

func (e *Entropy) saveImage(filename string, img image.Image) error {
	b := bytes.Buffer{}
	err := jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b.Bytes(), 0644)
	return err
}
