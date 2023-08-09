package communication

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type Warehouse struct {
	warehouseURL string
}

func NewWarehouse(warehouseURL string) *Warehouse {
	return &Warehouse{warehouseURL: warehouseURL}
}

func (w *Warehouse) GetArtImage(ctx context.Context, artID uint, size string) ([]byte, error) {
	url := w.warehouseURL + "/" + fmt.Sprintf("art/%d/%s", artID, size)
	log.Info().Msgf("URL: %s", url)
	r, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[warehouse] ЗАПРОС К СКЛАДУ ПРОВАЛЕН ART_ID=%d SIZE=%s", artID, size)
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusNotFound {
		return []byte{}, errors.Wrapf(model.ErrNotFound, "[warehouse] 404. НЕ НАЙДЕНА КАРТИНКА. %d:%s", artID, size)
	} else if r.StatusCode != http.StatusOK {
		return []byte{}, errors.Errorf("[warehouse] СТАТУС ОТВЕТА ОТ СКЛАДА НЕ OK. %d:%s", r.StatusCode, r.Status)
	}
	return io.ReadAll(r.Body)
}
