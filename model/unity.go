package model

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
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
	UnityStatePreUnified    = "pre-unified"   // частично заполненное единство. Внутри него написаны еще не все картины. Коллаж частичный.
	UnityStateReunification = "reunification" // специальный статус, который указывает Архитектору перезаполнить единство
	// Когда коллаж единства нужно обновить, то ставится статус reunification и в следующем цикле Архитектор его пересоберёт.

	// размер сетки -  NxN элементов
	CollageSize100K = 7 * 7
	CollageSize10K  = 6 * 6
	CollageSize1K   = 5 * 5
	CollageSize100  = 4 * 4
)

type Unity struct {
	Mask      string    `gorm:"primaryKey" json:"mask"`
	Parent    string    `json:"parent"` // маска родительского единства
	Rank      uint      `json:"rank"`   // тип единства
	MinID     uint      `json:"minID"`
	MaxID     uint      `json:"maxID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	State     string    `json:"state"`
	Leads     string    `json:"leads"`   // массив ID картин, которые попали в коллаж в виде строки [100, 121, 110, 0, 130, 0, 0, 100...]. Нули - пустые места (на картине чёрным)
	Version   uint      `json:"version"` // при пересборке единства версия повышается (чтобы старые картинки не кешировались)
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

// GetRoot - получение всех корневых единств
func (up *UnityPile) GetRoot(ctx context.Context) ([]Unity, error) {
	var unities []Unity
	err := up.db.WithContext(ctx).Where("rank = ?", Unity100K).Order("mask ASC").Find(&unities).Error
	return unities, err
}

func (up *UnityPile) Create(ctx context.Context, mask, state string, rank, min, max uint) (Unity, error) {
	unity := Unity{
		Mask:    mask,
		Parent:  getParentMask(mask),
		Rank:    rank,
		MinID:   min,
		MaxID:   max,
		State:   state,
		Leads:   "[]",
		Version: 0,
	}
	err := up.db.WithContext(ctx).Save(&unity).Error
	return unity, err
}

func (up *UnityPile) Save(ctx context.Context, unity Unity) (Unity, error) {
	err := up.db.WithContext(ctx).Save(&unity).Error
	return unity, err
}

func (up *UnityPile) GetNextUnityForReunification(ctx context.Context) (Unity, error) {
	var unity Unity
	err := up.db.WithContext(ctx).
		Where("state = ?", UnityStateReunification).
		Order("rank DESC, mask ASC"). // Первыми в процесс объединения попадают большие единства
		Limit(1).
		First(&unity).
		Error
	return unity, err
}

func (up *UnityPile) GetChildren(ctx context.Context, unity Unity) ([]Unity, error) {
	var unities []Unity
	err := up.db.WithContext(ctx).
		Where("parent = ?", unity.Mask).
		Order("mask ASC").
		Find(&unities).
		Error
	return unities, err
}

// getParentMask - 91XXXX превращает в 9XXXXX (добавляется один X).
// Если родителя выше уже нет (верхний уровень единства), то возвращается пустая строка
func getParentMask(mask string) string {
	for i := len(mask) - 1; i >= 0; i-- {
		if string(mask[i]) == "X" {
			continue
		}
		if i == 0 { // Ведущее число маски не может быть заменено на X. Это уже самое верхнее единство
			return ""
		}
		return fmt.Sprintf("%sX%s", mask[:i], mask[i+1:])
	}
	log.Fatal().Msgf("[model:unity] ОШИБКА С ПРЕОБРАЗОВАНИЕМ МАСКИ %s В РОДИТЕЛЬСКУЮ", mask)
	return ""
}
