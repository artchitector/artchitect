package portals

import (
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type ArtImageRequest struct {
	ID   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"`
}

type UnityImageRequest struct {
	Mask    string `uri:"mask" binding:"required"`
	Version uint   `uri:"version" binding:"numeric"`
	Size    string `uri:"size" binding:"required"`
}

type ImagePortal struct {
	warehouse warehouse
}

func NewImagePortal(warehouse warehouse) *ImagePortal {
	return &ImagePortal{warehouse: warehouse}
}

func (ip *ImagePortal) HandleArtImage(c *gin.Context) {
	var request ArtImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	var img []byte
	var err error
	if request.Size == model.SizeOrigin {
		img, err = ip.warehouse.DownloadArtOrigin(c, request.ID)
	} else {
		img, err = ip.warehouse.DownloadArtImage(c, request.ID, request.Size)
	}

	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}
	c.Data(http.StatusOK, "image/jpeg", img)
}

func (ip *ImagePortal) HandleUnityImage(c *gin.Context) {
	var request UnityImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	img, err := ip.warehouse.DownloadUnityImage(c, request.Mask, request.Version, request.Size)
	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	c.Data(http.StatusOK, "image/jpeg", img)
}
