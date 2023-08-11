package model

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"math"
	"time"
)

const (
	// Version1 Odin: Версия алгоритма создания картины. Сейчас тут только одна.
	// Version1 Odin: в версии v1 тут у нас допотопный StableDiffusion v1.5, но он рисует ярко, пусть и не аккуратно
	Version1 = "v1"
)

// Art - сущность "картина"
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

/*
GetUnityMask - смотри документацию тут (/model/unity.go): Unity100K, Unity10K, Unity1K, Unity100.
Odin: Если картина #012345 будет сделана, то она попадёт в такие единства:
Unity100K = U0XXXXX (префикс U тут добавлен для наглядности, но в методе GetUnityMask не добавляется)
Unity10K = U01XXXX
Unity1K = U012XXX
Unity100 = U0123XX
Всего во множестве 6 символов, вплоть до картины #999999
*/
func (a Art) GetUnityMask(unityType int) string {
	if a.ID > 999999 {
		// Числа более 999999 не поддерживаются. Artchitect не собирается рисовать 1млн.
		// краш Artchitect. Дальше #999999 уже ничего не будет нарисовано, пока тут не исправить.
		log.Fatal().Msgf("[ART] ID %d IS TOO LONG", a.ID)
	}
	normalized := int(math.Floor(float64(a.ID) / float64(unityType)))
	switch unityType {
	case Unity100K:
		return fmt.Sprintf("%dXXXXX", normalized)
	case Unity10K:
		return fmt.Sprintf("%02dXXXX", normalized)
	case Unity1K:
		return fmt.Sprintf("%03dXXX", normalized)
	case Unity100:
		return fmt.Sprintf("%04dXX", normalized)
	default:
		log.Fatal().Msgf("[ART] UNITY TYPE %d UNKNOWN", unityType)
	}
	return ""
}

// ### pile - куча - репозиторий

type ArtPile struct {
	db *gorm.DB
}

func NewArtPile(db *gorm.DB) *ArtPile {
	return &ArtPile{db: db}
}

func (ap *ArtPile) GetArt(ctx context.Context, ID uint) (Art, error) {
	var art Art
	err := ap.db.WithContext(ctx).Where("id = ?", ID).Limit(1).First(&art).Error
	return art, err
}

func (ap *ArtPile) GetArtRecursive(ctx context.Context, ID uint) (Art, error) {
	var art Art
	err := ap.db.
		WithContext(ctx).
		Preload("Idea").
		Preload("Idea.Words").
		Where("id = ?", ID).
		Limit(1).
		First(&art).
		Error
	return art, err
}

func (ap *ArtPile) GetMaxArtID(ctx context.Context) (uint, error) {
	var id uint
	err := ap.db.
		WithContext(ctx).
		Select("case when max(id) is null then 0 else max(id) end as max_id").
		Model(&Art{}).
		Scan(&id).
		Error
	return id, err
}

func (ap *ArtPile) GetNextArtID(ctx context.Context) (uint, error) {
	id, err := ap.GetMaxArtID(ctx)
	return id + 1, err
}

func (ap *ArtPile) GetLastArts(ctx context.Context, last uint) ([]Art, error) {
	arts := make([]Art, 0, last)
	err := ap.db.WithContext(ctx).
		Preload("Idea").
		Order("id desc").
		Limit(int(last)).
		Find(&arts).Error
	return arts, err
}

func (ap *ArtPile) SaveArt(ctx context.Context, artID uint, art Art, idea Idea) (Art, error) {
	art.ID = artID
	idea.ArtID = artID
	db := ap.db.WithContext(ctx)
	if err := db.Save(&art).Error; err != nil {
		return Art{}, err
	}
	words := idea.Words
	idea.Words = nil
	if err := db.Save(&idea).Error; err != nil {
		return Art{}, err
	}
	for _, w := range words {
		w.IdeaID = artID
		if err := db.Save(&w).Error; err != nil {
			return Art{}, err
		}
	}

	idea.Words = words
	art.Idea = idea
	return art, nil
}

func (ap *ArtPile) GetLastPaintTime(ctx context.Context) (uint, error) {
	id, err := ap.GetMaxArtID(ctx)
	if err != nil {
		return 0, err
	}
	art, err := ap.GetArt(ctx, id)
	if err != nil {
		return 0, err
	}
	return art.PaintTime, nil
}
