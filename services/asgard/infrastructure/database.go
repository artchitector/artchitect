package infrastructure

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(ctx context.Context, dns string) *gorm.DB {
	pg := postgres.Open(dns)
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Warn, // logger.Info
		},
	)
	db, err := gorm.Open(pg, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		zlog.Fatal().Err(errors.Wrap(err, "[database] ОШИБКА ПОДКЛЮЧЕНИЯ"))
	}

	if err := db.AutoMigrate(
		&model.Art{},
		&model.Idea{},
		&model.Word{},
		&model.Unity{},
		&model.Like{},
	); err != nil {
		zlog.Fatal().Err(errors.Wrap(err, "[database] ОШИБКА МИГРАЦИИ"))
	}

	return db
}
