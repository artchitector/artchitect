package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

// ### entities

type Art struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`

	Version     string `json:"version"`     // version of art-generation-algorithm (different dictionary or settings...)
	Seed        uint   `json:"seed"`        // entropy-generated seed number, which is seed for Stable Diffusion input
	SeedEntropy string `json:"seedEntropy"` // unique identifier of original entropy image

	TotalTime uint `json:"totalTime"` // seconds, how long whole generation process taken
	PaintTime uint `json:"paintTime"` // seconds, how long stable diffusion generate image

	Tags  []ArtTag `json:"tags"`  // set of tags
	Likes ArtLikes `json:"likes"` // current likes-summary
}

type ArtTag struct {
	ID        uint      `json:"-" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	ArtID     uint      `json:"-"`
	Keyword   string    `json:"keyword"` // keyword which will be passed into Stable Diffusion as part of prompt
	Entropy   string    `json:"entropy"` // unique identifier of original entropy image
}

type ArtLikes struct {
	ArtID uint `json:"-" gorm:"primaryKey"`
	Likes uint `json:"likes"` // likes total amount
}

// ### repository

type ArtRepository struct {
	db      *gorm.DB
	entropy entropy
}

func NewArtRepository(db *gorm.DB, entropy entropy) *ArtRepository {
	return &ArtRepository{db: db, entropy: entropy}
}

func (ar *ArtRepository) GetArt(ctx context.Context, ID uint) (Art, error) {
	return Art{}, errors.New("fake method GetArt")
}

func (ar *ArtRepository) GetMaxArtID(ctx context.Context) (uint, error) {
	var id uint
	err := ar.db.WithContext(ctx).Select("case when max(id) is null then 0 else max(id) end as max_id").Model(&Art{}).Scan(&id).Error
	return id, err
}

func (ar *ArtRepository) GetNextArtID(ctx context.Context) (uint, error) {
	id, err := ar.GetMaxArtID(ctx)
	return id + 1, err
}
