package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// Odin: все картины Artchitect делятся на группы - "единства"
// Odin: Единства соединяет моя супруга pantheon.Frigg (смотри код Frigg для понимания единств)

const (
	// Ранг единства
	Unity100K = 100000
	Unity10K  = 10000
	Unity1K   = 1000
	Unity100  = 100

	// UpdateInterval100K - Как часто будет обновляться коллаж единства. Раз в сколько картин.
	UpdateInterval100K = 1000
	UpdateInterval10K  = 100
	UpdateInterval1K   = 50
	UpdateInterval100  = 10

	UnityStateEmpty         = "empty"         // пустое единство. Только создано. Коллаж еще не создавался.
	UnityStateUnified       = "unified"       // окончательно сформированное единство, где уже все картины на писаны. Больше не изменяется.
	UnityStatePartial       = "partial"       // частично заполненное единство. Внутри него написаны еще не все картины. Коллаж частичный.
	UnityStateReunification = "reunification" // специальный статус, который указывает Архитектору перезаполнить единство
	// Когда коллаж единства нужно обновить, то ставится статус reunification и в следующем цикле Архитектор его пересоберёт.
)

type Unity struct {
	Mask      string `gorm:"primaryKey"`
	Rank      uint   // тип единства
	MinID     uint
	MaxID     uint
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
	Leads     string // массив ID картин, которые попали в коллаж в виде строки [100, 121, 110, 0, 130, 0, 0, 100...]. Нули - пустые места (на картине чёрным)
	Version   int    // при пересборке единства версия повышается (чтобы старые картинки не кешировались)
}

type UnityPile struct {
	db *gorm.DB
}

func NewUnityPile(db *gorm.DB) *UnityPile {
	return &UnityPile{db: db}
}

func (up *UnityPile) Get(ctx context.Context, mask string) (Unity, error) {
	var unity Unity
	err := up.db.WithContext(ctx).Where("mask = ?", mask).Limit(1).First(&unity).Error
	return unity, err
}

func (up *UnityPile) Create(ctx context.Context, mask string, rank, min, max uint) (Unity, error) {
	unity := Unity{
		Mask:    mask,
		Rank:    rank,
		MinID:   min,
		MaxID:   max,
		State:   UnityStateEmpty,
		Leads:   "[]",
		Version: 0,
	}
	err := up.db.Save(&unity).Error
	return unity, err
}

func (up *UnityPile) Save(ctx context.Context, unity Unity) (Unity, error) {
	err := up.db.Save(&unity).Error
	return unity, err
}
