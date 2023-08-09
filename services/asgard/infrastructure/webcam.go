package infrastructure

import (
	"bytes"
	"context"
	"github.com/blackjack/webcam"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const ReadTimeoutSeconds = 5
const V4L2_PIX_FMT_YUYV = 0x56595559 // https://github.com/blackjack/webcam/blob/master/examples/http_mjpeg_streamer/webcam.go

// Webcam - камера, которая считывает кадры пространства. Технически камера и есть пустой глаз Одина в Artchitect.
// pantheon.Odin: мне приходится смотреть в эту ограниченную электронную коробочку с разрешением 640х480,
// pantheon.Odin: которая еще и направлена на маленькую иконку Иисуса. Безобразие! Издевательство надо мной!
// pantheon.Odin: Это проделки Локи, не иначе - его чувство юмора. Каждый раз думаю, зачем я на это согласился...
type Webcam struct {
	deviceID   string
	resolution string
	started    bool
	mutex      sync.Mutex
	// rawSingleReceivers - разовое получение изначального кадра внешним подписчиком
	rawSingleReceivers []chan image.Image
}

func NewWebcam(deviceID string, resolution string) *Webcam {
	return &Webcam{deviceID: deviceID, resolution: resolution}
}

func (w *Webcam) Start(ctx context.Context, outputCh chan image.Image) error {
	if strings.Index(w.deviceID, "http") == 0 {
		//w.StartHttpTransfer(ctx, outputCh)
		return nil
	}

	if w.started == true {
		return errors.Errorf("[webcam] ДОСТУПЕН ЛИШЬ 1 ПОТОК КИНОКАМЕРЫ")
	}

	// инициализация камеры
	cam, err := webcam.Open(w.deviceID)
	if err != nil {
		return errors.Wrap(err, "[webcam] КИНОКАМЕРА - НЕТ СВЯЗИ")
	}
	defer func() {
		if err := cam.Close(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	// настройки формата и разрешения
	format, width, height, err := w.getFormatAndSize(cam)
	if err != nil {
		return errors.Wrap(err, "[webcam] НЕТ НАСТРОЕК КАДРА")
	}
	camFormat, camWidth, camHeight, err := cam.SetImageFormat(format, uint32(width), uint32(height))
	if err != nil {
		return errors.Wrap(err, "[webcam] КИНОКАМЕРА НЕ НАСТРОЕНА")
	}

	log.Info().Msgf("[webcam] СТАРТ ПОТОКА КИНОКАМЕРЫ")
	// начало стрима  кадров с веб-камеры
	if err := cam.StartStreaming(); err != nil {
		return errors.Wrap(err, "[webcam] НЕТ ПОТОКА КИНОКАМЕРЫ")
	}

	// для остановки чтения с камеры по контекст.Done
	done := false
	go func() {
		<-ctx.Done()
		log.Debug().Msgf("[webcam] ОСТАНОВКА ПОТОКА КИНОКАМЕРЫ")
		done = true
	}()

	w.started = true

	// цикл чтения кадров с камеры
	for {
		if done {
			break
		}
		err := cam.WaitForFrame(uint32(ReadTimeoutSeconds))
		switch err.(type) {
		case nil:
		// ok
		case *webcam.Timeout:
			log.Error().Err(err).Msgf("[webcam] ВРЕМЯ ЧТЕНИЯ КАДРА ВЫШЛО")
			continue
		default:
			return errors.Wrap(err, "[webcam] ОШИБКА ОШИДАНИЯ КАДРА")
		}

		if done {
			break
		}
		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			img, err := w.encodeFrame(frame, camWidth, camHeight, camFormat)
			if err != nil {
				log.Error().Err(err).Msgf("[webcam] ОШИБКА ДЕКОДИРОВАНИЯ КАДРА")
				continue
			}

			go func(img image.Image) {
				// отправляем кадр в singleReceiver
				if len(w.rawSingleReceivers) == 0 {
					return
				}
				w.mutex.Lock()
				currentSingleReceiver := w.rawSingleReceivers[0]
				w.rawSingleReceivers = w.rawSingleReceivers[1:]
				w.mutex.Unlock()
				select {
				case <-time.After(time.Millisecond * 10):
					log.Warn().Msgf("[webcam] ЗАКРЫТИЕ КАНАЛА ЕДИНИЧНОГО ПОЛУЧАТЕЛЯ ПО ТАЙМАУТУ")
					close(currentSingleReceiver)
					return
				case currentSingleReceiver <- img:
					// ok
				}
			}(img)

			// отправка кадра получателю
			select {
			case <-ctx.Done():
				log.Debug().Msgf("[webcam] ОСТАНАВЛИВАЮ КАМЕРУ")
				break
			case <-time.After(time.Millisecond * 10):
				// ok
			case outputCh <- img:
				// ok
			}

		} else if err != nil {
			log.Error().Err(err).Msgf("[webcam] ОШИБКА ЧТЕНИЯ КАДРА")
		}
	}

	log.Debug().Msgf("[webcam] КАМЕРА ОСТАНОВЛЕНА")
	return nil
}

/*
StartHttpTransfer - для локальной разработки. Кадры для энтропии вытаскиваются из ручки, которую держит production-Асгард
*/
func (w *Webcam) StartHttpTransfer(ctx context.Context, outputCh chan image.Image) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 10):
			func() {
				r, err := http.Get(w.deviceID)
				if err != nil {
					log.Error().Err(err).Msgf("[webcam] НЕ МОГУ ЗАГРУЗИТЬ КАДР ПО HTTP")
					return
				}
				defer r.Body.Close()
				if r.StatusCode != http.StatusOK {
					log.Error().Err(err).Msgf("[webcam] HTTP-СТАТУС ОТВЕТА НЕ ОК %d:%s", r.StatusCode, r.Status)
					return
				}
				imgData, err := io.ReadAll(r.Body)
				if err != nil {
					log.Error().Err(err).Msgf("[webcam] НЕВОЗМОЖНО ПРОЧИТАТЬ ТЕЛО ОТВЕТА")
					return
				}
				b := bytes.NewReader(imgData)
				img, err := jpeg.Decode(b)
				if err != nil {
					log.Error().Err(err).Msgf("[webcam] НЕВОЗМОЖНО РАСПОЗНАТЬ JPEG")
					return
				}
				outputCh <- img
			}()
		}
	}
}

