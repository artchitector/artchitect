package model

const (
	SizeF      = "f"      // 1024x1536
	SizeM      = "m"      // 512x768
	SizeS      = "s"      // 256x384
	SizeXS     = "xs"     // 128x192
	SizeOrigin = "origin" //исходник (jpeg-изображение в большом разрешении, качество 100%)

	HeightToWidth = float64(3.0 / 2.0) // Все картинки Artchitect в соотношении 2:3. Умножь это на Width и получишь Height
	WidthF        = 1024               // height=1536
	WidthM        = 512                // height=768
	WidthS        = 256                // height=384
	WidthXS       = 128                // height=192

	QualityTransfer = 100 // передача между серверами в jpeg со 100% качеством
	QualityF        = 90
	QualityM        = 80
	QualityS        = 75
	QualityXS       = 75
)
