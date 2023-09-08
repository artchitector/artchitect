package model

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Like struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ArtID     uint      `gorm:"uniqueIndex:user_card_uq" json:"-"`
	Art       Art       `json:"-"`
	UserID    uint      `gorm:"uniqueIndex:user_card_uq" json:"-"`
}

type LikePile struct {
	db *gorm.DB
}

func NewLikePile(db *gorm.DB) *LikePile {
	return &LikePile{db: db}
}

func (lp *LikePile) Get(ctx context.Context, userID uint, artID uint) (Like, error) {
	var like Like
	err := lp.db.WithContext(ctx).Where("user_id = ? and art_id = ?", userID, artID).Limit(1).First(&like).Error
	return like, err
}

func (lp *LikePile) GetList(ctx context.Context, userID uint) ([]Like, error) {
	var likes []Like
	err := lp.db.WithContext(ctx).Where("user_id = ?", userID).Find(&likes).Error
	return likes, err
}

func (lp *LikePile) Set(ctx context.Context, userID uint, artID uint, liked bool) error {
	var like Like
	if !liked {
		result := lp.db.WithContext(ctx).Where("user_id = ?", userID).Where("art_id = ?", artID).Delete(&like)
		if err := result.Error; err != nil {
			return err
		}
		if result.RowsAffected > 0 {
			log.Info().Msgf("[like_pile] УДАЛЁН ЛАЙК u=%d art=%d", userID, artID)
		} else {
			log.Info().Msgf("[like_pile] ЛАЙК НЕ БЫЛ УСТАНОВЛЕН, u=%d art=%d", userID, artID)
		}
		return nil
	} else {
		like = Like{
			ArtID:  artID,
			UserID: userID,
		}
		err := lp.db.Create(&like).Error
		return err
	}
}
