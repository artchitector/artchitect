module github.com/artchitector/artchitect/services/asgard

go 1.20

require (
	github.com/artchitector/artchitect/libraries/resizer v0.0.0-20230818151811-c8f78c5df613
	github.com/artchitector/artchitect/libraries/warehouse v0.0.0-20230819093634-3a18e65db5f3
	github.com/artchitector/artchitect/model v0.0.0-20230819093634-3a18e65db5f3
	github.com/blackjack/webcam v0.0.0-20230509180125-87693b3f29dc
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.30.0
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b
	golang.org/x/image v0.11.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.3
)

require (
	github.com/bytedance/sonic v1.10.0-rc3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.0 // indirect
	github.com/go-telegram/bot v0.7.15 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.4.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/artchitector/artchitect/model => ../../model

replace github.com/artchitector/artchitect/libraries/resizer => ../../libraries/resizer

replace github.com/artchitector/artchitect/libraries/warehouse => ../../libraries/warehouse
