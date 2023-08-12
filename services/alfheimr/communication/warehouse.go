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
	artWarehouseURL    string
	originWarehouseURL string
}

func NewWarehouse(artWarehouseURL string, originWarehouseURL string) *Warehouse {
	return &Warehouse{artWarehouseURL: artWarehouseURL, originWarehouseURL: originWarehouseURL}
}

func (w *Warehouse) GetArtImage(ctx context.Context, artID uint, size string) ([]byte, error) {
	return w.get(ctx, w.artWarehouseURL, artID, size)
}

func (w *Warehouse) GetArtOrigin(ctx context.Context, artID uint) ([]byte, error) {
	// Odin: исходник, подлинник, в png-формате
	return w.get(ctx, w.originWarehouseURL, artID, "origin")
}

func (w *Warehouse) get(ctx context.Context, warehouseURL string, artID uint, size string) ([]byte, error) {
	url := warehouseURL + "/" + fmt.Sprintf("art/%d/%s", artID, size)
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
