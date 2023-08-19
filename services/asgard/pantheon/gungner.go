package pantheon

import (
	"bytes"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"os"
)

/*
Gungner - копьё Одина.
Копьё было изготовлено двумя карликами-свартальвами, сыновьями Ивальди, чтобы показать богам мастерство подземного народа.
Odin: Так как Я истинный автор каждой картины, то я буду оставлять на ней свою подпись.
Odin: Я буду высекать символы своим копьём Gungner.
Odin: На каждую картину я нанесу отпечаток с её номером.
*/
type Gungner struct {
	font   *truetype.Font
	catImg image.Image
}

func NewGungner() *Gungner {
	return &Gungner{}
}

func (g *Gungner) MakeArtWatermark(img image.Image, artID uint) (image.Image, error) {
	// Odin: ХОЧУ, чтобы на каждой картине была подпись с номером картины, а рядом с ним был КОТ!
	// Freyja: Есть подходящий кот на картине "Есть ли кошачий Бог?" за авторством Artchitect опытной первой версии.
	// Odin: именно этого кота я и хотел увидеть.
	/*
		[artchitector]: для любопытствующих смотреть файл services/asgard/files/images/is_there_cat_god.jpg.
		Это была первая отладка Artchitect, рисовалось всё без энтропии. Это было самое начало проекта Artchitect.
		И с этой картины был вырезан кот, который появляется в углу каждой картины.
		Card #1294.
		Created: 2023 Jan 6 18:31
		Seed: 4091966908
		Words: intricate,cat,Sun,galactic,nuclear,symmetrical,Allah,girl,stunning beautiful,europe,dynamic lighting,
			greek,darkblue,art,sadness,light,fantastically beautiful,red,Gothic,train,john constable,textured,yellow,
			tribal patterns,hyper,high details,electricity
	*/
	if err := g.loadResources(); err != nil {
		return nil, errors.Wrap(err, "[gungner] ТРЕБУЕМЫЕ РЕСУРСЫ НЕ ЗАГРУЖЕНЫ")
	}
	text := fmt.Sprintf("#%d", artID)
	return g.addWatermark(img, text), nil
}

func (g *Gungner) MakeUnityWatermark(img image.Image, mask string) (image.Image, error) {
	// Odin: Логика для нанесения номера единства повторяет нанесение номера картины
	if err := g.loadResources(); err != nil {
		return nil, errors.Wrap(err, "[gungner] ТРЕБУЕМЫЕ РЕСУРСЫ НЕ ЗАГРУЖЕНЫ")
	}
	text := fmt.Sprintf("U%s", mask)
	return g.addWatermark(img, text), nil
}

func (g *Gungner) addWatermark(img image.Image, text string) image.Image {
	// Odin: начинаю наносить подпись на картину своим копьём!
	//rgba := img.(*image.RGBA)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// Odin: готовлю свою подпись
	watermark := g.prepareWatermark(rgba.Bounds(), text)
	// Odin: выбираю место, куда нанесу подпись
	padding := fixed.I(rgba.Bounds().Max.X / 30) // Odin: отступ = 1/30 от ширины холста
	rightBottomPoint := image.Point{
		X: -1 * (fixed.I(rgba.Bounds().Max.X) - padding - fixed.I(watermark.Bounds().Dx())).Ceil(),
		Y: -1 * (fixed.I(rgba.Bounds().Max.Y) - padding - fixed.I(watermark.Bounds().Dy())).Ceil(),
	}
	// Odin: Наношу подпись поверх картины!
	draw.Draw(rgba, rgba.Bounds(), watermark, rightBottomPoint, draw.Over)
	// Odin: Работа готова! Восхваляйте Всеотца вашего за это чудо!
	return rgba
}

func (g *Gungner) prepareWatermark(bounds image.Rectangle, text string) image.Image {
	// Odin: создаю маленькую картинку с подписью, которая нанесётся поверх картины
	size := 86.0 // Odin: размер шрифта
	isSmallImage := bounds.Dx() <= model.WidthF
	if isSmallImage {
		size /= 2 // Odin: для картинок малых размеров (а это не оригинал картины, а коллаж множества) и размер подписи уменьшается
	}

	fontDrawer := font.Drawer{
		Src: image.NewUniform(color.RGBA{200, 200, 200, 255}), // серый шрифт
		Face: truetype.NewFace(g.font, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		}),
	}

	// Odin: рассчитываю границы текста
	textBounds, _ := fontDrawer.BoundString(text)
	textWidth := fontDrawer.MeasureString(text)
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	catHeight := textHeight

	// Odin: готовлю картинку милого котика (слева в подписи)
	catImgResized := image.NewRGBA(image.Rect(0, 0, catHeight.Ceil(), catHeight.Ceil()))
	draw.NearestNeighbor.Scale(catImgResized, catImgResized.Rect, g.catImg, g.catImg.Bounds(), draw.Over, nil)
	watermark := image.NewRGBA(image.Rect(
		0,
		0,
		(textWidth + textHeight*7/6).Ceil(), // Odin: запамятовал уже, зачем такие числа...
		textHeight.Ceil(),
	))
	// Odin: Подготовка завершена, начинаю компоновать подпись
	// Odin: О МОЙ GUNGNER, НАЧНЁМ НАШ ТРУД!
	// Odin: первым слоем на подпись накладывает полупрозрачный чёрный оттеняющий фон
	blackTransparentClr := color.RGBA{0, 0, 0, 128}
	draw.Draw(watermark, watermark.Bounds(), &image.Uniform{blackTransparentClr}, image.Point{}, draw.Src)
	// Odin: следом наношу текст на изображение
	fontDrawer.Dst = watermark
	fontDrawer.Dot = fixed.Point26_6{
		X: catHeight * 7 / 6, // Odin: Отклоняю своё копьё, пропуская место для изображения КОТА с отступом, и наношу надпись
		Y: fixed.I(watermark.Bounds().Max.Y),
	}
	fontDrawer.DrawString(text)
	// Odin: теперь на пустое место слева наношу КОТИКА верхним слоем
	draw.Draw(watermark, watermark.Bounds(), catImgResized, image.Point{X: 0, Y: 0}, draw.Over)
	// Odin: делу время, потехе час! эта часть работы окончена. вернёмся теперь к картине.
	return watermark
}

func (g *Gungner) loadResources() error {
	if g.font == nil {
		fontFile := "./files/font/conso.ttf"
		fontData, err := os.ReadFile(fontFile)
		if err != nil {
			return errors.Wrapf(err, "[gungner] ШРИФТ %s НЕ ЗАГРУЖЕН!", fontFile)
		}
		fontFace, err := freetype.ParseFont(fontData)
		if err != nil {
			return errors.Wrapf(err, "[gungner] ШРИФТ %s ПОВРЕЖДЁН!", fontFile)
		}
		g.font = fontFace
	}
	if g.catImg == nil {
		catImgData, err := os.ReadFile("./files/images/watermark.png")
		if err != nil {
			return errors.Wrapf(err, "[gungner] КОТ НЕ ЗАГРУЖЕН!")
		}
		r := bytes.NewReader(catImgData)
		catImg, err := png.Decode(r)
		if err != nil {
			return errors.Wrapf(err, "[gungner] КОТ ПОВРЕЖДЁН!")
		}
		g.catImg = catImg
	}
	return nil
}
