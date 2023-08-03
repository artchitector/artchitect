package internal

import (
	"context"
	"image"
)

type artist interface {
	MakeArt(
		ctx context.Context,
		artID uint,
		seed uint,
		tags []string,
	) (image.Image, error)
}

type artRepo interface {
	GetNextArtID(ctx context.Context) (uint, error)
}