func (w *Webcam) GetNextFrame(ctx context.Context) (image.Image, error) {
	singleReceiver := make(chan image.Image)
	w.mutex.Lock()
	w.rawSingleReceivers = append(w.rawSingleReceivers, singleReceiver)
	w.mutex.Unlock()
	select {
	case <-ctx.Done():
		return nil, errors.Errorf("[webcam] ОСТАНОВ")
	case <-time.After(time.Second):
		return nil, errors.Errorf("[webcam] ТАЙМАУТ")
	case img := <-singleReceiver:
		return img, nil
	}
}

func (w *Webcam) getFormatAndSize(cam *webcam.Webcam) (webcam.PixelFormat, uint32, uint32, error) {
	// получим формат и размеры изображения с веб-камеры
	var selectedFormat webcam.PixelFormat
	var selectedFormatDesc string
	supportedFormats := cam.GetSupportedFormats()
	if len(supportedFormats) == 0 {
		return 0, 0, 0, errors.New("[webcam] НЕТ ДОСТУПНЫХ ФОРМАТОВ ВИДЕО")
	}
	for f, s := range supportedFormats {
		// Берём первый попавшийся формат для кадра. В итоге он будет превращён в png картинку
		selectedFormat = f
		selectedFormatDesc = s
		log.Debug().Msgf("[webcam] ИСПОЛЬЗУЮ ФОРМАТ - %s", s)
		break
	}

	var selectedFrameSize *webcam.FrameSize
	frames := cam.GetSupportedFrameSizes(selectedFormat)
	for _, f := range frames {
		if f.GetString() == w.resolution {
			selectedFrameSize = &f
			log.Debug().Msgf("[webcam] ИСПОЛЬЗУЮ РАЗМЕР - %s", f.GetString())
			break
		}
	}

	log.Debug().Msgf("[webcam] НАСТРОЙКИ КИНОКАМЕРЫ: %s %dx%d", selectedFormatDesc, selectedFrameSize.MaxWidth, selectedFrameSize.MaxHeight)
	return selectedFormat, selectedFrameSize.MaxWidth, selectedFrameSize.MaxHeight, nil
}

func (w *Webcam) encodeFrame(frame []byte, width, height uint32, format webcam.PixelFormat) (image.Image, error) {
	var img image.Image
	switch format {
	case V4L2_PIX_FMT_YUYV:
		yuyv := image.NewYCbCr(image.Rect(0, 0, int(width), int(height)), image.YCbCrSubsampleRatio422)
		for i := range yuyv.Cb {
			ii := i * 4
			yuyv.Y[i*2] = frame[ii]
			yuyv.Y[i*2+1] = frame[ii+2]
			yuyv.Cb[i] = frame[ii+1]
			yuyv.Cr[i] = frame[ii+3]

		}
		img = yuyv
	default:
		return nil, errors.Errorf("[webcam] ФОРМАТ НЕ ПОДДЕРЖИВАЕТСЯ %v", format)
	}

	return img, nil
}
