package communication

import (
	"context"
	"image"
)

type webcam interface {
	GetNextFrame(ctx context.Context) (image.Image, error)
}
