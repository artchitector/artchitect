package frigg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/libraries/resizer"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

type warehouse interface {
	GetAndDecodeArtImage(ctx context.Context, artID uint, size string) (image.Image, error)
	SaveUnityCollage(ctx context.Context, mask string, version uint, img image.Image) error
}

// gungner - копьё Одина для нанесения водяных знаков (подписи) на картины.
// Frigg: Воспользуюсь силой копья Одина для того, чтобы подписать коллаж единства
type gungner interface {
	MakeUnityWatermark(img image.Image, mask string) (image.Image, error)
}

type Collage struct {
	warehouse warehouse
	gungner   gungner
}

func NewFriggCollage(warehouse warehouse, gungner gungner) *Collage {
	return &Collage{warehouse: warehouse, gungner: gungner}
}

// MakeCollage - собираем изображение NxN и переданных лидеров. Файлы загружаются из warehouse.
func (c *Collage) MakeCollage(ctx context.Context, mask string, leads []uint, maxID uint) (image.Image, error) {
	size := model.SizeS
	if len(leads) > model.CollageSize100 {
		size = model.SizeXS // Frigg: для больших коллажей так загружаться будут лишь небольшие файлы и процессор не будет сильно нагружен
	}
	var imgs []image.Image
	for _, lead := range leads {
		if lead > maxID {
			// Frigg: эта картина еще не написана. Вместо неё будет использован чёрный фон.
			if img, err := c.getBlackImage(size); err != nil {
				return nil, errors.Wrapf(err, "[frigg_collage] АВАРИЯ ЧЁРНОЙ КАРТИНКИ %s", size)
			} else {
				imgs = append(imgs, img)
			}
		} else {
			if img, err := c.warehouse.GetAndDecodeArtImage(ctx, lead, size); err != nil {
				return nil, errors.Wrapf(err, "[frigg_collage] СБОЙ ЗАГРУЗКИ ИЗОБРАЖЕНИЯ #%d/%s", lead, size)
			} else {
				imgs = append(imgs, img)
			}
		}
	}
	// Frigg: все изображение лидеров получены. Теперь их надо собрать в одну картину.
	collage := c.combineCollage(imgs)

	var err error
	collage, err = c.gungner.MakeUnityWatermark(collage, mask)
	if err != nil {
		return nil, errors.Wrapf(err, "[frigg_collage] ГУНГНИР ПОВРЕЖДЁН. МАСКА U%s", mask)
	}

	return collage, nil

}

func (c *Collage) SaveCollage(ctx context.Context, unity model.Unity, img image.Image) error {
	return c.warehouse.SaveUnityCollage(ctx, unity.Mask, unity.Version, img)
}

func (c *Collage) MakeNSaveBlankCollage(ctx context.Context, unity model.Unity) error {
	img := image.NewRGBA(image.Rect(0, 0, model.WidthF, int(model.WidthF*model.HeightToWidth)))
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)
	unityImage, err := c.gungner.MakeUnityWatermark(img, unity.Mask)
	if err != nil {
		return errors.Wrapf(err, "[frigg_collage] ПРОБЛЕМКИ С НАНЕСЕНИЕМ ПОДПИСИ НА ПУСТУЮ КАРТИНУ")
	}
	return c.SaveCollage(ctx, unity, unityImage)
}

func (c *Collage) getBlackImage(size string) (image.Image, error) {
	if dt, err := os.ReadFile(fmt.Sprintf("./files/images/black/black-%s.jpg", size)); err != nil {
		return nil, errors.Wrapf(err, "[frigg_collage] ОШИБКА ПОИСКА ЧЁРНОЙ-%s КАРТИНКИ", size)
	} else {
		img, err := jpeg.Decode(bytes.NewReader(dt))
		return img, err
	}
}

// combineCollage - сборка малых изображений в коллаж
func (c *Collage) combineCollage(imgs []image.Image) image.Image {
	size := int(math.Sqrt(float64(len(imgs))))
	first := imgs[0]
	// Frigg: формируется холст под все изображения
	collage := image.NewRGBA(
		image.Rect(0, 0, first.Bounds().Dx()*size, first.Bounds().Dy()*size),
	)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			idx := y*size + x
			point := image.Point{-x * imgs[idx].Bounds().Dx(), -y * imgs[idx].Bounds().Dy()}
			draw.Draw(collage, collage.Bounds(), imgs[idx], point, draw.Over)
		}
	}

	return resizer.ResizeImage(collage, model.SizeF)
}
