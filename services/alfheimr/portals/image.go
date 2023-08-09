package portals

import (
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type ImageRequest struct {
	ID   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"`
}

type ImagePortal struct {
	warehouse warehouse
}

func NewImagePortal(warehouse warehouse) *ImagePortal {
	return &ImagePortal{warehouse: warehouse}
}

func (ip *ImagePortal) HandleImage(c *gin.Context) {
	var request ImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	img, err := ip.warehouse.GetArtImage(c, request.ID, request.Size)
	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}
	c.Data(http.StatusOK, "image/jpeg", img)
}
