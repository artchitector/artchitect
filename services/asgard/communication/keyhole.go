package communication

import (
	"bytes"
	"context"
	"image/jpeg"
	"net/http"

	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Keyhole - замочная скважина (http-ручка)
// Через эту замочную скважину другая внешняя система может получать кадры с веб-камеры
// Веб-камеру единолично занимает Asgard, поэтому он и выставляет эту ручку
type Keyhole struct {
	httpPort string
	webcam   webcam
}

func NewKeyhole(httpPort string, webcam webcam) *Keyhole {
	return &Keyhole{httpPort: httpPort, webcam: webcam}
}

func (kh *Keyhole) StartHttpServer(ctx context.Context) error {
	if kh.httpPort == "0" {
		log.Info().Msgf("[keyhole] НЕ ЗАПУСКАЮ HTTP-СЕРВЕР. ПОРТ=0")
		return nil
	}
	r := gin.Default()
	r.GET("/frame", kh.handle)
	if err := r.Run("0.0.0.0:" + kh.httpPort); err != nil {
		return errors.Wrap(err, "[keyhole] HTTP-СЕРВЕР УПАЛ")
	}
	return nil
}

func (kh *Keyhole) handle(c *gin.Context) {
	img, err := kh.webcam.GetNextFrame(c)
	if err != nil {
		log.Error().Err(err).Send()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	b := new(bytes.Buffer)
	if err := jpeg.Encode(b, img, &jpeg.Options{Quality: model.QualityTransfer}); err != nil {
		log.Error().Err(err).Send()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Data(http.StatusOK, "image/jpeg", b.Bytes())
}
