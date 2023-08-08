package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

const (
	// Version1 Odin: Версия алгоритма создания картины. Сейчас тут только одна.
	// Version1 Odin: в версии v1 тут у нас допотопный StableDiffusion v1.5, но он рисует ярко, пусть и не аккуратно
	Version1 = "v1"
)

// ### entities

type Art struct {
	// ID не автоинкрементное поле. Автоинкремент сделан в коде вручную.
	// Odin: Все номера картин должны идти подряд без пропусков, поэтому тут не используется sequence/autoincrement
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	// Version - Odin: может быть несколько вариантов генерации картины (разные словари, разные версии StableDiffusion)
	Version string `json:"version"`
	// Odin: это идея картины. Картина может быть воссоздана по этой идее на той же версии ИИ без изменений сколько угодно раз
	// Odin: саму картинку можно репродуцировать, а вот идея пришла "из космоса", и её не повторить
	// Loki: кстати, если в настройках ИИ поменять разрешение или другой параметр,
	// Loki: то из этой же идеи будет нарисована похожая, но другая картина
	// Odin: раз так, то общая уникальность изображения заключена в составном ключе seed+words+AIversion+AIsettings
	// Odin: используя идею и повторив эти настройки можно воссоздать ТУ ЖЕ картину и на другой машине
	Idea Idea `json:"idea"`

	TotalTime          uint `json:"totalTime"`          // мс, сколько заняло рисование картины от начала до конца
	IdeaGenerationTime uint `json:"ideaGenerationTime"` // мс. сколько заняла сборка идеи
	PaintTime          uint `json:"paintTime"`          // мс. сколько заняло рисование в недрах ИИ
}

// ### pile - куча - репозиторий

type ArtPile struct {
	db *gorm.DB
}

func NewArtPile(db *gorm.DB) *ArtPile {
	return &ArtPile{db: db}
}

func (ap *ArtPile) GetArt(ctx context.Context, ID uint) (Art, error) {
	return Art{}, errors.New("fake method GetArt")
}

func (ap *ArtPile) GetMaxArtID(ctx context.Context) (uint, error) {
	var id uint
	err := ap.db.WithContext(ctx).Select("case when max(id) is null then 0 else max(id) end as max_id").Model(&Art{}).Scan(&id).Error
	return id, err
}

func (ap *ArtPile) GetNextArtID(ctx context.Context) (uint, error) {
	id, err := ap.GetMaxArtID(ctx)
	return id + 1, err
}

func (ap *ArtPile) SaveArt(ctx context.Context, artID uint, art Art, idea Idea) (Art, error) {
	art.ID = artID
	idea.ArtID = artID
	if err := ap.db.Save(&art).Error; err != nil {
		return Art{}, err
	}
	if err := ap.db.Save(&idea).Error; err != nil {
		return Art{}, err
	}
	return art, nil
}
