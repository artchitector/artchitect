package external

import "context"

type ai interface {
	GenerateImage(ctx context.Context, seed uint, prompt string) ([]byte, error)
}
