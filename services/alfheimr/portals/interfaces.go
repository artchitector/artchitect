package portals

import (
	"context"

	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
)

type radio interface {
	ListenRadio(subscribeCtx context.Context) chan model.Radiogram
}

type authService interface {
	GetUserIdFromContext(c *gin.Context) uint
}

type artPile interface {
	GetArt(ctx context.Context, ID uint) (model.Art, error)
	GetArtRecursive(ctx context.Context, ID uint) (model.Art, error)
	GetLastArts(ctx context.Context, last uint) ([]model.Art, error)
	GetMaxArtID(ctx context.Context) (uint, error)
	GetArtsInterval(ctx context.Context, min, max uint) ([]model.Art, error)
}

type unityPile interface {
	Get(ctx context.Context, mask string) (model.Unity, error)
	GetRoot(ctx context.Context) ([]model.Unity, error)
	GetChildren(ctx context.Context, unity model.Unity) ([]model.Unity, error)
}

type likePile interface {
	GetList(ctx context.Context, userID uint) ([]model.Like, error)
	Get(ctx context.Context, userID uint, artID uint) (model.Like, error)
	Set(ctx context.Context, userID uint, artID uint, liked bool) error
}

type warehouse interface {
	DownloadArtImage(ctx context.Context, artID uint, size string) ([]byte, error)
	DownloadArtOrigin(ctx context.Context, artID uint) ([]byte, error)
	DownloadUnityImage(ctx context.Context, mask string, version uint, size string) ([]byte, error)
}

type harbour interface {
	SendCrownWaitCargo(ctx context.Context, request string) (string, error)
}
