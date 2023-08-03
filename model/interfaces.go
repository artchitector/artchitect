package model

import "context"

type entropy interface {
	SelectOf(ctx context.Context, totalVariant uint) (uint, error)
}
