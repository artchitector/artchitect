module github.com/artchitector/artchitect2/services/asgard

go 1.20

require (
	github.com/artchitector/artchitect2/model v0.0.0-20230808134650-4a673e65f228
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.30.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.2
)

require (
	github.com/artchitector/artchitect2/libraries/resizer v0.0.0-20230808145946-07f8c6ba781d // indirect
	github.com/blackjack/webcam v0.0.0-20230509180125-87693b3f29dc // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/artchitector/artchitect2/model => ../../model
replace github.com/artchitector/artchitect2/libraries/resizer => ../../libraries/resizer
