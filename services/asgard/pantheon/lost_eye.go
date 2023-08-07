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
	EntropyImageSide        = 8      // размер изображения в пикселях (энтропия 8х8 пикселей, всего 64 байта, это под uint64 число)
)

// Ворон Huginn подписывается на данные с глаза LostEye и получает энтропию по подписке на go-канал
// Других подписчиков нет, но множество их поддерживается в LostEye
// этот же subscriber использует и Huginn для похожего механизма
type subscriber struct {
	ctx context.Context
	ch  chan model.EntropyPack
}

// LostEye - утраченный глаз Odin-а, пустой. По волшебству он всё равно видит, но видит ткань мироздания, энтропию пространства.
// Odin с помощью своих воронов Muninn и Huginn может посмотреть на эту энтропию, и найти в ней смысл - вспомнить идею будущей картины.
type LostEye struct {
	previousFrame image.Image
	currentFrame  image.Image
	mutex         sync.Mutex
	sMutex        sync.Mutex
	subscribers   []*subscriber
}

func NewLostEye() *LostEye {
	return &LostEye{
		previousFrame: nil,
		currentFrame:  nil,
		mutex:         sync.Mutex{},
		sMutex:        sync.Mutex{},
		subscribers:   make([]*subscriber, 0),
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
func (le *LostEye) Subscribe(subscriberCtx context.Context) chan model.EntropyPack {
	ch := make(chan model.EntropyPack)
	sub := subscriber{
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

func (le *LostEye) unsubscribe(sub *subscriber) {
	le.sMutex.Lock()
	defer le.sMutex.Unlock()

	idx := slices.IndexFunc(le.subscribers, func(s *subscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[lost_eye] ПОЛУЧАТЕЛЬ ИСЧЕЗ. ПРОБЛЕМА")
		return
	}

	le.subscribers = append(le.subscribers[:idx], le.subscribers[idx+1:]...)
	log.Debug().Msgf("[lost_eye] ПОЛУЧАТЕЛЬ %d УДАЛЁН. УСПЕХ.", idx)
}

func (le *LostEye) notifyListeners(ctx context.Context, entropy image.Image, inverted image.Image) {
	le.sMutex.Lock()
	subscribers := le.subscribers[:]
	le.sMutex.Unlock()

	pack := model.EntropyPack{
		Timestamp: time.Now(),
		Entropy: model.Entropy{
			Image: entropy,
		},
		Choice: model.Entropy{
			Image: inverted,
		},
	}

	for _, sub := range subscribers {
		// отправка энтропии всем слушателям
		go func(s *subscriber) {
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

	entropy, choice, ready, err := le.extractEntropy(ctx)
	if err != nil {
		return false, err
	} else if !ready {
		log.Debug().Msgf("[lost_eye] ЭНТРОПИЯ НЕ ГОТОВА")
		return false, nil
	}

	le.mutex.Lock()
	defer le.mutex.Unlock()

	le.notifyListeners(ctx, entropy, choice)
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

func (le *LostEye) extractEntropy(ctx context.Context) (image.Image, image.Image, bool, error) {
	noise, ready, err := le.extractNoise()
	if err != nil {
		return nil, nil, false, err
	}
	if !ready {
		log.Debug().Msgf("[lost_eye] ДАЙТЕ ШУМА!")
		return nil, nil, false, nil
	}

	entropyImage := le.noiseToEntropy(noise)
	invertedEntropyImage := le.invertEntropy(entropyImage)

	return entropyImage, invertedEntropyImage, true, nil
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

func (le *LostEye) noiseToEntropy(noise image.Image) image.Image {
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

	return entropyImage
}

// invertEntropy - шумовая энтропия есть плавно изменяющая величина
// Но для работы Архитектора нужно иметь очень разнообразный выбор
// (выбирать разнообразные слова из словаря для одной картины, то есть даже 2 соседних кадра должны давать большой разброс значений)
// Чтобы превратить "плавную" энтропию в "резкий" выбор - картинка энтропии бинарно инвертируется (попиксельно)
func (le *LostEye) invertEntropy(entropy image.Image) image.Image {
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

	return choiceImage
}

func (le *LostEye) makeEntropyStruct(img image.Image, val uint64) model.Entropy {
	return model.Entropy{
		IntValue:   val,
		FloatValue: float64(val) / float64(math.MaxUint64),
		Image:      img,
	}
}
