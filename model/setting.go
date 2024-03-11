package model

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const SettingOdinActive = "odin_active"
const OdinActive = "odin_active"
const OdinDisactive = "odin_disactive"

type Setting struct {
	SettingID string `gorm:"primaryKey"`
	Value     string
}

type SettingPile struct {
	db *gorm.DB
}

func NewSettingPile(db *gorm.DB) *SettingPile {
	return &SettingPile{db: db}
}

func (sp *SettingPile) SetValue(ctx context.Context, name string, value string) (Setting, error) {
	setting := Setting{
		SettingID: name,
		Value:     value,
	}
	if err := sp.db.Save(&setting).Error; err != nil {
		return Setting{}, fmt.Errorf("[setting] ПОПЫТКА СОХРАНИТЬ НАСТРОЙКУ ПРОВАЛЕНА")
	}
	return setting, nil
}

func (sp *SettingPile) GetValue(ctx context.Context, name string) (string, error) {
	var setting Setting
	err := sp.db.WithContext(ctx).Where("setting_id = ?", name).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("[setting] ПОПЫТКА ПОЛУЧИТЬ НАСТРОЙКУ %s ПРОВАЛИЛАСЬ: %w", name)
	}
	return setting.Value, nil
}
