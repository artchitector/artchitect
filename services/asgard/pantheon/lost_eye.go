package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"image"
	"image/color"
	"math"
	"math/bits"
	"sync"
	"time"
)

const (
	SquareSize              = 64 * 7 // центральная область кадра, которая пойдёт в работу. Лишнее отбрасывается
	NoiseAmplifierRatio int = 30     // Odin: чтобы люди тоже могли увидеть, как шумит ткань пространства, я для них усилю ощущения
	// размер изображения в пикселях (энтропия 8х8 пикселей, всего 64 байта, это под uint64 число)
)

// Ворон Huginn подписывается на данные с глаза LostEye и получает энтропию по подписке на go-канал
// Других подписчиков нет, но множество их поддерживается в LostEye
// этот же subscriber использует и Huginn для похожего механизма
type leSubscriber struct {
	ctx context.Context
	ch  chan model.EntropyPackExtended
}

// LostEye - утраченный глаз Odin-а, пустой. По волшебству он всё равно видит, но видит ткань мироздания, энтропию пространства.
// Odin с помощью своих воронов Muninn и Huginn может посмотреть на эту энтропию, и найти в ней смысл - вспомнить идею будущей картины.
type LostEye struct {
	previousFrame image.Image
	currentFrame  image.Image
	mutex         sync.Mutex
	sMutex        sync.Mutex
	subscribers   []*leSubscriber
}

func NewLostEye() *LostEye {
	return &LostEye{
		previousFrame: nil,
		currentFrame:  nil,
		mutex:         sync.Mutex{},
		sMutex:        sync.Mutex{},
		subscribers:   make([]*leSubscriber, 0),
	}
}

// StartEntropyDecode - запуск процесс непрерывной расшифровки энтропии
// Odin Я не постоянно придумываю идеи (чаще жду), но мой пустой LostEye смотрит непрерывно
func (le *LostEye) StartEntropyDecode(ctx context.Context, stream chan image.Image) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[lost_eye] ПРОСМОТР ЭНТРОПИИ ОСТАНОВЛЕН")
			return
		case img := <-stream:
			saved, err := le.handleFrame(ctx, img)
			if err != nil {
				log.Error().Err(err).Msgf("[lost_eye] ГЛАЗ ДАЛ СБОЙ")
			} else if saved {
				// ok?
			}
		}
	}
}

// Subscribe - отдаёт канал, из которого подписчик читает сообщения.
// Если подписчик закрывает контекст, то отправка прерывается.
func (le *LostEye) Subscribe(subscriberCtx context.Context) chan model.EntropyPackExtended {
	ch := make(chan model.EntropyPackExtended)
	sub := leSubscriber{
		ctx: subscriberCtx,
		ch:  ch,
	}
	le.sMutex.Lock()
	defer le.sMutex.Unlock()

	le.subscribers = append(le.subscribers, &sub)
	go func() {
		<-subscriberCtx.Done()
		le.unsubscribe(&sub)
	}()
	return ch
}

func (le *LostEye) unsubscribe(sub *leSubscriber) {
	le.sMutex.Lock()
	defer le.sMutex.Unlock()

	idx := slices.IndexFunc(le.subscribers, func(s *leSubscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[lost_eye] ПОЛУЧАТЕЛЬ ИСЧЕЗ. ПРОБЛЕМА")
		return
	}

	le.subscribers = append(le.subscribers[:idx], le.subscribers[idx+1:]...)
	log.Debug().Msgf("[lost_eye] ПОЛУЧАТЕЛЬ %d УДАЛЁН. УСПЕХ.", idx)
}

func (le *LostEye) notifyListeners(
	ctx context.Context,
	pack model.EntropyPackExtended,
) {
	le.sMutex.Lock()
	subscribers := le.subscribers[:]
	le.sMutex.Unlock()

	for _, sub := range subscribers {
		// отправка энтропии всем слушателям
		go func(s *leSubscriber) {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
				log.Error().Msgf("[lost_eye] ОТПРАВКА ЗАВИСЛА, ЭНТРОПИЯ ПОТЕРЯНА")
			case s.ch <- pack:
				//log.Debug().Msgf("[lost_eye] ЭНТРОПИЯ ОТПРАВЛЕНА")
			}
		}(sub)
	}
}

func (le *LostEye) handleFrame(ctx context.Context, img image.Image) (bool, error) {
	if err := le.saveFrame(ctx, img); err != nil {
		return false, errors.Wrap(err, "[lost_eye] СОХРАНЕНИЕ КАДРА. АВАРИЯ.")
	}

	pack, ready, err := le.extractEntropy(ctx)
	if err != nil {
		return false, err
	} else if !ready {
		log.Debug().Msgf("[lost_eye] ЭНТРОПИЯ НЕ ГОТОВА")
		return false, nil
	}

	le.mutex.Lock()
	defer le.mutex.Unlock()

	le.notifyListeners(ctx, pack)
	return true, nil
}

// saveFrame - сохраняет кадр в текущее состояние LostEye, чтобы воспользоваться им далее
func (le *LostEye) saveFrame(ctx context.Context, img image.Image) error {
	var err error
	if img, err = le.extractSquare(img); err != nil {
		return err
	}

	if le.currentFrame != nil {
		le.previousFrame = le.currentFrame
	}
	le.currentFrame = img

	return nil
}

