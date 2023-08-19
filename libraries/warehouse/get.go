package warehouse

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"io"
	"net/http"
)

func (w *Warehouse) DownloadArtImage(ctx context.Context, artID uint, size string) ([]byte, error) {
	url := w.artWarehouseURL + "/" + fmt.Sprintf("art/%d/%s", artID, size)
	return w.get(ctx, url)
}

func (w *Warehouse) DownloadArtOrigin(ctx context.Context, artID uint) ([]byte, error) {
	// Odin: исходник, подлинник
	url := w.originWarehouseURL + "/" + fmt.Sprintf("art/%d/origin", artID)
	return w.get(ctx, url)
}

func (w *Warehouse) DownloadUnityImage(ctx context.Context, mask string, version uint, size string) ([]byte, error) {
	url := w.artWarehouseURL + "/" + fmt.Sprintf("unity/%s/%d/%s", mask, version, size)
	return w.get(ctx, url)
}

func (w *Warehouse) GetAndDecodeArtImage(ctx context.Context, artID uint, size string) (image.Image, error) {
	b, err := w.DownloadArtImage(ctx, artID, size)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ЗАГРУЗКИ ART-ИЗОБРАЖЕНИЯ. #%d/%s", artID, size)
	}
	r := bytes.NewReader(b)
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ДЕКОДИРОВАНИЯ ART-ИЗОБРАЖЕНИЯ. #%d/%s", artID, size)
	}
	return img, nil
}
func (w *Warehouse) GetAndDecodeArtOrigin(ctx context.Context, artID uint) (image.Image, error) {
	b, err := w.DownloadArtOrigin(ctx, artID)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ЗАГРУЗКИ ORIGIN-ИЗОБРАЖЕНИЯ. #%d", artID)
	}
	r := bytes.NewReader(b)
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ДЕКОДИРОВАНИЯ ORIGIN-ИЗОБРАЖЕНИЯ. #%d", artID)
	}
	return img, nil
}
func (w *Warehouse) GetAndDecodeUnityImage(ctx context.Context, mask string, version uint, size string) (image.Image, error) {
	b, err := w.DownloadUnityImage(ctx, mask, version, size)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ЗАГРУЗКИ UNITY-ИЗОБРАЖЕНИЯ. U%s/%d/%s", mask, version, size)
	}
	r := bytes.NewReader(b)
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, errors.Wrapf(err, "[warehouse] ОШИБКА ДЕКОДИРОВАНИЯ UNITY-ИЗОБРАЖЕНИЯ. U%s/%d/%s", mask, version, size)
	}
	return img, nil
}

func (w *Warehouse) get(ctx context.Context, url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[warehouse] ЗАПРОС К СКЛАДУ ПРОВАЛЕН: %s", url)
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusNotFound {
		return []byte{}, errors.Wrapf(model.ErrNotFound, "[warehouse] 404. НЕ НАЙДЕНА КАРТИНКА. %s", url)
	} else if r.StatusCode != http.StatusOK {
		return []byte{}, errors.Errorf("[warehouse] СТАТУС ОТВЕТА ОТ СКЛАДА НЕ OK. %d:%s", r.StatusCode, r.Status)
	}
	return io.ReadAll(r.Body)
}
