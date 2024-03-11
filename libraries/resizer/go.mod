module github.com/artchitector/artchitect/libraries/resizer

go 1.20

require (
	github.com/artchitector/artchitect/model v0.0.0-20230808134650-4a673e65f228
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rs/zerolog v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	gorm.io/gorm v1.25.2 // indirect
)

replace github.com/artchitector/artchitect/model => ../../model