// extractSquare - вырезает ровный квадрат из кадра. С него будет считываться шум
func (le *LostEye) extractSquare(frame image.Image) (image.Image, error) {
	oldBounds := frame.Bounds()
	if oldBounds.Dx() < SquareSize || oldBounds.Dy() < SquareSize {
		return nil, errors.Errorf("[lost_eye] КВАДРАТ МАЛ ДЛЯ ВЫРЕЗА. РАЗМЕР %dх%d", oldBounds.Dx(), oldBounds.Dy())
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

func (le *LostEye) extractEntropy(ctx context.Context) (model.EntropyPackExtended, bool, error) {
	noise, ready, err := le.extractNoise()
	if err != nil {
		return model.EntropyPackExtended{}, false, err
	}
	if !ready {
		log.Debug().Msgf("[lost_eye] ДАЙТЕ ШУМА!")
		return model.EntropyPackExtended{}, false, nil
	}

	entropyMatrix := le.noiseToEntropy(noise)
	choiceMatrix := le.invertEntropy(entropyMatrix)

	pack := model.EntropyPackExtended{
		Timestamp: time.Now(),
		Entropy: model.Entropy{
			Matrix: entropyMatrix,
		},
		Choice: model.Entropy{
			Matrix: choiceMatrix,
		},
		ImageFrame: le.currentFrame,
		ImageNoise: noise,
	}

	return pack, true, nil
}

// extractNoise - вычитание цветов двух соседних кадров с целью изъять шум из этой разницы
func (le *LostEye) extractNoise() (image.Image, bool, error) {
	if le.previousFrame == nil || le.currentFrame == nil {
		return nil, false, nil
	}

	bounds := le.currentFrame.Bounds()
	noiseImage := image.NewRGBA(bounds)

	for x := 0; x <= bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			oldColor := le.previousFrame.At(x, y)
			newColor := le.currentFrame.At(x, y)
			if _, ok := oldColor.(color.RGBA); !ok {
				return nil, false, errors.New("[lost_eye] СТАРЫЙ ЦВЕТ НЕ RGBA")
			}
			if _, ok := newColor.(color.RGBA); !ok {
				return nil, false, errors.New("[lost_eye] НОВЫЙ ЦВЕТ НЕ RGBA")
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

			// ВАЖНО! Усиление и изменение цветов на шумовой картине не влияет на результат.
			// Дальше шум проходит нормализацию,
			// Абсолютные значение не так важны, всё строится на относительной светимости пикселей.
			// Цвет можно выбирать по своему вкусу и дизайну

			// Odin: Я хочу чтобы и вы это увидели, как шумит изнанка пространства, мои дороги мидгардцы!!

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

func (le *LostEye) noiseToEntropy(noise image.Image) model.EntropyMatrix {
	noiseBounds := noise.Bounds()

	proportion := noiseBounds.Dx() / model.EntropySize
	// power - интенсивность пикселя, а min и max - это шкала для нормализации всего шума внутри этой шкалы.
	// Все 64 пикселя энтропии будут нормализованы по шкале 0 до 255 (так усиление с предыдущего шага и исчезнет)
	var minPixelPower int64 = math.MaxInt64
	var maxPixelPower int64 = math.MinInt64
	powers := make([][]int64, 0, 8)

	// собираем силу каждого пикселя.
	// размер шума - 448х448, а энтропии 8х8. В каждом пикселе энтропии находится 56х56 пикселей исходного шума.
	// Среднее значение зоны 56х56 пикселей будет сжато в 1 пиксель энтропии
	for x := 0; x < model.EntropySize; x++ {
		powers = append(powers, make([]int64, 8))
		for y := 0; y < model.EntropySize; y++ {
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

	scale := maxPixelPower - minPixelPower
	matrix := model.EntropyMatrix{}

	for x := 0; x < model.EntropySize; x++ {
		for y := 0; y < model.EntropySize; y++ {
			// нормализация - из энергии текущего пикселя вычитаем энергию самого слабого пикселя,
			//так лишняя энергия будет погашена во всех пикселях
			powerOfPixel := powers[x][y] - minPixelPower
			scaledPower := math.Round(float64(powerOfPixel) / float64(scale) * math.MaxUint8)
			matrix.Set(x, y, uint8(scaledPower))
		}
	}

	// Huginn: теперь эта матрица содержит 64 числа. Это 64 пикселя на картинке энтропии и 64 бита в будущем uint64-числе
	return matrix
}

// invertEntropy - шумовая энтропия есть плавно изменяющая величина
// Но для работы Архитектора нужно иметь очень разнообразный выбор
// (выбирать разнообразные слова из словаря для одной картины, то есть даже 2 соседних кадра должны давать большой разброс значений)
// Чтобы превратить "плавную" энтропию в "резкий" выбор - картинка энтропии бинарно инвертируется (попиксельно)
func (le *LostEye) invertEntropy(entropyMatrix model.EntropyMatrix) model.EntropyMatrix {
	choiceMatrix := model.EntropyMatrix{}

	for x := 0; x < model.EntropySize; x++ {
		for y := 0; y < model.EntropySize; y++ {
			power := entropyMatrix.Get(x, y)
			power = bits.Reverse8(power)  // теперь сила этого пикселя инвертирована бинарно
			choiceMatrix.Set(x, y, power) // сохраняем эту инверсию в матрицу
		}
	}

	return choiceMatrix
}
